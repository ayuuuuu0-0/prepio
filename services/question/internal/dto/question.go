package dto

// DailyPaperResponse is returned by GET /api/v1/questions/daily.
type DailyPaperResponse struct {
	SessionID         string             `json:"session_id"`
	Date              string             `json:"date"`
	Questions         []QuestionResponse `json:"questions"`
	MinimumToStreak   int                `json:"minimum_to_streak"`
}

// QuestionResponse is the public question shape (no answer guide).
type QuestionResponse struct {
	ID          string   `json:"id"`
	Body        string   `json:"body"`
	RoundType   string   `json:"round_type"`
	Difficulty  string   `json:"difficulty"`
	CompanyTags []string `json:"company_tags"`
	IsWeekend   bool     `json:"is_weekend"`
}

// SubmitRequest is the body for POST /api/v1/questions/{id}/submit.
type SubmitRequest struct {
	SessionID        string `json:"session_id"`
	Answer           string `json:"answer"`
	TimeSpentSeconds int    `json:"time_spent_seconds"`
	SubmittedAt      string `json:"submitted_at"`
}

// SubmitResponse is returned after answer submission.
type SubmitResponse struct {
	Correct       bool   `json:"correct"`
	XPAwarded     int    `json:"xp_awarded"`
	GemsAwarded   int    `json:"gems_awarded"`
	StreakUpdated bool   `json:"streak_updated"`
	Feedback      string `json:"feedback"`
}
