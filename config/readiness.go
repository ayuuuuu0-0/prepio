package config

// DifficultyMultiplierByDifficulty scales mastery contribution by lesson difficulty.
var DifficultyMultiplierByDifficulty = map[string]float64{
	"easy":   0.90,
	"medium": 1.00,
	"hard":   1.10,
}

// MaxSkillMastery is the upper bound for per-skill mastery scores.
const MaxSkillMastery = 100

// DifficultyMultiplier returns the mastery multiplier for a difficulty band.
func DifficultyMultiplier(difficulty string) float64 {
	if mult, ok := DifficultyMultiplierByDifficulty[difficulty]; ok {
		return mult
	}
	return DifficultyMultiplierByDifficulty["medium"]
}
