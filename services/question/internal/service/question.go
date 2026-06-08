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
	users       *store.UserStore
	redis       *redis.Client
	publisher   EventPublisher
}

// NewQuestionService creates a QuestionService.
func NewQuestionService(
	questions *store.QuestionStore,
	dailyPapers *store.DailyPaperStore,
	history *store.HistoryStore,
	users *store.UserStore,
	redisClient *redis.Client,
	publisher EventPublisher,
) *QuestionService {
	return &QuestionService{
		questions:   questions,
		dailyPapers: dailyPapers,
		history:     history,
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

	correct := EvaluateAnswer(req.Answer, question.AnswerGuide)
	if err := s.history.Insert(ctx, userID, questionID, req.SessionID, correct, submittedAt); err != nil {
		return nil, err
	}

	event := events.QuestionAnswered{
		EventID:     uuid.NewString(),
		UserID:      userID,
		QuestionID:  questionID,
		RoundType:   question.RoundType,
		Difficulty:  question.Difficulty,
		CompanyTags: question.CompanyTags,
		Correct:     correct,
		SubmittedAt: submittedAt,
		SessionID:   req.SessionID,
	}
	if err := s.publisher.Publish(ctx, events.TopicQuestionAnswered, userID, event); err != nil {
		return nil, fmt.Errorf("publish question answered: %w", err)
	}

	return &dto.SubmitResponse{
		Correct:       correct,
		XPAwarded:     0,
		GemsAwarded:   0,
		StreakUpdated: false,
		Feedback:      FeedbackFor(correct),
	}, nil
}

// ListCompanies returns available company tags.
func (s *QuestionService) ListCompanies(ctx context.Context) ([]string, error) {
	return s.questions.ListCompanies(ctx)
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
