package dto

// ProgressResponse is returned by GET /api/v1/progress/me.
type ProgressResponse struct {
	TotalXP       int `json:"total_xp"`
	CurrentLevel  int `json:"current_level"`
	GemBalance    int `json:"gem_balance"`
	XPToNextLevel int `json:"xp_to_next_level"`
}

// DeductGemsRequest is the body for internal gem deduction.
type DeductGemsRequest struct {
	Amount int    `json:"amount"`
	Reason string `json:"reason"`
}

// DeductGemsResponse is returned after a successful gem deduction.
type DeductGemsResponse struct {
	GemBalance int `json:"gem_balance"`
}

// MasteryChangeResponse explains how one skill's mastery moved.
type MasteryChangeResponse struct {
	SkillSlug string  `json:"skill_slug"`
	SkillName string  `json:"skill_name"`
	TopicSlug string  `json:"topic_slug,omitempty"`
	TopicName string  `json:"topic_name,omitempty"`
	Before    int     `json:"before"`
	After     int     `json:"after"`
	Delta     int     `json:"delta"`
	Accuracy  float64 `json:"accuracy"`
}

// AttemptRewardsResponse is returned by GET /api/v1/progress/attempts/{attemptID}/rewards.
type AttemptRewardsResponse struct {
	AttemptID       string                  `json:"attempt_id"`
	FirstCompletion bool                    `json:"first_completion"`
	XPAwarded       int                     `json:"xp_awarded"`
	GemsAwarded     int                     `json:"gems_awarded"`
	MasteryChanges  []MasteryChangeResponse `json:"mastery_changes"`
	// AchievementsUnlocked lists achievements this completion earned (empty for most lessons).
	AchievementsUnlocked []AchievementResponse `json:"achievements_unlocked"`
}
