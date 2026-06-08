package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// FreezeStore handles streak_freeze_inventory queries.
type FreezeStore struct {
	pool *pgxpool.Pool
}

// NewFreezeStore creates a FreezeStore.
func NewFreezeStore(pool *pgxpool.Pool) *FreezeStore {
	return &FreezeStore{pool: pool}
}

// Count returns the user's freeze inventory count.
func (s *FreezeStore) Count(ctx context.Context, userID string) (int, error) {
	const q = `SELECT count FROM streak_freeze_inventory WHERE user_id = $1`
	var count int
	err := s.pool.QueryRow(ctx, q, userID).Scan(&count)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, nil
	}
	if err != nil {
		return 0, fmt.Errorf("get freeze count: %w", err)
	}
	return count, nil
}

// SetCount upserts the freeze inventory count.
func (s *FreezeStore) SetCount(ctx context.Context, userID string, count int) error {
	const q = `
		INSERT INTO streak_freeze_inventory (user_id, count)
		VALUES ($1, $2)
		ON CONFLICT (user_id) DO UPDATE SET count = EXCLUDED.count`

	_, err := s.pool.Exec(ctx, q, userID, count)
	if err != nil {
		return fmt.Errorf("set freeze count: %w", err)
	}
	return nil
}
