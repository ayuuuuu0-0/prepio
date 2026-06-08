package constants

import "fmt"

// TokenBlacklistKey returns the Redis key for a blacklisted access token JTI.
func TokenBlacklistKey(jti string) string {
	return fmt.Sprintf("token_blacklist:%s", jti)
}

// StreakKey returns the Redis cache key for a user's streak.
func StreakKey(userID string) string {
	return fmt.Sprintf("streak:%s", userID)
}

// GemsKey returns the Redis read cache key for a user's gem balance.
func GemsKey(userID string) string {
	return fmt.Sprintf("gems:%s", userID)
}

// NotifCapKey returns the Redis key for daily notification count.
func NotifCapKey(userID, dateYYYYMMDD string) string {
	return fmt.Sprintf("notif_cap:%s:%s", userID, dateYYYYMMDD)
}

// SessionKey returns the Redis key for a question session.
func SessionKey(sessionID string) string {
	return fmt.Sprintf("session:%s", sessionID)
}
