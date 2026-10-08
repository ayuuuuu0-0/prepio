package integration_test

import (
	"context"
	"os"
	"path/filepath"
	"sort"
	"strings"
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

	// Every migration from 000035 (the lesson platform) up to the newest is rolled back newest
	// first, then re-applied oldest first. Later migrations build on these (focus topics, leagues,
	// world topics), so the list is read from the directory and grows with every new migration.
	entries, err := os.ReadDir(filepath.Join("..", "..", "migrations"))
	require.NoError(t, err)
	var l2 []string
	for _, e := range entries {
		name := e.Name()
		if strings.HasSuffix(name, ".up.sql") && name >= "000035" {
			l2 = append(l2, strings.TrimSuffix(name, ".up.sql"))
		}
	}
	sort.Strings(l2)
	require.GreaterOrEqual(t, len(l2), 5)
	reversed := make([]string, len(l2))
	for i, name := range l2 {
		reversed[len(l2)-1-i] = name
	}
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
		"lesson_step_results", "node_prerequisites", "mastery_ledger", "lesson_rewards", "user_focus_topics"} {
		require.True(t, tableExists(table), table)
	}

	exec("down", reversed)
	for _, table := range []string{"topics", "lessons", "lesson_steps", "lesson_skills", "lesson_attempts",
		"lesson_step_results", "node_prerequisites", "mastery_ledger", "lesson_rewards", "user_focus_topics"} {
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
