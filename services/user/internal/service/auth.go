package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/prepio/prepio/constants"
	"github.com/prepio/prepio/services/user/internal/dto"
	"github.com/prepio/prepio/services/user/internal/store"
	"github.com/prepio/prepio/shared/jwt"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

// AuthService handles registration, login, logout, and token refresh.
type AuthService struct {
	users         *store.UserStore
	refreshTokens *store.RefreshTokenStore
	signer        *jwt.Signer
	redis         *redis.Client
}

// NewAuthService creates an AuthService.
func NewAuthService(
	users *store.UserStore,
	refreshTokens *store.RefreshTokenStore,
	signer *jwt.Signer,
	redisClient *redis.Client,
) *AuthService {
	return &AuthService{
		users:         users,
		refreshTokens: refreshTokens,
		signer:        signer,
		redis:         redisClient,
	}
}

// Register creates a new user and returns auth tokens.
func (s *AuthService) Register(ctx context.Context, req dto.RegisterRequest) (*dto.AuthResponse, error) {
	if err := validateCredentials(req.Email, req.Username, req.Password); err != nil {
		return nil, err
	}

	timezone := req.Timezone
	if len(timezone) == 0 {
		timezone = constants.DefaultTimezone
	}

	existing, err := s.users.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrEmailTaken
	}

	existing, err = s.users.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrUsernameTaken
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	user, err := s.users.Create(ctx, req.Email, req.Username, string(hash), timezone)
	if err != nil {
		return nil, err
	}

	if err := s.users.UnlockCharacter(ctx, user.ID, constants.DefaultCharacterID); err != nil {
		return nil, err
	}

	activeChar := constants.DefaultCharacterID
	user, err = s.users.SetActiveCharacter(ctx, user.ID, activeChar)
	if err != nil {
		return nil, err
	}

	if err := s.users.CreateNotificationPreferences(ctx, user.ID); err != nil {
		return nil, err
	}

	return s.issueTokens(ctx, user)
}

// Login authenticates a user and returns auth tokens.
func (s *AuthService) Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error) {
	if len(req.Email) == 0 || len(req.Password) == 0 {
		return nil, ErrInvalidCredentials
	}

	user, err := s.users.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	return s.issueTokens(ctx, user)
}

// Refresh rotates a refresh token and returns new auth tokens.
func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (*dto.AuthResponse, error) {
	if len(refreshToken) == 0 {
		return nil, ErrRefreshTokenInvalid
	}

	claims, err := s.signer.Verify(refreshToken)
	if err != nil {
		return nil, ErrRefreshTokenInvalid
	}
	if !jwt.IsRefreshToken(claims) {
		return nil, ErrRefreshTokenInvalid
	}

	hash := hashToken(refreshToken)
	record, err := s.refreshTokens.GetByHash(ctx, hash)
	if err != nil {
		return nil, err
	}
	if record == nil {
		return nil, ErrRefreshTokenInvalid
	}

	if err := s.refreshTokens.MarkUsed(ctx, record.ID); err != nil {
		return nil, err
	}

	user, err := s.users.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	return s.issueTokens(ctx, user)
}

// Logout blacklists the access token until it expires.
func (s *AuthService) Logout(ctx context.Context, accessToken string) error {
	if len(accessToken) == 0 {
		return ErrInvalidToken
	}

	claims, err := s.signer.Verify(accessToken)
	if err != nil {
		return ErrInvalidToken
	}
	if !jwt.IsAccessToken(claims) {
		return ErrInvalidToken
	}

	ttl := time.Until(claims.ExpiresAt.Time)
	if ttl <= 0 {
		return nil
	}

	key := constants.TokenBlacklistKey(claims.ID)
	if err := s.redis.Set(ctx, key, "1", ttl).Err(); err != nil {
		return fmt.Errorf("blacklist token: %w", err)
	}

	return nil
}

func (s *AuthService) issueTokens(ctx context.Context, user *store.User) (*dto.AuthResponse, error) {
	accessToken, _, _, err := s.signer.SignAccessToken(user.ID)
	if err != nil {
		return nil, err
	}

	refreshToken, _, expiresAt, err := s.signer.SignRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	if err := s.refreshTokens.Create(ctx, user.ID, hashToken(refreshToken), expiresAt); err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         toUserResponse(user),
	}, nil
}

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func validateCredentials(email, username, password string) error {
	if len(email) == 0 || len(username) == 0 || len(password) == 0 {
		return ErrInvalidRequest
	}
	if len(password) < 8 {
		return ErrInvalidRequest
	}
	return nil
}
