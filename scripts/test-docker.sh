#!/usr/bin/env bash
# Runs the Go test suite inside a Linux container.
#
# Why: on Windows machines with Smart App Control / Application Control, every freshly
# compiled (unsigned) test .exe can be blocked. Inside Docker no Windows binary is run.
# Module and build caches persist in named volumes, so repeat runs are fast.
#
#   scripts/test-docker.sh                          # go test ./... -count=1
#   scripts/test-docker.sh ./services/question/...  # specific packages
#   scripts/test-docker.sh ./test/integration/ -run TestLessonPlatform -v
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
if command -v cygpath >/dev/null 2>&1; then ROOT="$(cygpath -m "$ROOT")"; fi
IMAGE="${GO_TEST_IMAGE:-golang:1.25}"
export MSYS_NO_PATHCONV=1

# The embedded Postgres refuses to run as root, so tests run as uid 1000.
docker volume create prepio-go-mod >/dev/null
docker volume create prepio-go-build >/dev/null
docker volume create prepio-test-home >/dev/null
docker run --rm -v prepio-go-mod:/gomod -v prepio-go-build:/gobuild -v prepio-test-home:/home/tester \
  busybox chown -R 1000:1000 /gomod /gobuild /home/tester

ARGS=("$@")
if [ ${#ARGS[@]} -eq 0 ]; then ARGS=(./... -count=1 -timeout 30m); fi

exec docker run --rm --user 1000:1000 \
  -e HOME=/home/tester -e GOMODCACHE=/gomod -e GOCACHE=/gobuild -e GOFLAGS=-buildvcs=false \
  -v prepio-go-mod:/gomod -v prepio-go-build:/gobuild -v prepio-test-home:/home/tester \
  -v "$ROOT":/src -w /src \
  "$IMAGE" go test "${ARGS[@]}"
