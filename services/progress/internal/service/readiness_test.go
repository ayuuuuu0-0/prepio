package service_test

import (
	"testing"

	"github.com/prepio/prepio/config"
	"github.com/prepio/prepio/services/progress/internal/service"
	"github.com/prepio/prepio/services/progress/internal/store"
	"github.com/stretchr/testify/require"
)

func TestMasteryContribution(t *testing.T) {
	contribution := store.QuestionSkillContribution{
		SkillWeight:     1.0,
		ReadinessWeight: 1.0,
		Difficulty:      "medium",
	}
	delta := service.MasteryContribution(80, contribution)
	require.InDelta(t, 0.80, delta, 0.001)
}

func TestApplyMasteryDeltaSmoothing(t *testing.T) {
	contribution := store.QuestionSkillContribution{
		SkillWeight:     1.0,
		ReadinessWeight: 1.0,
		Difficulty:      "medium",
	}
	delta := service.MasteryContribution(100, contribution)
	newMastery := service.ApplyMasteryDelta(0, delta)
	expected := int(100.0 * config.MasterySmoothingFactor)
	require.Equal(t, expected, newMastery)
}

func TestComputeCompanyReadiness(t *testing.T) {
	weights := []store.CompanySkillWeight{
		{SkillID: "skill-a", Weight: 50},
		{SkillID: "skill-b", Weight: 50},
	}
	mastery := map[string]int{
		"skill-a": 80,
		"skill-b": 60,
	}
	score := service.ComputeCompanyReadiness(weights, mastery)
	require.Equal(t, 70, score)
}

func TestComputeCompanyReadinessCapsAtMax(t *testing.T) {
	weights := []store.CompanySkillWeight{
		{SkillID: "skill-a", Weight: 100},
	}
	mastery := map[string]int{
		"skill-a": 100,
	}
	score := service.ComputeCompanyReadiness(weights, mastery)
	require.Equal(t, config.MaxCompanyReadiness, score)
}

func TestComputeCompanyReadinessMissingSkill(t *testing.T) {
	weights := []store.CompanySkillWeight{
		{SkillID: "skill-a", Weight: 100},
	}
	score := service.ComputeCompanyReadiness(weights, map[string]int{})
	require.Equal(t, 0, score)
}
