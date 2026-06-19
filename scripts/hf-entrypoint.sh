#!/usr/bin/env bash
set -euo pipefail

echo "=== Starting Prepio Hugging Face Setup ==="

# Define directories using /tmp which is guaranteed to be writable
PGDATA="/tmp/postgres_data"
PGSOCKETS="/tmp"
REDISDATA="/tmp"

# --- 1. Start Redis ---
echo "Starting Redis..."
redis-server --port 6379 --dir "$REDISDATA" --daemonize yes

# --- 2. Start PostgreSQL ---
echo "Initializing Postgres DB..."
if [ ! -d "$PGDATA" ]; then
  initdb -D "$PGDATA" --auth-host=trust --auth-local=trust
fi

echo "Starting Postgres..."
# -k specifies the socket directory. -F disables fsync for faster startup/performance in transient container environments.
pg_ctl -D "$PGDATA" -o "-F -p 5432 -k $PGSOCKETS" start

echo "Waiting for Postgres to accept connections..."
until pg_isready -h localhost -p 5432; do
  sleep 1
done

# Create the user and database if they do not exist
echo "Setting up Postgres database and user..."
psql -h localhost -p 5432 -U "$(whoami)" -d postgres -c "CREATE ROLE prepio WITH LOGIN PASSWORD 'prepio' SUPERUSER;" || true
psql -h localhost -p 5432 -U "$(whoami)" -d postgres -c "CREATE DATABASE prepio OWNER prepio;" || true

# --- 3. Run Migrations ---
echo "Running database migrations..."
migrate -path /app/migrations -database "postgres://prepio:prepio@localhost:5432/prepio?sslmode=disable" up

# --- 4. Start Microservices ---
export DATABASE_URL="postgres://prepio:prepio@localhost:5432/prepio?sslmode=disable"
export REDIS_ADDR="localhost:6379"
export DEV_SYNC_EVENTS="true"
export JWT_SECRET="${JWT_SECRET:-huggingface-portfolio-prepio-secret-key-32-chars-long}"

export USER_SERVICE_URL="http://localhost:8081"
export QUESTION_SERVICE_URL="http://localhost:8082"
export STREAK_SERVICE_URL="http://localhost:8083"
export PROGRESS_SERVICE_URL="http://localhost:8084"
export NOTIFICATION_SERVICE_URL="http://localhost:8085"

mkdir -p /tmp/prepio_logs

echo "Starting user service..."
/app/user > /tmp/prepio_logs/user.log 2>&1 &

echo "Starting question service..."
/app/question > /tmp/prepio_logs/question.log 2>&1 &

echo "Starting streak service..."
/app/streak > /tmp/prepio_logs/streak.log 2>&1 &

echo "Starting progress service..."
/app/progress > /tmp/prepio_logs/progress.log 2>&1 &

echo "Starting notification service..."
/app/notification > /tmp/prepio_logs/notification.log 2>&1 &

# Wait for internal services to spin up
sleep 3

# --- 5. Start Gateway in the foreground ---
# Hugging Face routes traffic to the port specified in GATEWAY_PORT (7860)
export GATEWAY_PORT="7860"
echo "Starting gateway service on port $GATEWAY_PORT..."
exec /app/gateway
