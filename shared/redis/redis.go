package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// Config holds Redis connection settings.
type Config struct {
	Addr     string
	Password string
	DB       int
}

// New creates a Redis client and verifies connectivity.
func New(ctx context.Context, cfg Config) (*redis.Client, error) {
	if len(cfg.Addr) == 0 {
		return nil, fmt.Errorf("redis addr is required")
	}

	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		_ = client.Close()
		return nil, fmt.Errorf("ping redis: %w", err)
	}

	return client, nil
}
