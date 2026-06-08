package store

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// TargetStore handles user_targets table queries.
type TargetStore struct {
	pool *pgxpool.Pool
}

// NewTargetStore creates a TargetStore.
func NewTargetStore(pool *pgxpool.Pool) *TargetStore {
	return &TargetStore{pool: pool}
}

// Replace sets the user's target companies, replacing any existing rows.
func (s *TargetStore) Replace(ctx context.Context, userID string, companies []string) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, `DELETE FROM user_targets WHERE user_id = $1`, userID); err != nil {
		return fmt.Errorf("delete targets: %w", err)
	}

	for _, company := range companies {
		if _, err := tx.Exec(ctx,
			`INSERT INTO user_targets (user_id, company) VALUES ($1, $2)`,
			userID, company,
		); err != nil {
			return fmt.Errorf("insert target: %w", err)
		}
	}

	return tx.Commit(ctx)
}

// List returns target companies for a user.
func (s *TargetStore) List(ctx context.Context, userID string) ([]string, error) {
	const q = `SELECT company FROM user_targets WHERE user_id = $1 ORDER BY company`

	rows, err := s.pool.Query(ctx, q, userID)
	if err != nil {
		return nil, fmt.Errorf("list targets: %w", err)
	}
	defer rows.Close()

	var companies []string
	for rows.Next() {
		var company string
		if err := rows.Scan(&company); err != nil {
			return nil, fmt.Errorf("scan target: %w", err)
		}
		companies = append(companies, company)
	}
	return companies, rows.Err()
}
