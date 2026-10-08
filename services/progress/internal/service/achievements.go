package service

import (
	"context"
	"time"

	"github.com/prepio/prepio/config"
	"github.com/prepio/prepio/services/progress/internal/dto"
	"github.com/prepio/prepio/services/progress/internal/store"
)

// AchievementService serves the achievement catalog with the learner's progress through it.
type AchievementService struct {
	achievements *store.AchievementStore
}

// NewAchievementService creates an AchievementService.
func NewAchievementService(achievements *store.AchievementStore) *AchievementService {
	return &AchievementService{achievements: achievements}
}

// List returns every achievement in catalog order, marking the ones the learner has earned.
func (s *AchievementService) List(ctx context.Context, userID string) ([]dto.AchievementResponse, error) {
	if !validIDs(userID) {
		return nil, ErrInvalidRequest
	}
	earned, err := s.achievements.List(ctx, userID)
	if err != nil {
		return nil, err
	}
	at := make(map[string]time.Time, len(earned))
	for _, e := range earned {
		at[e.Slug] = e.UnlockedAt
	}
	out := make([]dto.AchievementResponse, 0, len(config.Achievements))
	for _, a := range config.Achievements {
		r := dto.AchievementResponse{Slug: a.Slug, Name: a.Name, Description: a.Description}
		if t, ok := at[a.Slug]; ok {
			stamp := t.UTC().Format(time.RFC3339)
			r.Unlocked, r.UnlockedAt = true, &stamp
		}
		out = append(out, r)
	}
	return out, nil
}

// achievementResponses describes newly unlocked achievement slugs. Unknown slugs (removed
// from the catalog) are skipped.
func achievementResponses(slugs []string) []dto.AchievementResponse {
	out := make([]dto.AchievementResponse, 0, len(slugs))
	for _, slug := range slugs {
		if a, ok := config.AchievementBySlug(slug); ok {
			out = append(out, dto.AchievementResponse{Slug: a.Slug, Name: a.Name, Description: a.Description, Unlocked: true})
		}
	}
	return out
}
