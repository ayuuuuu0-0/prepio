package dto

// AchievementResponse is one achievement in the catalog. Description says how to earn it;
// UnlockedAt is set once the learner has earned it.
type AchievementResponse struct {
	Slug        string  `json:"slug"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Unlocked    bool    `json:"unlocked"`
	UnlockedAt  *string `json:"unlocked_at,omitempty"`
}
