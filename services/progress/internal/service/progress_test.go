package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/prepio/prepio/config"
	"github.com/prepio/prepio/services/progress/internal/service"
	"github.com/prepio/prepio/services/progress/internal/store"
	"github.com/prepio/prepio/shared/events"
	"github.com/prepio/prepio/test/factories"
	"github.com/prepio/prepio/test/fakes"
	"github.com/prepio/prepio/test/testdb"
	"github.com/stretchr/testify/require"
)

func setupProgress(t *testing.T) (*service.ProgressService, *fakes.KafkaProducer, string) {
	t.Helper()
	pool, _ := testdb.Start(t)
	testdb.Migrate(t, pool)
	publisher := &fakes.KafkaProducer{}
	svc := service.NewProgressService(store.NewProgressStore(pool), store.NewLedgerStore(pool), publisher)

	ctx := context.Background()
	var userID string
	require.NoError(t, pool.QueryRow(ctx, `
		INSERT INTO users (email, username, password_hash)
		VALUES ('p@test.com', 'puser', 'hash') RETURNING id`).Scan(&userID))

	_, err := pool.Exec(ctx, `
		INSERT INTO user_progress (user_id, total_xp, current_level, gem_balance)
		VALUES ($1, 0, 1, 500)`, userID)
	require.NoError(t, err)

	return svc, publisher, userID
}

func TestQuestionAnsweredAwardsXPAndGems(t *testing.T) {
	svc, publisher, userID := setupProgress(t)
	ctx := context.Background()

	event := factories.QuestionAnsweredEvent(userID, time.Now())
	event.Correct = true
	event.Difficulty = "medium"
	event.CompanyTags = nil
	require.NoError(t, svc.ProcessQuestionAnswered(ctx, event))

	resp, err := svc.GetMe(ctx, userID)
	require.NoError(t, err)
	require.Equal(t, config.XPByDifficulty["medium"], resp.TotalXP)
	require.Equal(t, 500+config.GemsByDifficulty["medium"], resp.GemBalance)
	require.NotNil(t, publisher.Last())
}

func TestGemDeductionOnStreakFreezePurchase(t *testing.T) {
	svc, _, userID := setupProgress(t)
	ctx := context.Background()

	balance, err := svc.DeductGems(ctx, userID, config.StreakFreezeGemCost, "streak_freeze_purchase")
	require.NoError(t, err)
	require.Equal(t, 500-config.StreakFreezeGemCost, balance)
}

func TestGemDeductionFailsWithInsufficientBalance(t *testing.T) {
	svc, _, userID := setupProgress(t)
	ctx := context.Background()

	_, err := svc.DeductGems(ctx, userID, 501, "streak_freeze_purchase")
	require.Error(t, err)
	require.True(t, errors.Is(err, service.ErrInsufficientGems))
}

func TestStreakIncrementAwardsBonusGems(t *testing.T) {
	svc, _, userID := setupProgress(t)
	ctx := context.Background()

	event := events.StreakUpdated{
		EventID: uuid.NewString(), UserID: userID,
		PreviousStreak: 1, CurrentStreak: 2,
	}
	require.NoError(t, svc.ProcessStreakUpdated(ctx, event))

	resp, err := svc.GetMe(ctx, userID)
	require.NoError(t, err)
	require.Equal(t, 500+config.StreakIncrementGemBonus, resp.GemBalance)
}
