package dto

// ProgressResponse is returned by GET /api/v1/progress/me.
type ProgressResponse struct {
	TotalXP        int `json:"total_xp"`
	CurrentLevel   int `json:"current_level"`
	GemBalance     int `json:"gem_balance"`
	XPToNextLevel  int `json:"xp_to_next_level"`
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
