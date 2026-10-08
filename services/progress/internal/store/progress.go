package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prepio/prepio/config"
)

// ProgressState is a row from user_progress.
type ProgressState struct {
	UserID       string
	TotalXP      int
	CurrentLevel int
	GemBalance   int
}

// ProgressStore handles user_progress queries. Every write locks the user's row and
// changes balances relative to what is stored, so concurrent consumers (lesson and
// streak events) never overwrite each other's XP or gems.
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

// AwardGems adds gems for an event. It is idempotent per (event, reason): a redelivered
// event changes nothing and returns applied=false.
func (s *ProgressStore) AwardGems(ctx context.Context, userID string, amount int, reason, eventID string) (state *ProgressState, applied bool, err error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, false, fmt.Errorf("begin award gems: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	state, err = lockProgress(ctx, tx, userID)
	if err != nil {
		return nil, false, err
	}

	var seen bool
	if err := tx.QueryRow(ctx,
		`SELECT EXISTS (SELECT 1 FROM gem_ledger WHERE user_id = $1 AND source_event_id = $2 AND reason = $3)`,
		userID, eventID, reason).Scan(&seen); err != nil {
		return nil, false, fmt.Errorf("check gem award: %w", err)
	}
	if seen {
		return state, false, nil
	}

	state.GemBalance += amount
	if _, err := tx.Exec(ctx, `UPDATE user_progress SET gem_balance = $2 WHERE user_id = $1`, userID, state.GemBalance); err != nil {
		return nil, false, fmt.Errorf("award gems: %w", err)
	}
	if _, err := tx.Exec(ctx,
		`INSERT INTO gem_ledger (user_id, amount, reason, source_event_id) VALUES ($1, $2, $3, $4)`,
		userID, amount, reason, eventID); err != nil {
		return nil, false, fmt.Errorf("insert gem ledger: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, false, fmt.Errorf("commit award gems: %w", err)
	}
	return state, true, nil
}

// DeductGems removes gems if the balance covers them and records the spend.
func (s *ProgressStore) DeductGems(ctx context.Context, userID string, amount int, reason string) (int, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("begin deduct gems: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	state, err := lockProgress(ctx, tx, userID)
	if err != nil {
		return 0, err
	}
	if state.GemBalance < amount {
		return 0, ErrInsufficientGems
	}

	state.GemBalance -= amount
	if _, err := tx.Exec(ctx, `UPDATE user_progress SET gem_balance = $2 WHERE user_id = $1`, userID, state.GemBalance); err != nil {
		return 0, fmt.Errorf("deduct gems: %w", err)
	}
	if _, err := tx.Exec(ctx,
		`INSERT INTO gem_ledger (user_id, amount, reason, source_event_id) VALUES ($1, $2, $3, $4)`,
		userID, -amount, reason, uuid.NewString()); err != nil {
		return 0, fmt.Errorf("insert gem ledger: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("commit deduct gems: %w", err)
	}
	return state.GemBalance, nil
}

// lockProgress makes sure the user's progress row exists and locks it for the rest of
// the transaction. Creating the row first matters: FOR UPDATE on a missing row locks
// nothing, which would let two first writes race.
func lockProgress(ctx context.Context, tx pgx.Tx, userID string) (*ProgressState, error) {
	if _, err := tx.Exec(ctx,
		`INSERT INTO user_progress (user_id, total_xp, current_level, gem_balance) VALUES ($1, 0, $2, 0)
		 ON CONFLICT (user_id) DO NOTHING`, userID, config.CurrentLevel(0)); err != nil {
		return nil, fmt.Errorf("ensure user progress: %w", err)
	}
	state := &ProgressState{UserID: userID}
	if err := tx.QueryRow(ctx,
		`SELECT total_xp, current_level, gem_balance FROM user_progress WHERE user_id = $1 FOR UPDATE`,
		userID).Scan(&state.TotalXP, &state.CurrentLevel, &state.GemBalance); err != nil {
		return nil, fmt.Errorf("lock user progress: %w", err)
	}
	return state, nil
}
