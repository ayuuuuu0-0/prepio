package service

import (
	"fmt"
	"strings"

	"github.com/prepio/prepio/constants"
)

// EvaluationResult holds scored rubric output for an answer submission.
type EvaluationResult struct {
	Score     int
	Correct   bool
	Strengths []string
	Gaps      []string
	Summary   string
}

// EvaluateAnswer scores the answer against required concepts in the answer guide.
func EvaluateAnswer(answer, guide string) EvaluationResult {
	trimmed := strings.TrimSpace(answer)
	if len(trimmed) < constants.MinAnswerLength {
		return EvaluationResult{
			Score: 0,
			Gaps: []string{
				fmt.Sprintf("Write at least %d characters — a meaningful technical answer needs more detail", constants.MinAnswerLength),
			},
			Summary: "Your answer is too brief to evaluate fairly.",
		}
	}

	concepts := parseRequiredConcepts(guide)
	if len(concepts) == 0 {
		return EvaluationResult{
			Score:   0,
			Gaps:    []string{"Question rubric is missing required concepts"},
			Summary: "Unable to evaluate this question — please try again later.",
		}
	}

	answerLower := strings.ToLower(trimmed)
	matched := make([]string, 0, len(concepts))
	missed := make([]string, 0, len(concepts))

	for _, concept := range concepts {
		if conceptMatched(answerLower, concept) {
			matched = append(matched, concept)
		} else {
			missed = append(missed, concept)
		}
	}

	score := (len(matched) * 100) / len(concepts)
	correct := score >= constants.MinEvaluationScore

	strengths := matched
	if len(strengths) > 3 {
		strengths = strengths[:3]
	}

	gaps := missed
	if len(gaps) > 3 {
		gaps = gaps[:3]
	}

	summary := summaryFor(score, correct)

	return EvaluationResult{
		Score:     score,
		Correct:   correct,
		Strengths: strengths,
		Gaps:      gaps,
		Summary:   summary,
	}
}

// parseRequiredConcepts extracts pipe-separated concepts from a concepts: prefixed guide.
func parseRequiredConcepts(guide string) []string {
	guide = strings.TrimSpace(guide)
	lower := strings.ToLower(guide)
	const prefix = "concepts:"
	if !strings.HasPrefix(lower, prefix) {
		return nil
	}

	rest := strings.TrimSpace(guide[len(prefix):])
	if idx := strings.Index(rest, "\n"); idx >= 0 {
		rest = rest[:idx]
	}

	parts := strings.Split(rest, "|")
	concepts := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if len(part) > 0 {
			concepts = append(concepts, part)
		}
	}
	return concepts
}

// conceptMatched reports whether the answer contains the required concept phrase.
func conceptMatched(answerLower, concept string) bool {
	concept = strings.ToLower(strings.TrimSpace(concept))
	if len(concept) == 0 {
		return false
	}
	return strings.Contains(answerLower, concept)
}

func summaryFor(score int, correct bool) string {
	if correct && score >= 90 {
		return "Excellent — you covered the core concepts clearly."
	}
	if correct {
		return "Good work — key concepts covered, with room to go deeper."
	}
	if score >= 40 {
		return "You're on the right track — review the missed concepts and try again."
	}
	return "Keep going — focus on the approach, complexity, and tradeoffs."
}

// FeedbackFor returns a short legacy feedback string.
func FeedbackFor(correct bool) string {
	if correct {
		return "solid answer — key concepts covered"
	}
	return "keep going — review the core concepts and try again tomorrow"
}
