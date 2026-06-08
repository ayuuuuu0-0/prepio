package constants

// AuthenticatedRateLimitPerMinute is the max requests per minute for logged-in users.
const AuthenticatedRateLimitPerMinute = 300

// UnauthenticatedRateLimitPerMinute is the max requests per minute per IP.
const UnauthenticatedRateLimitPerMinute = 20

// RateLimitWindow is the sliding window duration for rate limiting.
const RateLimitWindowSeconds = 60
