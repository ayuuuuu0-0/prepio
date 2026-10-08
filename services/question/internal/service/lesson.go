package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/prepio/prepio/config"
	"github.com/prepio/prepio/services/question/internal/dto"
	"github.com/prepio/prepio/services/question/internal/grading"
	"github.com/prepio/prepio/services/question/internal/lesson"
	"github.com/prepio/prepio/services/question/internal/store"
	"github.com/prepio/prepio/shared/events"
)

// EventPublisher publishes domain events.
type EventPublisher interface {
	Publish(ctx context.Context, topic, key string, payload any) error
}

// LessonService owns the Journey domain's lesson runtime: the path, attempts,
// server-side grading, and completion. It is the only grader.
type LessonService struct {
	lessons   *store.LessonStore
	evaluator Evaluator
	publisher EventPublisher
}

// NewLessonService creates a LessonService. evaluator grades prose steps only.
func NewLessonService(lessons *store.LessonStore, evaluator Evaluator, publisher EventPublisher) *LessonService {
	return &LessonService{lessons: lessons, evaluator: evaluator, publisher: publisher}
}

// GetPath returns the learner's path: worlds, nodes in order, status, and lesson previews.
// Worlds whose topic is one of focus (the learner's focus topic slugs, in priority order)
// come first, so the single current node respects focus. Order never changes what is
// locked: unlocks depend only on prerequisites.
func (s *LessonService) GetPath(ctx context.Context, userID string, focus []string) (*dto.PathResponse, error) {
	if len(userID) == 0 {
		return nil, ErrInvalidRequest
	}
	nodes, err := s.lessons.ListPath(ctx, userID)
	if err != nil {
		return nil, err
	}
	nodes = orderByFocus(nodes, focus)

	byID := make(map[string]store.PathNode, len(nodes))
	for _, n := range nodes {
		byID[n.NodeID] = n
	}

	resp := &dto.PathResponse{Worlds: []dto.PathWorld{}}
	currentSet := false
	for _, n := range nodes {
		unmet := unmetPrerequisites(n, byID)
		status := nodeStatus(n, unmet, &currentSet)

		node := dto.PathNode{
			ID:         n.NodeID,
			Slug:       n.NodeSlug,
			Label:      n.NodeLabel,
			NodeType:   n.NodeType,
			Status:     status,
			InProgress: n.InProgress && !n.Completed,
			LessonID:   n.Lesson.ID,
			LessonSlug: n.Lesson.Slug,
			Title:      n.Lesson.Title,
			Takeaways:  n.Lesson.Summary,
			Kind:       n.Lesson.Kind,
			Difficulty: n.Lesson.Difficulty,
			EstMinutes: n.Lesson.EstMinutes,
			XPPreview:  config.LessonXPPreview(n.Lesson.Kind, n.Lesson.Difficulty),
		}
		if status == "locked" {
			node.UnlockHint = unlockHint(unmet)
		}

		last := len(resp.Worlds) - 1
		if last < 0 || resp.Worlds[last].ID != n.WorldID {
			resp.Worlds = append(resp.Worlds, dto.PathWorld{
				ID: n.WorldID, Slug: n.WorldSlug, Name: n.WorldName,
				Description: n.WorldDescription, Theme: n.WorldTheme,
				Topic: n.WorldTopic, Focused: focusRank(n.WorldTopic, focus) < len(focus),
				Nodes: []dto.PathNode{},
			})
			last++
		}
		resp.Worlds[last].Nodes = append(resp.Worlds[last].Nodes, node)
	}
	return resp, nil
}

// orderByFocus moves the nodes of worlds whose topic is a focus topic to the front, in
// focus priority order. Everything else keeps its authored order (the sort is stable, and
// nodes arrive grouped by world).
func orderByFocus(nodes []store.PathNode, focus []string) []store.PathNode {
	if len(focus) == 0 {
		return nodes
	}
	out := append([]store.PathNode(nil), nodes...)
	sort.SliceStable(out, func(i, j int) bool {
		return focusRank(out[i].WorldTopic, focus) < focusRank(out[j].WorldTopic, focus)
	})
	return out
}

// focusRank is the topic's position in focus, or len(focus) when it is not a focus topic.
func focusRank(topic string, focus []string) int {
	if topic != "" {
		for i, f := range focus {
			if f == topic {
				return i
			}
		}
	}
	return len(focus)
}

// unmetPrerequisites returns the nodes the user still has to finish before this one.
// A prerequisite that is not part of the published path cannot block anyone.
func unmetPrerequisites(n store.PathNode, byID map[string]store.PathNode) []store.PathNode {
	var unmet []store.PathNode
	for _, id := range n.Requires {
		req, ok := byID[id]
		if ok && !req.Completed {
			unmet = append(unmet, req)
		}
	}
	return unmet
}

func nodeStatus(n store.PathNode, unmet []store.PathNode, currentSet *bool) string {
	switch {
	case n.Completed:
		return "done"
	case len(unmet) > 0:
		return "locked"
	case !*currentSet:
		*currentSet = true
		return "current"
	default:
		return "available"
	}
}

func unlockHint(unmet []store.PathNode) string {
	names := make([]string, 0, len(unmet))
	for _, n := range unmet {
		names = append(names, n.NodeLabel)
	}
	switch len(names) {
	case 0:
		return ""
	case 1:
		return "Complete " + names[0] + " to unlock"
	default:
		return "Complete " + strings.Join(names[:len(names)-1], ", ") + " and " + names[len(names)-1] + " to unlock"
	}
}

// StartAttempt starts a lesson attempt, or resumes the one already in progress.
func (s *LessonService) StartAttempt(ctx context.Context, userID, lessonID string) (*dto.AttemptResponse, error) {
	if len(userID) == 0 || !validUUID(lessonID) {
		return nil, ErrLessonNotFound
	}
	l, err := s.lessons.GetLesson(ctx, lessonID)
	if err != nil {
		return nil, err
	}
	if l == nil || l.Status != lesson.StatusPublished {
		return nil, ErrLessonNotFound
	}

	attempt, err := s.lessons.GetInProgressAttempt(ctx, userID, l.ID)
	if err != nil {
		return nil, err
	}
	resumed := attempt != nil
	if !resumed {
		unlocked, err := s.lessons.PrerequisitesDone(ctx, userID, l.NodeID)
		if err != nil {
			return nil, err
		}
		if !unlocked {
			return nil, ErrLessonLocked
		}
		attempt, err = s.lessons.StartAttempt(ctx, userID, l.ID, l.Version)
		if err != nil {
			return nil, err
		}
		if attempt == nil {
			return nil, fmt.Errorf("start attempt: no attempt returned")
		}
	}

	steps, err := s.loadSteps(ctx, attempt.LessonID, attempt.LessonVersion)
	if err != nil {
		return nil, err
	}
	client := make([]lesson.ClientStep, 0, len(steps))
	for _, st := range steps {
		projected, err := lesson.Project(st.Key, st.Type, st.Position, st.payload, attempt.ID)
		if err != nil {
			return nil, err
		}
		client = append(client, projected)
	}

	results, err := s.lessons.ListResults(ctx, attempt.ID)
	if err != nil {
		return nil, err
	}

	return &dto.AttemptResponse{
		AttemptID: attempt.ID,
		Resumed:   resumed,
		Lesson:    lessonInfo(l),
		Steps:     client,
		Progress:  stepProgress(steps, results),
	}, nil
}

// SubmitAnswer grades one try of one step. Repeating the same (attempt, step, try)
// returns the stored result without grading again.
func (s *LessonService) SubmitAnswer(ctx context.Context, userID, attemptID, stepKey string, req dto.AnswerRequest) (*dto.AnswerResponse, error) {
	attempt, err := s.ownedAttempt(ctx, userID, attemptID)
	if err != nil {
		return nil, err
	}
	if attempt.Status != "in_progress" {
		return nil, ErrAttemptNotInProgress
	}
	if req.Try < 1 {
		return nil, ErrInvalidTry
	}

	steps, err := s.loadSteps(ctx, attempt.LessonID, attempt.LessonVersion)
	if err != nil {
		return nil, err
	}
	step, ok := findStep(steps, stepKey)
	if !ok {
		return nil, ErrStepNotFound
	}

	results, err := s.lessons.ListResults(ctx, attempt.ID)
	if err != nil {
		return nil, err
	}
	tries, solved := stepState(results, step.ID)

	answerJSON, err := json.Marshal(req.Answer)
	if err != nil {
		return nil, fmt.Errorf("encode answer: %w", err)
	}

	if req.Try <= tries {
		return s.replay(ctx, attempt.ID, step, req.Try, answerJSON, steps, results)
	}
	if req.Try != tries+1 || solved {
		return nil, ErrInvalidTry
	}

	resp, err := s.grade(step, req)
	if err != nil {
		return nil, err
	}
	resp.Try = req.Try
	resp.StepDone = resp.Correct

	feedback, err := json.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("encode feedback: %w", err)
	}
	inserted, err := s.lessons.InsertResult(ctx, attempt.ID, step.ID, req.Try, answerJSON, resp.Correct, feedback)
	if err != nil {
		return nil, err
	}
	if !inserted {
		// A concurrent request recorded this try first; its result is authoritative.
		return s.replay(ctx, attempt.ID, step, req.Try, answerJSON, steps, nil)
	}

	results, err = s.lessons.ListResults(ctx, attempt.ID)
	if err != nil {
		return nil, err
	}
	resp.AllDone = allGradedSolved(steps, results)
	return resp, nil
}

func (s *LessonService) replay(ctx context.Context, attemptID string, step loadedStep, try int, answerJSON []byte, steps []loadedStep, results []store.StepResult) (*dto.AnswerResponse, error) {
	stored, err := s.lessons.GetResult(ctx, attemptID, step.ID, try)
	if err != nil {
		return nil, err
	}
	if stored == nil {
		return nil, ErrInvalidTry
	}
	if !sameAnswer(stored.Answer, answerJSON) {
		return nil, ErrInvalidTry
	}
	var resp dto.AnswerResponse
	if err := json.Unmarshal(stored.Feedback, &resp); err != nil {
		return nil, fmt.Errorf("decode stored feedback: %w", err)
	}
	if results == nil {
		results, err = s.lessons.ListResults(ctx, attemptID)
		if err != nil {
			return nil, err
		}
	}
	resp.Replayed = true
	resp.AllDone = allGradedSolved(steps, results)
	return &resp, nil
}

// grade judges one answer. Deterministic types go to the grading package; prose goes to the evaluator.
func (s *LessonService) grade(step loadedStep, req dto.AnswerRequest) (*dto.AnswerResponse, error) {
	switch p := step.payload.(type) {
	case lesson.IntroPayload:
		return &dto.AnswerResponse{Correct: true}, nil
	case lesson.ProsePayload:
		return s.gradeProse(p, req.Answer)
	default:
		res, err := grading.Grade(step.payload, req.Answer)
		if errors.Is(err, grading.ErrInvalidAnswer) {
			return nil, fmt.Errorf("%w: %v", ErrInvalidAnswer, err)
		}
		if err != nil {
			return nil, err
		}
		return &dto.AnswerResponse{
			Correct:       res.Correct,
			Explanation:   res.Explanation,
			WhyNot:        res.WhyNot,
			CorrectAnswer: res.CorrectAnswer,
		}, nil
	}
}

func (s *LessonService) gradeProse(p lesson.ProsePayload, a grading.Answer) (*dto.AnswerResponse, error) {
	if a.Text == nil || len(strings.TrimSpace(*a.Text)) < p.MinChars {
		return nil, fmt.Errorf("%w: write at least %d characters", ErrInvalidAnswer, p.MinChars)
	}
	guide := AnswerGuide{Concepts: make([]Concept, 0, len(p.Rubric.Concepts))}
	for _, c := range p.Rubric.Concepts {
		guide.Concepts = append(guide.Concepts, Concept{Name: c.Name, Aliases: c.Aliases, Required: c.Required})
	}
	rawGuide, err := json.Marshal(guide)
	if err != nil {
		return nil, fmt.Errorf("encode rubric: %w", err)
	}
	eval := s.evaluator.Evaluate(*a.Text, string(rawGuide))
	score := eval.Score
	return &dto.AnswerResponse{
		Correct:     eval.Correct,
		Explanation: p.Explanation,
		Feedback:    eval.Summary,
		Score:       &score,
		Strengths:   eval.Strengths,
		Gaps:        eval.Gaps,
	}, nil
}

// Complete finishes an attempt. It succeeds only when every graded step has a
// correct result. It is idempotent: completing again returns the same summary and
// re-publishes lesson.completed only if the first publish never succeeded.
func (s *LessonService) Complete(ctx context.Context, userID, attemptID string) (*dto.CompleteResponse, error) {
	attempt, err := s.ownedAttempt(ctx, userID, attemptID)
	if err != nil {
		return nil, err
	}

	l, err := s.lessons.GetLesson(ctx, attempt.LessonID)
	if err != nil {
		return nil, err
	}
	if l == nil {
		return nil, ErrLessonNotFound
	}

	if attempt.Status == "in_progress" {
		steps, err := s.loadSteps(ctx, attempt.LessonID, attempt.LessonVersion)
		if err != nil {
			return nil, err
		}
		results, err := s.lessons.ListResults(ctx, attempt.ID)
		if err != nil {
			return nil, err
		}
		if !allGradedSolved(steps, results) {
			return nil, ErrAttemptIncomplete
		}
		graded, firstTry, totalTries := summarize(steps, results)

		completed, err := s.lessons.CompleteAttempt(ctx, attempt.ID, uuid.NewString(), graded, firstTry, totalTries)
		if err != nil {
			return nil, err
		}
		if completed == nil {
			// Another request completed it first; continue with the stored state.
			completed, err = s.lessons.GetAttempt(ctx, attempt.ID)
			if err != nil {
				return nil, err
			}
		}
		attempt = completed
	}

	if attempt.EventPublishedAt == nil {
		if err := s.publishCompleted(ctx, attempt, l); err != nil {
			return nil, err
		}
	}

	unlocked, err := s.unlockedBy(ctx, userID, l.NodeID)
	if err != nil {
		return nil, err
	}
	accuracy := 0.0
	if attempt.GradedSteps > 0 {
		accuracy = float64(attempt.FirstTryCorrect) / float64(attempt.GradedSteps)
	}
	return &dto.CompleteResponse{
		AttemptID:       attempt.ID,
		Lesson:          lessonInfo(l),
		Takeaways:       l.Summary,
		GradedSteps:     attempt.GradedSteps,
		FirstTryCorrect: attempt.FirstTryCorrect,
		TotalTries:      attempt.TotalTries,
		Accuracy:        accuracy,
		NodeID:          l.NodeID,
		UnlockedNodes:   unlocked,
	}, nil
}

func (s *LessonService) publishCompleted(ctx context.Context, attempt *store.Attempt, l *store.Lesson) error {
	skills, err := s.lessons.ListSkills(ctx, attempt.LessonID, attempt.LessonVersion)
	if err != nil {
		return err
	}
	weights := make([]events.LessonSkillWeight, 0, len(skills))
	for _, sk := range skills {
		weights = append(weights, events.LessonSkillWeight{SkillID: sk.SkillID, SkillSlug: sk.SkillSlug, Weight: sk.Weight})
	}

	completedAt := time.Now().UTC()
	if attempt.CompletedAt != nil {
		completedAt = attempt.CompletedAt.UTC()
	}
	event := events.LessonCompleted{
		EventID:         attempt.CompletionEventID,
		UserID:          attempt.UserID,
		LessonID:        attempt.LessonID,
		LessonSlug:      l.Slug,
		AttemptID:       attempt.ID,
		Kind:            l.Kind,
		Difficulty:      l.Difficulty,
		Skills:          weights,
		GradedSteps:     attempt.GradedSteps,
		FirstTryCorrect: attempt.FirstTryCorrect,
		TotalTries:      attempt.TotalTries,
		CompletedAt:     completedAt,
	}
	if err := s.publisher.Publish(ctx, events.TopicLessonCompleted, attempt.UserID, event); err != nil {
		return fmt.Errorf("publish lesson completed: %w", err)
	}
	return s.lessons.MarkEventPublished(ctx, attempt.ID)
}

// unlockedBy lists the nodes that depend on the completed node and are now open.
func (s *LessonService) unlockedBy(ctx context.Context, userID, nodeID string) ([]dto.UnlockedNode, error) {
	path, err := s.lessons.ListPath(ctx, userID)
	if err != nil {
		return nil, err
	}
	byID := make(map[string]store.PathNode, len(path))
	for _, n := range path {
		byID[n.NodeID] = n
	}
	out := []dto.UnlockedNode{}
	for _, n := range path {
		if n.Completed || !contains(n.Requires, nodeID) {
			continue
		}
		if len(unmetPrerequisites(n, byID)) > 0 {
			continue
		}
		out = append(out, dto.UnlockedNode{ID: n.NodeID, Slug: n.NodeSlug, Label: n.NodeLabel, LessonID: n.Lesson.ID})
	}
	return out, nil
}

func (s *LessonService) ownedAttempt(ctx context.Context, userID, attemptID string) (*store.Attempt, error) {
	if len(userID) == 0 || !validUUID(attemptID) {
		return nil, ErrAttemptNotFound
	}
	attempt, err := s.lessons.GetAttempt(ctx, attemptID)
	if err != nil {
		return nil, err
	}
	if attempt == nil || attempt.UserID != userID {
		return nil, ErrAttemptNotFound
	}
	return attempt, nil
}

type loadedStep struct {
	store.LessonStep
	payload any
}

func (s *LessonService) loadSteps(ctx context.Context, lessonID string, version int) ([]loadedStep, error) {
	rows, err := s.lessons.ListSteps(ctx, lessonID, version)
	if err != nil {
		return nil, err
	}
	steps := make([]loadedStep, 0, len(rows))
	for _, row := range rows {
		payload, err := lesson.ParsePayload(row.Type, row.Payload)
		if err != nil {
			return nil, fmt.Errorf("step %s: %w", row.Key, err)
		}
		steps = append(steps, loadedStep{LessonStep: row, payload: payload})
	}
	return steps, nil
}

func findStep(steps []loadedStep, key string) (loadedStep, bool) {
	for _, st := range steps {
		if st.Key == key {
			return st, true
		}
	}
	return loadedStep{}, false
}

// stepState returns how many tries a step has and whether any try was correct.
func stepState(results []store.StepResult, stepID string) (tries int, solved bool) {
	for _, r := range results {
		if r.StepID != stepID {
			continue
		}
		tries++
		if r.Correct {
			solved = true
		}
	}
	return tries, solved
}

func stepProgress(steps []loadedStep, results []store.StepResult) []dto.StepProgress {
	out := make([]dto.StepProgress, 0, len(steps))
	for _, st := range steps {
		tries, solved := stepState(results, st.ID)
		if tries == 0 {
			continue
		}
		out = append(out, dto.StepProgress{StepID: st.Key, Tries: tries, Done: solved})
	}
	return out
}

func allGradedSolved(steps []loadedStep, results []store.StepResult) bool {
	for _, st := range steps {
		if !lesson.IsGraded(st.Type) {
			continue
		}
		if _, solved := stepState(results, st.ID); !solved {
			return false
		}
	}
	return true
}

// summarize counts graded steps, steps correct on the first try, and total tries.
func summarize(steps []loadedStep, results []store.StepResult) (graded, firstTry, totalTries int) {
	graded = 0
	gradedIDs := map[string]bool{}
	for _, st := range steps {
		if lesson.IsGraded(st.Type) {
			graded++
			gradedIDs[st.ID] = true
		}
	}
	for _, r := range results {
		if !gradedIDs[r.StepID] {
			continue
		}
		totalTries++
		if r.Try == 1 && r.Correct {
			firstTry++
		}
	}
	return graded, firstTry, totalTries
}

func sameAnswer(stored, submitted []byte) bool {
	var a, b grading.Answer
	if err := json.Unmarshal(stored, &a); err != nil {
		return false
	}
	if err := json.Unmarshal(submitted, &b); err != nil {
		return false
	}
	ra, errA := json.Marshal(a)
	rb, errB := json.Marshal(b)
	return errA == nil && errB == nil && bytes.Equal(ra, rb)
}

func lessonInfo(l *store.Lesson) dto.LessonInfo {
	return dto.LessonInfo{ID: l.ID, Slug: l.Slug, Title: l.Title, Kind: l.Kind, Difficulty: l.Difficulty, EstMinutes: l.EstMinutes}
}

func validUUID(v string) bool {
	_, err := uuid.Parse(v)
	return err == nil
}

func contains(list []string, v string) bool {
	for _, item := range list {
		if item == v {
			return true
		}
	}
	return false
}
