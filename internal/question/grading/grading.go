// Package grading evaluates learner answers against step definitions.
// It is the only place mcq, true_false, fill_blank, arrange, and code_snippet answers are judged.
package grading

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/prepio/prepio/internal/question/lesson"
)

var (
	// ErrInvalidAnswer is returned when an answer is structurally malformed
	// (e.g. wrong type, missing required field, negative index).
	ErrInvalidAnswer = errors.New("invalid answer")
)

// ErrNotDeterministic is returned when a step type cannot be graded without external services (e.g. prose).
var ErrNotDeterministic = errors.New("step type is not deterministically graded")

// Answer is a learner's submission. Which field is read depends on the step type.
type Answer struct {
	Choice   *int     `json:"choice,omitempty"`   // mcq
	Value    *bool    `json:"value,omitempty"`    // true_false
	Blanks   []string `json:"blanks,omitempty"`   // fill_blank
	Order    []int    `json:"order,omitempty"`    // arrange
	Text     *string  `json:"text,omitempty"`     // prose
	Code     *string  `json:"code,omitempty"`     // code_snippet
	Language *string  `json:"language,omitempty"` // code_snippet language
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

// CodeSnippetTestResult is the outcome of one assertion in a code_snippet step.
type CodeSnippetTestResult struct {
	Name   string `json:"name"`
	Passed bool   `json:"passed"`
}

// CodeSnippetResult is the outcome of grading a code_snippet step.
// It embeds Result for the top-level correct/explanation fields, and adds
// per-test details shown in the feedback tray.
type CodeSnippetResult struct {
	Result
	Tests []CodeSnippetTestResult `json:"test_results"`
}

// Grade judges an answer against a stored step payload.
// Returns ErrNotDeterministic for prose steps.
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

// GradeCodeSnippet judges a code_snippet answer against its test assertions.
// It is separate from Grade() because it returns CodeSnippetResult, not Result.
func GradeCodeSnippet(p lesson.CodeSnippetPayload, a Answer) (CodeSnippetResult, error) {
	if a.Code == nil {
		return CodeSnippetResult{}, fmt.Errorf("%w: code is required", ErrInvalidAnswer)
	}
	code := *a.Code

	tests := p.Tests
	solution := p.Solution
	if a.Language != nil && len(*a.Language) > 0 && len(p.Languages) > 0 {
		if variant, ok := p.Languages[*a.Language]; ok {
			tests = variant.Tests
			if variant.Solution != "" {
				solution = variant.Solution
			}
		}
	}

	var results []CodeSnippetTestResult
	allPassed := true

	for _, t := range tests {
		passed := runSnippetTest(t, code)
		results = append(results, CodeSnippetTestResult{Name: t.Name, Passed: passed})
		if !passed {
			allPassed = false
		}
	}

	var correctAnswer *Answer
	if !allPassed && len(solution) > 0 {
		correctAnswer = &Answer{Code: &solution}
	}

	return CodeSnippetResult{
		Result: Result{
			Correct:       allPassed,
			Explanation:   p.Explanation,
			CorrectAnswer: correctAnswer,
		},
		Tests: results,
	}, nil
}

// runSnippetTest executes one assertion against the submitted code.
func runSnippetTest(t lesson.CodeSnippetTest, code string) bool {
	switch t.Type {
	case "contains":
		return strings.Contains(code, t.Value)
	case "not_contains":
		return !strings.Contains(code, t.Value)
	case "pattern":
		re, err := regexp.Compile(t.Value)
		if err != nil {
			// Invalid regex — authoring error that should have been caught by validation.
			// Fail safe: treat as not matching rather than panicking.
			return false
		}
		return re.MatchString(code)
	case "structure":
		// All values must appear in the code in order; each searched from after the previous match.
		pos := 0
		for _, v := range t.Values {
			idx := strings.Index(code[pos:], v)
			if idx == -1 {
				return false
			}
			pos += idx + len(v)
		}
		return true
	}
	return false
}

// MCQ grades multiple-choice answers.
func MCQ(p lesson.MCQPayload, a Answer) (Result, error) {
	if a.Choice == nil {
		return Result{}, fmt.Errorf("%w: choice is required", ErrInvalidAnswer)
	}
	c := *a.Choice
	if c < 0 || c >= len(p.Options) {
		return Result{}, fmt.Errorf("%w: choice %d out of range [0, %d)", ErrInvalidAnswer, c, len(p.Options))
	}
	correct := c == p.Answer
	res := Result{Correct: correct}
	if correct {
		res.Explanation = p.Explanations.Correct
	} else {
		res.Explanation = p.Explanations.Correct
		res.WhyNot = p.Explanations.Wrong[c]
		res.CorrectAnswer = &Answer{Choice: &p.Answer}
	}
	return res, nil
}

// TrueFalse grades boolean statement answers.
func TrueFalse(p lesson.TrueFalsePayload, a Answer) (Result, error) {
	if a.Value == nil {
		return Result{}, fmt.Errorf("%w: value is required", ErrInvalidAnswer)
	}
	correct := *a.Value == p.Answer
	res := Result{Correct: correct, Explanation: p.Explanation}
	if !correct {
		res.CorrectAnswer = &Answer{Value: &p.Answer}
	}
	return res, nil
}

// FillBlank grades ordered token blank fills against the exact answer string.
func FillBlank(p lesson.FillBlankPayload, a Answer) (Result, error) {
	if len(a.Blanks) != len(p.Answers) {
		return Result{}, fmt.Errorf("%w: expected %d blanks, got %d", ErrInvalidAnswer, len(p.Answers), len(a.Blanks))
	}
	correct := true
	for i := range p.Answers {
		if strings.TrimSpace(a.Blanks[i]) != strings.TrimSpace(p.Answers[i]) {
			correct = false
			break
		}
	}
	res := Result{Correct: correct, Explanation: p.Explanation}
	if !correct {
		res.CorrectAnswer = &Answer{Blanks: append([]string(nil), p.Answers...)}
	}
	return res, nil
}

// Arrange grades sequence permutation against the target order index array.
func Arrange(p lesson.ArrangePayload, a Answer) (Result, error) {
	if len(a.Order) != len(p.Order) {
		return Result{}, fmt.Errorf("%w: expected %d items, got %d", ErrInvalidAnswer, len(p.Order), len(a.Order))
	}
	seen := make(map[int]bool, len(p.Order))
	for _, idx := range a.Order {
		if idx < 0 || idx >= len(p.Items) || seen[idx] {
			return Result{}, fmt.Errorf("%w: invalid or duplicate index %d", ErrInvalidAnswer, idx)
		}
		seen[idx] = true
	}
	correct := true
	for i := range p.Order {
		if a.Order[i] != p.Order[i] {
			correct = false
			break
		}
	}
	res := Result{Correct: correct, Explanation: p.Explanation}
	if !correct {
		res.CorrectAnswer = &Answer{Order: append([]int(nil), p.Order...)}
	}
	return res, nil
}
