package store

import (
	"context"
	"fmt"
	"time"

	"github.com/prepio/prepio/constants"
	"github.com/redis/go-redis/v9"
)

// CacheStore mirrors streak state in Redis.
type CacheStore struct {
	redis *redis.Client
}

// NewCacheStore creates a CacheStore.
func NewCacheStore(redisClient *redis.Client) *CacheStore {
	return &CacheStore{redis: redisClient}
}

// WriteThrough caches streak fields after a PostgreSQL write.
func (s *CacheStore) WriteThrough(ctx context.Context, state StreakState, freezeCount int) error {
	key := constants.StreakKey(state.UserID)
	last := ""
	if state.LastActivityDate != nil {
		last = state.LastActivityDate.Format("2006-01-02")
	}

	err := s.redis.HSet(ctx, key, map[string]any{
		"current_streak":     state.CurrentStreak,
		"longest_streak":     state.LongestStreak,
		"last_activity_date": last,
		"freeze_count":       freezeCount,
	}).Err()
	if err != nil {
		return fmt.Errorf("cache streak: %w", err)
	}
	return s.redis.Expire(ctx, key, constants.StreakCacheTTL).Err()
}

// Get reads cached streak state. Returns nil if cache miss.
func (s *CacheStore) Get(ctx context.Context, userID string) (*StreakState, int, error) {
	key := constants.StreakKey(userID)
	values, err := s.redis.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, 0, fmt.Errorf("get cached streak: %w", err)
	}
	if len(values) == 0 {
		return nil, 0, nil
	}

	state := &StreakState{UserID: userID}
	var freezeCount int
	_, _ = fmt.Sscanf(values["current_streak"], "%d", &state.CurrentStreak)
	_, _ = fmt.Sscanf(values["longest_streak"], "%d", &state.LongestStreak)
	_, _ = fmt.Sscanf(values["freeze_count"], "%d", &freezeCount)

	if len(values["last_activity_date"]) > 0 {
		parsed, err := time.Parse("2006-01-02", values["last_activity_date"])
		if err == nil {
			state.LastActivityDate = &parsed
		}
	}
	return state, freezeCount, nil
}
