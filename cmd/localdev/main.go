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

	// 1. Embedded Postgres
	port := uint32(5433)
	dataDir := filepath.Join(os.TempDir(), "prepio-local-pg")
	_ = os.RemoveAll(dataDir)

	db := embeddedpostgres.NewDatabase(embeddedpostgres.DefaultConfig().
		RuntimePath(dataDir).
		Port(port).
		Database("prepio").
		Username("prepio").
		Password("prepio"))

	log.Printf("[1/4] Starting embedded PostgreSQL on port %d...", port)
	if err := db.Start(); err != nil {
		log.Fatalf("Failed to start embedded postgres: %v", err)
	}
	defer func() {
		log.Println("Stopping PostgreSQL...")
		_ = db.Stop()
		_ = os.RemoveAll(dataDir)
	}()

	dsn := fmt.Sprintf("postgres://prepio:prepio@localhost:%d/prepio?sslmode=disable", port)

	// 2. Connect and run migrations
	log.Println("[2/4] Connecting and applying migrations...")
	pool, err := connectWithRetry(dsn, 30)
	if err != nil {
		log.Fatalf("Failed to connect to postgres: %v", err)
	}

	if err := applyMigrations(pool, "migrations"); err != nil {
		log.Fatalf("Failed to apply migrations: %v", err)
	}
	pool.Close()
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

	// 4. Start backend microservices
	log.Println("[4/4] Starting backend services...")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	baseEnv := append(os.Environ(),
		"DATABASE_URL="+dsn,
		"REDIS_ADDR=localhost:6379",
		"JWT_SECRET=dev-secret-change-in-production-use-32-chars",
		"DEV_SYNC_EVENTS=true",
		"USER_SERVICE_URL=http://localhost:8081",
		"QUESTION_SERVICE_URL=http://localhost:8082",
		"STREAK_SERVICE_URL=http://localhost:8083",
		"PROGRESS_SERVICE_URL=http://localhost:8084",
		"GATEWAY_PORT=8080",
	)

	services := []struct {
		name string
		path string
	}{
		{"user", "./services/user/cmd"},
		{"question", "./services/question/cmd"},
		{"streak", "./services/streak/cmd"},
		{"progress", "./services/progress/cmd"},
		{"gateway", "./services/gateway/cmd"},
	}

	var cmds []*exec.Cmd
	for _, s := range services {
		cmd := exec.CommandContext(ctx, "go", "run", s.path)
		cmd.Env = baseEnv
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Start(); err != nil {
			log.Fatalf("Failed to start %s service: %v", s.name, err)
		}
		cmds = append(cmds, cmd)
		log.Printf("  Started %s service (PID %d)", s.name, cmd.Process.Pid)
		time.Sleep(1 * time.Second)
	}

	log.Println("\n=======================================================")
	log.Println(" PREPIO LOCAL BACKEND READY AT http://localhost:8080")
	log.Println(" API Gateway: http://localhost:8080/api/v1")
	log.Println(" Press Ctrl+C to stop all services.")
	log.Println("=======================================================\n")

	// Wait for gateway to be responsive
	for i := 0; i < 15; i++ {
		resp, err := http.Get("http://localhost:8080/api/v1/companions")
		if err == nil && resp.StatusCode == http.StatusOK {
			resp.Body.Close()
			log.Println("✓ Gateway healthcheck passed: http://localhost:8080 is answering!")
			break
		}
		time.Sleep(1 * time.Second)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Println("\nShutting down all backend services...")
	cancel()
	for _, cmd := range cmds {
		if cmd.Process != nil {
			_ = cmd.Process.Kill()
		}
	}
	log.Println("All services stopped.")
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
