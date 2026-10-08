package config

import (
	"testing"
	"time"
)

func TestLeagueWeekStart(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"2026-10-05T00:00:00Z", "2026-10-05"},      // Monday midnight
		{"2026-10-07T13:45:00Z", "2026-10-05"},      // Wednesday
		{"2026-10-11T23:59:59Z", "2026-10-05"},      // Sunday last second
		{"2026-10-12T00:00:00Z", "2026-10-12"},      // next Monday
		{"2026-10-12T03:00:00+05:30", "2026-10-05"}, // still Sunday in UTC
	}
	for _, c := range cases {
		in, err := time.Parse(time.RFC3339, c.in)
		if err != nil {
			t.Fatal(err)
		}
		if got := LeagueWeekStart(in).Format("2006-01-02"); got != c.want {
			t.Errorf("LeagueWeekStart(%s) = %s, want %s", c.in, got, c.want)
		}
	}
}

func TestLeagueOutcome(t *testing.T) {
	top := len(LeagueTiers) - 1
	cases := []struct {
		name             string
		rank, size, tier int
		wantOutcome      string
		wantNext         int
	}{
		{"top rank promotes", 1, 30, 0, LeaguePromoted, 1},
		{"last promotion rank", LeaguePromoteCount, 30, 3, LeaguePromoted, 4},
		{"just below promotion stays", LeaguePromoteCount + 1, 30, 3, LeagueStayed, 3},
		{"bottom rank demotes", 30, 30, 3, LeagueDemoted, 2},
		{"first demotion rank", 30 - LeagueDemoteCount + 1, 30, 3, LeagueDemoted, 2},
		{"just above demotion stays", 30 - LeagueDemoteCount, 30, 3, LeagueStayed, 3},
		{"no demotion from the lowest tier", 30, 30, 0, LeagueStayed, 0},
		{"no promotion past the top tier", 1, 30, top, LeagueStayed, top},
		{"small cohort never demotes", 5, 5, 3, LeaguePromoted, 4},
		{"cohort at the threshold never demotes", 12, LeaguePromoteCount + LeagueDemoteCount, 3, LeagueStayed, 3},
		{"one above the threshold demotes its last", 13, LeaguePromoteCount + LeagueDemoteCount + 1, 3, LeagueDemoted, 2},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			outcome, next := LeagueOutcome(c.rank, c.size, c.tier)
			if outcome != c.wantOutcome || next != c.wantNext {
				t.Errorf("LeagueOutcome(%d, %d, %d) = %s, %d; want %s, %d", c.rank, c.size, c.tier, outcome, next, c.wantOutcome, c.wantNext)
			}
		})
	}
}

func TestClampLeagueTier(t *testing.T) {
	if ClampLeagueTier(-1) != 0 || ClampLeagueTier(len(LeagueTiers)+3) != len(LeagueTiers)-1 || ClampLeagueTier(2) != 2 {
		t.Fatal("ClampLeagueTier does not keep tiers on the ladder")
	}
}
