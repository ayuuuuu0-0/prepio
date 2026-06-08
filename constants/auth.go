package constants

// RefreshTokenCookie is the httpOnly cookie name for the JWT refresh token (web clients).
const RefreshTokenCookie = "prepio_refresh_token"

// RefreshTokenCookieMaxAge is refresh token cookie lifetime in seconds (7 days).
const RefreshTokenCookieMaxAge = 7 * 24 * 60 * 60
