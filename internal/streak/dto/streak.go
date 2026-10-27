package dto

// StreakResponse is returned by GET /api/v1/streaks/me.
type StreakResponse struct {
	CurrentStreak     int    `json:"current_streak"`
	LongestStreak     int    `json:"longest_streak"`
	FreezeCount       int    `json:"freeze_count"`
	LastActivityDate  string `json:"last_activity_date,omitempty"`
	StreakActiveToday bool   `json:"streak_active_today"`
}

// FreezePurchaseResponse is returned after buying a streak freeze.
type FreezePurchaseResponse struct {
	FreezeCount int `json:"freeze_count"`
	GemsSpent   int `json:"gems_spent"`
}
