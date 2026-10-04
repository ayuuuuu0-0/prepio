# === Stage 1: Build static Go binaries ===
FROM golang:1.24-alpine AS builder
WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Compile each microservice as a static, lightweight binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o user ./services/user/cmd
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o question ./services/question/cmd
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o streak ./services/streak/cmd
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o progress ./services/progress/cmd
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o notification ./services/notification/cmd
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o gateway ./services/gateway/cmd

# === Stage 2: Extract migrate binary ===
FROM migrate/migrate:v4.18.1 AS migrate-bin

# === Stage 3: Minimal production runner ===
FROM alpine:3.20

# Install runtime utilities: postgresql-client (for pg_isready), bash, curl, ca-certificates
RUN apk add --no-cache bash curl ca-certificates postgresql-client

WORKDIR /app

# Copy the migration tool
COPY --from=migrate-bin /migrate /usr/local/bin/migrate

# Copy compiled binaries
COPY --from=builder /build/user /app/user
COPY --from=builder /build/question /app/question
COPY --from=builder /build/streak /app/streak
COPY --from=builder /build/progress /app/progress
COPY --from=builder /build/notification /app/notification
COPY --from=builder /build/gateway /app/gateway

# Copy database migrations
COPY migrations /app/migrations

# Copy entrypoint script
COPY scripts/prod-entrypoint.sh /app/prod-entrypoint.sh
RUN chmod +x /app/prod-entrypoint.sh

EXPOSE 8080

ENTRYPOINT ["/app/prod-entrypoint.sh"]
