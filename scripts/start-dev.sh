#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

if [ ! -f .env ]; then
  cp .env.example .env
fi

set -a
source .env
set +a

export JWT_SECRET="${JWT_SECRET:-dev-secret-change-in-production-use-32-chars}"
export DATABASE_URL="${DATABASE_URL:-postgres://prepio:prepio@localhost:5432/prepio?sslmode=disable}"
export REDIS_ADDR="${REDIS_ADDR:-localhost:6379}"
export DEV_SYNC_EVENTS="${DEV_SYNC_EVENTS:-true}"

echo "Starting infrastructure (postgres + redis)..."
docker compose up -d postgres redis

echo "Waiting for postgres..."
for i in $(seq 1 30); do
  if docker compose exec -T postgres pg_isready -U prepio >/dev/null 2>&1; then
    break
  fi
  sleep 2
done

echo "Applying migrations if needed..."
if command -v migrate >/dev/null 2>&1; then
  migrate -path migrations -database "$DATABASE_URL" up 2>/dev/null || true
else
  applied=$(docker compose exec -T postgres psql -U prepio -d prepio -tAc "SELECT COUNT(*) FROM schema_migrations" 2>/dev/null || echo "0")
  if [ "${applied:-0}" = "0" ]; then
    for f in migrations/*.up.sql; do
      echo "  applying $(basename "$f")"
      docker compose exec -T postgres psql -U prepio -d prepio -f - < "$f" >/dev/null 2>&1 || true
    done
  fi
fi

mkdir -p "$ROOT/.run"
PIDS_FILE="$ROOT/.run/pids"

stop_services() {
  if [ -f "$PIDS_FILE" ]; then
    while read -r pid; do
      kill "$pid" 2>/dev/null || true
    done < "$PIDS_FILE"
    rm -f "$PIDS_FILE"
  fi
}
trap stop_services EXIT

stop_services

start_service() {
  local name="$1"
  local path="$2"
  echo "Starting $name..."
  env JWT_SECRET="$JWT_SECRET" DATABASE_URL="$DATABASE_URL" REDIS_ADDR="$REDIS_ADDR" \
    DEV_SYNC_EVENTS="$DEV_SYNC_EVENTS" KAFKA_BROKERS="$KAFKA_BROKERS" \
    USER_SERVICE_URL="${USER_SERVICE_URL:-http://localhost:8081}" \
    QUESTION_SERVICE_URL="${QUESTION_SERVICE_URL:-http://localhost:8082}" \
    STREAK_SERVICE_URL="${STREAK_SERVICE_URL:-http://localhost:8083}" \
    PROGRESS_SERVICE_URL="${PROGRESS_SERVICE_URL:-http://localhost:8084}" \
    NOTIFICATION_SERVICE_URL="${NOTIFICATION_SERVICE_URL:-http://localhost:8085}" \
    go run "$path" >"$ROOT/.run/$name.log" 2>&1 &
  echo $! >> "$PIDS_FILE"
  sleep 1
}

start_service user         ./services/user/cmd
start_service question     ./services/question/cmd
start_service streak       ./services/streak/cmd
start_service progress     ./services/progress/cmd
start_service notification ./services/notification/cmd
start_service gateway      ./services/gateway/cmd

echo ""
echo "Backend running (DEV_SYNC_EVENTS=$DEV_SYNC_EVENTS):"
echo "  Gateway:      http://localhost:8080"
echo "  User:         http://localhost:8081"
echo "  Question:     http://localhost:8082"
echo "  Streak:       http://localhost:8083"
echo "  Progress:     http://localhost:8084"
echo "  Notification: http://localhost:8085"
echo ""
echo "Logs in .run/*.log"
echo "Press Ctrl+C to stop all services"

wait
