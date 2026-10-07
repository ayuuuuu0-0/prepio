package events

import "time"

// Topic names for the Kafka event bus.
const (
	TopicQuestionAnswered       = "question.answered"
	TopicStreakUpdated          = "streak.updated"
	TopicProgressUpdated        = "progress.updated"
	TopicNotificationsDispatch  = "notifications.dispatch"
	TopicLessonCompleted        = "lesson.completed"
)

// LessonSkillWeight is one skill a lesson teaches, with its weight (weights sum to 1.0).
type LessonSkillWeight struct {
	SkillID   string  `json:"skill_id"`
	SkillSlug string  `json:"skill_slug"`
	Weight    float64 `json:"weight"`
}

// LessonCompleted is emitted when a user finishes a lesson attempt.
// It carries performance facts only; XP and mastery are computed by Progress.
type LessonCompleted struct {
	EventID         string              `json:"event_id"`
	UserID          string              `json:"user_id"`
	LessonID        string              `json:"lesson_id"`
	LessonSlug      string              `json:"lesson_slug"`
	AttemptID       string              `json:"attempt_id"`
	Kind            string              `json:"kind"`
	Difficulty      string              `json:"difficulty"`
	Skills          []LessonSkillWeight `json:"skills"`
	GradedSteps     int                 `json:"graded_steps"`
	FirstTryCorrect int                 `json:"first_try_correct"`
	TotalTries      int                 `json:"total_tries"`
	CompletedAt     time.Time           `json:"completed_at"`
}

// QuestionAnswered is emitted when a user submits an answer.
type QuestionAnswered struct {
	EventID     string    `json:"event_id"`
	UserID      string    `json:"user_id"`
	QuestionID  string    `json:"question_id"`
	RoundType   string    `json:"round_type"`
	Difficulty  string    `json:"difficulty"`
	Correct     bool      `json:"correct"`
	Score       int       `json:"score"`
	XPAwarded   int       `json:"xp_awarded"`
	GemsAwarded int       `json:"gems_awarded"`
	SubmittedAt time.Time `json:"submitted_at"`
	SessionID   string    `json:"session_id"`
}

// StreakUpdated is emitted when a user's streak state changes.
type StreakUpdated struct {
	EventID        string    `json:"event_id"`
	UserID         string    `json:"user_id"`
	PreviousStreak int       `json:"previous_streak"`
	CurrentStreak  int       `json:"current_streak"`
	StreakBroken   bool      `json:"streak_broken"`
	FreezeConsumed bool      `json:"freeze_consumed"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// ProgressUpdated is emitted when XP, gems, or level changes.
type ProgressUpdated struct {
	EventID     string    `json:"event_id"`
	UserID      string    `json:"user_id"`
	XPAwarded   int       `json:"xp_awarded"`
	GemsAwarded int       `json:"gems_awarded"`
	TotalXP     int       `json:"total_xp"`
	TotalGems   int       `json:"total_gems"`
	LevelBefore int       `json:"level_before"`
	LevelAfter  int       `json:"level_after"`
	LeveledUp   bool      `json:"leveled_up"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// NotificationsDispatch triggers the notification service to send a push.
type NotificationsDispatch struct {
	EventID          string         `json:"event_id"`
	UserID           string         `json:"user_id"`
	NotificationType string         `json:"notification_type"`
	Metadata         map[string]any `json:"metadata"`
	TriggeredAt      time.Time      `json:"triggered_at"`
}

// Notification types consumed by the notification service.
const (
	NotificationStreakReminder          = "streak_reminder"
	NotificationStreakBroken            = "streak_broken"
	NotificationLevelUp                 = "level_up"
	NotificationLeaguePositionChange    = "league_position_change"
	NotificationWeekendChallengeAvail   = "weekend_challenge_available"
	NotificationStreakFreezeLow         = "streak_freeze_low"
)
