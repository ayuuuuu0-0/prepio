package store

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// NotificationStore handles notification_log queries.
type NotificationStore struct {
	pool *pgxpool.Pool
}

// NewNotificationStore creates a NotificationStore.
func NewNotificationStore(pool *pgxpool.Pool) *NotificationStore {
	return &NotificationStore{pool: pool}
}

// Insert logs a sent notification.
func (s *NotificationStore) Insert(ctx context.Context, userID, notificationType, channel string) error {
	const q = `INSERT INTO notification_log (user_id, notification_type, channel) VALUES ($1, $2, $3)`
	_, err := s.pool.Exec(ctx, q, userID, notificationType, channel)
	if err != nil {
		return fmt.Errorf("insert notification log: %w", err)
	}
	return nil
}

// CountToday returns notifications sent today for the user.
func (s *NotificationStore) CountToday(ctx context.Context, userID, date string) (int, error) {
	const q = `
		SELECT COUNT(*) FROM notification_log
		WHERE user_id = $1 AND sent_at::date = $2::date`

	var count int
	err := s.pool.QueryRow(ctx, q, userID, date).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count notifications: %w", err)
	}
	return count, nil
}
