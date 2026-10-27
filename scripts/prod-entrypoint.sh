#!/usr/bin/env bash
set -euo pipefail

echo "=== Starting Prepio Modular Monolith Backend ==="

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
# A failed migration must stop the deploy (migrate exits 0 when there is nothing to apply).
migrate -path /app/migrations -database "$DATABASE_URL" up

# Load authored worlds and lessons (idempotent). Invalid content must stop the deploy.
echo "Syncing lesson content..."
CONTENT_DIR=/app/content /app/content-sync

echo "Starting unified Prepio modular monolith on port ${GATEWAY_PORT:-8080}..."
exec /app/server
