package store

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// RefreshToken is a row from the refresh_tokens table.
type RefreshToken struct {
	ID        string
	UserID    string
	TokenHash string
	ExpiresAt time.Time
	UsedAt    *time.Time
	CreatedAt time.Time
}

// RefreshTokenStore handles refresh_tokens table queries.
type RefreshTokenStore struct {
	pool *pgxpool.Pool
}

// NewRefreshTokenStore creates a RefreshTokenStore.
func NewRefreshTokenStore(pool *pgxpool.Pool) *RefreshTokenStore {
	return &RefreshTokenStore{pool: pool}
}

// Create inserts a new refresh token hash.
func (s *RefreshTokenStore) Create(ctx context.Context, userID, tokenHash string, expiresAt time.Time) error {
	const q = `
		INSERT INTO refresh_tokens (user_id, token_hash, expires_at)
		VALUES ($1, $2, $3)`

	_, err := s.pool.Exec(ctx, q, userID, tokenHash, expiresAt)
	if err != nil {
		return fmt.Errorf("insert refresh token: %w", err)
	}
	return nil
}

// GetByHash returns an unused, unexpired refresh token by hash.
func (s *RefreshTokenStore) GetByHash(ctx context.Context, tokenHash string) (*RefreshToken, error) {
	const q = `
		SELECT id, user_id, token_hash, expires_at, used_at, created_at
		FROM refresh_tokens
		WHERE token_hash = $1 AND used_at IS NULL AND expires_at > now()`

	row := s.pool.QueryRow(ctx, q, tokenHash)
	var t RefreshToken
	err := row.Scan(&t.ID, &t.UserID, &t.TokenHash, &t.ExpiresAt, &t.UsedAt, &t.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get refresh token: %w", err)
	}
	return &t, nil
}

// MarkUsed marks a refresh token as consumed (single-use rotation).
func (s *RefreshTokenStore) MarkUsed(ctx context.Context, id string) error {
	const q = `UPDATE refresh_tokens SET used_at = now() WHERE id = $1`
	_, err := s.pool.Exec(ctx, q, id)
	if err != nil {
		return fmt.Errorf("mark refresh token used: %w", err)
	}
	return nil
}
