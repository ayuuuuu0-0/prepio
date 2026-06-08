package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/prepio/prepio/config"
	"github.com/prepio/prepio/constants"
	"github.com/prepio/prepio/services/streak/internal/dto"
	"github.com/prepio/prepio/services/streak/internal/store"
	"github.com/prepio/prepio/shared/events"
)

// EventPublisher publishes streak events.
type EventPublisher interface {
	Publish(ctx context.Context, topic, key string, payload any) error
}

// GemDeductor deducts gems via the progress service internal API.
type GemDeductor interface {
	DeductGems(ctx context.Context, userID string, amount int, reason string) error
}

// StreakService owns all streak logic.
type StreakService struct {
	streaks   *store.StreakStore
	freezes   *store.FreezeStore
	cache     *store.CacheStore
	publisher EventPublisher
	gems      GemDeductor
}

// NewStreakService creates a StreakService.
func NewStreakService(
	streaks *store.StreakStore,
	freezes *store.FreezeStore,
	cache *store.CacheStore,
	publisher EventPublisher,
	gems GemDeductor,
) *StreakService {
	return &StreakService{
		streaks:   streaks,
		freezes:   freezes,
		cache:     cache,
		publisher: publisher,
		gems:      gems,
	}
}

// ProcessQuestionAnswered applies streak rules for a submitted answer.
func (s *StreakService) ProcessQuestionAnswered(ctx context.Context, event events.QuestionAnswered) error {
	timezone, err := s.streaks.Timezone(ctx, event.UserID)
	if err != nil {
		return err
	}
	if len(timezone) == 0 {
		timezone = constants.DefaultTimezone
	}

	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return fmt.Errorf("load timezone: %w", err)
	}

	activityDate := calendarDateIn(event.SubmittedAt.In(loc))
	state, err := s.streaks.Get(ctx, event.UserID)
	if err != nil {
		return err
	}

	freezeCount, err := s.freezes.Count(ctx, event.UserID)
	if err != nil {
		return err
	}

	previous := state.CurrentStreak
	updated, consumedFreeze, changed := applyActivity(state, activityDate, freezeCount)

	if !changed {
		return nil
	}

	if consumedFreeze {
		freezeCount--
		if err := s.freezes.SetCount(ctx, event.UserID, freezeCount); err != nil {
			return err
		}
	}

	if err := s.streaks.Upsert(ctx, updated); err != nil {
		return err
	}
	if err := s.cache.WriteThrough(ctx, updated, freezeCount); err != nil {
		return err
	}

	streakEvent := events.StreakUpdated{
		EventID:        uuid.NewString(),
		UserID:         event.UserID,
		PreviousStreak: previous,
		CurrentStreak:  updated.CurrentStreak,
		StreakBroken:   updated.CurrentStreak < previous,
		FreezeConsumed: consumedFreeze,
		UpdatedAt:      time.Now().UTC(),
	}
	return s.publisher.Publish(ctx, events.TopicStreakUpdated, event.UserID, streakEvent)
}

// GetMe returns the authenticated user's streak summary.
func (s *StreakService) GetMe(ctx context.Context, userID, timezone string) (*dto.StreakResponse, error) {
	if len(timezone) == 0 {
		var err error
		timezone, err = s.streaks.Timezone(ctx, userID)
		if err != nil {
			return nil, err
		}
		if len(timezone) == 0 {
			timezone = constants.DefaultTimezone
		}
	}

	state, err := s.streaks.Get(ctx, userID)
	if err != nil {
		return nil, err
	}
	freezeCount, err := s.freezes.Count(ctx, userID)
	if err != nil {
		return nil, err
	}

	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return nil, fmt.Errorf("load timezone: %w", err)
	}
	today := calendarDateIn(time.Now().In(loc))

	var lastActivity string
	activeToday := false
	if state.LastActivityDate != nil {
		lastActivity = state.LastActivityDate.Format("2006-01-02")
		activeToday = calendarDateFromDB(*state.LastActivityDate).Equal(today)
	}

	return &dto.StreakResponse{
		CurrentStreak:    state.CurrentStreak,
		LongestStreak:    state.LongestStreak,
		FreezeCount:      freezeCount,
		LastActivityDate: lastActivity,
		StreakActiveToday: activeToday,
	}, nil
}

// PurchaseFreeze buys a streak freeze if inventory and gems allow.
func (s *StreakService) PurchaseFreeze(ctx context.Context, userID string) (*dto.FreezePurchaseResponse, error) {
	count, err := s.freezes.Count(ctx, userID)
	if err != nil {
		return nil, err
	}
	if count >= config.MaxStreakFreezes {
		return nil, ErrFreezeMaxHeld
	}

	if err := s.gems.DeductGems(ctx, userID, config.StreakFreezeGemCost, "streak_freeze_purchase"); err != nil {
		return nil, err
	}

	count++
	if err := s.freezes.SetCount(ctx, userID, count); err != nil {
		return nil, err
	}

	state, err := s.streaks.Get(ctx, userID)
	if err != nil {
		return nil, err
	}
	if err := s.cache.WriteThrough(ctx, *state, count); err != nil {
		return nil, err
	}

	return &dto.FreezePurchaseResponse{
		FreezeCount: count,
		GemsSpent:   config.StreakFreezeGemCost,
	}, nil
}

func applyActivity(state *store.StreakState, activityDate time.Time, freezeCount int) (store.StreakState, bool, bool) {
	updated := *state
	consumedFreeze := false
	changed := false

	if state.LastActivityDate != nil && calendarDateFromDB(*state.LastActivityDate).Equal(activityDate) {
		return updated, false, false
	}

	changed = true
	if state.LastActivityDate == nil {
		updated.CurrentStreak = 1
	} else {
		last := calendarDateFromDB(*state.LastActivityDate)
		diff := int(activityDate.Sub(last).Hours() / 24)

		switch {
		case diff == 1:
			updated.CurrentStreak++
		case diff == 2 && freezeCount > 0:
			consumedFreeze = true
			updated.CurrentStreak++
		default:
			updated.CurrentStreak = 1
		}
	}

	if updated.CurrentStreak > updated.LongestStreak {
		updated.LongestStreak = updated.CurrentStreak
	}
	day := activityDate
	updated.LastActivityDate = &day
	return updated, consumedFreeze, changed
}

// calendarDateIn normalizes a local instant to a UTC calendar date.
func calendarDateIn(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}

// calendarDateFromDB normalizes a PostgreSQL DATE value for comparison.
func calendarDateFromDB(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}
