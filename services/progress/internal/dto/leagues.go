package dto

// LeagueTierResponse is one rung of the league ladder.
type LeagueTierResponse struct {
	Index int    `json:"index"`
	Slug  string `json:"slug"`
	Name  string `json:"name"`
}

// LeagueStandingResponse is one row of a weekly leaderboard. Zone is what the rank would
// earn if the week ended now: "promote", "stay", or "demote".
type LeagueStandingResponse struct {
	Rank     int    `json:"rank"`
	UserID   string `json:"user_id"`
	WeeklyXP int    `json:"weekly_xp"`
	Zone     string `json:"zone"`
	IsMe     bool   `json:"is_me"`
}

// LeagueResultResponse is how the learner's most recent finished week ended.
type LeagueResultResponse struct {
	WeekStart  string             `json:"week_start"`
	Rank       int                `json:"rank"`
	CohortSize int                `json:"cohort_size"`
	Outcome    string             `json:"outcome"`
	FromTier   LeagueTierResponse `json:"from_tier"`
	ToTier     LeagueTierResponse `json:"to_tier"`
}

// LeagueResponse is the learner's league this week. Until they earn XP this week they have
// not joined: Standings is empty and Tier is the tier they will join.
type LeagueResponse struct {
	WeekStart string                   `json:"week_start"`
	EndsAt    string                   `json:"ends_at"`
	Joined    bool                     `json:"joined"`
	Tier      LeagueTierResponse       `json:"tier"`
	Tiers     []LeagueTierResponse     `json:"tiers"`
	MyRank    int                      `json:"my_rank"`
	Standings []LeagueStandingResponse `json:"standings"`
	Last      *LeagueResultResponse    `json:"last_result"`
}
