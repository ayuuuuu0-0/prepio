package store

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// StreakState is a row from user_streaks.
type StreakState struct {
	UserID           string
	CurrentStreak    int
	LongestStreak    int
	LastActivityDate *time.Time
}

// StreakStore handles user_streaks queries.
type StreakStore struct {
	pool *pgxpool.Pool
}

// NewStreakStore creates a StreakStore.
func NewStreakStore(pool *pgxpool.Pool) *StreakStore {
	return &StreakStore{pool: pool}
}

// Get returns streak state, creating a zero row if missing.
func (s *StreakStore) Get(ctx context.Context, userID string) (*StreakState, error) {
	const q = `
		SELECT user_id, current_streak, longest_streak, last_activity_date
		FROM user_streaks WHERE user_id = $1`

	var state StreakState
	var last *time.Time
	err := s.pool.QueryRow(ctx, q, userID).Scan(
		&state.UserID, &state.CurrentStreak, &state.LongestStreak, &last,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return &StreakState{UserID: userID}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get streak: %w", err)
	}
	state.LastActivityDate = last
	return &state, nil
}

// Upsert writes streak state to PostgreSQL.
func (s *StreakStore) Upsert(ctx context.Context, state StreakState) error {
	const q = `
		INSERT INTO user_streaks (user_id, current_streak, longest_streak, last_activity_date)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id) DO UPDATE SET
			current_streak = EXCLUDED.current_streak,
			longest_streak = EXCLUDED.longest_streak,
			last_activity_date = EXCLUDED.last_activity_date`

	_, err := s.pool.Exec(ctx, q, state.UserID, state.CurrentStreak, state.LongestStreak, state.LastActivityDate)
	if err != nil {
		return fmt.Errorf("upsert streak: %w", err)
	}
	return nil
}

// Timezone returns the user's timezone.
func (s *StreakStore) Timezone(ctx context.Context, userID string) (string, error) {
	const q = `SELECT timezone FROM users WHERE id = $1`
	var timezone string
	err := s.pool.QueryRow(ctx, q, userID).Scan(&timezone)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("get timezone: %w", err)
	}
	return timezone, nil
}
