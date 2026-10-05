package service_test

import (
	"testing"

	"github.com/prepio/prepio/services/progress/internal/dto"
	"github.com/prepio/prepio/services/progress/internal/service"
	"github.com/stretchr/testify/require"
)

func TestTopAndWeakestSkills(t *testing.T) {
	summaries := []dto.SkillSummary{
		{SkillSlug: "arrays", SkillName: "Arrays", Mastery: 85, Attempts: 3},
		{SkillSlug: "trees", SkillName: "Trees", Mastery: 41, Attempts: 2},
		{SkillSlug: "graphs", SkillName: "Graphs", Mastery: 0, Attempts: 0},
	}

	top := service.TopSkills(summaries, 2)
	require.Len(t, top, 2)
	require.Equal(t, "arrays", top[0].SkillSlug)

	weakest := service.WeakestSkills(summaries, 2)
	require.Len(t, weakest, 2)
	require.Equal(t, "trees", weakest[0].SkillSlug)
}

func TestComputeSkillGapScore(t *testing.T) {
	require.Equal(t, 100, service.ComputeSkillGapScore(0))
	require.Equal(t, 50, service.ComputeSkillGapScore(50))
	require.Equal(t, 0, service.ComputeSkillGapScore(100))
}

func TestBuildSkillGaps(t *testing.T) {
	weakest := []dto.SkillSummary{
		{SkillSlug: "arrays", SkillName: "Arrays", Mastery: 85, Attempts: 3},
		{SkillSlug: "trees", SkillName: "Trees", Mastery: 41, Attempts: 2},
	}

	gaps := service.BuildSkillGaps(weakest)
	require.Len(t, gaps, 1)
	require.Equal(t, "trees", gaps[0].SkillSlug)
	require.Equal(t, 59, gaps[0].GapScore)
}
