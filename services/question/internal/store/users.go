package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// UserStore provides read-only user lookups for the question service.
type UserStore struct {
	pool *pgxpool.Pool
}

// NewUserStore creates a UserStore.
func NewUserStore(pool *pgxpool.Pool) *UserStore {
	return &UserStore{pool: pool}
}

// Timezone returns the user's timezone or empty if not found.
func (s *UserStore) Timezone(ctx context.Context, userID string) (string, error) {
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

// Level returns the user's current level, defaulting to 1.
func (s *UserStore) Level(ctx context.Context, userID string) (int, error) {
	const q = `SELECT current_level FROM user_progress WHERE user_id = $1`
	var level int
	err := s.pool.QueryRow(ctx, q, userID).Scan(&level)
	if errors.Is(err, pgx.ErrNoRows) {
		return 1, nil
	}
	if err != nil {
		return 0, fmt.Errorf("get level: %w", err)
	}
	return level, nil
}
