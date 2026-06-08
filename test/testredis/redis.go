package testredis

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis/v2"
	redisclient "github.com/prepio/prepio/shared/redis"
	"github.com/redis/go-redis/v9"
)

// New starts an in-memory Redis server for tests.
func New(t *testing.T) (*redis.Client, *miniredis.Miniredis) {
	t.Helper()

	srv, err := miniredis.Run()
	if err != nil {
		t.Fatalf("start miniredis: %v", err)
	}
	t.Cleanup(srv.Close)

	client, err := redisclient.New(context.Background(), redisclient.Config{Addr: srv.Addr()})
	if err != nil {
		t.Fatalf("redis client: %v", err)
	}
	t.Cleanup(func() { _ = client.Close() })

	return client, srv
}
