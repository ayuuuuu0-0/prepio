package service_test

import (
	"strings"
	"testing"

	"github.com/prepio/prepio/constants"
	"github.com/prepio/prepio/services/question/internal/service"
	"github.com/stretchr/testify/require"
)

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
