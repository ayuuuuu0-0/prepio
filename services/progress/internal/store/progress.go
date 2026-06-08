package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ProgressState is a row from user_progress.
type ProgressState struct {
	UserID       string
	TotalXP      int
	CurrentLevel int
	GemBalance   int
}

// ProgressStore handles user_progress queries.
type ProgressStore struct {
	pool *pgxpool.Pool
}

// NewProgressStore creates a ProgressStore.
func NewProgressStore(pool *pgxpool.Pool) *ProgressStore {
	return &ProgressStore{pool: pool}
}

// Get returns progress state, defaulting to level 1 if missing.
func (s *ProgressStore) Get(ctx context.Context, userID string) (*ProgressState, error) {
	const q = `SELECT user_id, total_xp, current_level, gem_balance FROM user_progress WHERE user_id = $1`
	var state ProgressState
	err := s.pool.QueryRow(ctx, q, userID).Scan(&state.UserID, &state.TotalXP, &state.CurrentLevel, &state.GemBalance)
	if errors.Is(err, pgx.ErrNoRows) {
		return &ProgressState{UserID: userID, CurrentLevel: 1}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get progress: %w", err)
	}
	return &state, nil
}

// Upsert writes progress state.
func (s *ProgressStore) Upsert(ctx context.Context, state ProgressState) error {
	const q = `
		INSERT INTO user_progress (user_id, total_xp, current_level, gem_balance)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id) DO UPDATE SET
			total_xp = EXCLUDED.total_xp,
			current_level = EXCLUDED.current_level,
			gem_balance = EXCLUDED.gem_balance`

	_, err := s.pool.Exec(ctx, q, state.UserID, state.TotalXP, state.CurrentLevel, state.GemBalance)
	if err != nil {
		return fmt.Errorf("upsert progress: %w", err)
	}
	return nil
}

// DeductGems atomically deducts gems with optimistic locking on balance.
func (s *ProgressStore) DeductGems(ctx context.Context, userID string, amount int) (int, error) {
	const q = `
		UPDATE user_progress
		SET gem_balance = gem_balance - $2
		WHERE user_id = $1 AND gem_balance >= $2
		RETURNING gem_balance`

	var balance int
	err := s.pool.QueryRow(ctx, q, userID, amount).Scan(&balance)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, ErrInsufficientGems
	}
	if err != nil {
		return 0, fmt.Errorf("deduct gems: %w", err)
	}
	return balance, nil
}
