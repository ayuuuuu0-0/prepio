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
	"github.com/prepio/prepio/services/progress/internal/consumer"
	"github.com/prepio/prepio/services/progress/internal/handler"
	"github.com/prepio/prepio/services/progress/internal/service"
	"github.com/prepio/prepio/services/progress/internal/store"
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

	producer, err := kafka.NewEventPublisher(kafka.ProducerConfig{
		Brokers: strings.Split(envOrDefault("KAFKA_BROKERS", "localhost:9092"), ","),
	})
	if err != nil {
		log.Fatalf("event publisher: %v", err)
	}
	defer producer.Close()

	progressService := service.NewProgressService(
		store.NewProgressStore(pool),
		store.NewLedgerStore(pool),
		producer,
	)

	if !devSyncEnabled() {
		brokers := strings.Split(envOrDefault("KAFKA_BROKERS", "localhost:9092"), ",")
		questionConsumer, err := kafka.NewConsumer(kafka.ConsumerConfig{
			Brokers: brokers, Topic: events.TopicQuestionAnswered, GroupID: "progress-service",
		})
		if err != nil {
			log.Fatalf("question consumer: %v", err)
		}
		defer questionConsumer.Close()

		streakConsumer, err := kafka.NewConsumer(kafka.ConsumerConfig{
			Brokers: brokers, Topic: events.TopicStreakUpdated, GroupID: "progress-service",
		})
		if err != nil {
			log.Fatalf("streak consumer: %v", err)
		}
		defer streakConsumer.Close()

		eventHandler := consumer.NewHandler(progressService)
		go func() {
			if err := consumer.RunQuestionAnswered(ctx, questionConsumer, eventHandler); err != nil && ctx.Err() == nil {
				log.Printf("question consumer: %v", err)
			}
		}()
		go func() {
			if err := consumer.RunStreakUpdated(ctx, streakConsumer, eventHandler); err != nil && ctx.Err() == nil {
				log.Printf("streak consumer: %v", err)
			}
		}()
	} else {
		log.Printf("DEV_SYNC_EVENTS=true — kafka consumers disabled")
	}

	redisClient, err := redisclient.New(ctx, redisclient.Config{Addr: envOrDefault("REDIS_ADDR", "localhost:6379")})
	if err != nil {
		log.Fatalf("redis: %v", err)
	}
	defer redisClient.Close()

	signer, err := jwt.NewSigner(envOrDefault("JWT_SECRET", "dev-secret-change-in-production"))
	if err != nil {
		log.Fatalf("jwt: %v", err)
	}

	progressHandler := handler.NewProgressHandler(progressService)
	r := chi.NewRouter()
	r.Use(chimw.Recoverer)
	r.Get("/internal/progress/{userID}/gems", progressHandler.InternalGetGems)
	r.Post("/internal/progress/{userID}/gems/deduct", progressHandler.InternalDeductGems)
	r.Post("/internal/events/question-answered", progressHandler.InternalQuestionAnswered)
	r.Post("/internal/events/streak-updated", progressHandler.InternalStreakUpdated)
	r.Route("/api/v1", func(r chi.Router) {
		r.Use(middleware.Auth(signer, redisClient))
		r.Get("/progress/me", progressHandler.GetMe)
	})

	port := envOrDefault("PROGRESS_SERVICE_PORT", "8084")
	srv := &http.Server{Addr: ":" + port, Handler: r, ReadHeaderTimeout: 10 * time.Second}

	go func() {
		log.Printf("progress service listening on :%s", port)
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
