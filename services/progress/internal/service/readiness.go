package service

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/prepio/prepio/config"
	"github.com/prepio/prepio/services/progress/internal/dto"
	"github.com/prepio/prepio/services/progress/internal/store"
	"github.com/prepio/prepio/shared/events"
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

// ProcessQuestionAnswered updates skill mastery from an answer event.
func (s *ReadinessService) ProcessQuestionAnswered(ctx context.Context, event events.QuestionAnswered) error {
	if len(event.UserID) == 0 || len(event.QuestionID) == 0 {
		return fmt.Errorf("user id and question id are required")
	}
	if event.Score <= 0 {
		return nil
	}

	contributions, err := s.readiness.ListQuestionSkillContributions(ctx, event.QuestionID)
	if err != nil {
		return err
	}
	if len(contributions) == 0 {
		return nil
	}

	practicedAt := event.SubmittedAt
	if practicedAt.IsZero() {
		practicedAt = time.Now().UTC()
	}

	for _, contribution := range contributions {
		delta := MasteryContribution(event.Score, contribution)
		if delta <= 0 {
			continue
		}

		existing, err := s.readiness.GetUserSkillScore(ctx, event.UserID, contribution.SkillID)
		if err != nil {
			return err
		}

		currentMastery := 0
		attempts := 0
		if existing != nil {
			currentMastery = existing.Mastery
			attempts = existing.Attempts
		}

		newMastery := ApplyMasteryDelta(currentMastery, delta)
		if err := s.readiness.UpsertUserSkillScore(
			ctx,
			event.UserID,
			contribution.SkillID,
			newMastery,
			attempts+1,
			practicedAt,
			config.ReadinessSourceLive,
		); err != nil {
			return err
		}
	}
	return nil
}

// MasteryContribution calculates the raw mastery delta from one answer for a skill mapping.
func MasteryContribution(score int, contribution store.QuestionSkillContribution) float64 {
	if score <= 0 {
		return 0
	}
	return (float64(score) / 100.0) *
		contribution.ReadinessWeight *
		contribution.SkillWeight *
		config.DifficultyMultiplier(contribution.Difficulty)
}

// ApplyMasteryDelta smooths a contribution into the current mastery score.
func ApplyMasteryDelta(currentMastery int, contribution float64) int {
	smoothed := float64(currentMastery) + contribution*100.0*config.MasterySmoothingFactor
	return int(math.Min(float64(config.MaxSkillMastery), math.Round(smoothed)))
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
