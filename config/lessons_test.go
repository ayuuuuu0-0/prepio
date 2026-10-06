package config_test

import (
	"testing"

	"github.com/prepio/prepio/config"
	"github.com/stretchr/testify/require"
)

func TestLessonXP(t *testing.T) {
	tests := []struct {
		name       string
		kind       string
		difficulty string
		accuracy   float64
		want       int
	}{
		{"perfect medium", "lesson", "medium", 1, 35},
		{"zero accuracy still earns the floor", "lesson", "medium", 0, 21},
		{"easy two thirds", "lesson", "easy", 2.0 / 3.0, 17},
		{"hard perfect", "lesson", "hard", 1, 50},
		{"boss multiplies", "boss", "hard", 1, 75},
		{"unknown difficulty falls back to medium", "lesson", "odd", 1, 35},
		{"accuracy is clamped", "lesson", "medium", 7, 35},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, config.LessonXP(tt.kind, tt.difficulty, tt.accuracy))
		})
	}
}

func TestLessonGems(t *testing.T) {
	require.Equal(t, 10, config.LessonGems("medium", 1))
	require.Equal(t, 5, config.LessonGems("medium", 0.5))
	require.Equal(t, 5, config.LessonGems("unknown-but-medium", 0.1), "unknown difficulty falls back to medium, half gems below the threshold")
}

func TestLessonMasteryDelta(t *testing.T) {
	tests := []struct {
		name     string
		current  int
		kind     string
		diff     string
		accuracy float64
		weight   float64
		want     int
	}{
		{"fresh skill, perfect, full weight", 0, "lesson", "medium", 1, 1, 12},
		{"partial weight", 0, "lesson", "medium", 1, 0.5, 6},
		{"half accuracy", 0, "lesson", "medium", 0.5, 1, 6},
		{"gains shrink near the top", 80, "lesson", "medium", 1, 1, 2},
		{"already maxed", 100, "lesson", "medium", 1, 1, 0},
		{"zero accuracy moves nothing", 10, "lesson", "medium", 0, 1, 0},
		{"zero weight moves nothing", 10, "lesson", "medium", 1, 0, 0},
		{"hard is worth more", 0, "lesson", "hard", 1, 1, 13},
		{"boss is capped", 0, "boss", "hard", 1, 1, 15},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := config.LessonMasteryDelta(tt.current, tt.kind, tt.diff, tt.accuracy, tt.weight)
			require.Equal(t, tt.want, got)
			require.LessOrEqual(t, got, config.MaxMasteryDeltaPerLesson)
			require.LessOrEqual(t, tt.current+got, config.MaxSkillMastery)
		})
	}
}
