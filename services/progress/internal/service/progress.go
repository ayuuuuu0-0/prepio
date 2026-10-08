package service

import (
	"context"
	"errors"
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
	publisher EventPublisher
}

// NewProgressService creates a ProgressService.
func NewProgressService(progress *store.ProgressStore, publisher EventPublisher) *ProgressService {
	return &ProgressService{progress: progress, publisher: publisher}
}

// ProcessStreakUpdated awards streak bonus gems when a streak increments. Redelivery is safe.
func (s *ProgressService) ProcessStreakUpdated(ctx context.Context, event events.StreakUpdated) error {
	if event.StreakBroken || event.CurrentStreak <= event.PreviousStreak {
		return nil
	}
	if !validIDs(event.EventID, event.UserID) {
		return ErrInvalidRequest
	}

	gems := config.StreakIncrementGemBonus
	state, applied, err := s.progress.AwardGems(ctx, event.UserID, gems, "streak_increment", event.EventID)
	if err != nil || !applied {
		return err
	}
	return s.emitUpdated(ctx, event.UserID, 0, gems, state, config.CurrentLevel(state.TotalXP))
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
	if amount <= 0 || !validIDs(userID) {
		return 0, ErrInvalidRequest
	}
	balance, err := s.progress.DeductGems(ctx, userID, amount, reason)
	if errors.Is(err, store.ErrInsufficientGems) {
		return 0, ErrInsufficientGems
	}
	return balance, err
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
