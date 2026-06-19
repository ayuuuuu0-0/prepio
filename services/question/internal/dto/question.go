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
	Correct         bool     `json:"correct"`
	Score           int      `json:"score"`
	XPAwarded       int      `json:"xp_awarded"`
	GemsAwarded     int      `json:"gems_awarded"`
	StreakUpdated   bool     `json:"streak_updated"`
	ReadinessDelta  int      `json:"readiness_delta"`
	Feedback        string   `json:"feedback"`
	Strengths       []string `json:"strengths"`
	Gaps            []string `json:"gaps"`
}

// JourneyWorldResponse describes a journey world.
type JourneyWorldResponse struct {
	ID          string `json:"id"`
	Slug        string `json:"slug"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Theme       string `json:"theme"`
}

// JourneyNodeResponse is a node on the journey map with live status.
type JourneyNodeResponse struct {
	ID         string              `json:"id"`
	Slug       string              `json:"slug,omitempty"`
	Label      string              `json:"label"`
	NodeType   string              `json:"node_type"`
	Status     string              `json:"status"`
	QuestionID string              `json:"question_id,omitempty"`
	SortOrder  int                 `json:"sort_order"`
	Skills     []NodeSkillResponse `json:"skills,omitempty"`
	Pools      []NodePoolResponse  `json:"pools,omitempty"`
}

// JourneyResponse is returned by GET /api/v1/journey.
type JourneyResponse struct {
	World     JourneyWorldResponse  `json:"world"`
	Nodes     []JourneyNodeResponse `json:"nodes"`
	SessionID string                `json:"session_id"`
}

// HistoryEntry is a prior answer submission for the authenticated user.
type HistoryEntry struct {
	QuestionID  string `json:"question_id"`
	SessionID   string `json:"session_id"`
	Correct     bool   `json:"correct"`
	Score       int    `json:"score"`
	SubmittedAt string `json:"submitted_at"`
}

// CompanyStats summarizes answer performance for a target company.
type CompanyStats struct {
	Company   string `json:"company"`
	Answered  int    `json:"answered"`
	Correct   int    `json:"correct"`
	ScoreAvg  int    `json:"score_avg"`
}

// ReadinessStats aggregates performance used for readiness computation.
type ReadinessStats struct {
	ByCompany []CompanyStats `json:"by_company"`
}
