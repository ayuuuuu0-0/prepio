package dto

import (
	"github.com/prepio/prepio/services/question/internal/grading"
	"github.com/prepio/prepio/services/question/internal/lesson"
)

// PathResponse is returned by GET /api/v1/path.
type PathResponse struct {
	Worlds []PathWorld `json:"worlds"`
}

// PathWorld is a world with its nodes in path order.
type PathWorld struct {
	ID          string `json:"id"`
	Slug        string `json:"slug"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Theme       string `json:"theme"`
	// Topic is the topic slug the world builds (empty if none). Focused is true when it is
	// one of the learner's focus topics; focused worlds come first.
	Topic   string     `json:"topic,omitempty"`
	Focused bool       `json:"focused"`
	Nodes   []PathNode `json:"nodes"`
}

// PathNode is one node with its lesson preview. Status is locked, current,
// available, or done; locked nodes say exactly what unlocks them.
type PathNode struct {
	ID         string   `json:"id"`
	Slug       string   `json:"slug"`
	Label      string   `json:"label"`
	NodeType   string   `json:"node_type"`
	Status     string   `json:"status"`
	UnlockHint string   `json:"unlock_hint,omitempty"`
	InProgress bool     `json:"in_progress"`
	LessonID   string   `json:"lesson_id"`
	LessonSlug string   `json:"lesson_slug"`
	Title      string   `json:"title"`
	Takeaways  []string `json:"takeaways"`
	Kind       string   `json:"kind"`
	Difficulty string   `json:"difficulty"`
	EstMinutes int      `json:"est_minutes"`
	XPPreview  int      `json:"xp_preview"`
}

// LessonInfo describes the lesson an attempt belongs to.
type LessonInfo struct {
	ID         string `json:"id"`
	Slug       string `json:"slug"`
	Title      string `json:"title"`
	Kind       string `json:"kind"`
	Difficulty string `json:"difficulty"`
	EstMinutes int    `json:"est_minutes"`
}

// StepProgress is the state of one step in a resumed attempt.
type StepProgress struct {
	StepID string `json:"step_id"`
	Tries  int    `json:"tries"`
	Done   bool   `json:"done"`
}

// AttemptResponse is returned when an attempt is started or resumed.
// Steps never contain answers, explanations, or rubrics.
type AttemptResponse struct {
	AttemptID string              `json:"attempt_id"`
	Resumed   bool                `json:"resumed"`
	Lesson    LessonInfo          `json:"lesson"`
	Steps     []lesson.ClientStep `json:"steps"`
	Progress  []StepProgress      `json:"progress"`
}

// AnswerRequest is the body for POST /api/v1/attempts/{id}/steps/{stepId}/answer.
// Try is 1-based and makes the call idempotent: repeating a try returns the stored result.
type AnswerRequest struct {
	Try    int            `json:"try"`
	Answer grading.Answer `json:"answer"`
}

// AnswerResponse is the server's grading of one try.
// CorrectAnswer is present only after an incorrect answer.
type AnswerResponse struct {
	Correct       bool            `json:"correct"`
	Try           int             `json:"try"`
	Explanation   string          `json:"explanation,omitempty"`
	WhyNot        string          `json:"why_not,omitempty"`
	CorrectAnswer *grading.Answer `json:"correct_answer,omitempty"`
	Feedback      string          `json:"feedback,omitempty"`
	Score         *int            `json:"score,omitempty"`
	Strengths     []string        `json:"strengths,omitempty"`
	Gaps          []string        `json:"gaps,omitempty"`
	StepDone      bool            `json:"step_done"`
	AllDone       bool            `json:"all_done"`
	Replayed      bool            `json:"replayed"`
}

// UnlockedNode is a node that became available because of a completion.
type UnlockedNode struct {
	ID       string `json:"id"`
	Slug     string `json:"slug"`
	Label    string `json:"label"`
	LessonID string `json:"lesson_id"`
}

// CompleteResponse is the performance summary of a finished attempt.
// The gateway adds rewards (XP, mastery changes) from Progress.
type CompleteResponse struct {
	AttemptID       string         `json:"attempt_id"`
	Lesson          LessonInfo     `json:"lesson"`
	Takeaways       []string       `json:"takeaways"`
	GradedSteps     int            `json:"graded_steps"`
	FirstTryCorrect int            `json:"first_try_correct"`
	TotalTries      int            `json:"total_tries"`
	Accuracy        float64        `json:"accuracy"`
	NodeID          string         `json:"node_id"`
	UnlockedNodes   []UnlockedNode `json:"unlocked_nodes"`
}
