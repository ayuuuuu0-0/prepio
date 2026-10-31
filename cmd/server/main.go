package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/prepio/prepio/constants"
	gatewaydashboard "github.com/prepio/prepio/internal/gateway/dashboard"
	gatewayleagues "github.com/prepio/prepio/internal/gateway/leagues"
	gatewaylessons "github.com/prepio/prepio/internal/gateway/lessons"
	notificationhandler "github.com/prepio/prepio/internal/notification/handler"
	notificationservice "github.com/prepio/prepio/internal/notification/service"
	notificationstore "github.com/prepio/prepio/internal/notification/store"
	progresshandler "github.com/prepio/prepio/internal/progress/handler"
	progressservice "github.com/prepio/prepio/internal/progress/service"
	progressstore "github.com/prepio/prepio/internal/progress/store"
	questionhandler "github.com/prepio/prepio/internal/question/handler"
	questionservice "github.com/prepio/prepio/internal/question/service"
	questionstore "github.com/prepio/prepio/internal/question/store"
	streakhandler "github.com/prepio/prepio/internal/streak/handler"
	streakservice "github.com/prepio/prepio/internal/streak/service"
	streakstore "github.com/prepio/prepio/internal/streak/store"
	userhandler "github.com/prepio/prepio/internal/user/handler"
	usermailer "github.com/prepio/prepio/internal/user/mailer"
	userservice "github.com/prepio/prepio/internal/user/service"
	userstore "github.com/prepio/prepio/internal/user/store"
	"github.com/prepio/prepio/shared/devsync"
	"github.com/prepio/prepio/shared/events"
	"github.com/prepio/prepio/shared/jwt"
	"github.com/prepio/prepio/shared/kafka"
	"github.com/prepio/prepio/shared/middleware"
	"github.com/prepio/prepio/shared/postgres"
	redisclient "github.com/prepio/prepio/shared/redis"
	"github.com/prepio/prepio/shared/response"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Println("Starting Prepio Modular Monolith Server...")

	// 1. Single shared PostgreSQL connection pool
	dbURL := envOrDefault("DATABASE_URL", "postgres://prepio:prepio@localhost:5432/prepio?sslmode=disable")
	pool, err := postgres.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("postgres pool: %v", err)
	}
	defer pool.Close()
	log.Println("✓ Connected to PostgreSQL")

	// 2. Single shared Redis client
	redisClient, err := redisclient.New(ctx, redisclient.Config{
		Addr:     envOrDefault("REDIS_ADDR", "localhost:6379"),
		Password: os.Getenv("REDIS_PASSWORD"),
	})
	if err != nil {
		log.Fatalf("redis client: %v", err)
	}
	defer redisClient.Close()
	log.Println("✓ Connected to Redis")

	// 3. JWT Signer
	signer, err := jwt.NewSigner(envOrDefault("JWT_SECRET", "dev-secret-change-in-production-use-32-chars"))
	if err != nil {
		log.Fatalf("jwt signer: %v", err)
	}

	port := envOrDefault("GATEWAY_PORT", envOrDefault("PORT", "8080"))
	serverBaseURL := "http://localhost:" + port

	// 4. Unified Event Publisher
	var producer kafka.EventPublisher
	if devSyncEnabled() {
		producer = devsync.NewPublisher(devsync.Config{
			StreakURL:       serverBaseURL,
			ProgressURL:     serverBaseURL,
			NotificationURL: serverBaseURL,
		})
		log.Println("✓ Event publisher: dev-sync (in-process loopback)")
	} else {
		producer, err = kafka.NewEventPublisher(kafka.ProducerConfig{
			Brokers: strings.Split(envOrDefault("KAFKA_BROKERS", "localhost:9092"), ","),
		})
		if err != nil {
			log.Fatalf("kafka event publisher: %v", err)
		}
		defer producer.Close()
		log.Println("✓ Event publisher: Kafka")
	}

	// --------------------------------------------------------------------------
	// 5. Module Initialization
	// --------------------------------------------------------------------------

	// --- User Module ---
	userStore := userstore.NewUserStore(pool)
	refreshStore := userstore.NewRefreshTokenStore(pool)
	deviceStore := userstore.NewUserDeviceStore(pool)
	characterStore := userstore.NewCharacterStore(pool)
	mailClient := usermailer.NewMailerFromEnv()

	userAuthService := userservice.NewAuthService(userStore, refreshStore, signer, redisClient, mailClient)
	userService := userservice.NewUserService(userStore, deviceStore)
	userOnboardingService := userservice.NewOnboardingService(userStore, characterStore)

	userAuthHandler := userhandler.NewAuthHandler(userAuthService)
	userHandler := userhandler.NewUserHandler(userService)
	userOnboardingHandler := userhandler.NewOnboardingHandler(userOnboardingService)

	// --- Question Module ---
	questionSkillStore := questionstore.NewSkillStore(pool)
	questionLessonStore := questionstore.NewLessonStore(pool)

	geminiClient := questionservice.NewGeminiClient(os.Getenv("GEMINI_API_KEY"))
	var llm questionservice.LLMEvaluator
	if geminiClient != nil {
		llm = geminiClient
		log.Println("evaluator: gemini stage 3 enabled")
	} else {
		log.Println("evaluator: structural-only (stages 1+2)")
	}
	evaluator := questionservice.NewPipelineEvaluator(llm)

	questionLessonService := questionservice.NewLessonService(questionLessonStore, evaluator, producer)
	questionSkillService := questionservice.NewSkillService(questionSkillStore)

	questionSkillHandler := questionhandler.NewSkillHandler(questionSkillService)
	questionLessonHandler := questionhandler.NewLessonHandler(questionLessonService)

	// --- Streak Module ---
	streakGemClient := streakservice.NewProgressGemClient(serverBaseURL)
	streakService := streakservice.NewStreakService(
		streakstore.NewStreakStore(pool),
		streakstore.NewFreezeStore(pool),
		streakstore.NewCacheStore(redisClient),
		producer,
		streakGemClient,
	)
	streakHandler := streakhandler.NewStreakHandler(streakService)

	// --- Progress Module ---
	progressAchievementStore := progressstore.NewAchievementStore(pool)
	progressService := progressservice.NewProgressService(
		progressstore.NewProgressStore(pool),
		progressAchievementStore,
		producer,
	)
	progressReadinessStore := progressstore.NewReadinessStore(pool)
	progressReadinessService := progressservice.NewReadinessService(progressReadinessStore)
	progressLessonService := progressservice.NewLessonService(progressstore.NewLessonStore(pool), producer)
	progressLeagueService := progressservice.NewLeagueService(progressstore.NewLeagueStore(pool))

	progressHandler := progresshandler.NewProgressHandler(progressService, progressReadinessService)
	progressReadinessHandler := progresshandler.NewReadinessHandler(progressReadinessService)
	progressLessonHandler := progresshandler.NewLessonHandler(progressLessonService)
	progressLeagueHandler := progresshandler.NewLeagueHandler(progressLeagueService)
	progressAchievementHandler := progresshandler.NewAchievementHandler(progressservice.NewAchievementService(progressAchievementStore))

	// --- Notification Module ---
	notificationStore := notificationstore.NewNotificationStore(pool)
	notificationService := notificationservice.NewNotificationService(notificationStore, redisClient)
	notificationHandler := notificationhandler.NewNotificationHandler(notificationService)

	// --- Gateway Aggregators ---
	dashboardService := gatewaydashboard.NewService(serverBaseURL, serverBaseURL, serverBaseURL, serverBaseURL)
	dashboardHandler := gatewaydashboard.NewHandler(dashboardService)

	leaguesService := gatewayleagues.NewService(serverBaseURL, serverBaseURL)
	leaguesHandler := gatewayleagues.NewHandler(leaguesService)

	lessonsHandler := gatewaylessons.NewHandler(serverBaseURL+"/internal/question", serverBaseURL)

	// --------------------------------------------------------------------------
	// 6. Router Setup & Middleware
	// --------------------------------------------------------------------------
	r := chi.NewRouter()
	r.Use(middleware.CORS)
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)

	// Healthcheck
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		response.Data(w, http.StatusOK, map[string]string{
			"status": "healthy",
			"mode":   "modular_monolith",
		})
	})

	// Internal inter-module endpoints (for events & gem transactions)
	r.Get("/internal/progress/{userID}/gems", progressHandler.InternalGetGems)
	r.Post("/internal/progress/{userID}/gems/deduct", progressHandler.InternalDeductGems)
	r.Post("/internal/events/streak-updated", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			response.Error(w, http.StatusBadRequest, constants.ErrInvalidRequest, "invalid body")
			return
		}
		// 1. Notify Progress
		r1 := r.Clone(r.Context())
		r1.Body = io.NopCloser(bytes.NewReader(body))
		progressHandler.InternalStreakUpdated(w, r1)

		// 2. Notify Notification service
		var event events.StreakUpdated
		if err := json.Unmarshal(body, &event); err == nil {
			_ = notificationService.HandleStreakUpdated(r.Context(), event)
		}
	})
	r.Post("/internal/events/progress-updated", func(w http.ResponseWriter, r *http.Request) {
		var event events.ProgressUpdated
		if err := json.NewDecoder(r.Body).Decode(&event); err == nil {
			_ = notificationService.HandleProgressUpdated(r.Context(), event)
		}
		response.Data(w, http.StatusOK, map[string]bool{"ok": true})
	})
	r.Post("/internal/events/lesson-completed", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			response.Error(w, http.StatusBadRequest, constants.ErrInvalidRequest, "invalid body")
			return
		}
		// 1. Notify Streak
		rStreak := r.Clone(r.Context())
		rStreak.Body = io.NopCloser(bytes.NewReader(body))
		streakHandler.InternalLessonCompleted(w, rStreak)

		// 2. Notify Progress
		rProgress := r.Clone(r.Context())
		rProgress.Body = io.NopCloser(bytes.NewReader(body))
		progressLessonHandler.InternalLessonCompleted(w, rProgress)
	})
	r.Post("/internal/streak/events/lesson-completed", streakHandler.InternalLessonCompleted)
	r.Post("/internal/progress/events/lesson-completed", progressLessonHandler.InternalLessonCompleted)
	r.With(middleware.Auth(signer, redisClient)).Post("/internal/question/api/v1/attempts/{id}/complete", questionLessonHandler.Complete)
	r.Post("/internal/notification/events/progress-updated", notificationHandler.InternalProgressUpdated)
	r.Post("/internal/notification/events/streak-updated", notificationHandler.InternalStreakUpdated)

	// API v1 Routes
	r.Route("/api/v1", func(r chi.Router) {
		// Public routes (Rate limited by IP)
		r.Group(func(r chi.Router) {
			r.Use(middleware.RateLimit(redisClient, constants.UnauthenticatedRateLimitPerMinute, middleware.RateLimitKeyByIP))

			// Auth public
			r.Post("/auth/otp/send", userAuthHandler.SendOTP)
			r.Post("/auth/otp/verify", userAuthHandler.VerifyOTP)
			r.Post("/auth/register", userAuthHandler.Register)
			r.Post("/auth/login", userAuthHandler.Login)
			r.Post("/auth/refresh", userAuthHandler.Refresh)

			// Companions catalog
			r.Get("/companions", userOnboardingHandler.ListCompanions)
		})

		// Protected routes (Authenticated + Rate limited by User)
		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(signer, redisClient))
			r.Use(middleware.RateLimit(redisClient, constants.AuthenticatedRateLimitPerMinute, middleware.RateLimitKeyByUser))

			// Auth logout
			r.Post("/auth/logout", userAuthHandler.Logout)

			// Users
			r.Get("/users/me", userHandler.GetMe)
			r.Patch("/users/me", userHandler.UpdateMe)
			r.Get("/users/profile", userOnboardingHandler.GetProfile)
			r.Post("/users/onboarding", userOnboardingHandler.CompleteOnboarding)
			r.Post("/users/public-cards", userHandler.PublicCards)
			r.Post("/users/me/devices", userHandler.RegisterDevice)
			r.Delete("/users/me/devices/{deviceID}", userHandler.DeleteDevice)
			r.Get("/companions/reaction", userOnboardingHandler.GetReaction)

			// Gateway aggregates
			r.Get("/dashboard/home", dashboardHandler.GetHome)
			r.Get("/league", leaguesHandler.GetLeague)

			// Journey & Lessons
			r.Get("/path", questionLessonHandler.GetPath)
			r.Get("/topics", questionSkillHandler.ListTopics)
			r.Get("/skills", questionSkillHandler.ListSkills)
			r.Get("/skills/{slug}", questionSkillHandler.GetSkill)
			r.Post("/lessons/{id}/attempts", questionLessonHandler.StartAttempt)
			r.Post("/attempts/{id}/steps/{stepId}/answer", questionLessonHandler.SubmitAnswer)
			r.Post("/attempts/{id}/complete", lessonsHandler.Complete)

			// Streaks
			r.Get("/streaks/me", streakHandler.GetMe)
			r.Post("/streaks/me/freeze/purchase", streakHandler.PurchaseFreeze)

			// Progress
			r.Get("/progress/me", progressHandler.GetMe)
			r.Get("/progress/topics", progressReadinessHandler.GetTopicMastery)
			r.Get("/progress/league", progressLeagueHandler.GetLeague)
			r.Get("/progress/achievements", progressAchievementHandler.List)
			r.Get("/progress/attempts/{attemptID}/rewards", progressLessonHandler.GetAttemptRewards)
		})
	})

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           r,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		log.Printf("🚀 Prepio Modular Monolith listening on http://localhost:%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Println("Shutting down modular monolith gracefully...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown error: %v", err)
	}
	log.Println("Server stopped.")
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); len(v) > 0 {
		return v
	}
	return fallback
}

func devSyncEnabled() bool {
	v := os.Getenv("DEV_SYNC_EVENTS")
	return v == "" || v == "true" || v == "1"
}
