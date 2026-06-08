package config

// XP awarded per question by difficulty.
var XPByDifficulty = map[string]int{
	"easy":   20,
	"medium": 50,
	"hard":   80,
}

// XP multiplier applied when a question has a top-tier company tag.
const TopTierCompanyXPMultiplier = 1.5

// TopTierCompanies are FAANG+ tier companies that earn the XP multiplier.
var TopTierCompanies = map[string]bool{
	"google": true,
	"meta":   true,
	"amazon": true,
	"apple":  true,
	"netflix": true,
}

// Gems awarded per correct answer by difficulty.
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

// DailyPaperMaxQuestions is the maximum questions in a daily paper.
const DailyPaperMaxQuestions = 5

// MinimumAnswersForStreak is the minimum submitted answers to qualify for a streak.
const MinimumAnswersForStreak = 1

// WeekendChallengeGemBonus is the gem bonus for completing a weekend challenge.
const WeekendChallengeGemBonus = 50

// DefaultGemBalance is the starting gem balance for new users.
const DefaultGemBalance = 0
