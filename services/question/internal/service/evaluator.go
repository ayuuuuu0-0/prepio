package service

import (
	"strings"
	"unicode"
)

// EvaluateAnswer checks the answer against the rubric guide using keyword overlap.
func EvaluateAnswer(answer, guide string) bool {
	if len(strings.TrimSpace(answer)) < 10 {
		return false
	}

	answerLower := strings.ToLower(answer)
	matchCount := 0
	for _, word := range strings.FieldsFunc(guide, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	}) {
		if len(word) < 4 {
			continue
		}
		if strings.Contains(answerLower, strings.ToLower(word)) {
			matchCount++
		}
	}

	return matchCount >= 2
}

// FeedbackFor returns user-facing feedback for a submission result.
func FeedbackFor(correct bool) string {
	if correct {
		return "solid answer — key concepts covered"
	}
	return "keep going — review the core concepts and try again tomorrow"
}
