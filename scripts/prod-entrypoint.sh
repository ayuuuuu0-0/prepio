#!/usr/bin/env bash
set -euo pipefail

echo "=== Starting Prepio Production Backend Services ==="

# Wait for PostgreSQL to become ready
echo "Waiting for PostgreSQL database to be ready..."
DB_HOST="${DB_HOST:-postgres}"
DB_PORT="${DB_PORT:-5432}"
DB_USER="${POSTGRES_USER:-prepio}"

until pg_isready -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER"; do
  echo "PostgreSQL is not ready yet - retrying in 1s..."
  sleep 1
done
echo "PostgreSQL is ready and accepting connections."

# Execute database migrations
echo "Executing database migrations..."
migrate -path /app/migrations -database "$DATABASE_URL" up || {
  echo "Migration failed or already up to date."
}

# Internal service communication configuration
export USER_SERVICE_URL="${USER_SERVICE_URL:-http://localhost:8081}"
export QUESTION_SERVICE_URL="${QUESTION_SERVICE_URL:-http://localhost:8082}"
export STREAK_SERVICE_URL="${STREAK_SERVICE_URL:-http://localhost:8083}"
export PROGRESS_SERVICE_URL="${PROGRESS_SERVICE_URL:-http://localhost:8084}"
export NOTIFICATION_SERVICE_URL="${NOTIFICATION_SERVICE_URL:-http://localhost:8085}"
export DEV_SYNC_EVENTS="${DEV_SYNC_EVENTS:-true}"
export GATEWAY_PORT="${GATEWAY_PORT:-8080}"

mkdir -p /tmp/prepio_logs

PIDS=()

cleanup() {
  echo "Received shutdown signal. Stopping all services gracefully..."
  for pid in "${PIDS[@]}"; do
    if kill -0 "$pid" 2>/dev/null; then
      kill "$pid" 2>/dev/null || true
    fi
  done
  wait
  echo "All services stopped."
  exit 0
}

trap cleanup SIGINT SIGTERM

echo "Starting user service (:8081)..."
/app/user > /tmp/prepio_logs/user.log 2>&1 &
PIDS+=($!)

echo "Starting question service (:8082)..."
/app/question > /tmp/prepio_logs/question.log 2>&1 &
PIDS+=($!)

echo "Starting streak service (:8083)..."
/app/streak > /tmp/prepio_logs/streak.log 2>&1 &
PIDS+=($!)

echo "Starting progress service (:8084)..."
/app/progress > /tmp/prepio_logs/progress.log 2>&1 &
PIDS+=($!)

echo "Starting notification service (:8085)..."
/app/notification > /tmp/prepio_logs/notification.log 2>&1 &
PIDS+=($!)

# Brief pause to ensure internal microservices bind to ports
sleep 2

echo "Starting gateway on port $GATEWAY_PORT..."
/app/gateway &
PIDS+=($!)

# Wait for the gateway or any critical service to exit
wait -n "${PIDS[@]}"
