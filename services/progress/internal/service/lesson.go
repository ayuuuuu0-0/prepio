package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/prepio/prepio/services/progress/internal/dto"
	"github.com/prepio/prepio/services/progress/internal/store"
	"github.com/prepio/prepio/shared/events"
)

// LessonService turns lesson.completed facts into mastery, XP, and gems.
type LessonService struct {
	lessons   *store.LessonStore
	publisher EventPublisher
}

// NewLessonService creates a LessonService.
func NewLessonService(lessons *store.LessonStore, publisher EventPublisher) *LessonService {
	return &LessonService{lessons: lessons, publisher: publisher}
}

// ProcessLessonCompleted applies a lesson.completed event. Redelivery is safe.
func (s *LessonService) ProcessLessonCompleted(ctx context.Context, event events.LessonCompleted) error {
	if !validIDs(event.EventID, event.UserID, event.LessonID, event.AttemptID) || len(event.Skills) == 0 {
		return fmt.Errorf("lesson completed: %w", ErrInvalidRequest)
	}
	for _, sk := range event.Skills {
		if !validIDs(sk.SkillID) || sk.Weight <= 0 {
			return fmt.Errorf("lesson completed: invalid skill weight: %w", ErrInvalidRequest)
		}
	}
	if event.CompletedAt.IsZero() {
		event.CompletedAt = time.Now().UTC()
	}

	outcome, applied, err := s.lessons.ApplyLessonCompleted(ctx, event)
	if err != nil {
		return err
	}
	if !applied || (outcome.XPAwarded == 0 && outcome.GemsAwarded == 0) {
		return nil
	}

	return s.publisher.Publish(ctx, events.TopicProgressUpdated, event.UserID, events.ProgressUpdated{
		EventID:     uuid.NewString(),
		UserID:      event.UserID,
		XPAwarded:   outcome.XPAwarded,
		GemsAwarded: outcome.GemsAwarded,
		TotalXP:     outcome.TotalXP,
		TotalGems:   outcome.TotalGems,
		LevelBefore: outcome.LevelBefore,
		LevelAfter:  outcome.LevelAfter,
		LeveledUp:   outcome.LevelAfter > outcome.LevelBefore,
		UpdatedAt:   time.Now().UTC(),
	})
}

// GetAttemptRewards returns what the user's attempt earned. found is false until
// Progress has processed the completion.
func (s *LessonService) GetAttemptRewards(ctx context.Context, userID, attemptID string) (resp *dto.AttemptRewardsResponse, found bool, err error) {
	if !validIDs(userID, attemptID) {
		return nil, false, ErrInvalidRequest
	}
	outcome, err := s.lessons.GetAttemptRewards(ctx, userID, attemptID)
	if err != nil {
		return nil, false, err
	}
	if outcome == nil {
		return nil, false, nil
	}
	resp = &dto.AttemptRewardsResponse{
		AttemptID:            outcome.AttemptID,
		FirstCompletion:      outcome.FirstCompletion,
		XPAwarded:            outcome.XPAwarded,
		GemsAwarded:          outcome.GemsAwarded,
		MasteryChanges:       make([]dto.MasteryChangeResponse, 0, len(outcome.Mastery)),
		AchievementsUnlocked: achievementResponses(outcome.Achievements),
	}
	for _, m := range outcome.Mastery {
		resp.MasteryChanges = append(resp.MasteryChanges, dto.MasteryChangeResponse{
			SkillSlug: m.SkillSlug, SkillName: m.SkillName,
			TopicSlug: m.TopicSlug, TopicName: m.TopicName,
			Before: m.Before, After: m.After, Delta: m.Delta, Accuracy: m.Accuracy,
		})
	}
	return resp, true, nil
}

func validIDs(ids ...string) bool {
	for _, id := range ids {
		if _, err := uuid.Parse(id); err != nil {
			return false
		}
	}
	return true
}
