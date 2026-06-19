package store_test

import (
	"context"
	"testing"

	"github.com/prepio/prepio/constants"
	"github.com/prepio/prepio/services/question/internal/store"
	"github.com/prepio/prepio/test/testdb"
	"github.com/stretchr/testify/require"
)

func TestContentStoreFoundationForestBindings(t *testing.T) {
	pool, _ := testdb.Start(t)
	testdb.Migrate(t, pool)

	contentStore := store.NewContentStore(pool)
	journeyStore := store.NewJourneyStore(pool)
	ctx := context.Background()

	world, err := journeyStore.GetWorldBySlug(ctx, constants.FoundationForestWorldSlug)
	require.NoError(t, err)
	require.NotNil(t, world)

	nodes, err := journeyStore.ListNodesByWorld(ctx, world.ID)
	require.NoError(t, err)
	require.Len(t, nodes, 5)

	require.Equal(t, "arrays-basics", nodes[0].Slug)
	require.Equal(t, "forest-boss", nodes[4].Slug)

	arraysContent, err := contentStore.GetNodeContent(ctx, nodes[0].ID)
	require.NoError(t, err)
	require.Len(t, arraysContent.Skills, 1)
	require.Equal(t, "arrays", arraysContent.Skills[0].SkillSlug)
	require.True(t, arraysContent.Skills[0].IsPrimary)
	require.Len(t, arraysContent.Pools, 1)
	require.Equal(t, "foundation-arrays-beginner", arraysContent.Pools[0].PoolSlug)

	questionPool, err := contentStore.GetPoolBySlug(ctx, "foundation-arrays-beginner")
	require.NoError(t, err)
	require.NotNil(t, questionPool)

	questionIDs, err := contentStore.ListPoolQuestionIDs(ctx, questionPool.ID)
	require.NoError(t, err)
	require.Len(t, questionIDs, 1)
	require.Equal(t, "b0000000-0000-4000-8000-000000000001", questionIDs[0])

	bossContent, err := contentStore.GetNodeContent(ctx, nodes[4].ID)
	require.NoError(t, err)
	require.Len(t, bossContent.Skills, 3)
	require.Len(t, bossContent.Pools, 3)
}

func TestContentStoreListPoolsBySkillID(t *testing.T) {
	pool, _ := testdb.Start(t)
	testdb.Migrate(t, pool)

	contentStore := store.NewContentStore(pool)
	skillStore := store.NewSkillStore(pool)
	ctx := context.Background()

	skill, err := skillStore.GetSkillBySlug(ctx, "arrays")
	require.NoError(t, err)
	require.NotNil(t, skill)

	pools, err := contentStore.ListPoolsBySkillID(ctx, skill.ID)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(pools), 1)
	require.Equal(t, "arrays", pools[0].SkillSlug)
}
