package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/prepio/prepio/services/streak/internal/consumer"
	"github.com/prepio/prepio/services/streak/internal/handler"
	"github.com/prepio/prepio/services/streak/internal/service"
	"github.com/prepio/prepio/services/streak/internal/store"
	"github.com/prepio/prepio/shared/events"
	"github.com/prepio/prepio/shared/jwt"
	"github.com/prepio/prepio/shared/kafka"
	"github.com/prepio/prepio/shared/middleware"
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

	redisClient, err := redisclient.New(ctx, redisclient.Config{
		Addr:     envOrDefault("REDIS_ADDR", "localhost:6379"),
		Password: os.Getenv("REDIS_PASSWORD"),
	})
	if err != nil {
		log.Fatalf("redis: %v", err)
	}
	defer redisClient.Close()

	producer, err := kafka.NewEventPublisher(kafka.ProducerConfig{
		Brokers: strings.Split(envOrDefault("KAFKA_BROKERS", "localhost:9092"), ","),
	})
	if err != nil {
		log.Fatalf("event publisher: %v", err)
	}
	defer producer.Close()

	gemClient := service.NewProgressGemClient(envOrDefault("PROGRESS_SERVICE_URL", "http://localhost:8084"))

	streakService := service.NewStreakService(
		store.NewStreakStore(pool),
		store.NewFreezeStore(pool),
		store.NewCacheStore(redisClient),
		producer,
		gemClient,
	)

	if !devSyncEnabled() {
		kafkaConsumer, err := kafka.NewConsumer(kafka.ConsumerConfig{
			Brokers: strings.Split(envOrDefault("KAFKA_BROKERS", "localhost:9092"), ","),
			Topic:   events.TopicQuestionAnswered,
			GroupID: "streak-service",
		})
		if err != nil {
			log.Fatalf("kafka consumer: %v", err)
		}
		defer kafkaConsumer.Close()

		go func() {
			if err := consumer.Run(ctx, kafkaConsumer, consumer.NewQuestionAnsweredConsumer(streakService)); err != nil && ctx.Err() == nil {
				log.Printf("consumer stopped: %v", err)
			}
		}()
	} else {
		log.Printf("DEV_SYNC_EVENTS=true — kafka consumer disabled")
	}

	signer, err := jwt.NewSigner(envOrDefault("JWT_SECRET", "dev-secret-change-in-production"))
	if err != nil {
		log.Fatalf("jwt: %v", err)
	}

	streakHandler := handler.NewStreakHandler(streakService)
	r := chi.NewRouter()
	r.Use(chimw.Recoverer)
	r.Post("/internal/events/question-answered", streakHandler.InternalQuestionAnswered)
	r.Route("/api/v1", func(r chi.Router) {
		r.Use(middleware.Auth(signer, redisClient))
		r.Get("/streaks/me", streakHandler.GetMe)
		r.Post("/streaks/me/freeze/purchase", streakHandler.PurchaseFreeze)
	})

	port := envOrDefault("STREAK_SERVICE_PORT", "8083")
	srv := &http.Server{Addr: ":" + port, Handler: r, ReadHeaderTimeout: 10 * time.Second}

	go func() {
		log.Printf("streak service listening on :%s", port)
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

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); len(v) > 0 {
		return v
	}
	return fallback
}

func devSyncEnabled() bool {
	return strings.EqualFold(os.Getenv("DEV_SYNC_EVENTS"), "true")
}
