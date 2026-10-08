package kafka

import (
	"context"
	"time"
)

// retryPolicy bounds how a consumer retries a failing handler before it skips
// the message: capped exponential backoff, a fixed number of attempts.
type retryPolicy struct {
	attempts  int
	baseDelay time.Duration
	maxDelay  time.Duration
}

// defaultRetryPolicy retries for roughly ten seconds in total (200ms, 400ms,
// 800ms, 1.6s between five attempts; capped at 5s).
var defaultRetryPolicy = retryPolicy{attempts: 5, baseDelay: 200 * time.Millisecond, maxDelay: 5 * time.Second}

// delay is the wait before retry number n (n >= 1).
func (p retryPolicy) delay(n int) time.Duration {
	d := p.baseDelay
	for i := 1; i < n; i++ {
		d *= 2
		if d >= p.maxDelay {
			return p.maxDelay
		}
	}
	if d > p.maxDelay {
		return p.maxDelay
	}
	return d
}

// handleWithRetry runs handler until it succeeds or the attempts run out, and
// returns the last handler error (nil on success). It stops early with the
// context's error when ctx is cancelled.
func handleWithRetry(ctx context.Context, p retryPolicy, handler MessageHandler, key, value []byte) error {
	var err error
	for attempt := 1; attempt <= p.attempts; attempt++ {
		if err = handler(ctx, key, value); err == nil {
			return nil
		}
		if ctx.Err() != nil {
			return ctx.Err()
		}
		if attempt == p.attempts {
			break
		}
		timer := time.NewTimer(p.delay(attempt))
		select {
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()
		case <-timer.C:
		}
	}
	return err
}
