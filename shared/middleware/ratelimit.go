package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/prepio/prepio/constants"
	"github.com/prepio/prepio/shared/response"
	"github.com/redis/go-redis/v9"
)

// RateLimit enforces a per-key request cap within a fixed window using Redis.
func RateLimit(redisClient *redis.Client, limit int, keyFn func(*http.Request) string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := keyFn(r)
			if len(key) == 0 {
				response.Error(w, http.StatusInternalServerError, constants.ErrInternal, "rate limit key missing")
				return
			}

			redisKey := fmt.Sprintf("ratelimit:%s", key)
			count, err := redisClient.Incr(r.Context(), redisKey).Result()
			if err != nil {
				response.Error(w, http.StatusInternalServerError, constants.ErrInternal, "rate limit check failed")
				return
			}

			if count == 1 {
				redisClient.Expire(r.Context(), redisKey, time.Duration(constants.RateLimitWindowSeconds)*time.Second)
			}

			if int(count) > limit {
				response.Error(w, http.StatusTooManyRequests, constants.ErrRateLimited, "rate limit exceeded")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RateLimitKeyByIP returns a rate limit key scoped to client IP and current minute.
func RateLimitKeyByIP(r *http.Request) string {
	ip := r.RemoteAddr
	if forwarded := r.Header.Get("X-Forwarded-For"); len(forwarded) > 0 {
		ip = forwarded
	}
	minute := time.Now().UTC().Format("200601021504")
	return fmt.Sprintf("ip:%s:%s", ip, minute)
}

// RateLimitKeyByUser returns a rate limit key scoped to authenticated user and current minute.
func RateLimitKeyByUser(r *http.Request) string {
	userID, ok := UserIDFromContext(r.Context())
	if !ok || len(userID) == 0 {
		return RateLimitKeyByIP(r)
	}
	minute := time.Now().UTC().Format("200601021504")
	return fmt.Sprintf("user:%s:%s", userID, minute)
}
