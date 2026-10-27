package store

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prepio/prepio/config"
	"github.com/prepio/prepio/test/testdb"
	"github.com/stretchr/testify/require"
)

// week1 is a Monday; league weeks start Monday 00:00 UTC.
var week1 = time.Date(2026, 10, 5, 0, 0, 0, 0, time.UTC)

func leaguePool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	pool, _ := testdb.Start(t)
	testdb.Migrate(t, pool)
	return pool
}

func leagueUser(t *testing.T, pool *pgxpool.Pool, n int) string {
	t.Helper()
	var id string
	require.NoError(t, pool.QueryRow(context.Background(), `
		INSERT INTO users (email, username, password_hash) VALUES ($1, $2, 'hash') RETURNING id`,
		fmt.Sprintf("league%d@test.com", n), fmt.Sprintf("league%d", n)).Scan(&id))
	return id
}

func earn(t *testing.T, pool *pgxpool.Pool, userID string, xp int, at time.Time) {
	t.Helper()
	ctx := context.Background()
	tx, err := pool.Begin(ctx)
	require.NoError(t, err)
	defer func() { _ = tx.Rollback(ctx) }()
	_, err = lockProgress(ctx, tx, userID)
	require.NoError(t, err)
	require.NoError(t, addLeagueXP(ctx, tx, userID, xp, at))
	require.NoError(t, tx.Commit(ctx))
}

func TestLeagueJoinRankAndPromotion(t *testing.T) {
	pool := leaguePool(t)
	ls := NewLeagueStore(pool)
	ctx := context.Background()

	a, b, c := leagueUser(t, pool, 1), leagueUser(t, pool, 2), leagueUser(t, pool, 3)

	// Not joined until XP is earned; zero XP never joins.
	earn(t, pool, a, 0, week1.Add(time.Hour))
	view, err := ls.Get(ctx, a, week1.Add(time.Hour))
	require.NoError(t, err)
	require.False(t, view.Joined)
	require.Equal(t, 0, view.Tier)
	require.Empty(t, view.Standings)

	// b and c tie on XP; c got there first, so c ranks above b.
	earn(t, pool, a, 20, week1.Add(2*time.Hour))
	earn(t, pool, c, 15, week1.Add(3*time.Hour))
	earn(t, pool, b, 10, week1.Add(4*time.Hour))
	earn(t, pool, b, 5, week1.Add(5*time.Hour))
	earn(t, pool, a, 10, week1.Add(6*time.Hour))

	view, err = ls.Get(ctx, b, week1.Add(7*time.Hour))
	require.NoError(t, err)
	require.True(t, view.Joined)
	require.Equal(t, []LeagueStanding{{a, 30}, {c, 15}, {b, 15}}, view.Standings)
	require.Nil(t, view.Last)

	// Next week: before b plays they see last week's result and the tier they'll join.
	week2 := week1.AddDate(0, 0, 7)
	view, err = ls.Get(ctx, b, week2.Add(time.Hour))
	require.NoError(t, err)
	require.False(t, view.Joined)
	require.NotNil(t, view.Last)
	require.Equal(t, 3, view.Last.Rank)
	require.Equal(t, 3, view.Last.CohortSize)
	require.Equal(t, config.LeaguePromoted, view.Last.Outcome) // a small cohort promotes its top ranks
	require.Equal(t, 1, view.Tier)

	// Earning XP joins week 2 at the promoted tier.
	earn(t, pool, b, 12, week2.Add(2*time.Hour))
	view, err = ls.Get(ctx, b, week2.Add(3*time.Hour))
	require.NoError(t, err)
	require.True(t, view.Joined)
	require.Equal(t, 1, view.Tier)
	require.Equal(t, []LeagueStanding{{b, 12}}, view.Standings)
}

func TestLeagueSkippedWeekKeepsEarnedTier(t *testing.T) {
	pool := leaguePool(t)
	ls := NewLeagueStore(pool)
	ctx := context.Background()
	a := leagueUser(t, pool, 1)

	earn(t, pool, a, 10, week1.Add(time.Hour)) // promoted to tier 1 at week end
	week4 := week1.AddDate(0, 0, 21)
	earn(t, pool, a, 10, week4.Add(time.Hour))

	view, err := ls.Get(ctx, a, week4.Add(2*time.Hour))
	require.NoError(t, err)
	require.Equal(t, 1, view.Tier)
	require.Equal(t, week1, view.Last.WeekStart.UTC())
}

func TestLeagueCohortsFillThenSplit(t *testing.T) {
	pool := leaguePool(t)
	ctx := context.Background()

	for i := 0; i < config.LeagueCohortSize+1; i++ {
		earn(t, pool, leagueUser(t, pool, i), 5, week1.Add(time.Duration(i)*time.Minute))
	}

	var cohorts []int
	rows, err := pool.Query(ctx, `
		SELECT count(*) FROM league_memberships GROUP BY cohort_id ORDER BY count(*) DESC`)
	require.NoError(t, err)
	defer rows.Close()
	for rows.Next() {
		var n int
		require.NoError(t, rows.Scan(&n))
		cohorts = append(cohorts, n)
	}
	require.NoError(t, rows.Err())
	require.Equal(t, []int{config.LeagueCohortSize, 1}, cohorts)
}
