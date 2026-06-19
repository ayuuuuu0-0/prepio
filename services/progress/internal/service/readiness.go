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

// ReadinessService computes and persists skill-based readiness scores.
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

// ComputeCompanyReadiness derives a weighted readiness score for one company.
func ComputeCompanyReadiness(weights []store.CompanySkillWeight, masteryBySkill map[string]int) int {
	if len(weights) == 0 {
		return 0
	}

	totalWeight := 0
	weightedSum := 0
	for _, weight := range weights {
		mastery, ok := masteryBySkill[weight.SkillID]
		if !ok {
			mastery = 0
		}
		totalWeight += weight.Weight
		weightedSum += mastery * weight.Weight
	}
	if totalWeight == 0 {
		return 0
	}

	score := weightedSum / totalWeight
	if score > config.MaxCompanyReadiness {
		score = config.MaxCompanyReadiness
	}
	return score
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

	targets, err := s.readiness.ListUserTargetCompanies(ctx, userID)
	if err != nil {
		return nil, err
	}
	masteryBySkill := masteryMapFromScores(scores)
	allGaps := make([]dto.SkillGap, 0)
	for _, company := range targets {
		weights, err := s.readiness.ListCompanySkillWeights(ctx, company)
		if err != nil {
			return nil, err
		}
		allGaps = append(allGaps, BuildCompanySkillGaps(company, weights, masteryBySkill)...)
	}

	explanations := []dto.ReadinessExplanation{
		BuildSkillMasteryExplanation(overall, weakestSkills),
	}

	return &dto.SkillReadinessResponse{
		Skills:        entries,
		Overall:       overall,
		TopSkills:     topSkills,
		WeakestSkills: weakestSkills,
		SkillGaps:     MergeSkillGaps(allGaps),
		Explanations:  explanations,
		Version:       readinessVersionV2,
	}, nil
}

// GetCompanyReadiness returns company readiness for the user's target companies.
func (s *ReadinessService) GetCompanyReadiness(ctx context.Context, userID string) (*dto.CompanyReadinessResponse, error) {
	if len(userID) == 0 {
		return nil, ErrInvalidRequest
	}

	targets, err := s.readiness.ListUserTargetCompanies(ctx, userID)
	if err != nil {
		return nil, err
	}

	scores, err := s.readiness.ListUserSkillScores(ctx, userID)
	if err != nil {
		return nil, err
	}

	masteryBySkill := masteryMapFromScores(scores)
	summaries := BuildSkillSummaries(scores)

	entries := make([]dto.CompanyReadinessEntry, 0, len(targets))
	allGaps := make([]dto.SkillGap, 0)
	explanations := make([]dto.ReadinessExplanation, 0, len(targets))
	totalReadiness := 0

	for _, company := range targets {
		weights, err := s.readiness.ListCompanySkillWeights(ctx, company)
		if err != nil {
			return nil, err
		}

		contributions := make([]dto.SkillContribution, 0, len(weights))
		companySummaries := make([]dto.SkillSummary, 0, len(weights))
		for _, weight := range weights {
			mastery := masteryBySkill[weight.SkillID]
			contributions = append(contributions, dto.SkillContribution{
				SkillSlug: weight.SkillSlug,
				SkillName: weight.SkillName,
				Mastery:   mastery,
				Weight:    weight.Weight,
			})
			companySummaries = append(companySummaries, dto.SkillSummary{
				SkillSlug: weight.SkillSlug,
				SkillName: weight.SkillName,
				Mastery:   mastery,
				Attempts:  attemptsForSkill(scores, weight.SkillID),
			})
		}

		readinessScore := ComputeCompanyReadiness(weights, masteryBySkill)
		gaps := BuildCompanySkillGaps(company, weights, masteryBySkill)
		explanation := BuildCompanyExplanation(company, readinessScore, gaps)

		entries = append(entries, dto.CompanyReadinessEntry{
			Company:            company,
			Readiness:          readinessScore,
			SkillContributions: contributions,
			TopSkills:          TopSkills(companySummaries, maxTopWeakestSkills),
			WeakestSkills:      WeakestSkills(companySummaries, maxTopWeakestSkills),
			SkillGaps:          gaps,
			Explanation:        explanation,
		})
		allGaps = append(allGaps, gaps...)
		explanations = append(explanations, explanation)
		totalReadiness += readinessScore
	}

	overall := 0
	if len(entries) > 0 {
		overall = totalReadiness / len(entries)
	}

	return &dto.CompanyReadinessResponse{
		Companies:     entries,
		Overall:       overall,
		TopSkills:     TopSkills(summaries, maxTopWeakestSkills),
		WeakestSkills: WeakestSkills(summaries, maxTopWeakestSkills),
		SkillGaps:     MergeSkillGaps(allGaps),
		Explanations:  explanations,
		Version:       readinessVersionV2,
	}, nil
}

// GetReadinessDashboard aggregates skill and company readiness for validation tooling.
func (s *ReadinessService) GetReadinessDashboard(ctx context.Context, userID string) (*dto.ReadinessDashboardResponse, error) {
	skillMastery, err := s.GetSkillReadiness(ctx, userID)
	if err != nil {
		return nil, err
	}
	companyReadiness, err := s.GetCompanyReadiness(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &dto.ReadinessDashboardResponse{
		SkillMastery:     *skillMastery,
		CompanyReadiness: *companyReadiness,
	}, nil
}

func masteryMapFromScores(scores []store.UserSkillScore) map[string]int {
	masteryBySkill := make(map[string]int, len(scores))
	for _, score := range scores {
		masteryBySkill[score.SkillID] = score.Mastery
	}
	return masteryBySkill
}

func attemptsForSkill(scores []store.UserSkillScore, skillID string) int {
	for _, score := range scores {
		if score.SkillID == skillID {
			return score.Attempts
		}
	}
	return 0
}
