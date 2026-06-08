package store

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// UserDevice is a row from the user_devices table.
type UserDevice struct {
	ID        string
	UserID    string
	FCMToken  string
	Platform  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// UserDeviceStore handles user_devices table queries.
type UserDeviceStore struct {
	pool *pgxpool.Pool
}

// NewUserDeviceStore creates a UserDeviceStore.
func NewUserDeviceStore(pool *pgxpool.Pool) *UserDeviceStore {
	return &UserDeviceStore{pool: pool}
}

// Upsert inserts or updates a device token for the user and platform.
func (s *UserDeviceStore) Upsert(ctx context.Context, userID, fcmToken, platform string) (*UserDevice, error) {
	const q = `
		INSERT INTO user_devices (user_id, fcm_token, platform)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, platform)
		DO UPDATE SET fcm_token = EXCLUDED.fcm_token
		RETURNING id, user_id, fcm_token, platform, created_at, updated_at`

	row := s.pool.QueryRow(ctx, q, userID, fcmToken, platform)
	var d UserDevice
	err := row.Scan(&d.ID, &d.UserID, &d.FCMToken, &d.Platform, &d.CreatedAt, &d.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("upsert user device: %w", err)
	}
	return &d, nil
}

// Delete removes a device by ID scoped to the user.
func (s *UserDeviceStore) Delete(ctx context.Context, userID, deviceID string) (bool, error) {
	const q = `DELETE FROM user_devices WHERE id = $1 AND user_id = $2`
	tag, err := s.pool.Exec(ctx, q, deviceID, userID)
	if err != nil {
		return false, fmt.Errorf("delete user device: %w", err)
	}
	return tag.RowsAffected() > 0, nil
}

// GetByID returns a device owned by the user.
func (s *UserDeviceStore) GetByID(ctx context.Context, userID, deviceID string) (*UserDevice, error) {
	const q = `
		SELECT id, user_id, fcm_token, platform, created_at, updated_at
		FROM user_devices WHERE id = $1 AND user_id = $2`

	row := s.pool.QueryRow(ctx, q, deviceID, userID)
	var d UserDevice
	err := row.Scan(&d.ID, &d.UserID, &d.FCMToken, &d.Platform, &d.CreatedAt, &d.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get user device: %w", err)
	}
	return &d, nil
}
