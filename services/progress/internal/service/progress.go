package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/prepio/prepio/config"
	"github.com/prepio/prepio/services/progress/internal/dto"
	"github.com/prepio/prepio/services/progress/internal/store"
	"github.com/prepio/prepio/shared/events"
)

// EventPublisher publishes progress events.
type EventPublisher interface {
	Publish(ctx context.Context, topic, key string, payload any) error
}

// ProgressService owns XP, gems, and level state.
type ProgressService struct {
	progress  *store.ProgressStore
	ledger    *store.LedgerStore
	publisher EventPublisher
}

// NewProgressService creates a ProgressService.
func NewProgressService(progress *store.ProgressStore, ledger *store.LedgerStore, publisher EventPublisher) *ProgressService {
	return &ProgressService{progress: progress, ledger: ledger, publisher: publisher}
}

// ProcessQuestionAnswered awards XP and gems for a correct answer.
func (s *ProgressService) ProcessQuestionAnswered(ctx context.Context, event events.QuestionAnswered) error {
	if !event.Correct {
		return nil
	}

	xp := config.XPByDifficulty[event.Difficulty]
	for _, company := range event.CompanyTags {
		if config.TopTierCompanies[company] {
			xp = int(float64(xp) * config.TopTierCompanyXPMultiplier)
			break
		}
	}
	gems := config.GemsByDifficulty[event.Difficulty]

	state, err := s.progress.Get(ctx, event.UserID)
	if err != nil {
		return err
	}

	levelBefore := config.CurrentLevel(state.TotalXP)
	state.TotalXP += xp
	state.GemBalance += gems
	state.CurrentLevel = config.CurrentLevel(state.TotalXP)

	if err := s.progress.Upsert(ctx, *state); err != nil {
		return err
	}
	if err := s.ledger.InsertXP(ctx, event.UserID, xp, "question_answered", event.EventID); err != nil {
		return err
	}
	if gems > 0 {
		if err := s.ledger.InsertGem(ctx, event.UserID, gems, "question_answered", event.EventID); err != nil {
			return err
		}
	}

	return s.emitUpdated(ctx, event.UserID, xp, gems, state, levelBefore)
}

// ProcessStreakUpdated awards streak bonus gems when a streak increments.
func (s *ProgressService) ProcessStreakUpdated(ctx context.Context, event events.StreakUpdated) error {
	if event.StreakBroken || event.CurrentStreak <= event.PreviousStreak {
		return nil
	}

	gems := config.StreakIncrementGemBonus
	state, err := s.progress.Get(ctx, event.UserID)
	if err != nil {
		return err
	}

	levelBefore := config.CurrentLevel(state.TotalXP)
	state.GemBalance += gems
	state.CurrentLevel = config.CurrentLevel(state.TotalXP)

	if err := s.progress.Upsert(ctx, *state); err != nil {
		return err
	}
	if err := s.ledger.InsertGem(ctx, event.UserID, gems, "streak_increment", event.EventID); err != nil {
		return err
	}

	return s.emitUpdated(ctx, event.UserID, 0, gems, state, levelBefore)
}

// GetMe returns the authenticated user's progress summary.
func (s *ProgressService) GetMe(ctx context.Context, userID string) (*dto.ProgressResponse, error) {
	state, err := s.progress.Get(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &dto.ProgressResponse{
		TotalXP:       state.TotalXP,
		CurrentLevel:  state.CurrentLevel,
		GemBalance:    state.GemBalance,
		XPToNextLevel: config.XPToNextLevel(state.TotalXP),
	}, nil
}

// GetGems returns gem balance for internal API.
func (s *ProgressService) GetGems(ctx context.Context, userID string) (int, error) {
	state, err := s.progress.Get(ctx, userID)
	if err != nil {
		return 0, err
	}
	return state.GemBalance, nil
}

// DeductGems deducts gems for transactional operations like streak freeze purchase.
func (s *ProgressService) DeductGems(ctx context.Context, userID string, amount int, reason string) (int, error) {
	if amount <= 0 {
		return 0, ErrInvalidRequest
	}

	state, err := s.progress.Get(ctx, userID)
	if err != nil {
		return 0, err
	}
	if state.GemBalance < amount {
		// ensure row exists before deduct attempt
		if err := s.progress.Upsert(ctx, *state); err != nil {
			return 0, err
		}
	}

	balance, err := s.progress.DeductGems(ctx, userID, amount)
	if err != nil {
		if err == store.ErrInsufficientGems {
			return 0, ErrInsufficientGems
		}
		return 0, err
	}

	eventID := uuid.NewString()
	if err := s.ledger.InsertGem(ctx, userID, -amount, reason, eventID); err != nil {
		return 0, err
	}
	return balance, nil
}

func (s *ProgressService) emitUpdated(ctx context.Context, userID string, xp, gems int, state *store.ProgressState, levelBefore int) error {
	levelAfter := config.CurrentLevel(state.TotalXP)
	event := events.ProgressUpdated{
		EventID:     uuid.NewString(),
		UserID:      userID,
		XPAwarded:   xp,
		GemsAwarded: gems,
		TotalXP:     state.TotalXP,
		TotalGems:   state.GemBalance,
		LevelBefore: levelBefore,
		LevelAfter:  levelAfter,
		LeveledUp:   levelAfter > levelBefore,
		UpdatedAt:   time.Now().UTC(),
	}
	return s.publisher.Publish(ctx, events.TopicProgressUpdated, userID, event)
}
