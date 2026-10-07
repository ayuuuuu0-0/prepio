package integration_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/prepio/prepio/test/testdb"
	"github.com/stretchr/testify/require"
)

// TestLessonMigrationsRollBack proves the Phase L2 migrations apply, roll back
// cleanly in reverse order, and apply again on the same database.
func TestLessonMigrationsRollBack(t *testing.T) {
	pool, _ := testdb.Start(t)
	testdb.Migrate(t, pool)
	ctx := context.Background()

	l2 := []string{"000035_create_topics", "000036_create_lessons", "000037_create_mastery_ledger"}
	exec := func(suffix string, names []string) {
		t.Helper()
		for _, name := range names {
			raw, err := os.ReadFile(filepath.Join("..", "..", "migrations", name+"."+suffix+".sql"))
			require.NoError(t, err)
			_, err = pool.Exec(ctx, string(raw))
			require.NoError(t, err, "%s.%s", name, suffix)
		}
	}
	tableExists := func(name string) bool {
		var ok bool
		require.NoError(t, pool.QueryRow(ctx, `SELECT to_regclass('public.'||$1) IS NOT NULL`, name).Scan(&ok))
		return ok
	}

	for _, table := range []string{"topics", "lessons", "lesson_steps", "lesson_skills", "lesson_attempts",
		"lesson_step_results", "node_prerequisites", "mastery_ledger", "lesson_rewards"} {
		require.True(t, tableExists(table), table)
	}

	exec("down", []string{l2[2], l2[1], l2[0]})
	for _, table := range []string{"topics", "lessons", "lesson_steps", "lesson_skills", "lesson_attempts",
		"lesson_step_results", "node_prerequisites", "mastery_ledger", "lesson_rewards"} {
		require.False(t, tableExists(table), table)
	}
	var backend int
	require.NoError(t, pool.QueryRow(ctx, `SELECT count(*) FROM skills WHERE slug LIKE 'backend-%'`).Scan(&backend))
	require.Zero(t, backend, "rollback removes the skills it added")

	exec("up", l2)
	require.True(t, tableExists("lessons"))
	require.NoError(t, pool.QueryRow(ctx, `SELECT count(*) FROM skills WHERE slug LIKE 'backend-%'`).Scan(&backend))
	require.Equal(t, 4, backend)

	var topics, mapped int
	require.NoError(t, pool.QueryRow(ctx, `SELECT count(*) FROM topics`).Scan(&topics))
	require.NoError(t, pool.QueryRow(ctx, `SELECT count(DISTINCT topic_id) FROM skill_categories WHERE topic_id IS NOT NULL`).Scan(&mapped))
	require.Equal(t, 4, topics)
	require.Equal(t, 4, mapped, "every topic has at least one skill category")
}

// TestLegacyQuestionLoopMigrationRollsBack proves the Phase L3 drop applies, that every
// legacy table is gone afterwards, and that its rollback restores the structure.
func TestLegacyQuestionLoopMigrationRollsBack(t *testing.T) {
	pool, _ := testdb.Start(t)
	testdb.Migrate(t, pool)
	ctx := context.Background()

	legacy := []string{
		"questions", "question_skills", "question_pools", "pool_questions", "node_pools",
		"node_skills", "user_question_history", "daily_papers", "daily_paper_questions", "user_journey_progress",
	}
	kept := []string{"skills", "subskills", "user_skill_scores", "worlds", "journey_nodes", "lessons", "mastery_ledger"}

	exists := func(name string) bool {
		var ok bool
		require.NoError(t, pool.QueryRow(ctx, `SELECT to_regclass('public.'||$1) IS NOT NULL`, name).Scan(&ok))
		return ok
	}
	run := func(suffix string) {
		t.Helper()
		raw, err := os.ReadFile(filepath.Join("..", "..", "migrations", "000038_drop_legacy_question_loop."+suffix+".sql"))
		require.NoError(t, err)
		_, err = pool.Exec(ctx, string(raw))
		require.NoError(t, err, suffix)
	}

	for _, table := range legacy {
		require.False(t, exists(table), "%s must be dropped by the migration", table)
	}
	for _, table := range kept {
		require.True(t, exists(table), "%s must survive", table)
	}

	run("down")
	for _, table := range legacy {
		require.True(t, exists(table), "%s must be restored by the rollback", table)
	}

	run("up")
	for _, table := range legacy {
		require.False(t, exists(table), "%s must be dropped again", table)
	}
}
