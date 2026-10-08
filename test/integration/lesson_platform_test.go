package integration_test

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prepio/prepio/config"
	progresstest "github.com/prepio/prepio/services/progress/testing"
	questiontest "github.com/prepio/prepio/services/question/testing"
	streaktest "github.com/prepio/prepio/services/streak/testing"
	"github.com/prepio/prepio/shared/events"
	"github.com/prepio/prepio/test/fakes"
	"github.com/prepio/prepio/test/testdb"
	"github.com/prepio/prepio/test/testredis"
	"github.com/stretchr/testify/require"
)

// secretExplanation must never appear in a payload sent before the answer is graded.
const secretExplanation = "SECRET-EXPLANATION-SENTINEL"

func intp(v int) *int       { return &v }
func boolp(v bool) *bool    { return &v }
func strp(v string) *string { return &v }

// testContent builds fixtures for two chained lessons. Fixtures live only in tests.
func testContent() *questiontest.Content {
	return &questiontest.Content{
		Worlds: []questiontest.World{{
			Slug: "test-world", Name: "Test World", Description: "d", Theme: "t", Order: 1,
			Nodes: []questiontest.Node{
				{Slug: "n-start", Label: "Start Node", Type: "lesson"},
				{Slug: "n-next", Label: "Next Node", Type: "lesson", Requires: []string{"n-start"}},
			},
		}},
		Lessons: []questiontest.Lesson{
			{
				Slug: "lesson-one", Title: "Lesson One", Node: "n-start", Kind: "lesson", Status: questiontest.StatusPublished,
				Difficulty: "medium", EstMinutes: 4,
				Skills: []questiontest.SkillRef{
					{Skill: "system-design-scaling", Weight: 0.6},
					{Skill: "backend-databases", Weight: 0.4},
				},
				Summary: []string{"takeaway one", "takeaway two"},
				Steps: []questiontest.Step{
					{ID: "intro", Type: questiontest.StepIntro, Beats: []questiontest.Beat{{Text: "hello world", Emphasis: "world", DurationMs: 3000}}},
					{ID: "q1", Type: questiontest.StepMCQ, Prompt: "p1", Options: []string{"a", "b", "c"}, Answer: 1,
						Explanations: &questiontest.Explanations{Correct: secretExplanation + " correct", Wrong: map[int]string{0: "not a", 2: "not c"}}},
					{ID: "q2", Type: questiontest.StepTrueFalse, Statement: "s2", Answer: true, Explanation: secretExplanation + " tf"},
					{ID: "q3", Type: questiontest.StepFillBlank, Code: "x = ___1___", Bank: []string{"a", "b", "c"}, Answers: []string{"a"}, Explanation: secretExplanation + " fb"},
				},
			},
			{
				Slug: "lesson-two", Title: "Lesson Two", Node: "n-next", Kind: "lesson", Status: questiontest.StatusPublished,
				Difficulty: "easy", EstMinutes: 4,
				Skills:  []questiontest.SkillRef{{Skill: "system-design-scaling", Weight: 1.0}},
				Summary: []string{"takeaway one", "takeaway two"},
				Steps: []questiontest.Step{
					{ID: "m1", Type: questiontest.StepMCQ, Prompt: "p", Options: []string{"a", "b", "c"}, Answer: 0,
						Explanations: &questiontest.Explanations{Correct: "ok", Wrong: map[int]string{1: "no b", 2: "no c"}}},
					{ID: "a1", Type: questiontest.StepArrange, Prompt: "order", Items: []string{"c", "a", "b"}, Order: []int{1, 2, 0}, Explanation: "e"},
					{ID: "pr1", Type: questiontest.StepProse, Prompt: "explain", MinChars: 100, Explanation: "worked answer",
						Rubric: &questiontest.Rubric{Concepts: []questiontest.Concept{
							{Name: "cache", Required: true},
							{Name: "stale", Required: true},
						}}},
				},
			},
		},
	}
}

type failingOncePublisher struct {
	fakes.KafkaProducer
	failNext bool
}

func (p *failingOncePublisher) Publish(ctx context.Context, topic, key string, payload any) error {
	if p.failNext {
		p.failNext = false
		return errors.New("broker unavailable")
	}
	return p.KafkaProducer.Publish(ctx, topic, key, payload)
}

func lessonEvents(p *fakes.KafkaProducer) []events.LessonCompleted {
	var out []events.LessonCompleted
	for _, m := range p.Messages {
		if m.Topic != events.TopicLessonCompleted {
			continue
		}
		var e events.LessonCompleted
		if err := json.Unmarshal(m.Payload, &e); err != nil {
			panic(err)
		}
		out = append(out, e)
	}
	return out
}

func newUser(t *testing.T, pool *pgxpool.Pool, name string) string {
	t.Helper()
	var id string
	require.NoError(t, pool.QueryRow(context.Background(), `
		INSERT INTO users (email, username, password_hash, timezone)
		VALUES ($1, $2, 'hash', 'Asia/Kolkata') RETURNING id`, name+"@test.com", name).Scan(&id))
	return id
}

func TestLessonPlatform(t *testing.T) {
	pool, _ := testdb.Start(t)
	testdb.Migrate(t, pool)
	redisClient, _ := testredis.New(t)
	ctx := context.Background()

	publisher := &failingOncePublisher{}
	lessons := questiontest.NewLessonService(pool, publisher)
	sync := questiontest.NewContentSync(pool)
	content := testContent()

	catalog, err := sync.Catalog(ctx)
	require.NoError(t, err)
	require.Empty(t, questiontest.Validate(content, catalog), "fixture content must be valid")

	userID := newUser(t, pool, "learner")
	otherID := newUser(t, pool, "intruder")

	t.Run("content sync is idempotent and never duplicates", func(t *testing.T) {
		first, err := sync.Sync(ctx, content)
		require.NoError(t, err)
		require.Equal(t, 2, first.LessonsCreated)
		require.Equal(t, 0, first.LessonsNewVersion)

		second, err := sync.Sync(ctx, content)
		require.NoError(t, err)
		require.Equal(t, 0, second.LessonsCreated)
		require.Equal(t, 0, second.LessonsNewVersion)
		require.Equal(t, 2, second.LessonsUnchanged)

		var lessonsN, stepsN, skillsN int
		require.NoError(t, pool.QueryRow(ctx, `SELECT count(*) FROM lessons`).Scan(&lessonsN))
		require.NoError(t, pool.QueryRow(ctx, `SELECT count(*) FROM lesson_steps`).Scan(&stepsN))
		require.NoError(t, pool.QueryRow(ctx, `SELECT count(*) FROM lesson_skills`).Scan(&skillsN))
		require.Equal(t, 2, lessonsN)
		require.Equal(t, 7, stepsN)
		require.Equal(t, 3, skillsN)
	})

	var lessonOneID, lessonTwoID string

	t.Run("path shows status, previews, and exactly what unlocks a locked node", func(t *testing.T) {
		path, err := lessons.GetPath(ctx, userID, nil)
		require.NoError(t, err)
		require.Len(t, path.Worlds, 1, "the legacy deprecated world must not appear")
		nodes := path.Worlds[0].Nodes
		require.Len(t, nodes, 2)

		require.Equal(t, "current", nodes[0].Status)
		require.Equal(t, "Lesson One", nodes[0].Title)
		require.Equal(t, []string{"takeaway one", "takeaway two"}, nodes[0].Takeaways)
		require.Equal(t, 4, nodes[0].EstMinutes)
		require.Equal(t, config.LessonXPPreview("lesson", "medium"), nodes[0].XPPreview)

		require.Equal(t, "locked", nodes[1].Status)
		require.Equal(t, "Complete Start Node to unlock", nodes[1].UnlockHint)
		lessonOneID, lessonTwoID = nodes[0].LessonID, nodes[1].LessonID
	})

	t.Run("a locked lesson cannot be started", func(t *testing.T) {
		_, err := lessons.StartAttempt(ctx, userID, lessonTwoID)
		require.ErrorIs(t, err, questiontest.ErrLessonLocked)
	})

	t.Run("unknown lesson ids are not found", func(t *testing.T) {
		_, err := lessons.StartAttempt(ctx, userID, "not-a-uuid")
		require.ErrorIs(t, err, questiontest.ErrLessonNotFound)
		_, err = lessons.StartAttempt(ctx, userID, "00000000-0000-4000-8000-000000000000")
		require.ErrorIs(t, err, questiontest.ErrLessonNotFound)
	})

	var attemptID string

	t.Run("client payload contains no answers", func(t *testing.T) {
		attempt, err := lessons.StartAttempt(ctx, userID, lessonOneID)
		require.NoError(t, err)
		require.False(t, attempt.Resumed)
		require.Len(t, attempt.Steps, 4)
		attemptID = attempt.AttemptID

		raw, err := json.Marshal(attempt)
		require.NoError(t, err)
		body := string(raw)
		require.NotContains(t, body, secretExplanation)
		for _, leaked := range []string{`"answer"`, `"answers"`, `"explanations"`, `"explanation"`, `"rubric"`, `"order"`, "not a", "not c"} {
			require.NotContains(t, body, leaked)
		}
	})

	t.Run("starting again resumes the same attempt", func(t *testing.T) {
		again, err := lessons.StartAttempt(ctx, userID, lessonOneID)
		require.NoError(t, err)
		require.True(t, again.Resumed)
		require.Equal(t, attemptID, again.AttemptID)
	})

	t.Run("other users cannot touch the attempt", func(t *testing.T) {
		_, err := lessons.SubmitAnswer(ctx, otherID, attemptID, "q1", questiontest.AnswerRequest{Try: 1, Answer: questiontest.Answer{Choice: intp(1)}})
		require.ErrorIs(t, err, questiontest.ErrAttemptNotFound)
		_, err = lessons.Complete(ctx, otherID, attemptID)
		require.ErrorIs(t, err, questiontest.ErrAttemptNotFound)
	})

	answer := func(step string, try int, a questiontest.Answer) (*questiontest.AnswerResponse, error) {
		return lessons.SubmitAnswer(ctx, userID, attemptID, step, questiontest.AnswerRequest{Try: try, Answer: a})
	}

	t.Run("invalid answers are rejected and not recorded", func(t *testing.T) {
		_, err := answer("q1", 1, questiontest.Answer{Choice: intp(9)})
		require.ErrorIs(t, err, questiontest.ErrInvalidAnswer)
		_, err = answer("q1", 1, questiontest.Answer{})
		require.ErrorIs(t, err, questiontest.ErrInvalidAnswer)
		_, err = answer("q3", 1, questiontest.Answer{Blanks: []string{"a", "b"}})
		require.ErrorIs(t, err, questiontest.ErrInvalidAnswer)
		_, err = answer("nope", 1, questiontest.Answer{Choice: intp(0)})
		require.ErrorIs(t, err, questiontest.ErrStepNotFound)

		var recorded int
		require.NoError(t, pool.QueryRow(ctx, `SELECT count(*) FROM lesson_step_results WHERE attempt_id = $1`, attemptID).Scan(&recorded))
		require.Zero(t, recorded)
	})

	t.Run("try numbers must be sequential and start at 1", func(t *testing.T) {
		_, err := answer("q1", 0, questiontest.Answer{Choice: intp(1)})
		require.ErrorIs(t, err, questiontest.ErrInvalidTry)
		_, err = answer("q1", 2, questiontest.Answer{Choice: intp(1)})
		require.ErrorIs(t, err, questiontest.ErrInvalidTry)
	})

	t.Run("a wrong answer reveals the correct one, and a replayed try is idempotent", func(t *testing.T) {
		wrong, err := answer("q1", 1, questiontest.Answer{Choice: intp(0)})
		require.NoError(t, err)
		require.False(t, wrong.Correct)
		require.False(t, wrong.StepDone)
		require.Equal(t, "not a", wrong.WhyNot)
		require.Contains(t, wrong.Explanation, secretExplanation)
		require.NotNil(t, wrong.CorrectAnswer)
		require.Equal(t, 1, *wrong.CorrectAnswer.Choice)

		replay, err := answer("q1", 1, questiontest.Answer{Choice: intp(0)})
		require.NoError(t, err)
		require.True(t, replay.Replayed)
		require.False(t, replay.Correct)
		require.Equal(t, wrong.Explanation, replay.Explanation)

		_, err = answer("q1", 1, questiontest.Answer{Choice: intp(1)})
		require.ErrorIs(t, err, questiontest.ErrInvalidTry, "same try with a different answer is a conflict, not a regrade")

		var rows int
		require.NoError(t, pool.QueryRow(ctx, `SELECT count(*) FROM lesson_step_results WHERE attempt_id = $1`, attemptID).Scan(&rows))
		require.Equal(t, 1, rows, "replays must not add rows")

		right, err := answer("q1", 2, questiontest.Answer{Choice: intp(1)})
		require.NoError(t, err)
		require.True(t, right.Correct)
		require.True(t, right.StepDone)
		require.Nil(t, right.CorrectAnswer, "the correct answer is only revealed after an incorrect answer")

		_, err = answer("q1", 3, questiontest.Answer{Choice: intp(1)})
		require.ErrorIs(t, err, questiontest.ErrInvalidTry, "a solved step takes no more tries")
	})

	t.Run("completion requires every graded step to be correct", func(t *testing.T) {
		_, err := lessons.Complete(ctx, userID, attemptID)
		require.ErrorIs(t, err, questiontest.ErrAttemptIncomplete)

		got, err := answer("q2", 1, questiontest.Answer{Value: boolp(false)})
		require.NoError(t, err)
		require.False(t, got.Correct)
		require.True(t, *got.CorrectAnswer.Value)
		got, err = answer("q2", 2, questiontest.Answer{Value: boolp(true)})
		require.NoError(t, err)
		require.True(t, got.Correct)

		_, err = lessons.Complete(ctx, userID, attemptID)
		require.ErrorIs(t, err, questiontest.ErrAttemptIncomplete, "q3 is still unanswered")
	})

	t.Run("a resumed attempt reports per-step progress", func(t *testing.T) {
		resumed, err := lessons.StartAttempt(ctx, userID, lessonOneID)
		require.NoError(t, err)
		require.True(t, resumed.Resumed)
		byStep := map[string]int{}
		for _, p := range resumed.Progress {
			require.True(t, p.Done)
			byStep[p.StepID] = p.Tries
		}
		require.Equal(t, map[string]int{"q1": 2, "q2": 2}, byStep)
	})

	var completion *questiontest.CompleteResp
	var firstEvent events.LessonCompleted

	t.Run("completing publishes lesson.completed once and is idempotent", func(t *testing.T) {
		last, err := answer("q3", 1, questiontest.Answer{Blanks: []string{"a"}})
		require.NoError(t, err)
		require.True(t, last.Correct)
		require.True(t, last.AllDone)

		completion, err = lessons.Complete(ctx, userID, attemptID)
		require.NoError(t, err)
		require.Equal(t, 3, completion.GradedSteps)
		require.Equal(t, 1, completion.FirstTryCorrect, "only q3 was right on the first try")
		require.Equal(t, 5, completion.TotalTries)
		require.InDelta(t, 1.0/3.0, completion.Accuracy, 0.001)
		require.Len(t, completion.UnlockedNodes, 1)
		require.Equal(t, "n-next", completion.UnlockedNodes[0].Slug)

		again, err := lessons.Complete(ctx, userID, attemptID)
		require.NoError(t, err)
		require.Equal(t, completion, again)

		published := lessonEvents(&publisher.KafkaProducer)
		require.Len(t, published, 1, "completion must publish exactly once")
		firstEvent = published[0]
		require.Equal(t, userID, firstEvent.UserID)
		require.Equal(t, attemptID, firstEvent.AttemptID)
		require.Equal(t, "lesson-one", firstEvent.LessonSlug)
		require.Equal(t, 1, firstEvent.FirstTryCorrect)
		require.Len(t, firstEvent.Skills, 2)

		_, err = answer("q3", 2, questiontest.Answer{Blanks: []string{"a"}})
		require.Error(t, err, "a completed attempt takes no more answers")
	})

	t.Run("path reflects completion and unlocks the next node", func(t *testing.T) {
		path, err := lessons.GetPath(ctx, userID, nil)
		require.NoError(t, err)
		nodes := path.Worlds[0].Nodes
		require.Equal(t, "done", nodes[0].Status)
		require.Equal(t, "current", nodes[1].Status)
		require.Empty(t, nodes[1].UnlockHint)
	})

	progressLessons := progresstest.NewLessonService(pool, &fakes.KafkaProducer{})

	t.Run("progress writes mastery, ledger, and XP once per attempt", func(t *testing.T) {
		require.NoError(t, progressLessons.ProcessLessonCompleted(ctx, firstEvent))

		rewards, found, err := progressLessons.GetAttemptRewards(ctx, userID, attemptID)
		require.NoError(t, err)
		require.True(t, found)
		require.True(t, rewards.FirstCompletion)

		accuracy := 1.0 / 3.0
		require.Equal(t, config.LessonXP("lesson", "medium", accuracy), rewards.XPAwarded)
		require.Equal(t, config.LessonGems("medium", accuracy), rewards.GemsAwarded)
		require.Len(t, rewards.MasteryChanges, 2)

		bySkill := map[string]int{}
		for _, m := range rewards.MasteryChanges {
			bySkill[m.SkillSlug] = m.Delta
			require.Equal(t, 0, m.Before)
			require.Equal(t, m.Delta, m.After)
			require.NotEmpty(t, m.TopicSlug, "every changed skill belongs to a topic")
		}
		require.Equal(t, config.LessonMasteryDelta(0, "lesson", "medium", accuracy, 0.6), bySkill["system-design-scaling"])
		require.Equal(t, config.LessonMasteryDelta(0, "lesson", "medium", accuracy, 0.4), bySkill["backend-databases"])
		require.Positive(t, bySkill["system-design-scaling"])

		var mastery, ledger, xpRows int
		require.NoError(t, pool.QueryRow(ctx, `
			SELECT mastery FROM user_skill_scores u JOIN skills s ON s.id = u.skill_id
			WHERE u.user_id = $1 AND s.slug = 'system-design-scaling'`, userID).Scan(&mastery))
		require.Equal(t, bySkill["system-design-scaling"], mastery)
		require.NoError(t, pool.QueryRow(ctx, `SELECT count(*) FROM mastery_ledger WHERE user_id = $1`, userID).Scan(&ledger))
		require.Equal(t, 2, ledger)
		require.NoError(t, pool.QueryRow(ctx, `SELECT count(*) FROM xp_ledger WHERE user_id = $1`, userID).Scan(&xpRows))
		require.Equal(t, 1, xpRows)

		// Redelivery of the same event changes nothing.
		require.NoError(t, progressLessons.ProcessLessonCompleted(ctx, firstEvent))
		require.NoError(t, pool.QueryRow(ctx, `SELECT count(*) FROM mastery_ledger WHERE user_id = $1`, userID).Scan(&ledger))
		require.Equal(t, 2, ledger)
		require.NoError(t, pool.QueryRow(ctx, `SELECT count(*) FROM xp_ledger WHERE user_id = $1`, userID).Scan(&xpRows))
		require.Equal(t, 1, xpRows)
		var totalXP int
		require.NoError(t, pool.QueryRow(ctx, `SELECT total_xp FROM user_progress WHERE user_id = $1`, userID).Scan(&totalXP))
		require.Equal(t, rewards.XPAwarded, totalXP)

		// Rewards are private to the user who earned them.
		_, found, err = progressLessons.GetAttemptRewards(ctx, otherID, attemptID)
		require.NoError(t, err)
		require.False(t, found)
	})

	t.Run("replaying a completed lesson grants no further mastery or XP", func(t *testing.T) {
		replay, err := lessons.StartAttempt(ctx, userID, lessonOneID)
		require.NoError(t, err)
		require.False(t, replay.Resumed)
		require.NotEqual(t, attemptID, replay.AttemptID)

		for _, s := range []struct {
			step string
			a    questiontest.Answer
		}{
			{"q1", questiontest.Answer{Choice: intp(1)}},
			{"q2", questiontest.Answer{Value: boolp(true)}},
			{"q3", questiontest.Answer{Blanks: []string{"a"}}},
		} {
			_, err := lessons.SubmitAnswer(ctx, userID, replay.AttemptID, s.step, questiontest.AnswerRequest{Try: 1, Answer: s.a})
			require.NoError(t, err)
		}
		_, err = lessons.Complete(ctx, userID, replay.AttemptID)
		require.NoError(t, err)

		published := lessonEvents(&publisher.KafkaProducer)
		require.Len(t, published, 2)
		require.NoError(t, progressLessons.ProcessLessonCompleted(ctx, published[1]))

		rewards, found, err := progressLessons.GetAttemptRewards(ctx, userID, replay.AttemptID)
		require.NoError(t, err)
		require.True(t, found)
		require.False(t, rewards.FirstCompletion)
		require.Zero(t, rewards.XPAwarded)
		require.Empty(t, rewards.MasteryChanges)

		var ledger int
		require.NoError(t, pool.QueryRow(ctx, `SELECT count(*) FROM mastery_ledger WHERE user_id = $1`, userID).Scan(&ledger))
		require.Equal(t, 2, ledger)
	})

	t.Run("streak counts a completed lesson", func(t *testing.T) {
		streaks := streaktest.NewService(pool, redisClient, &fakes.KafkaProducer{}, &fakeGems{balance: 100})
		require.NoError(t, streaks.ProcessLessonCompleted(ctx, firstEvent))
		streak, err := streaks.GetMe(ctx, userID, "Asia/Kolkata")
		require.NoError(t, err)
		require.Equal(t, 1, streak.CurrentStreak)
	})

	t.Run("prose and arrange steps are graded server-side", func(t *testing.T) {
		attempt, err := lessons.StartAttempt(ctx, userID, lessonTwoID)
		require.NoError(t, err)
		submit := func(step string, try int, a questiontest.Answer) (*questiontest.AnswerResponse, error) {
			return lessons.SubmitAnswer(ctx, userID, attempt.AttemptID, step, questiontest.AnswerRequest{Try: try, Answer: a})
		}

		_, err = submit("a1", 1, questiontest.Answer{Order: []int{0, 0, 1}})
		require.ErrorIs(t, err, questiontest.ErrInvalidAnswer)
		bad, err := submit("a1", 1, questiontest.Answer{Order: []int{0, 1, 2}})
		require.NoError(t, err)
		require.False(t, bad.Correct)
		require.Equal(t, []int{1, 2, 0}, bad.CorrectAnswer.Order)
		good, err := submit("a1", 2, questiontest.Answer{Order: []int{1, 2, 0}})
		require.NoError(t, err)
		require.True(t, good.Correct)

		_, err = submit("pr1", 1, questiontest.Answer{Text: strp("too short")})
		require.ErrorIs(t, err, questiontest.ErrInvalidAnswer)
		_, err = submit("pr1", 1, questiontest.Answer{})
		require.ErrorIs(t, err, questiontest.ErrInvalidAnswer)

		filler := strings.Repeat("lorem ipsum dolor sit amet ", 5)
		missing, err := submit("pr1", 1, questiontest.Answer{Text: strp(filler)})
		require.NoError(t, err)
		require.False(t, missing.Correct, "an answer covering none of the rubric must not pass")
		require.NotNil(t, missing.Score)

		hit, err := submit("pr1", 2, questiontest.Answer{Text: strp("A cache keeps data close, but entries can go stale until they expire. " + filler)})
		require.NoError(t, err)
		require.True(t, hit.Correct)
		require.Contains(t, hit.Strengths, "cache")

		_, err = submit("m1", 1, questiontest.Answer{Choice: intp(0)})
		require.NoError(t, err)
	})

	t.Run("a failed publish is retried on the next completion without losing the event", func(t *testing.T) {
		attempt, err := lessons.StartAttempt(ctx, userID, lessonTwoID)
		require.NoError(t, err)
		require.True(t, attempt.Resumed)

		publisher.failNext = true
		_, err = lessons.Complete(ctx, userID, attempt.AttemptID)
		require.Error(t, err)
		require.Len(t, lessonEvents(&publisher.KafkaProducer), 2, "nothing was published while the broker was down")

		done, err := lessons.Complete(ctx, userID, attempt.AttemptID)
		require.NoError(t, err)
		require.Equal(t, 3, done.GradedSteps)
		published := lessonEvents(&publisher.KafkaProducer)
		require.Len(t, published, 3)

		again, err := lessons.Complete(ctx, userID, attempt.AttemptID)
		require.NoError(t, err)
		require.Equal(t, done, again)
		require.Len(t, lessonEvents(&publisher.KafkaProducer), 3, "published exactly once after recovery")

		var published2 events.LessonCompleted = published[2]
		require.NoError(t, progressLessons.ProcessLessonCompleted(ctx, published2))
		rewards, found, err := progressLessons.GetAttemptRewards(ctx, userID, attempt.AttemptID)
		require.NoError(t, err)
		require.True(t, found)
		require.True(t, rewards.FirstCompletion)
		require.Len(t, rewards.MasteryChanges, 1)
	})

	t.Run("changed content gets a new version while in-flight attempts keep theirs", func(t *testing.T) {
		other := newUser(t, pool, "versioned")
		// Unlock lesson two for this user by finishing lesson one.
		a1, err := lessons.StartAttempt(ctx, other, lessonOneID)
		require.NoError(t, err)
		for _, s := range []struct {
			step string
			a    questiontest.Answer
		}{
			{"q1", questiontest.Answer{Choice: intp(1)}},
			{"q2", questiontest.Answer{Value: boolp(true)}},
			{"q3", questiontest.Answer{Blanks: []string{"a"}}},
		} {
			_, err := lessons.SubmitAnswer(ctx, other, a1.AttemptID, s.step, questiontest.AnswerRequest{Try: 1, Answer: s.a})
			require.NoError(t, err)
		}
		_, err = lessons.Complete(ctx, other, a1.AttemptID)
		require.NoError(t, err)

		inflight, err := lessons.StartAttempt(ctx, other, lessonTwoID)
		require.NoError(t, err)

		changed := testContent()
		changed.Lessons[1].Steps[1].Order = []int{2, 0, 1} // arrange answer changes in v2
		report, err := sync.Sync(ctx, changed)
		require.NoError(t, err)
		require.Equal(t, 1, report.LessonsNewVersion)
		require.Equal(t, 1, report.LessonsUnchanged)

		// The in-flight attempt still grades against version 1.
		old, err := lessons.SubmitAnswer(ctx, other, inflight.AttemptID, "a1", questiontest.AnswerRequest{Try: 1, Answer: questiontest.Answer{Order: []int{1, 2, 0}}})
		require.NoError(t, err)
		require.True(t, old.Correct)

		var attemptVersion, lessonVersion int
		require.NoError(t, pool.QueryRow(ctx, `SELECT lesson_version FROM lesson_attempts WHERE id = $1`, inflight.AttemptID).Scan(&attemptVersion))
		require.NoError(t, pool.QueryRow(ctx, `SELECT version FROM lessons WHERE id = $1`, lessonTwoID).Scan(&lessonVersion))
		require.Equal(t, 1, attemptVersion)
		require.Equal(t, 2, lessonVersion)

		var versions int
		require.NoError(t, pool.QueryRow(ctx, `SELECT count(DISTINCT lesson_version) FROM lesson_steps WHERE lesson_id = $1`, lessonTwoID).Scan(&versions))
		require.Equal(t, 2, versions, "old versions are kept, never deleted")
	})
}

type fakeGems struct{ balance int }

func (f *fakeGems) DeductGems(_ context.Context, _ string, amount int, _ string) error {
	if f.balance < amount {
		return errors.New("insufficient gems")
	}
	f.balance -= amount
	return nil
}
