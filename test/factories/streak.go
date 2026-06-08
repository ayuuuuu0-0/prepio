package factories

import (
	"time"

	"github.com/google/uuid"
	"github.com/prepio/prepio/shared/events"
)

// QuestionAnsweredEvent builds a question.answered Kafka event for tests.
func QuestionAnsweredEvent(userID string, submittedAt time.Time) events.QuestionAnswered {
	return events.QuestionAnswered{
		EventID:     uuid.NewString(),
		UserID:      userID,
		QuestionID:  "test-question",
		RoundType:   "dsa",
		Difficulty:  "easy",
		CompanyTags: []string{"google"},
		Correct:     true,
		SubmittedAt: submittedAt,
		SessionID:   "test-session",
	}
}
