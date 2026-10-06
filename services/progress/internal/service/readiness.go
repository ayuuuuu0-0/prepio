package service

import (
	"context"
	"time"

	"github.com/prepio/prepio/services/progress/internal/dto"
	"github.com/prepio/prepio/services/progress/internal/store"
)

const readinessVersionV2 = "v2"

// ReadinessService computes and persists skill-based mastery scores.
type ReadinessService struct {
	readiness *store.ReadinessStore
}

// NewReadinessService creates a ReadinessService.
func NewReadinessService(readiness *store.ReadinessStore) *ReadinessService {
	return &ReadinessService{readiness: readiness}
}

// GetSkillReadiness returns the user's skill mastery scores with analysis.
func (s *ReadinessService) GetSkillReadiness(ctx context.Context, userID string) (*dto.SkillReadinessResponse, error) {
	if len(userID) == 0 {
		return nil, ErrInvalidRequest
	}

	scores, err := s.readiness.ListUserSkillScores(ctx, userID)
	if err != nil {
		return nil, err
	}

	entries := make([]dto.SkillReadinessEntry, 0, len(scores))
	totalMastery := 0
	for _, score := range scores {
		entry := dto.SkillReadinessEntry{
			SkillSlug: score.SkillSlug,
			SkillName: score.SkillName,
			Mastery:   score.Mastery,
			Attempts:  score.Attempts,
		}
		if score.LastPracticedAt != nil {
			entry.LastPracticedAt = score.LastPracticedAt.UTC().Format(time.RFC3339)
		}
		entries = append(entries, entry)
		totalMastery += score.Mastery
	}

	summaries := BuildSkillSummaries(scores)
	topSkills := TopSkills(summaries, maxTopWeakestSkills)
	weakestSkills := WeakestSkills(summaries, maxTopWeakestSkills)

	overall := 0
	if len(scores) > 0 {
		overall = totalMastery / len(scores)
	}

	gaps := BuildSkillGaps(weakestSkills)
	explanations := []dto.ReadinessExplanation{
		BuildSkillMasteryExplanation(overall, weakestSkills),
	}

	return &dto.SkillReadinessResponse{
		Skills:        entries,
		Overall:       overall,
		TopSkills:     topSkills,
		WeakestSkills: weakestSkills,
		SkillGaps:     gaps,
		Explanations:  explanations,
		Version:       readinessVersionV2,
	}, nil
}
