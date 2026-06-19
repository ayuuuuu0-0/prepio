package store_test

import (
	"context"
	"testing"

	"github.com/prepio/prepio/services/question/internal/store"
	"github.com/prepio/prepio/test/testdb"
	"github.com/stretchr/testify/require"
)

func TestSkillStoreListCategories(t *testing.T) {
	pool, _ := testdb.Start(t)
	testdb.Migrate(t, pool)

	skillStore := store.NewSkillStore(pool)
	ctx := context.Background()

	categories, err := skillStore.ListCategories(ctx)
	require.NoError(t, err)
	require.Len(t, categories, 8)
	require.Equal(t, "programming-fundamentals", categories[0].Slug)
}

func TestSkillStoreGetSkillBySlug(t *testing.T) {
	pool, _ := testdb.Start(t)
	testdb.Migrate(t, pool)

	skillStore := store.NewSkillStore(pool)
	ctx := context.Background()

	skill, err := skillStore.GetSkillBySlug(ctx, "arrays")
	require.NoError(t, err)
	require.NotNil(t, skill)
	require.Equal(t, "Arrays", skill.Name)

	subskills, err := skillStore.ListSubskillsBySkillID(ctx, skill.ID)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(subskills), 5)
}

func TestSkillStoreListQuestionSkills(t *testing.T) {
	pool, _ := testdb.Start(t)
	testdb.Migrate(t, pool)

	skillStore := store.NewSkillStore(pool)
	ctx := context.Background()

	mappings, err := skillStore.ListQuestionSkills(ctx, "b0000000-0000-4000-8000-000000000001")
	require.NoError(t, err)
	require.Len(t, mappings, 2)

	totalWeight := 0.0
	for _, mapping := range mappings {
		totalWeight += mapping.Weight
		require.NotEmpty(t, mapping.SkillSlug)
		require.NotEmpty(t, mapping.SubskillSlug)
	}
	require.InDelta(t, 1.0, totalWeight, 0.001)
}

func TestSkillStoreGetQuestionContentMetadata(t *testing.T) {
	pool, _ := testdb.Start(t)
	testdb.Migrate(t, pool)

	skillStore := store.NewSkillStore(pool)
	ctx := context.Background()

	meta, err := skillStore.GetQuestionContentMetadata(ctx, "b0000000-0000-4000-8000-000000000001")
	require.NoError(t, err)
	require.NotNil(t, meta)
	require.Equal(t, "coding", meta.EvaluationType)
	require.NotEmpty(t, meta.Explanation)
	require.NotEmpty(t, meta.Solution)
	require.Greater(t, meta.ReadinessWeight, 0.0)
	require.Greater(t, meta.EstimatedTime, 0)
	require.NotEmpty(t, meta.Hints)
}
