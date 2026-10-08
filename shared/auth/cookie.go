package auth

import (
	"net/http"
	"strings"

	"github.com/prepio/prepio/constants"
)

// SetRefreshTokenCookie writes the refresh token as an httpOnly cookie.
// Over HTTPS the cookie is Secure and SameSite=None so a web client on another
// site (e.g. a separately hosted frontend) can send it with credentialed fetches.
func SetRefreshTokenCookie(w http.ResponseWriter, r *http.Request, refreshToken string) {
	if len(refreshToken) == 0 {
		return
	}
	secure, sameSite := cookiePolicy(r)
	http.SetCookie(w, &http.Cookie{
		Name:     constants.RefreshTokenCookie,
		Value:    refreshToken,
		Path:     "/api/v1/auth",
		MaxAge:   constants.RefreshTokenCookieMaxAge,
		HttpOnly: true,
		SameSite: sameSite,
		Secure:   secure,
	})
}

// ClearRefreshTokenCookie removes the refresh token cookie.
func ClearRefreshTokenCookie(w http.ResponseWriter, r *http.Request) {
	secure, sameSite := cookiePolicy(r)
	http.SetCookie(w, &http.Cookie{
		Name:     constants.RefreshTokenCookie,
		Value:    "",
		Path:     "/api/v1/auth",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: sameSite,
		Secure:   secure,
	})
}

// cookiePolicy returns Secure+SameSite=None for HTTPS requests (directly or via a
// TLS-terminating proxy) and insecure Lax for plain HTTP (local development).
func cookiePolicy(r *http.Request) (bool, http.SameSite) {
	if isHTTPS(r) {
		return true, http.SameSiteNoneMode
	}
	return false, http.SameSiteLaxMode
}

func isHTTPS(r *http.Request) bool {
	if r == nil {
		return false
	}
	if r.TLS != nil {
		return true
	}
	proto := r.Header.Get("X-Forwarded-Proto")
	if i := strings.IndexByte(proto, ','); i >= 0 {
		proto = proto[:i]
	}
	return strings.EqualFold(strings.TrimSpace(proto), "https")
}

// RefreshTokenFromRequest reads the refresh token from cookie or returns empty.
func RefreshTokenFromRequest(r *http.Request) string {
	c, err := r.Cookie(constants.RefreshTokenCookie)
	if err != nil || len(c.Value) == 0 {
		return ""
	}
	return c.Value
}
