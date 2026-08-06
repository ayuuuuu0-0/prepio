package service

import (
	"fmt"
	"strings"

	"github.com/prepio/prepio/constants"
)

// EvaluationResult holds scored rubric output for an answer submission.
type EvaluationResult struct {
	Score   int
	Correct bool

	// Top 3 matched concept names shown to the user as positive feedback.
	Strengths []string

	// Top 3 missed concept names shown to the user as improvement areas.
	Gaps []string

	Summary string

	// StructuralScore is the raw Stage 2 percentage (0–100) before any LLM adjustment.
	// Always populated; equals Score when the LLM is not used.
	StructuralScore int

	// EvaluatedBy records which stage produced the final score.
	// Values: "filter_rejected" | "structural" | "llm"
	EvaluatedBy string
}

// ConceptCoverage records the match result for a single concept.
type ConceptCoverage struct {
	Concept      string
	Required     bool
	Matched      bool
	MatchedAlias string // empty string means matched by primary name
}

// CoverageReport is the Stage 2 output: per-concept match details and aggregate scores.
// It is passed as grounding context to the Stage 3 LLM prompt.
type CoverageReport struct {
	Concepts        []ConceptCoverage
	RequiredMatched int
	RequiredTotal   int
	OptionalMatched int
	OptionalTotal   int
	// StructuralScore is (RequiredMatched + OptionalMatched) / total * 100.
	StructuralScore int
}

// Evaluator is the injection point for answer scoring.
// StructuralEvaluator satisfies this interface using Stages 1+2 only.
// PipelineEvaluator additionally calls the LLM (Stage 3) when configured.
type Evaluator interface {
	Evaluate(answer, rawGuide string) EvaluationResult
}

// --- Stage 1: keyword filter --------------------------------------------------

// keywordFilter checks whether the answer contains enough required concepts to
// be worth a full structural evaluation. Returns true if the ratio of matched
// required concepts meets the minimum threshold.
//
// This stage rejects obviously incomplete answers in microseconds so that
// neither Stage 2 nor (when wired) Stage 3 wastes time on them.
func keywordFilter(answerLower string, guide AnswerGuide) (pass bool, matchedRequired int) {
	required := guide.RequiredConcepts()
	if len(required) == 0 {
		// No required concepts declared — let Stage 2 decide.
		return true, 0
	}

	for _, c := range required {
		if conceptMatched(answerLower, c) {
			matchedRequired++
		}
	}

	ratio := float64(matchedRequired) / float64(len(required))
	return ratio >= constants.MinKeywordFilterRatio, matchedRequired
}

// --- Stage 2: structural coverage --------------------------------------------

// buildCoverageReport matches every concept (required + optional) in the guide
// against the answer and builds a detailed coverage report.
func buildCoverageReport(answerLower string, guide AnswerGuide) CoverageReport {
	coverages := make([]ConceptCoverage, 0, len(guide.Concepts))
	reqMatched, reqTotal, optMatched, optTotal := 0, 0, 0, 0

	for _, c := range guide.Concepts {
		cov := ConceptCoverage{
			Concept:  c.Name,
			Required: c.Required,
		}

		alias := matchedAlias(answerLower, c)
		if alias != "" {
			cov.Matched = true
			if alias != c.Name {
				cov.MatchedAlias = alias
			}
		}

		if c.Required {
			reqTotal++
			if cov.Matched {
				reqMatched++
			}
		} else {
			optTotal++
			if cov.Matched {
				optMatched++
			}
		}

		coverages = append(coverages, cov)
	}

	total := reqTotal + optTotal
	allMatched := reqMatched + optMatched
	score := 0
	if total > 0 {
		score = allMatched * 100 / total
	}

	return CoverageReport{
		Concepts:        coverages,
		RequiredMatched: reqMatched,
		RequiredTotal:   reqTotal,
		OptionalMatched: optMatched,
		OptionalTotal:   optTotal,
		StructuralScore: score,
	}
}

// coverageToResult converts a CoverageReport into an EvaluationResult.
func coverageToResult(report CoverageReport) EvaluationResult {
	score := report.StructuralScore
	correct := score >= constants.MinEvaluationScore

	var strengths, gaps []string
	for _, c := range report.Concepts {
		if c.Matched {
			strengths = append(strengths, c.Concept)
		} else {
			gaps = append(gaps, c.Concept)
		}
	}
	if len(strengths) > 3 {
		strengths = strengths[:3]
	}
	if len(gaps) > 3 {
		gaps = gaps[:3]
	}

	return EvaluationResult{
		Score:           score,
		Correct:         correct,
		Strengths:       strengths,
		Gaps:            gaps,
		Summary:         summaryFor(score, correct),
		StructuralScore: score,
		EvaluatedBy:     "structural",
	}
}

// --- StructuralEvaluator (Stages 1+2 only) -----------------------------------

// StructuralEvaluator runs Stage 1 (keyword filter) and Stage 2 (coverage report)
// and produces a deterministic result without any LLM call.
// It is used in all environments — as the sole evaluator when no LLM is configured,
// and as the fallback inside PipelineEvaluator when the LLM is unavailable.
type StructuralEvaluator struct{}

// Evaluate implements Evaluator.
func (e StructuralEvaluator) Evaluate(answer, rawGuide string) EvaluationResult {
	return EvaluateAnswer(answer, rawGuide)
}

// --- Public entry point (backward-compatible) ---------------------------------

// EvaluateAnswer scores the answer against the guide using Stage 1 + Stage 2.
// It is the canonical structural evaluation path and remains the public API so
// that existing call sites and tests require no changes.
func EvaluateAnswer(answer, rawGuide string) EvaluationResult {
	trimmed := strings.TrimSpace(answer)

	// Minimum length guard (pre-filter before concept parsing).
	if len(trimmed) < constants.MinAnswerLength {
		return EvaluationResult{
			Score:       0,
			EvaluatedBy: "filter_rejected",
			Gaps: []string{
				fmt.Sprintf("Write at least %d characters — a meaningful technical answer needs more detail", constants.MinAnswerLength),
			},
			Summary: "Your answer is too brief to evaluate fairly.",
		}
	}

	guide := ParseAnswerGuide(rawGuide)
	if len(guide.Concepts) == 0 {
		return EvaluationResult{
			Score:       0,
			EvaluatedBy: "filter_rejected",
			Gaps:        []string{"Question rubric is missing required concepts"},
			Summary:     "Unable to evaluate this question — please try again later.",
		}
	}

	answerLower := strings.ToLower(trimmed)

	// Stage 1: keyword filter.
	pass, _ := keywordFilter(answerLower, guide)
	if !pass {
		return EvaluationResult{
			Score:       0,
			Correct:     false,
			EvaluatedBy: "filter_rejected",
			Gaps:        topMissedRequired(answerLower, guide),
			Summary:     "Your answer is missing several key concepts — review the topic and try again.",
		}
	}

	// Stage 2: full structural coverage.
	report := buildCoverageReport(answerLower, guide)
	return coverageToResult(report)
}

// --- Helpers ------------------------------------------------------------------

// conceptMatched reports whether the answer contains the concept's primary name
// or any of its aliases.
func conceptMatched(answerLower string, c Concept) bool {
	return matchedAlias(answerLower, c) != ""
}

// matchedAlias returns the alias (or primary name) that matched, or empty string.
func matchedAlias(answerLower string, c Concept) string {
	primary := strings.ToLower(strings.TrimSpace(c.Name))
	if len(primary) > 0 && strings.Contains(answerLower, primary) {
		return c.Name
	}
	for _, alias := range c.Aliases {
		a := strings.ToLower(strings.TrimSpace(alias))
		if len(a) > 0 && strings.Contains(answerLower, a) {
			return alias
		}
	}
	return ""
}

// topMissedRequired returns up to 3 required concept names that the answer missed.
func topMissedRequired(answerLower string, guide AnswerGuide) []string {
	var missed []string
	for _, c := range guide.RequiredConcepts() {
		if !conceptMatched(answerLower, c) {
			missed = append(missed, c.Name)
		}
		if len(missed) == 3 {
			break
		}
	}
	return missed
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
