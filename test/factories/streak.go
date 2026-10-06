package factories

import (
	"time"

	"github.com/google/uuid"
	"github.com/prepio/prepio/shared/events"
)

// LessonCompletedEvent builds a lesson.completed event for tests.
func LessonCompletedEvent(userID string, completedAt time.Time) events.LessonCompleted {
	return events.LessonCompleted{
		EventID:         uuid.NewString(),
		UserID:          userID,
		LessonID:        uuid.NewString(),
		LessonSlug:      "test-lesson",
		AttemptID:       uuid.NewString(),
		Kind:            "lesson",
		Difficulty:      "easy",
		GradedSteps:     3,
		FirstTryCorrect: 3,
		TotalTries:      3,
		CompletedAt:     completedAt,
	}
}
