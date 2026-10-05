package service

import (
	"fmt"
	"sort"

	"github.com/prepio/prepio/config"
	"github.com/prepio/prepio/services/progress/internal/dto"
	"github.com/prepio/prepio/services/progress/internal/store"
)

const (
	readinessGapMasteryThreshold = 70
	maxTopWeakestSkills          = 5
)

// BuildSkillSummaries converts score rows into skill summaries.
func BuildSkillSummaries(scores []store.UserSkillScore) []dto.SkillSummary {
	summaries := make([]dto.SkillSummary, 0, len(scores))
	for _, score := range scores {
		summaries = append(summaries, dto.SkillSummary{
			SkillSlug: score.SkillSlug,
			SkillName: score.SkillName,
			Mastery:   score.Mastery,
			Attempts:  score.Attempts,
		})
	}
	return summaries
}

// TopSkills returns the highest mastery skills with at least one attempt.
func TopSkills(summaries []dto.SkillSummary, limit int) []dto.SkillSummary {
	if limit <= 0 {
		limit = maxTopWeakestSkills
	}
	filtered := make([]dto.SkillSummary, 0, len(summaries))
	for _, summary := range summaries {
		if summary.Attempts > 0 {
			filtered = append(filtered, summary)
		}
	}
	sort.Slice(filtered, func(i, j int) bool {
		if filtered[i].Mastery == filtered[j].Mastery {
			return filtered[i].SkillName < filtered[j].SkillName
		}
		return filtered[i].Mastery > filtered[j].Mastery
	})
	if len(filtered) > limit {
		filtered = filtered[:limit]
	}
	return filtered
}

// WeakestSkills returns the lowest mastery skills with at least one attempt.
func WeakestSkills(summaries []dto.SkillSummary, limit int) []dto.SkillSummary {
	if limit <= 0 {
		limit = maxTopWeakestSkills
	}
	filtered := make([]dto.SkillSummary, 0, len(summaries))
	for _, summary := range summaries {
		if summary.Attempts > 0 {
			filtered = append(filtered, summary)
		}
	}
	sort.Slice(filtered, func(i, j int) bool {
		if filtered[i].Mastery == filtered[j].Mastery {
			return filtered[i].SkillName < filtered[j].SkillName
		}
		return filtered[i].Mastery < filtered[j].Mastery
	})
	if len(filtered) > limit {
		filtered = filtered[:limit]
	}
	return filtered
}

// ComputeSkillGapScore measures distance from full mastery.
func ComputeSkillGapScore(mastery int) int {
	return config.MaxSkillMastery - mastery
}

// BuildSkillGaps identifies practiced skills holding back mastery.
func BuildSkillGaps(weakest []dto.SkillSummary) []dto.SkillGap {
	gaps := make([]dto.SkillGap, 0, len(weakest))
	for _, skill := range weakest {
		if skill.Mastery >= readinessGapMasteryThreshold {
			continue
		}
		gaps = append(gaps, dto.SkillGap{
			SkillSlug:   skill.SkillSlug,
			SkillName:   skill.SkillName,
			Mastery:     skill.Mastery,
			GapScore:    ComputeSkillGapScore(skill.Mastery),
			Explanation: fmt.Sprintf("%s mastery is %d after %d attempts", skill.SkillName, skill.Mastery, skill.Attempts),
		})
	}
	sort.Slice(gaps, func(i, j int) bool {
		if gaps[i].GapScore == gaps[j].GapScore {
			return gaps[i].SkillName < gaps[j].SkillName
		}
		return gaps[i].GapScore > gaps[j].GapScore
	})
	return gaps
}

// BuildSkillMasteryExplanation summarizes overall skill mastery state.
func BuildSkillMasteryExplanation(overall int, weakest []dto.SkillSummary) dto.ReadinessExplanation {
	details := make([]string, 0, len(weakest))
	for _, skill := range weakest {
		details = append(details, fmt.Sprintf("%s mastery is %d after %d attempts", skill.SkillName, skill.Mastery, skill.Attempts))
	}
	return dto.ReadinessExplanation{
		Scope:   "skills",
		Summary: fmt.Sprintf("Overall skill mastery average is %d across practiced skills", overall),
		Details: details,
	}
}
