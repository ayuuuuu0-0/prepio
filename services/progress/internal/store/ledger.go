package store

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// LedgerStore handles xp and gem ledger entries.
type LedgerStore struct {
	pool *pgxpool.Pool
}

// NewLedgerStore creates a LedgerStore.
func NewLedgerStore(pool *pgxpool.Pool) *LedgerStore {
	return &LedgerStore{pool: pool}
}

// InsertXP records an XP ledger entry.
func (s *LedgerStore) InsertXP(ctx context.Context, userID string, amount int, reason, eventID string) error {
	const q = `INSERT INTO xp_ledger (user_id, amount, reason, source_event_id) VALUES ($1, $2, $3, $4)`
	_, err := s.pool.Exec(ctx, q, userID, amount, reason, eventID)
	if err != nil {
		return fmt.Errorf("insert xp ledger: %w", err)
	}
	return nil
}

// InsertGem records a gem ledger entry.
func (s *LedgerStore) InsertGem(ctx context.Context, userID string, amount int, reason, eventID string) error {
	const q = `INSERT INTO gem_ledger (user_id, amount, reason, source_event_id) VALUES ($1, $2, $3, $4)`
	_, err := s.pool.Exec(ctx, q, userID, amount, reason, eventID)
	if err != nil {
		return fmt.Errorf("insert gem ledger: %w", err)
	}
	return nil
}
