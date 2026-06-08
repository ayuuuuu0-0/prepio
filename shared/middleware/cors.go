package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/prepio/prepio/constants"
)

// CORS adds cross-origin headers for web clients.
func CORS(next http.Handler) http.Handler {
	origins := parseOrigins(os.Getenv("CORS_ALLOWED_ORIGINS"))
	relaxed := strings.EqualFold(os.Getenv("DEV_SYNC_EVENTS"), "true")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if len(origin) > 0 && isAllowedOrigin(origin, origins, relaxed) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func parseOrigins(raw string) []string {
	if len(raw) == 0 {
		raw = constants.DefaultCORSOrigins
	}
	parts := strings.Split(raw, ",")
	origins := make([]string, 0, len(parts))
	for _, part := range parts {
		if trimmed := strings.TrimSpace(part); len(trimmed) > 0 {
			origins = append(origins, trimmed)
		}
	}
	return origins
}

func isAllowedOrigin(origin string, origins []string, relaxed bool) bool {
	for _, allowed := range origins {
		if allowed == origin {
			return true
		}
	}
	if !relaxed {
		return false
	}
	return strings.HasPrefix(origin, "http://localhost:") ||
		strings.HasPrefix(origin, "http://127.0.0.1:") ||
		strings.HasPrefix(origin, "http://192.168.")
}
