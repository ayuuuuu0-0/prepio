package testing

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prepio/prepio/internal/notification/service"
	"github.com/prepio/prepio/internal/notification/store"
	"github.com/redis/go-redis/v9"
)

// NewService wires the notification service for integration tests.
func NewService(pool *pgxpool.Pool, redisClient *redis.Client) *service.NotificationService {
	return service.NewNotificationService(store.NewNotificationStore(pool), redisClient)
}
