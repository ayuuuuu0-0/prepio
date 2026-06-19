# === Stage 1: Build the Go microservices ===
FROM golang:1.22-alpine AS builder
WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Compile each microservice as a static binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o user ./services/user/cmd
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o question ./services/question/cmd
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o streak ./services/streak/cmd
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o progress ./services/progress/cmd
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o notification ./services/notification/cmd
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o gateway ./services/gateway/cmd

# === Stage 2: Extract migrate binary ===
FROM migrate/migrate:v4.18.1 AS migrate-bin

# === Stage 3: Runner container ===
FROM alpine:3.19

# Install dependencies (PostgreSQL, Redis, Bash, curl)
RUN apk add --no-cache postgresql redis bash curl

# Copy the migration tool to system path
COPY --from=migrate-bin /migrate /usr/local/bin/migrate

WORKDIR /app

# Ensure we create a non-root user (UID 1000) match Hugging Face's environment
# If UID 1000 doesn't exist, create it.
RUN getent passwd 1000 >/dev/null || adduser -u 1000 -D -g "" user

# Copy build output from Stage 1
COPY --from=builder /build/user /app/user
COPY --from=builder /build/question /app/question
COPY --from=builder /build/streak /app/streak
COPY --from=builder /build/progress /app/progress
COPY --from=builder /build/notification /app/notification
COPY --from=builder /build/gateway /app/gateway

# Copy migrations files
COPY migrations /app/migrations

# Copy and configure the startup script
COPY scripts/hf-entrypoint.sh /app/hf-entrypoint.sh
RUN chmod +x /app/hf-entrypoint.sh

# Change ownership of /app to the non-root user
RUN chown -R 1000:1000 /app

# Switch to the non-root user
USER 1000

# Expose Hugging Face default port
EXPOSE 7860

ENTRYPOINT ["/app/hf-entrypoint.sh"]
