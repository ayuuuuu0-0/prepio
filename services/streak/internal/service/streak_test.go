package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prepio/prepio/config"
	"github.com/prepio/prepio/constants"
	"github.com/prepio/prepio/services/streak/internal/service"
	"github.com/prepio/prepio/services/streak/internal/store"
	"github.com/prepio/prepio/test/factories"
	"github.com/prepio/prepio/test/fakes"
	"github.com/prepio/prepio/test/testdb"
	"github.com/prepio/prepio/test/testredis"
	"github.com/stretchr/testify/require"
)

type fakeGems struct {
	balance int
}

func (f *fakeGems) DeductGems(ctx context.Context, userID string, amount int, reason string) error {
	if f.balance < amount {
		return service.ErrInsufficientGems
	}
	f.balance -= amount
	return nil
}

type testEnv struct {
	svc      *service.StreakService
	publisher *fakes.KafkaProducer
	userID   string
	gems     *fakeGems
}

func setupStreakTest(t *testing.T) *testEnv {
	t.Helper()
	pool, _ := testdb.Start(t)
	testdb.Migrate(t, pool)
	redisClient, _ := testredis.New(t)
	publisher := &fakes.KafkaProducer{}
	gems := &fakeGems{balance: config.StreakFreezeGemCost * 10}

	svc := service.NewStreakService(
		store.NewStreakStore(pool),
		store.NewFreezeStore(pool),
		store.NewCacheStore(redisClient),
		publisher,
		gems,
	)

	return &testEnv{
		svc:       svc,
		publisher: publisher,
		userID:    seedStreakUser(t, pool),
		gems:      gems,
	}
}

func TestFirstQuestionOfDayIncrementsStreak(t *testing.T) {
	env := setupStreakTest(t)
	ctx := context.Background()
	loc, _ := time.LoadLocation(constants.DefaultTimezone)
	submittedAt := time.Date(2026, 6, 9, 10, 0, 0, 0, loc)

	require.NoError(t, env.svc.ProcessQuestionAnswered(ctx, factories.QuestionAnsweredEvent(env.userID, submittedAt)))

	resp, err := env.svc.GetMe(ctx, env.userID, constants.DefaultTimezone)
	require.NoError(t, err)
	require.Equal(t, 1, resp.CurrentStreak)
	require.NotNil(t, env.publisher.Last())
}

func TestSecondQuestionSameDayDoesNotDoubleIncrement(t *testing.T) {
	env := setupStreakTest(t)
	ctx := context.Background()
	loc, _ := time.LoadLocation(constants.DefaultTimezone)
	day := time.Date(2026, 6, 9, 10, 0, 0, 0, loc)

	require.NoError(t, env.svc.ProcessQuestionAnswered(ctx, factories.QuestionAnsweredEvent(env.userID, day)))
	require.NoError(t, env.svc.ProcessQuestionAnswered(ctx, factories.QuestionAnsweredEvent(env.userID, day.Add(2*time.Hour))))

	resp, err := env.svc.GetMe(ctx, env.userID, constants.DefaultTimezone)
	require.NoError(t, err)
	require.Equal(t, 1, resp.CurrentStreak)
}

func TestMissedDayWithoutFreezeBreaksStreak(t *testing.T) {
	env := setupStreakTest(t)
	ctx := context.Background()
	loc, _ := time.LoadLocation(constants.DefaultTimezone)

	day1 := time.Date(2026, 6, 9, 10, 0, 0, 0, loc)
	day3 := time.Date(2026, 6, 11, 10, 0, 0, 0, loc)

	require.NoError(t, env.svc.ProcessQuestionAnswered(ctx, factories.QuestionAnsweredEvent(env.userID, day1)))
	require.NoError(t, env.svc.ProcessQuestionAnswered(ctx, factories.QuestionAnsweredEvent(env.userID, day3)))

	resp, err := env.svc.GetMe(ctx, env.userID, constants.DefaultTimezone)
	require.NoError(t, err)
	require.Equal(t, 1, resp.CurrentStreak)
	require.Equal(t, 1, resp.LongestStreak)
}

func TestMissedDayWithFreezeConsumesFreezeAndHoldsStreak(t *testing.T) {
	env := setupStreakTest(t)
	ctx := context.Background()
	loc, _ := time.LoadLocation(constants.DefaultTimezone)

	_, err := env.svc.PurchaseFreeze(ctx, env.userID)
	require.NoError(t, err)

	day1 := time.Date(2026, 6, 9, 10, 0, 0, 0, loc)
	day3 := time.Date(2026, 6, 11, 10, 0, 0, 0, loc)

	require.NoError(t, env.svc.ProcessQuestionAnswered(ctx, factories.QuestionAnsweredEvent(env.userID, day1)))
	require.NoError(t, env.svc.ProcessQuestionAnswered(ctx, factories.QuestionAnsweredEvent(env.userID, day3)))

	resp, err := env.svc.GetMe(ctx, env.userID, constants.DefaultTimezone)
	require.NoError(t, err)
	require.Equal(t, 2, resp.CurrentStreak)
	require.Equal(t, 0, resp.FreezeCount)
}

func TestFreezePurchaseDeductsGems(t *testing.T) {
	env := setupStreakTest(t)
	env.gems.balance = config.StreakFreezeGemCost
	ctx := context.Background()

	resp, err := env.svc.PurchaseFreeze(ctx, env.userID)
	require.NoError(t, err)
	require.Equal(t, 1, resp.FreezeCount)
	require.Equal(t, config.StreakFreezeGemCost, resp.GemsSpent)
	require.Equal(t, 0, env.gems.balance)
}

func TestFreezePurchaseFailsWithInsufficientGems(t *testing.T) {
	env := setupStreakTest(t)
	env.gems.balance = config.StreakFreezeGemCost - 1
	ctx := context.Background()

	_, err := env.svc.PurchaseFreeze(ctx, env.userID)
	require.Error(t, err)
	require.True(t, errors.Is(err, service.ErrInsufficientGems))
}

func TestCannotHoldMoreThanMaxStreakFreezes(t *testing.T) {
	env := setupStreakTest(t)
	ctx := context.Background()

	for i := 0; i < config.MaxStreakFreezes; i++ {
		_, err := env.svc.PurchaseFreeze(ctx, env.userID)
		require.NoError(t, err)
	}

	_, err := env.svc.PurchaseFreeze(ctx, env.userID)
	require.Error(t, err)
	require.True(t, errors.Is(err, service.ErrFreezeMaxHeld))
}

func seedStreakUser(t *testing.T, pool *pgxpool.Pool) string {
	t.Helper()
	ctx := context.Background()
	var userID string
	err := pool.QueryRow(ctx, `
		INSERT INTO users (email, username, password_hash, timezone)
		VALUES ('streak@test.com', 'streakuser', 'hash', $1)
		RETURNING id`, constants.DefaultTimezone).Scan(&userID)
	require.NoError(t, err)
	return userID
}
