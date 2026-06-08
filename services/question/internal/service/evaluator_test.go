package service_test

import (
	"testing"

	"github.com/prepio/prepio/services/question/internal/service"
	"github.com/stretchr/testify/require"
)

func TestEvaluateAnswer(t *testing.T) {
	guide := "hash map approach with O(n) time and O(n) space"
	require.True(t, service.EvaluateAnswer("use a hash map for O(n) time and O(n) space", guide))
	require.False(t, service.EvaluateAnswer("short", guide))
}
