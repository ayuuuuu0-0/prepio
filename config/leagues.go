package config

import "time"

// LeagueTier is one rung of the weekly league ladder. Index 0 is where everyone starts.
type LeagueTier struct {
	Slug string
	Name string
}

// LeagueTiers is the ladder, lowest first.
var LeagueTiers = []LeagueTier{
	{Slug: "bronze", Name: "Bronze"},
	{Slug: "silver", Name: "Silver"},
	{Slug: "gold", Name: "Gold"},
	{Slug: "platinum", Name: "Platinum"},
	{Slug: "sapphire", Name: "Sapphire"},
	{Slug: "ruby", Name: "Ruby"},
	{Slug: "emerald", Name: "Emerald"},
	{Slug: "amethyst", Name: "Amethyst"},
	{Slug: "obsidian", Name: "Obsidian"},
	{Slug: "diamond", Name: "Diamond"},
}

const (
	// LeagueCohortSize is the most learners grouped into one weekly leaderboard.
	LeagueCohortSize = 30
	// LeaguePromoteCount is how many of the top ranks move up a tier at week end.
	LeaguePromoteCount = 7
	// LeagueDemoteCount is how many of the bottom ranks move down a tier at week end.
	// Demotion only applies when the cohort is big enough that nobody can be both
	// promoted and demoted (more than LeaguePromoteCount+LeagueDemoteCount members).
	LeagueDemoteCount = 5
)

// League outcomes at the end of a week.
const (
	LeaguePromoted = "promoted"
	LeagueStayed   = "stayed"
	LeagueDemoted  = "demoted"
)

// LeagueWeekStart returns the start of the league week containing t: Monday 00:00 UTC.
func LeagueWeekStart(t time.Time) time.Time {
	t = t.UTC()
	day := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
	offset := (int(day.Weekday()) + 6) % 7 // Monday = 0
	return day.AddDate(0, 0, -offset)
}

// LeagueOutcome decides what a final rank (1-based) in a cohort of size members at the
// given tier earns, and the tier the learner plays in next.
func LeagueOutcome(rank, size, tier int) (outcome string, nextTier int) {
	top := len(LeagueTiers) - 1
	switch {
	case rank <= LeaguePromoteCount && tier < top:
		return LeaguePromoted, tier + 1
	case tier > 0 && size > LeaguePromoteCount+LeagueDemoteCount && rank > size-LeagueDemoteCount:
		return LeagueDemoted, tier - 1
	default:
		return LeagueStayed, tier
	}
}

// ClampLeagueTier keeps a stored tier on the configured ladder (the ladder may shrink).
func ClampLeagueTier(tier int) int {
	if tier < 0 {
		return 0
	}
	if tier >= len(LeagueTiers) {
		return len(LeagueTiers) - 1
	}
	return tier
}
