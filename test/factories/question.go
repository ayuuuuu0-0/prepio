package factories

import "github.com/google/uuid"

// QuestionFactory builds test question metadata.
type QuestionFactory struct {
	Body        string
	RoundType   string
	Difficulty  string
	AnswerGuide string
	Status      string
	IsWeekend   bool
	Source      string
	Companies   []string
}

// NewQuestionFactory returns a QuestionFactory with defaults.
func NewQuestionFactory() QuestionFactory {
	return QuestionFactory{
		Body:        "explain hash map time complexity for two sum",
		RoundType:   "dsa",
		Difficulty:  "easy",
		AnswerGuide: "concepts:hash map|O(n) time|O(n) space|two sum",
		Status:      "approved",
		IsWeekend:   false,
		Source:      "manual",
		Companies:   []string{"google"},
	}
}

// WithDifficulty overrides difficulty.
func (f QuestionFactory) WithDifficulty(difficulty string) QuestionFactory {
	f.Difficulty = difficulty
	return f
}

// NewQuestionID returns a random question UUID string.
func NewQuestionID() string {
	return uuid.NewString()
}
