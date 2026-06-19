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
	"github.com/prepio/prepio/services/question/internal/handler"
	"github.com/prepio/prepio/services/question/internal/service"
	"github.com/prepio/prepio/services/question/internal/store"
	"github.com/prepio/prepio/shared/jwt"
	"github.com/prepio/prepio/shared/kafka"
	"github.com/prepio/prepio/shared/middleware"
	"github.com/prepio/prepio/shared/postgres"
	redisclient "github.com/prepio/prepio/shared/redis"
)

func main() {
	ctx := context.Background()

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

	signer, err := jwt.NewSigner(envOrDefault("JWT_SECRET", "dev-secret-change-in-production"))
	if err != nil {
		log.Fatalf("jwt: %v", err)
	}

	skillStore := store.NewSkillStore(pool)
	contentStore := store.NewContentStore(pool)
	journeyStore := store.NewJourneyStore(pool)

	questionService := service.NewQuestionService(
		store.NewQuestionStore(pool),
		store.NewDailyPaperStore(pool),
		store.NewHistoryStore(pool),
		journeyStore,
		contentStore,
		store.NewUserStore(pool),
		redisClient,
		producer,
	)
	skillService := service.NewSkillService(skillStore)
	contentService := service.NewContentService(contentStore, journeyStore)

	questionHandler := handler.NewQuestionHandler(questionService, contentService)
	skillHandler := handler.NewSkillHandler(skillService)
	contentHandler := handler.NewContentHandler(contentService, questionService)

	r := chi.NewRouter()
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)

	r.Route("/api/v1", func(r chi.Router) {
		r.Use(middleware.Auth(signer, redisClient))
		r.Get("/journey", questionHandler.GetJourney)
		r.Get("/journey/nodes/{id}/content", contentHandler.GetNodeContent)
		r.Get("/skills", skillHandler.ListSkills)
		r.Get("/skills/{slug}", skillHandler.GetSkill)
		r.Get("/questions/daily", questionHandler.GetDaily)
		r.Get("/questions/history", questionHandler.GetHistory)
		r.Get("/questions/stats/readiness", questionHandler.GetReadinessStats)
		r.Get("/questions/{id}/skills", skillHandler.GetQuestionSkills)
		r.Post("/questions/{id}/submit", questionHandler.Submit)
		r.Get("/questions/companies", questionHandler.ListCompanies)
	})

	port := envOrDefault("QUESTION_SERVICE_PORT", "8082")
	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           r,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		log.Printf("question service listening on :%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
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
