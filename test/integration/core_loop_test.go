package integration_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/prepio/prepio/config"
	notiftest "github.com/prepio/prepio/services/notification/testing"
	progresstest "github.com/prepio/prepio/services/progress/testing"
	questiontest "github.com/prepio/prepio/services/question/testing"
	streaktest "github.com/prepio/prepio/services/streak/testing"
	"github.com/prepio/prepio/shared/events"
	"github.com/prepio/prepio/test/fakes"
	"github.com/prepio/prepio/test/testdb"
	"github.com/prepio/prepio/test/testredis"
	"github.com/stretchr/testify/require"
)

func TestCoreLoop(t *testing.T) {
	pool, _ := testdb.Start(t)
	testdb.Migrate(t, pool)
	redisClient, _ := testredis.New(t)
	ctx := context.Background()
	publisher := &fakes.KafkaProducer{}

	questionService := questiontest.NewService(pool, redisClient, publisher)
	streakService := streaktest.NewService(pool, redisClient, publisher, &fakeGems{balance: 1000})
	progressService := progresstest.NewService(pool, publisher)
	notificationService := notiftest.NewService(pool, redisClient)

	var userID string
	require.NoError(t, pool.QueryRow(ctx, `
		INSERT INTO users (email, username, password_hash, timezone)
		VALUES ('loop@test.com', 'loopuser', 'hash', 'Asia/Kolkata')
		RETURNING id`).Scan(&userID))

	_, err := pool.Exec(ctx, `
		INSERT INTO user_progress (user_id, total_xp, current_level, gem_balance)
		VALUES ($1, 3750, 9, 100)`, userID)
	require.NoError(t, err)

	daily, err := questionService.GetDailyPaper(ctx, userID, "Asia/Kolkata")
	require.NoError(t, err)
	require.NotEmpty(t, daily.Questions)

	submitResp, err := questionService.SubmitAnswer(ctx, userID, daily.Questions[0].ID, questiontest.SubmitRequest{
		SessionID:        daily.SessionID,
		Answer:           "hash map approach with O(n) time and O(n) space complexity",
		TimeSpentSeconds: 120,
	})
	require.NoError(t, err)
	require.True(t, submitResp.Correct)
	require.Equal(t, events.TopicQuestionAnswered, publisher.Messages[0].Topic)

	var answered events.QuestionAnswered
	require.NoError(t, json.Unmarshal(publisher.Messages[0].Payload, &answered))

	require.NoError(t, streakService.ProcessQuestionAnswered(ctx, answered))
	require.NoError(t, progressService.ProcessQuestionAnswered(ctx, answered))

	streak, err := streakService.GetMe(ctx, userID, "Asia/Kolkata")
	require.NoError(t, err)
	require.Equal(t, 1, streak.CurrentStreak)

	progress, err := progressService.GetMe(ctx, userID)
	require.NoError(t, err)
	require.Equal(t, config.XPByDifficulty["medium"], progress.TotalXP)

	var streakEvent events.StreakUpdated
	require.NoError(t, json.Unmarshal(publisher.Messages[1].Payload, &streakEvent))
	require.NoError(t, progressService.ProcessStreakUpdated(ctx, streakEvent))

	var progressEvent events.ProgressUpdated
	require.NoError(t, json.Unmarshal(publisher.Messages[2].Payload, &progressEvent))
	require.True(t, progressEvent.LeveledUp)
	require.NoError(t, notificationService.HandleProgressUpdated(ctx, progressEvent))

	var notifCount int
	require.NoError(t, pool.QueryRow(ctx, `
		SELECT COUNT(*) FROM notification_log WHERE user_id = $1`, userID).Scan(&notifCount))
	require.Equal(t, 1, notifCount)
}

type fakeGems struct{ balance int }

func (f *fakeGems) DeductGems(ctx context.Context, userID string, amount int, reason string) error {
	if f.balance < amount {
		return errInsufficientGems
	}
	f.balance -= amount
	return nil
}

var errInsufficientGems = &gemError{}

type gemError struct{}

func (e *gemError) Error() string { return "insufficient gems" }
