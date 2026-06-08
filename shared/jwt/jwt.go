package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	// AccessTokenTTL is the lifetime of an access token.
	AccessTokenTTL = 15 * time.Minute

	// RefreshTokenTTL is the lifetime of a refresh token.
	RefreshTokenTTL = 7 * 24 * time.Hour

	tokenTypeAccess  = "access"
	tokenTypeRefresh = "refresh"
)

// Claims holds JWT claims for Prepio tokens.
type Claims struct {
	UserID    string `json:"user_id"`
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}

// Signer issues and verifies JWT tokens.
type Signer struct {
	secret []byte
}

// NewSigner creates a Signer with the given HMAC secret.
func NewSigner(secret string) (*Signer, error) {
	if len(secret) == 0 {
		return nil, fmt.Errorf("jwt secret is required")
	}
	return &Signer{secret: []byte(secret)}, nil
}

// SignAccessToken issues a short-lived access token for the user.
func (s *Signer) SignAccessToken(userID string) (string, string, time.Time, error) {
	return s.sign(userID, tokenTypeAccess, AccessTokenTTL)
}

// SignRefreshToken issues a long-lived refresh token for the user.
func (s *Signer) SignRefreshToken(userID string) (string, string, time.Time, error) {
	return s.sign(userID, tokenTypeRefresh, RefreshTokenTTL)
}

func (s *Signer) sign(userID, tokenType string, ttl time.Duration) (string, string, time.Time, error) {
	if len(userID) == 0 {
		return "", "", time.Time{}, fmt.Errorf("user id is required")
	}

	jti := uuid.NewString()
	expiresAt := time.Now().UTC().Add(ttl)

	claims := Claims{
		UserID:    userID,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        jti,
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(s.secret)
	if err != nil {
		return "", "", time.Time{}, fmt.Errorf("sign token: %w", err)
	}

	return signed, jti, expiresAt, nil
}

// Verify parses and validates a token string, returning its claims.
func (s *Signer) Verify(tokenString string) (*Claims, error) {
	if len(tokenString) == 0 {
		return nil, fmt.Errorf("token is required")
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return s.secret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("parse token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

// IsAccessToken reports whether the claims represent an access token.
func IsAccessToken(claims *Claims) bool {
	return claims != nil && claims.TokenType == tokenTypeAccess
}

// IsRefreshToken reports whether the claims represent a refresh token.
func IsRefreshToken(claims *Claims) bool {
	return claims != nil && claims.TokenType == tokenTypeRefresh
}
