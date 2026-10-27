# === Stage 1: Build static Go binaries ===
FROM golang:1.25-alpine AS builder
WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Compile the modular monolith server and content-sync CLI
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o server ./cmd/server
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o content-sync ./services/question/cmd/content-sync

# === Stage 2: Extract migrate binary ===
FROM migrate/migrate:v4.18.1 AS migrate-bin

# === Stage 3: Minimal production runner ===
FROM alpine:3.20

# Install runtime utilities: postgresql-client (for pg_isready), bash, curl, ca-certificates, tzdata (time zones for streaks)
RUN apk add --no-cache bash curl ca-certificates postgresql-client tzdata

WORKDIR /app

# Copy the migration tool
COPY --from=migrate-bin /migrate /usr/local/bin/migrate

# Copy compiled binaries
COPY --from=builder /build/server /app/server
COPY --from=builder /build/content-sync /app/content-sync

# Copy database migrations
COPY migrations /app/migrations

# Copy authored lesson content (loaded into the database by content-sync on every deploy)
COPY content /app/content

# Copy entrypoint script
COPY scripts/prod-entrypoint.sh /app/prod-entrypoint.sh
# Normalise line endings so images built from a Windows checkout (CRLF) still run.
RUN tr -d '\015' < /app/prod-entrypoint.sh > /tmp/entrypoint.sh && mv /tmp/entrypoint.sh /app/prod-entrypoint.sh && chmod +x /app/prod-entrypoint.sh

EXPOSE 8080

ENTRYPOINT ["/app/prod-entrypoint.sh"]
