package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	log.Println("=== Starting Prepio Local Development Environment ===")

	// 1. Start embedded PostgreSQL on port 5433
	log.Println("[1/4] Starting embedded PostgreSQL on port 5433...")
	pgConfig := embeddedpostgres.DefaultConfig().
		Port(5433).
		Database("prepio").
		Username("prepio").
		Password("prepio")

	pg := embeddedpostgres.NewDatabase(pgConfig)
	if err := pg.Start(); err != nil {
		log.Printf("Note on PostgreSQL start: %v (checking if already running...)", err)
	}
	defer func() {
		log.Println("Stopping embedded PostgreSQL...")
		_ = pg.Stop()
	}()

	dsn := "postgres://prepio:prepio@localhost:5433/prepio?sslmode=disable"

	// 2. Apply migrations
	log.Println("[2/4] Connecting and applying migrations...")
	pool, err := connectWithRetry(dsn, 10)
	if err != nil {
		log.Fatalf("Failed to connect to local PostgreSQL: %v", err)
	}
	defer pool.Close()

	if err := applyMigrations(pool, "./migrations"); err != nil {
		log.Fatalf("Failed to apply migrations: %v", err)
	}
	log.Println("  Migrations applied successfully!")

	// 3. Sync lesson content
	log.Println("[3/4] Syncing lesson content...")
	contentCmd := exec.Command("go", "run", "./services/question/cmd/content-sync")
	contentCmd.Env = append(os.Environ(), "DATABASE_URL="+dsn)
	contentCmd.Stdout = os.Stdout
	contentCmd.Stderr = os.Stderr
	if err := contentCmd.Run(); err != nil {
		log.Printf("Warning: content-sync had notices: %v", err)
	}

	// 4. Start unified modular monolith server
	log.Println("[4/4] Starting Prepio modular monolith server...")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	baseEnv := append(os.Environ(),
		"DATABASE_URL="+dsn,
		"REDIS_ADDR=localhost:6379",
		"JWT_SECRET=dev-secret-change-in-production-use-32-chars",
		"DEV_SYNC_EVENTS=true",
		"GATEWAY_PORT=8080",
	)

	cmd := exec.CommandContext(ctx, "go", "run", "./cmd/server")
	cmd.Env = baseEnv
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	log.Printf("  Started modular monolith server (PID %d)", cmd.Process.Pid)

	log.Println("\n=======================================================")
	log.Println(" PREPIO MODULAR MONOLITH READY AT http://localhost:8080")
	log.Println(" API Gateway: http://localhost:8080/api/v1")
	log.Println(" Press Ctrl+C to stop.")
	log.Println("=======================================================\n")

	// Wait for server to be responsive
	for i := 0; i < 15; i++ {
		resp, err := http.Get("http://localhost:8080/api/v1/companions")
		if err == nil && resp.StatusCode == http.StatusOK {
			resp.Body.Close()
			log.Println("✓ Monolith healthcheck passed: http://localhost:8080 is answering!")
			break
		}
		time.Sleep(1 * time.Second)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Println("\nShutting down modular monolith...")
	cancel()
	if cmd.Process != nil {
		_ = cmd.Process.Kill()
	}
	log.Println("Server stopped.")
}

func connectWithRetry(dsn string, retries int) (*pgxpool.Pool, error) {
	ctx := context.Background()
	var pool *pgxpool.Pool
	var err error
	for i := 0; i < retries; i++ {
		pool, err = pgxpool.New(ctx, dsn)
		if err == nil && pool.Ping(ctx) == nil {
			return pool, nil
		}
		if pool != nil {
			pool.Close()
		}
		time.Sleep(300 * time.Millisecond)
	}
	return nil, err
}

func applyMigrations(pool *pgxpool.Pool, dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	ctx := context.Background()
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || !strings.HasSuffix(name, ".up.sql") {
			continue
		}
		sqlBytes, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			return fmt.Errorf("read %s: %w", name, err)
		}
		if _, err := pool.Exec(ctx, string(sqlBytes)); err != nil {
			return fmt.Errorf("execute %s: %w", name, err)
		}
	}
	return nil
}
