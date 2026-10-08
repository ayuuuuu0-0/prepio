package kafka

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRetryPolicyDelayIsCappedExponential(t *testing.T) {
	p := retryPolicy{attempts: 10, baseDelay: 200 * time.Millisecond, maxDelay: 5 * time.Second}
	want := []time.Duration{200, 400, 800, 1600, 3200, 5000, 5000}
	for i, w := range want {
		require.Equal(t, w*time.Millisecond, p.delay(i+1), "retry %d", i+1)
	}
}

func TestHandleWithRetry(t *testing.T) {
	fast := retryPolicy{attempts: 5, baseDelay: time.Millisecond, maxDelay: 2 * time.Millisecond}
	boom := errors.New("boom")

	counting := func(failures int, calls *int) MessageHandler {
		return func(ctx context.Context, key, value []byte) error {
			*calls++
			if *calls <= failures {
				return boom
			}
			return nil
		}
	}

	t.Run("success needs one call", func(t *testing.T) {
		calls := 0
		require.NoError(t, handleWithRetry(context.Background(), fast, counting(0, &calls), nil, nil))
		require.Equal(t, 1, calls)
	})

	t.Run("transient failure is retried until it succeeds", func(t *testing.T) {
		calls := 0
		require.NoError(t, handleWithRetry(context.Background(), fast, counting(3, &calls), nil, nil))
		require.Equal(t, 4, calls)
	})

	t.Run("persistent failure gives up after the attempts with the last error", func(t *testing.T) {
		calls := 0
		err := handleWithRetry(context.Background(), fast, counting(100, &calls), nil, nil)
		require.ErrorIs(t, err, boom)
		require.Equal(t, 5, calls)
	})

	t.Run("cancellation stops retrying", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		calls := 0
		handler := func(context.Context, []byte, []byte) error {
			calls++
			cancel()
			return boom
		}
		slow := retryPolicy{attempts: 5, baseDelay: time.Hour, maxDelay: time.Hour}
		err := handleWithRetry(ctx, slow, handler, nil, nil)
		require.ErrorIs(t, err, context.Canceled)
		require.Equal(t, 1, calls)
	})
}
