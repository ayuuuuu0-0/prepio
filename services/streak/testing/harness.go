package testing

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prepio/prepio/internal/streak/service"
	"github.com/prepio/prepio/internal/streak/store"
	"github.com/redis/go-redis/v9"
)

// NewService wires the streak service for integration tests.
func NewService(pool *pgxpool.Pool, redisClient *redis.Client, publisher service.EventPublisher, gems service.GemDeductor) *service.StreakService {
	return service.NewStreakService(
		store.NewStreakStore(pool),
		store.NewFreezeStore(pool),
		store.NewCacheStore(redisClient),
		publisher,
		gems,
	)
}
