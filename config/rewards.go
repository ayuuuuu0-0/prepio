package config

// Gems awarded for a first lesson completion, by difficulty (halved below the accuracy threshold).
var GemsByDifficulty = map[string]int{
	"easy":   5,
	"medium": 10,
	"hard":   15,
}

// Gems awarded when a streak increments for the day.
const StreakIncrementGemBonus = 5

// StreakFreezeGemCost is the gem price to purchase one streak freeze.
const StreakFreezeGemCost = 100

// MaxStreakFreezes is the maximum freezes a user can hold at once.
const MaxStreakFreezes = 2
