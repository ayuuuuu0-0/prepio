package constants

// API error codes returned in the error envelope.
const (
	ErrInvalidRequest          = "invalid_request"
	ErrUnauthorized            = "unauthorized"
	ErrForbidden               = "forbidden"
	ErrNotFound                = "not_found"
	ErrConflict                = "conflict"
	ErrInternal                = "internal_error"
	ErrRateLimited             = "rate_limited"

	// Auth
	ErrInvalidCredentials      = "invalid_credentials"
	ErrEmailTaken              = "email_taken"
	ErrUsernameTaken           = "username_taken"
	ErrInvalidToken            = "invalid_token"
	ErrTokenExpired            = "token_expired"
	ErrTokenRevoked            = "token_revoked"
	ErrRefreshTokenInvalid     = "refresh_token_invalid"
	ErrRefreshTokenExpired     = "refresh_token_expired"

	// Users
	ErrUserNotFound            = "user_not_found"
	ErrCharacterNotFound       = "character_not_found"
	ErrCharacterAlreadyUnlocked = "character_already_unlocked"
	ErrCharacterNotUnlocked    = "character_not_unlocked"
	ErrInsufficientGems        = "insufficient_gems"
	ErrDeviceNotFound          = "device_not_found"

	// Skills
	ErrSkillNotFound           = "skill_not_found"

	// Journey
	ErrJourneyNodeNotFound     = "journey_node_not_found"

	// Questions
	ErrQuestionNotFound        = "question_not_found"
	ErrQuestionNotInSession    = "question_not_in_session"
	ErrSessionNotFound         = "session_not_found"
	ErrSessionExpired          = "session_expired"
	ErrDailyPaperUnavailable   = "daily_paper_unavailable"
	ErrAnswerAlreadySubmitted  = "answer_already_submitted"

	// Streaks
	ErrStreakNotFound          = "streak_not_found"
	ErrStreakFreezeMaxHeld     = "streak_freeze_max_held"
	ErrStreakFreezeInsufficientGems = "streak_freeze_insufficient_gems"

	// Progress
	ErrProgressNotFound        = "progress_not_found"

	// Leaderboard
	ErrLeaderboardUnavailable  = "leaderboard_unavailable"
)
