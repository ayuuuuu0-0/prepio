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

// ComputeSkillGapScore measures urgency: weighted distance from full mastery.
func ComputeSkillGapScore(mastery, weight int) int {
	if weight <= 0 {
		return 0
	}
	return weight * (config.MaxSkillMastery - mastery) / config.MaxSkillMastery
}

// BuildCompanySkillGaps returns weighted skills limiting company readiness.
func BuildCompanySkillGaps(
	company string,
	weights []store.CompanySkillWeight,
	masteryBySkill map[string]int,
) []dto.SkillGap {
	gaps := make([]dto.SkillGap, 0, len(weights))
	for _, weight := range weights {
		mastery := masteryBySkill[weight.SkillID]
		gapScore := ComputeSkillGapScore(mastery, weight.Weight)
		if mastery >= readinessGapMasteryThreshold {
			continue
		}
		gaps = append(gaps, dto.SkillGap{
			Company:     company,
			SkillSlug:   weight.SkillSlug,
			SkillName:   weight.SkillName,
			Mastery:     mastery,
			Weight:      weight.Weight,
			GapScore:    gapScore,
			Explanation: fmt.Sprintf("%s mastery is %d — carries %d%% weight for %s readiness", weight.SkillName, mastery, weight.Weight, company),
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

// BuildCompanyExplanation summarizes company readiness from weighted skills.
func BuildCompanyExplanation(company string, readiness int, gaps []dto.SkillGap) dto.ReadinessExplanation {
	details := make([]string, 0, len(gaps))
	for i, gap := range gaps {
		if i >= 3 {
			break
		}
		details = append(details, gap.Explanation)
	}

	summary := fmt.Sprintf("%s readiness is %d based on weighted skill mastery", company, readiness)
	if len(gaps) == 0 {
		summary = fmt.Sprintf("%s readiness is %d — no major weighted skill gaps below %d", company, readiness, readinessGapMasteryThreshold)
	} else if len(gaps) > 0 {
		summary = fmt.Sprintf("%s readiness is %d — weakest weighted skill is %s (%d mastery)", company, readiness, gaps[0].SkillName, gaps[0].Mastery)
	}

	return dto.ReadinessExplanation{
		Scope:   company,
		Summary: summary,
		Details: details,
	}
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

// MergeSkillGaps combines company gaps sorted by urgency.
func MergeSkillGaps(gapGroups ...[]dto.SkillGap) []dto.SkillGap {
	total := 0
	for _, group := range gapGroups {
		total += len(group)
	}
	merged := make([]dto.SkillGap, 0, total)
	for _, group := range gapGroups {
		merged = append(merged, group...)
	}
	sort.Slice(merged, func(i, j int) bool {
		if merged[i].GapScore == merged[j].GapScore {
			return merged[i].SkillName < merged[j].SkillName
		}
		return merged[i].GapScore > merged[j].GapScore
	})
	if len(merged) > maxTopWeakestSkills*2 {
		merged = merged[:maxTopWeakestSkills*2]
	}
	return merged
}
