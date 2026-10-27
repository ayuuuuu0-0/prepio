package store

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStreakAchievementsAreIdempotent(t *testing.T) {
	pool := leaguePool(t)
	ctx := context.Background()
	user := leagueUser(t, pool, 1)
	s := NewAchievementStore(pool)

	got, err := s.AwardStreak(ctx, user, 2)
	require.NoError(t, err)
	require.Empty(t, got, "a 2-day streak earns nothing")

	got, err = s.AwardStreak(ctx, user, 3)
	require.NoError(t, err)
	require.Equal(t, []string{"warming-up"}, got)

	got, err = s.AwardStreak(ctx, user, 7)
	require.NoError(t, err)
	require.Equal(t, []string{"full-week"}, got, "only what is new is reported")

	got, err = s.AwardStreak(ctx, user, 7)
	require.NoError(t, err)
	require.Empty(t, got, "awarding again changes nothing")

	list, err := s.List(ctx, user)
	require.NoError(t, err)
	require.Len(t, list, 2)
}

func TestLessonFactsMatchTopicReadiness(t *testing.T) {
	pool := leaguePool(t)
	ctx := context.Background()
	user := leagueUser(t, pool, 1)

	// Two started System Design skills at 40 and 61 (mean 50.5 rounds to 51) and one
	// untouched skill that must not count.
	_, err := pool.Exec(ctx, `
		INSERT INTO user_skill_scores (user_id, skill_id, mastery, attempts)
		SELECT $1, id, CASE slug WHEN 'system-design-scaling' THEN 40 ELSE 61 END, 1
		FROM skills WHERE slug IN ('system-design-scaling', 'system-design-data')`, user)
	require.NoError(t, err)
	_, err = pool.Exec(ctx, `
		INSERT INTO user_skill_scores (user_id, skill_id, mastery, attempts)
		SELECT $1, id, 0, 0 FROM skills WHERE slug = 'system-design-reliability'`, user)
	require.NoError(t, err)

	facts, err := lessonFacts(ctx, pool, user, 3)
	require.NoError(t, err)
	require.Equal(t, 51, facts.BestTopicMastery)
	require.Equal(t, 3, facts.Level)
	require.Zero(t, facts.Lessons)
	require.False(t, facts.Promoted)
}
