package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/prepio/prepio/services/notification/internal/consumer"
	"github.com/prepio/prepio/services/notification/internal/handler"
	"github.com/prepio/prepio/services/notification/internal/service"
	"github.com/prepio/prepio/services/notification/internal/store"
	"github.com/prepio/prepio/shared/events"
	"github.com/prepio/prepio/shared/kafka"
	"github.com/prepio/prepio/shared/postgres"
	redisclient "github.com/prepio/prepio/shared/redis"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool, err := postgres.New(ctx, envOrDefault("DATABASE_URL", "postgres://prepio:prepio@localhost:5432/prepio?sslmode=disable"))
	if err != nil {
		log.Fatalf("postgres: %v", err)
	}
	defer pool.Close()

	redisClient, err := redisclient.New(ctx, redisclient.Config{Addr: envOrDefault("REDIS_ADDR", "localhost:6379")})
	if err != nil {
		log.Fatalf("redis: %v", err)
	}
	defer redisClient.Close()

	notificationService := service.NewNotificationService(store.NewNotificationStore(pool), redisClient)

	if !devSyncEnabled() {
		brokers := strings.Split(envOrDefault("KAFKA_BROKERS", "localhost:9092"), ",")
		startConsumer(ctx, brokers, events.TopicProgressUpdated, "notification-service", func(ctx context.Context, value []byte) error {
			var event events.ProgressUpdated
			if err := json.Unmarshal(value, &event); err != nil {
				return fmt.Errorf("decode progress updated: %w", err)
			}
			return notificationService.HandleProgressUpdated(ctx, event)
		})
		startConsumer(ctx, brokers, events.TopicStreakUpdated, "notification-service", func(ctx context.Context, value []byte) error {
			var event events.StreakUpdated
			if err := json.Unmarshal(value, &event); err != nil {
				return fmt.Errorf("decode streak updated: %w", err)
			}
			return notificationService.HandleStreakUpdated(ctx, event)
		})
		startConsumer(ctx, brokers, events.TopicNotificationsDispatch, "notification-service", func(ctx context.Context, value []byte) error {
			var event events.NotificationsDispatch
			if err := json.Unmarshal(value, &event); err != nil {
				return fmt.Errorf("decode notification dispatch: %w", err)
			}
			return notificationService.HandleDispatch(ctx, event)
		})
		log.Printf("notification service consuming on kafka %v", brokers)
	} else {
		log.Printf("DEV_SYNC_EVENTS=true — kafka consumers disabled")
	}

	notificationHandler := handler.NewNotificationHandler(notificationService)
	r := chi.NewRouter()
	r.Use(chimw.Recoverer)
	r.Post("/internal/events/progress-updated", notificationHandler.InternalProgressUpdated)
	r.Post("/internal/events/streak-updated", notificationHandler.InternalStreakUpdated)

	port := envOrDefault("NOTIFICATION_SERVICE_PORT", "8085")
	srv := &http.Server{Addr: ":" + port, Handler: r, ReadHeaderTimeout: 10 * time.Second}

	go func() {
		log.Printf("notification service listening on :%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown: %v", err)
	}
}

func startConsumer(ctx context.Context, brokers []string, topic, group string, handler func(context.Context, []byte) error) {
	c, err := kafka.NewConsumer(kafka.ConsumerConfig{Brokers: brokers, Topic: topic, GroupID: group})
	if err != nil {
		log.Fatalf("consumer %s: %v", topic, err)
	}
	go func() {
		_ = consumer.Run(ctx, c, handler)
		_ = c.Close()
	}()
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); len(v) > 0 {
		return v
	}
	return fallback
}

func devSyncEnabled() bool {
	return strings.EqualFold(os.Getenv("DEV_SYNC_EVENTS"), "true")
}
