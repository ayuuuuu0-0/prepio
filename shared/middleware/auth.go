package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/prepio/prepio/constants"
	"github.com/prepio/prepio/shared/ctxkeys"
	"github.com/prepio/prepio/shared/jwt"
	"github.com/prepio/prepio/shared/response"
	"github.com/redis/go-redis/v9"
)

// Auth validates JWT access tokens and injects the user ID into context.
func Auth(signer *jwt.Signer, redisClient *redis.Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := ExtractBearerToken(r)
			if len(token) == 0 {
				response.Error(w, http.StatusUnauthorized, constants.ErrUnauthorized, "authorization required")
				return
			}

			claims, err := signer.Verify(token)
			if err != nil {
				response.Error(w, http.StatusUnauthorized, constants.ErrInvalidToken, "invalid token")
				return
			}
			if !jwt.IsAccessToken(claims) {
				response.Error(w, http.StatusUnauthorized, constants.ErrInvalidToken, "invalid token type")
				return
			}

			blacklisted, err := redisClient.Exists(r.Context(), constants.TokenBlacklistKey(claims.ID)).Result()
			if err != nil {
				response.Error(w, http.StatusInternalServerError, constants.ErrInternal, "token validation failed")
				return
			}
			if blacklisted > 0 {
				response.Error(w, http.StatusUnauthorized, constants.ErrTokenRevoked, "token revoked")
				return
			}

			ctx := context.WithValue(r.Context(), ctxkeys.UserIDKey, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// UserIDFromContext returns the authenticated user ID from context.
func UserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(ctxkeys.UserIDKey).(string)
	return userID, ok && len(userID) > 0
}

// ExtractBearerToken parses the Bearer token from the Authorization header.
func ExtractBearerToken(r *http.Request) string {
	header := r.Header.Get("Authorization")
	if len(header) == 0 {
		return ""
	}
	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
		return ""
	}
	return strings.TrimSpace(parts[1])
}
