// Package grading holds the pure, deterministic graders for lesson steps.
// It is the only place mcq, true_false, fill_blank, and arrange answers are judged.
// Prose steps are graded by the rubric evaluator, never here.
package grading

import (
	"errors"
	"fmt"
	"strings"

	"github.com/prepio/prepio/services/question/internal/lesson"
)

// ErrInvalidAnswer means the submitted answer is malformed for the step type
// (wrong shape, out of range). It is not an incorrect answer and is not recorded as one.
var ErrInvalidAnswer = errors.New("invalid answer")

// ErrNotDeterministic is returned for step types this package does not grade.
var ErrNotDeterministic = errors.New("step type is not deterministically graded")

// Answer is a learner's submission. Which field is read depends on the step type.
type Answer struct {
	Choice *int     `json:"choice,omitempty"` // mcq
	Value  *bool    `json:"value,omitempty"`  // true_false
	Blanks []string `json:"blanks,omitempty"` // fill_blank
	Order  []int    `json:"order,omitempty"`  // arrange
	Text   *string  `json:"text,omitempty"`   // prose
}

// Result is the outcome of grading one deterministic step.
type Result struct {
	Correct     bool
	Explanation string
	// WhyNot explains the chosen option when an mcq answer is wrong.
	WhyNot string
	// CorrectAnswer is set only when the answer was incorrect.
	CorrectAnswer *Answer
}

// Grade judges an answer against a deterministic step payload.
func Grade(payload any, a Answer) (Result, error) {
	switch p := payload.(type) {
	case lesson.MCQPayload:
		return MCQ(p, a)
	case lesson.TrueFalsePayload:
		return TrueFalse(p, a)
	case lesson.FillBlankPayload:
		return FillBlank(p, a)
	case lesson.ArrangePayload:
		return Arrange(p, a)
	default:
		return Result{}, ErrNotDeterministic
	}
}

// MCQ grades a multiple-choice answer.
func MCQ(p lesson.MCQPayload, a Answer) (Result, error) {
	if a.Choice == nil {
		return Result{}, fmt.Errorf("%w: choice is required", ErrInvalidAnswer)
	}
	choice := *a.Choice
	if choice < 0 || choice >= len(p.Options) {
		return Result{}, fmt.Errorf("%w: choice %d is not an option", ErrInvalidAnswer, choice)
	}
	if choice == p.Answer {
		return Result{Correct: true, Explanation: p.Explanations.Correct}, nil
	}
	right := p.Answer
	return Result{
		Explanation:   p.Explanations.Correct,
		WhyNot:        p.Explanations.Wrong[choice],
		CorrectAnswer: &Answer{Choice: &right},
	}, nil
}

// TrueFalse grades a true/false answer.
func TrueFalse(p lesson.TrueFalsePayload, a Answer) (Result, error) {
	if a.Value == nil {
		return Result{}, fmt.Errorf("%w: value is required", ErrInvalidAnswer)
	}
	if *a.Value == p.Answer {
		return Result{Correct: true, Explanation: p.Explanation}, nil
	}
	right := p.Answer
	return Result{Explanation: p.Explanation, CorrectAnswer: &Answer{Value: &right}}, nil
}

// FillBlank grades a fill-in-the-blank answer. Every blank must match its answer
// after trimming surrounding whitespace.
func FillBlank(p lesson.FillBlankPayload, a Answer) (Result, error) {
	if len(a.Blanks) != len(p.Answers) {
		return Result{}, fmt.Errorf("%w: expected %d blanks, got %d", ErrInvalidAnswer, len(p.Answers), len(a.Blanks))
	}
	correct := true
	for i, want := range p.Answers {
		if strings.TrimSpace(a.Blanks[i]) != strings.TrimSpace(want) {
			correct = false
			break
		}
	}
	if correct {
		return Result{Correct: true, Explanation: p.Explanation}, nil
	}
	return Result{
		Explanation:   p.Explanation,
		CorrectAnswer: &Answer{Blanks: append([]string(nil), p.Answers...)},
	}, nil
}

// Arrange grades an arrange-in-order answer. The submission must be a permutation
// of the item indices; it is correct when it equals the authored order.
func Arrange(p lesson.ArrangePayload, a Answer) (Result, error) {
	if len(a.Order) != len(p.Items) {
		return Result{}, fmt.Errorf("%w: expected %d positions, got %d", ErrInvalidAnswer, len(p.Items), len(a.Order))
	}
	seen := make([]bool, len(p.Items))
	for _, idx := range a.Order {
		if idx < 0 || idx >= len(p.Items) || seen[idx] {
			return Result{}, fmt.Errorf("%w: order must be a permutation of the item indices", ErrInvalidAnswer)
		}
		seen[idx] = true
	}
	correct := true
	for i, want := range p.Order {
		if a.Order[i] != want {
			correct = false
			break
		}
	}
	if correct {
		return Result{Correct: true, Explanation: p.Explanation}, nil
	}
	return Result{
		Explanation:   p.Explanation,
		CorrectAnswer: &Answer{Order: append([]int(nil), p.Order...)},
	}, nil
}
