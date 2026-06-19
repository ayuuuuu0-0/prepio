package config

// ReadinessWeightByDifficulty maps question difficulty to default readiness weight.
var ReadinessWeightByDifficulty = map[string]float64{
	"easy":   0.80,
	"medium": 1.00,
	"hard":   1.20,
}

// EstimatedTimeMinutesByDifficulty maps question difficulty to default completion time in minutes.
var EstimatedTimeMinutesByDifficulty = map[string]int{
	"easy":   8,
	"medium": 15,
	"hard":   25,
}

// DifficultyMultiplierByDifficulty scales mastery contribution by question difficulty.
var DifficultyMultiplierByDifficulty = map[string]float64{
	"easy":   0.90,
	"medium": 1.00,
	"hard":   1.10,
}

// MasterySmoothingFactor controls how much each answer moves skill mastery (0–1).
const MasterySmoothingFactor = 0.15

// MaxSkillMastery is the upper bound for per-skill mastery scores.
const MaxSkillMastery = 100

// MaxCompanyReadiness caps company readiness to avoid implying certainty.
const MaxCompanyReadiness = 95

// ReadinessSourceLive marks scores updated from live answer events.
const ReadinessSourceLive = "live"

// ReadinessSourceBackfill marks scores populated from historical data.
const ReadinessSourceBackfill = "backfill"

// DifficultyMultiplier returns the mastery multiplier for a difficulty band.
func DifficultyMultiplier(difficulty string) float64 {
	if mult, ok := DifficultyMultiplierByDifficulty[difficulty]; ok {
		return mult
	}
	return DifficultyMultiplierByDifficulty["medium"]
}
