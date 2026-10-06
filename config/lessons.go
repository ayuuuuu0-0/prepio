package config

import "math"

// Lesson kinds.
const (
	LessonKindLesson = "lesson"
	LessonKindBoss   = "boss"
)

// LessonXPByDifficulty is the base XP for completing a lesson for the first time.
var LessonXPByDifficulty = map[string]int{
	"easy":   20,
	"medium": 35,
	"hard":   50,
}

// BossXPMultiplier scales lesson XP for boss lessons.
const BossXPMultiplier = 1.5

// LessonMinXPAccuracyFactor is the share of base XP earned at zero first-try accuracy.
// XP scales linearly from this floor up to the full base at perfect accuracy.
const LessonMinXPAccuracyFactor = 0.6

// LessonGemAccuracyThreshold is the first-try accuracy at or above which full gems are awarded.
const LessonGemAccuracyThreshold = 0.8

// LessonMasteryGainMax is the mastery points a perfect, full-weight, medium lesson would add at mastery 0.
const LessonMasteryGainMax = 12.0

// MaxMasteryDeltaPerLesson caps the mastery change any single skill can receive from one lesson.
const MaxMasteryDeltaPerLesson = 15

// BossMasteryMultiplier scales the mastery gain of boss lessons.
const BossMasteryMultiplier = 1.25

// LessonXP returns the XP awarded for a first completion.
func LessonXP(kind, difficulty string, accuracy float64) int {
	base, ok := LessonXPByDifficulty[difficulty]
	if !ok {
		base = LessonXPByDifficulty["medium"]
	}
	accuracy = clamp01(accuracy)
	xp := float64(base) * (LessonMinXPAccuracyFactor + (1-LessonMinXPAccuracyFactor)*accuracy)
	if kind == LessonKindBoss {
		xp *= BossXPMultiplier
	}
	return int(math.Round(xp))
}

// LessonGems returns the gems awarded for a first completion.
func LessonGems(difficulty string, accuracy float64) int {
	base, ok := GemsByDifficulty[difficulty]
	if !ok {
		base = GemsByDifficulty["medium"]
	}
	if clamp01(accuracy) < LessonGemAccuracyThreshold {
		base /= 2
	}
	if base < 1 {
		base = 1
	}
	return base
}

// LessonXPPreview is the maximum XP a first completion can earn, shown on the journey preview.
func LessonXPPreview(kind, difficulty string) int {
	return LessonXP(kind, difficulty, 1)
}

// LessonMasteryDelta returns the capped, smoothed mastery change for one skill.
// Gains shrink as mastery approaches the maximum so mastery never jumps.
func LessonMasteryDelta(current int, kind, difficulty string, accuracy, skillWeight float64) int {
	if current >= MaxSkillMastery || accuracy <= 0 || skillWeight <= 0 {
		return 0
	}
	gain := LessonMasteryGainMax * clamp01(accuracy) * skillWeight * DifficultyMultiplier(difficulty)
	if kind == LessonKindBoss {
		gain *= BossMasteryMultiplier
	}
	headroom := float64(MaxSkillMastery-current) / float64(MaxSkillMastery)
	delta := int(math.Round(gain * headroom))
	if delta > MaxMasteryDeltaPerLesson {
		delta = MaxMasteryDeltaPerLesson
	}
	if delta > MaxSkillMastery-current {
		delta = MaxSkillMastery - current
	}
	return delta
}

func clamp01(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}
