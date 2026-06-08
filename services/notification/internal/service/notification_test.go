package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/prepio/prepio/services/notification/internal/service"
	"github.com/prepio/prepio/services/notification/internal/store"
	"github.com/prepio/prepio/shared/events"
	"github.com/prepio/prepio/test/testdb"
	"github.com/prepio/prepio/test/testredis"
	"github.com/stretchr/testify/require"
)

func TestLevelUpCreatesNotificationLog(t *testing.T) {
	pool, _ := testdb.Start(t)
	testdb.Migrate(t, pool)
	redisClient, _ := testredis.New(t)

	ctx := context.Background()
	var userID string
	require.NoError(t, pool.QueryRow(ctx, `
		INSERT INTO users (email, username, password_hash)
		VALUES ('n@test.com', 'nuser', 'hash') RETURNING id`).Scan(&userID))

	svc := service.NewNotificationService(store.NewNotificationStore(pool), redisClient)
	require.NoError(t, svc.HandleProgressUpdated(ctx, events.ProgressUpdated{
		EventID: uuid.NewString(), UserID: userID,
		LeveledUp: true, LevelAfter: 2,
	}))

	var count int
	require.NoError(t, pool.QueryRow(ctx, `
		SELECT COUNT(*) FROM notification_log WHERE user_id = $1`, userID).Scan(&count))
	require.Equal(t, 1, count)
}

func TestDailyCapBlocksFourthNotification(t *testing.T) {
	pool, _ := testdb.Start(t)
	testdb.Migrate(t, pool)
	redisClient, _ := testredis.New(t)

	ctx := context.Background()
	var userID string
	require.NoError(t, pool.QueryRow(ctx, `
		INSERT INTO users (email, username, password_hash)
		VALUES ('c@test.com', 'cuser', 'hash') RETURNING id`).Scan(&userID))

	svc := service.NewNotificationService(store.NewNotificationStore(pool), redisClient)
	for i := 0; i < 4; i++ {
		_ = svc.HandleDispatch(ctx, events.NotificationsDispatch{
			EventID: uuid.NewString(), UserID: userID,
			NotificationType: events.NotificationStreakReminder,
			TriggeredAt:      time.Now().UTC(),
		})
	}

	var count int
	require.NoError(t, pool.QueryRow(ctx, `
		SELECT COUNT(*) FROM notification_log WHERE user_id = $1`, userID).Scan(&count))
	require.Equal(t, 3, count)
}
