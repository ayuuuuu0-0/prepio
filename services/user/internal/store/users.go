package store

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// User is a row from the users table.
type User struct {
	ID                  string
	Email               string
	Username            string
	PasswordHash        string
	Timezone            string
	ActiveCharID        *string
	ReminderTime        string
	ExperienceLevel     *string
	OnboardingCompleted bool
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

// UserStore handles users table queries.
type UserStore struct {
	pool *pgxpool.Pool
}

// NewUserStore creates a UserStore.
func NewUserStore(pool *pgxpool.Pool) *UserStore {
	return &UserStore{pool: pool}
}

// Create inserts a new user and returns the created row.
func (s *UserStore) Create(ctx context.Context, email, username, passwordHash, timezone string) (*User, error) {
	const q = `
		INSERT INTO users (email, username, password_hash, timezone)
		VALUES ($1, $2, $3, $4)
		RETURNING id, email, username, password_hash, timezone, active_char_id,
		          reminder_time::text, experience_level, onboarding_completed,
		          created_at, updated_at`

	row := s.pool.QueryRow(ctx, q, email, username, passwordHash, timezone)
	return scanUser(row)
}

// GetByID returns a user by primary key.
func (s *UserStore) GetByID(ctx context.Context, id string) (*User, error) {
	const q = `
		SELECT id, email, username, password_hash, timezone, active_char_id,
		       reminder_time::text, experience_level, onboarding_completed,
		       created_at, updated_at
		FROM users WHERE id = $1`

	row := s.pool.QueryRow(ctx, q, id)
	user, err := scanUser(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return user, err
}

// GetByEmail returns a user by email address.
func (s *UserStore) GetByEmail(ctx context.Context, email string) (*User, error) {
	const q = `
		SELECT id, email, username, password_hash, timezone, active_char_id,
		       reminder_time::text, experience_level, onboarding_completed,
		       created_at, updated_at
		FROM users WHERE email = $1`

	row := s.pool.QueryRow(ctx, q, email)
	user, err := scanUser(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return user, err
}

// GetByUsername returns a user by username.
func (s *UserStore) GetByUsername(ctx context.Context, username string) (*User, error) {
	const q = `
		SELECT id, email, username, password_hash, timezone, active_char_id,
		       reminder_time::text, experience_level, onboarding_completed,
		       created_at, updated_at
		FROM users WHERE username = $1`

	row := s.pool.QueryRow(ctx, q, username)
	user, err := scanUser(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return user, err
}

// UpdateProfile updates mutable user profile fields.
func (s *UserStore) UpdateProfile(ctx context.Context, id string, timezone, reminderTime string, activeCharID *string) (*User, error) {
	const q = `
		UPDATE users
		SET timezone = $2,
		    reminder_time = $3::time,
		    active_char_id = $4
		WHERE id = $1
		RETURNING id, email, username, password_hash, timezone, active_char_id,
		          reminder_time::text, experience_level, onboarding_completed,
		          created_at, updated_at`

	row := s.pool.QueryRow(ctx, q, id, timezone, reminderTime, activeCharID)
	user, err := scanUser(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return user, err
}

// UnlockCharacter records a character unlock for the user.
func (s *UserStore) UnlockCharacter(ctx context.Context, userID, characterID string) error {
	const q = `
		INSERT INTO character_unlocks (user_id, character_id)
		VALUES ($1, $2)
		ON CONFLICT (user_id, character_id) DO NOTHING`

	_, err := s.pool.Exec(ctx, q, userID, characterID)
	if err != nil {
		return fmt.Errorf("insert character unlock: %w", err)
	}
	return nil
}

// HasCharacterUnlock reports whether the user has unlocked a character.
func (s *UserStore) HasCharacterUnlock(ctx context.Context, userID, characterID string) (bool, error) {
	const q = `SELECT 1 FROM character_unlocks WHERE user_id = $1 AND character_id = $2`

	var one int
	err := s.pool.QueryRow(ctx, q, userID, characterID).Scan(&one)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("check character unlock: %w", err)
	}
	return true, nil
}

// SetActiveCharacter sets the user's active character.
func (s *UserStore) SetActiveCharacter(ctx context.Context, userID, characterID string) (*User, error) {
	const q = `
		UPDATE users SET active_char_id = $2 WHERE id = $1
		RETURNING id, email, username, password_hash, timezone, active_char_id,
		          reminder_time::text, experience_level, onboarding_completed,
		          created_at, updated_at`

	row := s.pool.QueryRow(ctx, q, userID, characterID)
	user, err := scanUser(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return user, err
}

// CompleteOnboarding stores onboarding choices and marks the profile complete.
func (s *UserStore) CompleteOnboarding(ctx context.Context, userID, experienceLevel, characterID string) (*User, error) {
	const q = `
		UPDATE users
		SET experience_level = $2,
		    active_char_id = $3,
		    onboarding_completed = true
		WHERE id = $1
		RETURNING id, email, username, password_hash, timezone, active_char_id,
		          reminder_time::text, experience_level, onboarding_completed,
		          created_at, updated_at`

	row := s.pool.QueryRow(ctx, q, userID, experienceLevel, characterID)
	user, err := scanUser(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return user, err
}

// CreateNotificationPreferences inserts default notification preferences.
func (s *UserStore) CreateNotificationPreferences(ctx context.Context, userID string) error {
	const q = `INSERT INTO notification_preferences (user_id) VALUES ($1) ON CONFLICT DO NOTHING`
	_, err := s.pool.Exec(ctx, q, userID)
	if err != nil {
		return fmt.Errorf("insert notification preferences: %w", err)
	}
	return nil
}

func scanUser(row pgx.Row) (*User, error) {
	var u User
	err := row.Scan(
		&u.ID,
		&u.Email,
		&u.Username,
		&u.PasswordHash,
		&u.Timezone,
		&u.ActiveCharID,
		&u.ReminderTime,
		&u.ExperienceLevel,
		&u.OnboardingCompleted,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("scan user: %w", err)
	}
	return &u, nil
}
