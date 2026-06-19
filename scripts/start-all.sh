#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

# Cleanup function to stop backend services and databases when exiting
cleanup() {
  echo ""
  echo "Shutting down all services..."
  if [ -n "${BACKEND_PID:-}" ]; then
    # Killing start-dev.sh will trigger its exit trap to kill all Go services
    kill "$BACKEND_PID" 2>/dev/null || true
  fi
}
trap cleanup EXIT INT TERM

echo "Starting Go backend microservices and databases..."
./scripts/start-dev.sh &
BACKEND_PID=$!

# Wait briefly for backend ports to start opening before launching frontend
sleep 2

echo "Starting Next.js Frontend..."
if [ -d "web" ]; then
  npm --prefix web run dev
else
  echo "Error: web directory not found"
  exit 1
fi
