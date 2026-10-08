package store

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prepio/prepio/config"
)

// querier is what both a pool and a transaction can do; league reads run on either.
type querier interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
}

// LeagueStanding is one learner's row on a weekly leaderboard.
type LeagueStanding struct {
	UserID   string
	WeeklyXP int
}

// LeagueResult is how a finished week ended for a learner.
type LeagueResult struct {
	WeekStart  time.Time
	Tier       int
	Rank       int
	CohortSize int
	Outcome    string
	NextTier   int
}

// LeagueView is a learner's league state for one week.
type LeagueView struct {
	WeekStart time.Time
	Joined    bool
	Tier      int
	Standings []LeagueStanding // ordered best first; empty until joined
	Last      *LeagueResult    // the most recent finished week, if any
}

// LeagueStore reads league standings.
type LeagueStore struct {
	pool *pgxpool.Pool
}

// NewLeagueStore creates a LeagueStore.
func NewLeagueStore(pool *pgxpool.Pool) *LeagueStore {
	return &LeagueStore{pool: pool}
}

// Get returns the learner's league for the week containing now.
func (s *LeagueStore) Get(ctx context.Context, userID string, now time.Time) (*LeagueView, error) {
	week := config.LeagueWeekStart(now)
	view := &LeagueView{WeekStart: week}

	last, err := lastResult(ctx, s.pool, userID, week)
	if err != nil {
		return nil, err
	}
	view.Last = last
	if last != nil {
		view.Tier = last.NextTier
	}

	var cohortID string
	var tier int
	err = s.pool.QueryRow(ctx,
		`SELECT cohort_id, tier FROM league_memberships WHERE user_id = $1 AND week_start = $2`,
		userID, week).Scan(&cohortID, &tier)
	if errors.Is(err, pgx.ErrNoRows) {
		return view, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get league membership: %w", err)
	}
	view.Joined = true
	view.Tier = config.ClampLeagueTier(tier)

	rows, err := s.pool.Query(ctx, `
		SELECT user_id, weekly_xp FROM league_memberships
		WHERE cohort_id = $1
		ORDER BY weekly_xp DESC, last_xp_at, user_id`, cohortID)
	if err != nil {
		return nil, fmt.Errorf("list league standings: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var st LeagueStanding
		if err := rows.Scan(&st.UserID, &st.WeeklyXP); err != nil {
			return nil, fmt.Errorf("scan league standing: %w", err)
		}
		view.Standings = append(view.Standings, st)
	}
	return view, rows.Err()
}

// addLeagueXP credits XP to the learner's league for the week containing at, joining
// them to a cohort first if this is their first XP of the week. It runs inside the
// transaction that granted the XP, after the user's progress row is locked, so one
// learner's joins are serialised and the XP is counted exactly once.
func addLeagueXP(ctx context.Context, tx pgx.Tx, userID string, xp int, at time.Time) error {
	if xp <= 0 {
		return nil
	}
	week := config.LeagueWeekStart(at)

	tag, err := tx.Exec(ctx, `
		UPDATE league_memberships SET weekly_xp = weekly_xp + $3, last_xp_at = $4
		WHERE user_id = $1 AND week_start = $2`, userID, week, xp, at)
	if err != nil {
		return fmt.Errorf("add league xp: %w", err)
	}
	if tag.RowsAffected() == 1 {
		return nil
	}

	tier := 0
	last, err := lastResult(ctx, tx, userID, week)
	if err != nil {
		return err
	}
	if last != nil {
		tier = last.NextTier
	}

	cohortID, err := openCohort(ctx, tx, week, tier)
	if err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, `
		INSERT INTO league_memberships (user_id, week_start, cohort_id, tier, weekly_xp, last_xp_at)
		VALUES ($1, $2, $3, $4, $5, $6)`, userID, week, cohortID, tier, xp, at); err != nil {
		return fmt.Errorf("join league: %w", err)
	}
	return nil
}

// openCohort returns the oldest cohort for the week and tier with room left, creating one
// when all are full. Joins to the same week and tier are serialised so a cohort never
// overfills.
func openCohort(ctx context.Context, tx pgx.Tx, week time.Time, tier int) (string, error) {
	if _, err := tx.Exec(ctx, `SELECT pg_advisory_xact_lock(hashtextextended($1, 0))`,
		fmt.Sprintf("league:%s:%d", week.Format("2006-01-02"), tier)); err != nil {
		return "", fmt.Errorf("lock league cohort: %w", err)
	}

	var cohortID string
	err := tx.QueryRow(ctx, `
		SELECT c.id FROM league_cohorts c
		WHERE c.week_start = $1 AND c.tier = $2
		  AND (SELECT count(*) FROM league_memberships m WHERE m.cohort_id = c.id) < $3
		ORDER BY c.created_at, c.id
		LIMIT 1`, week, tier, config.LeagueCohortSize).Scan(&cohortID)
	if err == nil {
		return cohortID, nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return "", fmt.Errorf("find league cohort: %w", err)
	}
	if err := tx.QueryRow(ctx,
		`INSERT INTO league_cohorts (week_start, tier) VALUES ($1, $2) RETURNING id`,
		week, tier).Scan(&cohortID); err != nil {
		return "", fmt.Errorf("create league cohort: %w", err)
	}
	return cohortID, nil
}

// lastResult is how the learner's most recent week before week ended. Past weeks no
// longer change (XP always lands in the week it is granted), so the result is stable.
// A learner who skipped weeks keeps the tier their last week earned.
func lastResult(ctx context.Context, q querier, userID string, week time.Time) (*LeagueResult, error) {
	var (
		res      LeagueResult
		cohortID string
		xp       int
		lastXPAt time.Time
	)
	err := q.QueryRow(ctx, `
		SELECT week_start, tier, cohort_id, weekly_xp, last_xp_at FROM league_memberships
		WHERE user_id = $1 AND week_start < $2
		ORDER BY week_start DESC LIMIT 1`, userID, week).Scan(&res.WeekStart, &res.Tier, &cohortID, &xp, &lastXPAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get last league week: %w", err)
	}

	if err := q.QueryRow(ctx, `
		SELECT 1 + count(*) FILTER (WHERE weekly_xp > $2
		                              OR (weekly_xp = $2 AND last_xp_at < $3)
		                              OR (weekly_xp = $2 AND last_xp_at = $3 AND user_id < $4)),
		       count(*)
		FROM league_memberships WHERE cohort_id = $1`,
		cohortID, xp, lastXPAt, userID).Scan(&res.Rank, &res.CohortSize); err != nil {
		return nil, fmt.Errorf("rank last league week: %w", err)
	}

	res.Tier = config.ClampLeagueTier(res.Tier)
	res.Outcome, res.NextTier = config.LeagueOutcome(res.Rank, res.CohortSize, res.Tier)
	return &res, nil
}
