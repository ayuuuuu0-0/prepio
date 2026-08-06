package service_test

import (
	"context"
	"strings"
	"testing"

	"github.com/prepio/prepio/constants"
	"github.com/prepio/prepio/services/question/internal/service"
	"github.com/stretchr/testify/require"
)

// --- Existing Stage 1+2 regression tests (unchanged behaviour) ---------------

func TestEvaluateAnswerConcepts(t *testing.T) {
	guide := `concepts:hash map|O(n) time|O(n) space|two sum|duplicate handling
expect hash map approach with O(n) time and O(n) space`

	longAnswer := strings.Repeat("a", constants.MinAnswerLength) +
		" use a hash map for O(n) time and O(n) space with two sum and duplicate handling"

	result := service.EvaluateAnswer(longAnswer, guide)
	require.True(t, result.Correct)
	require.GreaterOrEqual(t, result.Score, constants.MinEvaluationScore)
	require.NotEmpty(t, result.Strengths)
}

func TestEvaluateAnswerRejectsShortAnswer(t *testing.T) {
	guide := "concepts:hash map|O(n) time|O(n) space"
	result := service.EvaluateAnswer("hash map only", guide)
	require.False(t, result.Correct)
	require.Equal(t, 0, result.Score)
	require.NotEmpty(t, result.Gaps)
}

func TestEvaluateAnswerRejectsKeywordStuffing(t *testing.T) {
	guide := "concepts:dynamic programming|memoization|optimal substructure"
	padding := strings.Repeat("x ", constants.MinAnswerLength)
	result := service.EvaluateAnswer(padding+" arrays are good ", guide)
	require.False(t, result.Correct)
}

func TestEvaluateBinaryTreeGuide(t *testing.T) {
	guide := `concepts:recursive|iterative|dfs|bfs|maximum depth|time complexity|space complexity
recursive or iterative dfs/bfs acceptable`

	padding := strings.Repeat("detail ", 20)
	answer := padding + "Use recursive dfs with time complexity O(n) and space complexity O(h) for maximum depth."
	result := service.EvaluateAnswer(answer, guide)
	require.True(t, result.Correct)
	require.GreaterOrEqual(t, result.Score, 60)
}

// --- Stage 1: keyword filter tests -------------------------------------------

func TestStage1RejectsAnswerBelowFilterThreshold(t *testing.T) {
	// Guide has 5 required concepts; answer only mentions 1 (20% < 35% threshold).
	guide := "concepts:dynamic programming|memoization|optimal substructure|overlapping subproblems|base case"
	padding := strings.Repeat("word ", 25) // meets MinAnswerLength
	answer := padding + "dynamic programming is used"

	result := service.EvaluateAnswer(answer, guide)
	require.False(t, result.Correct)
	require.Equal(t, 0, result.Score)
	require.Equal(t, "filter_rejected", result.EvaluatedBy)
	require.NotEmpty(t, result.Gaps, "filter-rejected result should list missed concepts")
}

func TestStage1PassesAnswerAboveFilterThreshold(t *testing.T) {
	// Guide has 4 required concepts; answer mentions 3 (75% ≥ 35% threshold).
	guide := "concepts:hash map|O(n) time|O(n) space|two pointers"
	padding := strings.Repeat("word ", 20)
	answer := padding + "use a hash map for O(n) time and O(n) space lookup"

	result := service.EvaluateAnswer(answer, guide)
	// Stage 1 passed — result is from Stage 2, so EvaluatedBy should be structural.
	require.Equal(t, "structural", result.EvaluatedBy)
}

// --- Stage 2: structural coverage tests --------------------------------------

func TestStage2RecordsStructuralScore(t *testing.T) {
	guide := "concepts:binary search|sorted array|O(log n)|mid pointer"
	padding := strings.Repeat("word ", 20)
	answer := padding + "binary search on a sorted array uses a mid pointer to achieve O(log n)"

	result := service.EvaluateAnswer(answer, guide)
	require.Equal(t, result.Score, result.StructuralScore, "structural score must equal main score when no LLM is used")
	require.True(t, result.Correct)
}

func TestStage2PartialCoverageGivesCorrectScore(t *testing.T) {
	// 2 of 4 concepts covered → score = 50, below MinEvaluationScore (60).
	// Answer mentions only "stack" and "LIFO"; "O(1) push" and "O(1) pop" are absent.
	guide := "concepts:stack|O(1) push|O(1) pop|LIFO"
	padding := strings.Repeat("word ", 20)
	answer := padding + "stack is a data structure that follows LIFO ordering"

	result := service.EvaluateAnswer(answer, guide)
	require.False(t, result.Correct)
	require.Equal(t, 50, result.Score)
	require.Equal(t, "structural", result.EvaluatedBy)
}

// --- Guide parser tests -------------------------------------------------------

func TestParseAnswerGuideLegacyFormat(t *testing.T) {
	raw := "concepts:hash map|O(n) time|two sum"
	guide := service.ParseAnswerGuide(raw)
	require.Len(t, guide.Concepts, 3)
	for _, c := range guide.Concepts {
		require.True(t, c.Required, "legacy concepts should all be required")
	}
}

func TestParseAnswerGuideJSONFormat(t *testing.T) {
	raw := `{
		"concepts": [
			{"name": "hash map", "aliases": ["hashmap", "hash_map"], "required": true},
			{"name": "O(n) time", "required": true},
			{"name": "two pointers", "required": false}
		]
	}`
	guide := service.ParseAnswerGuide(raw)
	require.Len(t, guide.Concepts, 3)
	require.Equal(t, "hash map", guide.Concepts[0].Name)
	require.ElementsMatch(t, []string{"hashmap", "hash_map"}, guide.Concepts[0].Aliases)
	require.False(t, guide.Concepts[2].Required)
}

func TestParseAnswerGuideAliasMatches(t *testing.T) {
	// "hashmap" alias should match even though primary name is "hash map".
	raw := `{
		"concepts": [
			{"name": "hash map", "aliases": ["hashmap"], "required": true}
		]
	}`
	padding := strings.Repeat("word ", 20)
	answer := padding + "we use a hashmap to store visited nodes"

	result := service.EvaluateAnswer(answer, raw)
	require.Equal(t, "structural", result.EvaluatedBy)
	require.NotEmpty(t, result.Strengths)
}

func TestParseAnswerGuideFallsBackToLegacy(t *testing.T) {
	// Malformed JSON should silently fall back to legacy parser, which finds nothing.
	raw := `{invalid json`
	guide := service.ParseAnswerGuide(raw)
	// Legacy parser requires "concepts:" prefix; this input doesn't have it.
	require.Empty(t, guide.Concepts)
}

// --- Pipeline: LLM never called when Stage 1 rejects -------------------------

func TestPipelineDoesNotCallLLMWhenStage1Rejects(t *testing.T) {
	// callCount lets us verify the LLM was never invoked.
	spy := &llmSpy{}

	evaluator := service.NewPipelineEvaluator(spy)

	// Answer has correct length but mentions none of the 5 required concepts.
	guide := "concepts:dynamic programming|memoization|optimal substructure|overlapping subproblems|base case"
	padding := strings.Repeat("word ", 25)
	answer := padding + "arrays are a fundamental data structure"

	result := evaluator.Evaluate(answer, guide)

	require.Equal(t, "filter_rejected", result.EvaluatedBy)
	require.Zero(t, spy.calls, "LLM must never be called when Stage 1 rejects")
}

// llmSpy is a test double that counts how many times Evaluate is called.
type llmSpy struct {
	calls int
}

func (s *llmSpy) Evaluate(_ context.Context, _ service.LLMRequest) (service.LLMResponse, error) {
	s.calls++
	return service.LLMResponse{}, nil
}
