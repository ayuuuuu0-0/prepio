package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/prepio/prepio/constants"
	"github.com/prepio/prepio/internal/user/dto"
	"github.com/prepio/prepio/internal/user/mailer"
	"github.com/prepio/prepio/internal/user/store"
	"github.com/prepio/prepio/shared/jwt"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

// AuthService handles registration, login, logout, token refresh, and email OTP.
type AuthService struct {
	users         *store.UserStore
	refreshTokens *store.RefreshTokenStore
	signer        *jwt.Signer
	redis         *redis.Client
	mailer        mailer.Mailer
}

// NewAuthService creates an AuthService.
func NewAuthService(
	users *store.UserStore,
	refreshTokens *store.RefreshTokenStore,
	signer *jwt.Signer,
	redisClient *redis.Client,
	mailClient mailer.Mailer,
) *AuthService {
	if mailClient == nil {
		mailClient = &mailer.MockMailer{}
	}
	return &AuthService{
		users:         users,
		refreshTokens: refreshTokens,
		signer:        signer,
		redis:         redisClient,
		mailer:        mailClient,
	}
}

// SendOTP generates and emails a 6-digit verification code.
func (s *AuthService) SendOTP(ctx context.Context, req dto.SendOTPRequest) (*dto.SendOTPResponse, error) {
	email := strings.TrimSpace(strings.ToLower(req.Email))
	if len(email) == 0 || !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		return nil, ErrInvalidRequest
	}

	if req.Signup {
		existing, err := s.users.GetByEmail(ctx, email)
		if err != nil {
			return nil, err
		}
		if existing != nil {
			return nil, ErrEmailTaken
		}
		if req.Username != "" {
			existingUser, err := s.users.GetByUsername(ctx, strings.TrimSpace(req.Username))
			if err != nil {
				return nil, err
			}
			if existingUser != nil {
				return nil, ErrUsernameTaken
			}
		}
	}

	// Rate limiting: allow at most 1 request per 60 seconds per email
	rateKey := fmt.Sprintf("otp:rate:%s", email)
	exists, err := s.redis.Exists(ctx, rateKey).Result()
	if err == nil && exists > 0 {
		return nil, ErrOTPCooldown
	}

	// Generate secure 6-digit numeric OTP
	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return nil, fmt.Errorf("generate otp: %w", err)
	}
	code := fmt.Sprintf("%06d", n.Int64())

	codeKey := fmt.Sprintf("otp:code:%s", email)
	attemptsKey := fmt.Sprintf("otp:attempts:%s", email)

	if err := s.redis.Set(ctx, codeKey, code, 5*time.Minute).Err(); err != nil {
		return nil, fmt.Errorf("store otp: %w", err)
	}
	if err := s.redis.Set(ctx, rateKey, "1", 60*time.Second).Err(); err != nil {
		return nil, fmt.Errorf("set rate limit: %w", err)
	}
	_ = s.redis.Del(ctx, attemptsKey).Err()

	if err := s.mailer.SendOTP(ctx, email, code); err != nil {
		return nil, fmt.Errorf("send otp email: %w", err)
	}

	resp := &dto.SendOTPResponse{
		Message:   "Verification code sent to email",
		ExpiresIn: 300,
	}
	if s.mailer.IsDev() {
		resp.DevCode = code
	}
	return resp, nil
}

// VerifyOTP verifies the 6-digit code and signs in or creates the user.
func (s *AuthService) VerifyOTP(ctx context.Context, req dto.VerifyOTPRequest) (*dto.AuthResponse, error) {
	email := strings.TrimSpace(strings.ToLower(req.Email))
	code := strings.TrimSpace(req.Code)

	if len(email) == 0 || len(code) != 6 {
		return nil, ErrInvalidRequest
	}

	codeKey := fmt.Sprintf("otp:code:%s", email)
	attemptsKey := fmt.Sprintf("otp:attempts:%s", email)

	attempts, _ := s.redis.Get(ctx, attemptsKey).Int()
	if attempts >= 5 {
		_ = s.redis.Del(ctx, codeKey).Err()
		return nil, ErrTooManyAttempts
	}

	storedCode, err := s.redis.Get(ctx, codeKey).Result()
	if err == redis.Nil {
		return nil, ErrInvalidOTP
	}
	if err != nil {
		return nil, fmt.Errorf("read otp: %w", err)
	}

	if storedCode != code {
		_ = s.redis.Incr(ctx, attemptsKey).Err()
		_ = s.redis.Expire(ctx, attemptsKey, 5*time.Minute).Err()
		return nil, ErrInvalidOTP
	}

	// Clean up OTP keys on successful verification
	_ = s.redis.Del(ctx, codeKey).Err()
	_ = s.redis.Del(ctx, attemptsKey).Err()

	// Check if user already exists
	user, err := s.users.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		// New user: derive username and create account
		username := strings.TrimSpace(req.Username)
		if len(username) == 0 {
			username = deriveUsername(email)
		}
		uniqueUsername, err := s.ensureUniqueUsername(ctx, username)
		if err != nil {
			return nil, err
		}

		var hash []byte
		if len(req.Password) > 0 {
			hash, err = bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
			if err != nil {
				return nil, fmt.Errorf("hash password: %w", err)
			}
		} else {
			randomBytes := make([]byte, 32)
			_, _ = rand.Read(randomBytes)
			hash, err = bcrypt.GenerateFromPassword(randomBytes, bcrypt.DefaultCost)
			if err != nil {
				return nil, fmt.Errorf("hash password: %w", err)
			}
		}

		user, err = s.users.Create(ctx, email, uniqueUsername, string(hash), constants.DefaultTimezone)
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
	}

	return s.issueTokens(ctx, user)
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

// Login authenticates a user by email or username, and returns auth tokens.
func (s *AuthService) Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error) {
	identifier := strings.TrimSpace(req.Email)
	password := req.Password
	if len(identifier) == 0 || len(password) == 0 {
		return nil, ErrInvalidCredentials
	}

	var user *store.User
	var err error
	if strings.Contains(identifier, "@") {
		user, err = s.users.GetByEmail(ctx, strings.ToLower(identifier))
	} else {
		user, err = s.users.GetByUsername(ctx, identifier)
	}
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
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

func (s *AuthService) ensureUniqueUsername(ctx context.Context, base string) (string, error) {
	candidate := base
	for i := 0; i < 5; i++ {
		existing, err := s.users.GetByUsername(ctx, candidate)
		if err != nil {
			return "", err
		}
		if existing == nil {
			return candidate, nil
		}
		num, _ := rand.Int(rand.Reader, big.NewInt(9000))
		candidate = fmt.Sprintf("%s_%d", base, num.Int64()+1000)
	}
	return fmt.Sprintf("%s_%d", base, time.Now().UnixNano()%100000), nil
}

func deriveUsername(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) == 0 || len(parts[0]) == 0 {
		return fmt.Sprintf("learner_%d", time.Now().Unix()%10000)
	}
	var b strings.Builder
	for _, r := range parts[0] {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
		}
	}
	s := strings.ToLower(b.String())
	if len(s) < 3 {
		s = "learner"
	}
	if len(s) > 20 {
		s = s[:20]
	}
	return s
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
