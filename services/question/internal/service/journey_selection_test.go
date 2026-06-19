package service

import (
	"testing"

	"github.com/prepio/prepio/constants"
	"github.com/prepio/prepio/services/question/internal/store"
	"github.com/stretchr/testify/require"
)

func TestSelectSequentialPrefersUnseen(t *testing.T) {
	seen := map[string]bool{"q1": true}
	id := selectSequential([]string{"q1", "q2", "q3"}, seen)
	require.Equal(t, "q2", id)
}

func TestSelectSequentialFallsBackWhenAllSeen(t *testing.T) {
	seen := map[string]bool{"q1": true, "q2": true}
	id := selectSequential([]string{"q1", "q2"}, seen)
	require.Equal(t, "q1", id)
}

func TestSelectRandomUnseenPrefersUnseen(t *testing.T) {
	seen := map[string]bool{"q1": true, "q2": true}
	id := selectRandomUnseen([]string{"q1", "q2", "q3"}, seen, 42)
	require.Equal(t, "q3", id)
}

func TestSelectRandomUnseenFallsBackWhenAllSeen(t *testing.T) {
	seen := map[string]bool{"q1": true}
	id := selectRandomUnseen([]string{"q1"}, seen, 99)
	require.Equal(t, "q1", id)
}

func TestSelectFromPoolStrategies(t *testing.T) {
	pool := []string{"a", "b", "c"}
	seen := map[string]bool{}

	id, ok := selectFromPool(constants.PoolSelectionSequential, pool, seen, 1)
	require.True(t, ok)
	require.Equal(t, "a", id)

	id, ok = selectFromPool(constants.PoolSelectionRandomUnseen, pool, seen, 7)
	require.True(t, ok)
	require.NotEmpty(t, id)
}

func TestAllAnsweredInSession(t *testing.T) {
	answered := map[string]bool{"q1": true, "q2": true}
	require.True(t, allAnsweredInSession([]string{"q1", "q2"}, answered))
	require.False(t, allAnsweredInSession([]string{"q1", "q3"}, answered))
	require.False(t, allAnsweredInSession(nil, answered))
}

func TestDisplayQuestionID(t *testing.T) {
	answered := map[string]bool{"q1": true}
	require.Equal(t, "q2", displayQuestionID([]string{"q1", "q2"}, answered))
	require.Equal(t, "q1", displayQuestionID([]string{"q1"}, answered))
	require.Empty(t, displayQuestionID(nil, answered))
}

func TestNodeUnlocked(t *testing.T) {
	assignments := []nodeAssignment{
		{questionIDs: []string{"q1"}},
		{questionIDs: []string{"q2"}},
	}
	answered := map[string]bool{"q1": true}
	require.True(t, nodeUnlocked(0, assignments, answered))
	require.True(t, nodeUnlocked(1, assignments, answered))
	require.False(t, nodeUnlocked(1, assignments, map[string]bool{}))
}

func TestSelectionSeedIsStable(t *testing.T) {
	a := selectionSeed("user", "session", "node")
	b := selectionSeed("user", "session", "node")
	c := selectionSeed("user", "session", "other")
	require.Equal(t, a, b)
	require.NotEqual(t, a, c)
}

func TestIsBossMixed(t *testing.T) {
	require.True(t, isBossMixed([]store.NodePoolBinding{
		{SelectionStrategy: constants.PoolSelectionBossMixed},
	}))
	require.True(t, isBossMixed([]store.NodePoolBinding{
		{SelectionStrategy: constants.PoolSelectionRandomUnseen},
		{SelectionStrategy: constants.PoolSelectionRandomUnseen},
	}))
	require.False(t, isBossMixed([]store.NodePoolBinding{
		{SelectionStrategy: constants.PoolSelectionRandomUnseen},
	}))
}
