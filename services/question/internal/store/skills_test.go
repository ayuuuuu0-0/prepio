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
	require.Len(t, categories, 9)
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
