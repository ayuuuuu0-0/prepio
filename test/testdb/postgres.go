package testdb

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Start boots an embedded PostgreSQL instance and returns a connection pool.
func Start(t *testing.T) (*pgxpool.Pool, string) {
	t.Helper()

	port := uint32(15432 + (time.Now().UnixNano() % 1000))
	db := embeddedpostgres.NewDatabase(embeddedpostgres.DefaultConfig().
		Port(port).
		Database("prepio").
		Username("prepio").
		Password("prepio"))

	if err := db.Start(); err != nil {
		t.Skipf("embedded postgres unavailable: %v", err)
	}
	t.Cleanup(func() { _ = db.Stop() })

	dsn := fmt.Sprintf("postgres://prepio:prepio@localhost:%d/prepio?sslmode=disable", port)
	pool := connect(t, dsn)
	return pool, dsn
}

// Migrate runs all .up.sql files from the migrations directory in order.
func Migrate(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()

	root, err := findRepoRoot()
	if err != nil {
		t.Fatalf("find repo root: %v", err)
	}

	dir := filepath.Join(root, "migrations")
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read migrations: %v", err)
	}

	ctx := context.Background()
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || filepath.Ext(name) != ".sql" {
			continue
		}
		if len(name) < 7 || name[len(name)-7:] != ".up.sql" {
			continue
		}
		sqlBytes, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			t.Fatalf("read migration %s: %v", name, err)
		}
		if _, err := pool.Exec(ctx, string(sqlBytes)); err != nil {
			t.Fatalf("apply migration %s: %v", name, err)
		}
	}
}

func connect(t *testing.T, dsn string) *pgxpool.Pool {
	t.Helper()
	ctx := context.Background()
	var pool *pgxpool.Pool
	var err error
	for i := 0; i < 30; i++ {
		pool, err = pgxpool.New(ctx, dsn)
		if err == nil && pool.Ping(ctx) == nil {
			return pool
		}
		if pool != nil {
			pool.Close()
		}
		time.Sleep(200 * time.Millisecond)
	}
	t.Fatalf("connect postgres: %v", err)
	return nil
}

func findRepoRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("go.mod not found")
		}
		dir = parent
	}
}
