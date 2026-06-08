package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/prepio/prepio/config"
	"github.com/prepio/prepio/constants"
	"github.com/prepio/prepio/services/question/internal/dto"
	"github.com/prepio/prepio/services/question/internal/store"
	"github.com/prepio/prepio/shared/events"
	"github.com/redis/go-redis/v9"
)

// EventPublisher publishes domain events to Kafka.
type EventPublisher interface {
	Publish(ctx context.Context, topic, key string, payload any) error
}

// QuestionService handles daily papers and answer submission.
type QuestionService struct {
	questions   *store.QuestionStore
	dailyPapers *store.DailyPaperStore
	history     *store.HistoryStore
	journey     *store.JourneyStore
	users       *store.UserStore
	redis       *redis.Client
	publisher   EventPublisher
}

// NewQuestionService creates a QuestionService.
func NewQuestionService(
	questions *store.QuestionStore,
	dailyPapers *store.DailyPaperStore,
	history *store.HistoryStore,
	journey *store.JourneyStore,
	users *store.UserStore,
	redisClient *redis.Client,
	publisher EventPublisher,
) *QuestionService {
	return &QuestionService{
		questions:   questions,
		dailyPapers: dailyPapers,
		history:     history,
		journey:     journey,
		users:       users,
		redis:       redisClient,
		publisher:   publisher,
	}
}

// GetDailyPaper returns or generates today's paper for the user.
func (s *QuestionService) GetDailyPaper(ctx context.Context, userID, timezone string) (*dto.DailyPaperResponse, error) {
	if len(userID) == 0 {
		return nil, ErrInvalidRequest
	}
	if len(timezone) == 0 {
		stored, err := s.users.Timezone(ctx, userID)
		if err != nil {
			return nil, err
		}
		timezone = stored
		if len(timezone) == 0 {
			timezone = constants.DefaultTimezone
		}
	}

	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return nil, fmt.Errorf("load timezone: %w", err)
	}

	now := time.Now().In(loc)
	paperDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	weekendOnly := now.Weekday() == time.Saturday || now.Weekday() == time.Sunday

	existing, questions, err := s.dailyPapers.GetByUserAndDate(ctx, userID, paperDate)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return toDailyPaperResponse(existing, questions), nil
	}

	level, err := s.users.Level(ctx, userID)
	if err != nil {
		return nil, err
	}
	difficulty := config.DifficultyForLevel(level)

	selected, err := s.selectQuestions(ctx, userID, difficulty, weekendOnly)
	if err != nil {
		return nil, err
	}
	if len(selected) == 0 {
		return nil, fmt.Errorf("no questions available")
	}

	sessionID := uuid.NewString()
	paper, err := s.dailyPapers.Create(ctx, userID, sessionID, paperDate, selected)
	if err != nil {
		return nil, err
	}

	sessionKey := constants.SessionKey(sessionID)
	if err := s.redis.Set(ctx, sessionKey, userID, 24*time.Hour).Err(); err != nil {
		return nil, fmt.Errorf("cache session: %w", err)
	}

	return toDailyPaperResponse(paper, selected), nil
}

// SubmitAnswer evaluates and records an answer, then emits question.answered.
func (s *QuestionService) SubmitAnswer(ctx context.Context, userID, questionID string, req dto.SubmitRequest) (*dto.SubmitResponse, error) {
	if len(req.SessionID) == 0 || len(req.Answer) == 0 {
		return nil, ErrInvalidRequest
	}

	inSession, err := s.dailyPapers.QuestionInSession(ctx, req.SessionID, questionID)
	if err != nil {
		return nil, err
	}
	if !inSession {
		return nil, ErrQuestionNotInSession
	}

	exists, err := s.history.ExistsForSession(ctx, userID, questionID, req.SessionID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrAnswerAlreadySubmitted
	}

	question, err := s.questions.GetByID(ctx, questionID)
	if err != nil {
		return nil, err
	}
	if question == nil {
		return nil, ErrQuestionNotFound
	}

	submittedAt := time.Now().UTC()
	if len(req.SubmittedAt) > 0 {
		parsed, err := time.Parse(time.RFC3339, req.SubmittedAt)
		if err != nil {
			return nil, ErrInvalidRequest
		}
		submittedAt = parsed
	}

	hadAnswerToday, err := s.history.HasAnswerToday(ctx, userID)
	if err != nil {
		return nil, err
	}

	readinessBefore, _ := s.history.AvgScoreByUser(ctx, userID)

	eval := EvaluateAnswer(req.Answer, question.AnswerGuide)
	xpAwarded, gemsAwarded := computeRewards(question, eval)

	if err := s.history.Insert(ctx, userID, questionID, req.SessionID, eval.Correct, eval.Score, submittedAt); err != nil {
		return nil, err
	}

	readinessAfter, _ := s.history.AvgScoreByUser(ctx, userID)
	readinessDelta := readinessAfter - readinessBefore
	if readinessDelta == 0 && eval.Correct {
		readinessDelta = 1
	}

	streakUpdated := eval.Correct && !hadAnswerToday

	event := events.QuestionAnswered{
		EventID:     uuid.NewString(),
		UserID:      userID,
		QuestionID:  questionID,
		RoundType:   question.RoundType,
		Difficulty:  question.Difficulty,
		CompanyTags: question.CompanyTags,
		Correct:     eval.Correct,
		Score:       eval.Score,
		XPAwarded:   xpAwarded,
		GemsAwarded: gemsAwarded,
		SubmittedAt: submittedAt,
		SessionID:   req.SessionID,
	}
	if err := s.publisher.Publish(ctx, events.TopicQuestionAnswered, userID, event); err != nil {
		return nil, fmt.Errorf("publish question answered: %w", err)
	}

	return &dto.SubmitResponse{
		Correct:        eval.Correct,
		Score:          eval.Score,
		XPAwarded:      xpAwarded,
		GemsAwarded:    gemsAwarded,
		StreakUpdated:  streakUpdated,
		ReadinessDelta: readinessDelta,
		Feedback:       eval.Summary,
		Strengths:      eval.Strengths,
		Gaps:           eval.Gaps,
	}, nil
}

// ListCompanies returns available company tags.
func (s *QuestionService) ListCompanies(ctx context.Context) ([]string, error) {
	return s.questions.ListCompanies(ctx)
}

// GetSessionHistory returns answer history, optionally scoped to a daily session.
func (s *QuestionService) GetSessionHistory(ctx context.Context, userID, sessionID string) ([]dto.HistoryEntry, error) {
	if len(sessionID) == 0 {
		return []dto.HistoryEntry{}, nil
	}

	records, err := s.history.ListBySession(ctx, userID, sessionID)
	if err != nil {
		return nil, err
	}

	entries := make([]dto.HistoryEntry, 0, len(records))
	for _, rec := range records {
		entries = append(entries, dto.HistoryEntry{
			QuestionID:  rec.QuestionID,
			SessionID:   rec.SessionID,
			Correct:     rec.Correct,
			Score:       rec.Score,
			SubmittedAt: rec.SubmittedAt.UTC().Format(time.RFC3339),
		})
	}
	return entries, nil
}

// GetReadinessStats aggregates per-company answer performance for readiness.
func (s *QuestionService) GetReadinessStats(ctx context.Context, userID string) (*dto.ReadinessStats, error) {
	rows, err := s.history.CompanyPerformanceByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	stats := make([]dto.CompanyStats, 0, len(rows))
	for _, row := range rows {
		stats = append(stats, dto.CompanyStats{
			Company:  row.Company,
			Answered: row.Answered,
			Correct:  row.Correct,
			ScoreAvg: row.ScoreAvg,
		})
	}

	return &dto.ReadinessStats{ByCompany: stats}, nil
}

func (s *QuestionService) selectQuestions(ctx context.Context, userID, difficulty string, weekendOnly bool) ([]store.Question, error) {
	limit := config.DailyPaperMaxQuestions

	// priorities 1-2 require target companies; skipped when none are configured
	questions, err := s.questions.SelectUnseenByDifficulty(ctx, userID, difficulty, limit, weekendOnly)
	if err != nil {
		return nil, err
	}
	if len(questions) >= 1 {
		return questions, nil
	}

	questions, err = s.questions.SelectRandomApproved(ctx, userID, limit, weekendOnly)
	if err != nil {
		return nil, err
	}
	return questions, nil
}

func toDailyPaperResponse(paper *store.DailyPaper, questions []store.Question) *dto.DailyPaperResponse {
	resp := &dto.DailyPaperResponse{
		SessionID:       paper.SessionID,
		Date:            paper.PaperDate.Format("2006-01-02"),
		MinimumToStreak: config.MinimumAnswersForStreak,
	}
	for _, q := range questions {
		resp.Questions = append(resp.Questions, dto.QuestionResponse{
			ID:          q.ID,
			Body:        q.Body,
			RoundType:   q.RoundType,
			Difficulty:  q.Difficulty,
			CompanyTags: q.CompanyTags,
			IsWeekend:   q.IsWeekend,
		})
	}
	return resp
}
