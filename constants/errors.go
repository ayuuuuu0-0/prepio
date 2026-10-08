package constants

// API error codes returned in the error envelope.
const (
	ErrInvalidRequest = "invalid_request"
	ErrUnauthorized   = "unauthorized"
	ErrForbidden      = "forbidden"
	ErrNotFound       = "not_found"
	ErrConflict       = "conflict"
	ErrInternal       = "internal_error"
	ErrRateLimited    = "rate_limited"

	// Auth
	ErrInvalidCredentials  = "invalid_credentials"
	ErrEmailTaken          = "email_taken"
	ErrUsernameTaken       = "username_taken"
	ErrInvalidToken        = "invalid_token"
	ErrTokenRevoked        = "token_revoked"
	ErrRefreshTokenInvalid = "refresh_token_invalid"

	// Users
	ErrUserNotFound     = "user_not_found"
	ErrInsufficientGems = "insufficient_gems"
	ErrDeviceNotFound   = "device_not_found"

	// Skills
	ErrSkillNotFound = "skill_not_found"

	// Lessons
	ErrLessonNotFound       = "lesson_not_found"
	ErrLessonLocked         = "lesson_locked"
	ErrAttemptNotFound      = "attempt_not_found"
	ErrStepNotFound         = "step_not_found"
	ErrAttemptNotInProgress = "attempt_not_in_progress"
	ErrAttemptIncomplete    = "attempt_incomplete"
	ErrInvalidAnswer        = "invalid_answer"
	ErrInvalidTry           = "invalid_try"

	// Streaks
	ErrStreakFreezeMaxHeld          = "streak_freeze_max_held"
	ErrStreakFreezeInsufficientGems = "streak_freeze_insufficient_gems"

	// Progress
	ErrProgressNotFound = "progress_not_found"
)
