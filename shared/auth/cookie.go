package auth

import (
	"net/http"

	"github.com/prepio/prepio/constants"
)

// SetRefreshTokenCookie writes the refresh token as an httpOnly cookie.
func SetRefreshTokenCookie(w http.ResponseWriter, refreshToken string) {
	if len(refreshToken) == 0 {
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     constants.RefreshTokenCookie,
		Value:    refreshToken,
		Path:     "/api/v1/auth",
		MaxAge:   constants.RefreshTokenCookieMaxAge,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   false,
	})
}

// ClearRefreshTokenCookie removes the refresh token cookie.
func ClearRefreshTokenCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     constants.RefreshTokenCookie,
		Value:    "",
		Path:     "/api/v1/auth",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

// RefreshTokenFromRequest reads the refresh token from cookie or returns empty.
func RefreshTokenFromRequest(r *http.Request) string {
	c, err := r.Cookie(constants.RefreshTokenCookie)
	if err != nil || len(c.Value) == 0 {
		return ""
	}
	return c.Value
}
