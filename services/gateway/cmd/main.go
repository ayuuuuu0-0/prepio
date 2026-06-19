package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/prepio/prepio/constants"
	"github.com/prepio/prepio/services/gateway/internal/dashboard"
	"github.com/prepio/prepio/shared/proxy"
	"github.com/prepio/prepio/shared/jwt"
	"github.com/prepio/prepio/shared/middleware"
	redisclient "github.com/prepio/prepio/shared/redis"
)

func main() {
	ctx := context.Background()

	redisClient, err := redisclient.New(ctx, redisclient.Config{
		Addr:     envOrDefault("REDIS_ADDR", "localhost:6379"),
		Password: os.Getenv("REDIS_PASSWORD"),
	})
	if err != nil {
		log.Fatalf("redis: %v", err)
	}
	defer redisClient.Close()

	signer, err := jwt.NewSigner(envOrDefault("JWT_SECRET", "dev-secret-change-in-production"))
	if err != nil {
		log.Fatalf("jwt: %v", err)
	}

	userProxy, err := proxy.New(envOrDefault("USER_SERVICE_URL", "http://localhost:8081"))
	if err != nil {
		log.Fatalf("user proxy: %v", err)
	}

	questionProxy, err := proxy.New(envOrDefault("QUESTION_SERVICE_URL", "http://localhost:8082"))
	if err != nil {
		log.Fatalf("question proxy: %v", err)
	}

	streakProxy, err := proxy.New(envOrDefault("STREAK_SERVICE_URL", "http://localhost:8083"))
	if err != nil {
		log.Fatalf("streak proxy: %v", err)
	}

	progressProxy, err := proxy.New(envOrDefault("PROGRESS_SERVICE_URL", "http://localhost:8084"))
	if err != nil {
		log.Fatalf("progress proxy: %v", err)
	}

	dashboardService := dashboard.NewService(
		envOrDefault("USER_SERVICE_URL", "http://localhost:8081"),
		envOrDefault("PROGRESS_SERVICE_URL", "http://localhost:8084"),
		envOrDefault("STREAK_SERVICE_URL", "http://localhost:8083"),
		envOrDefault("QUESTION_SERVICE_URL", "http://localhost:8082"),
	)
	dashboardHandler := dashboard.NewHandler(dashboardService)

	r := chi.NewRouter()
	r.Use(middleware.CORS)
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)

	r.Route("/api/v1", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(middleware.RateLimit(redisClient, constants.UnauthenticatedRateLimitPerMinute, middleware.RateLimitKeyByIP))
			r.Post("/auth/register", userProxy.ServeHTTP)
			r.Post("/auth/login", userProxy.ServeHTTP)
			r.Post("/auth/refresh", userProxy.ServeHTTP)
			r.Get("/companions", userProxy.ServeHTTP)
		})

		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(signer, redisClient))
			r.Use(middleware.RateLimit(redisClient, constants.AuthenticatedRateLimitPerMinute, middleware.RateLimitKeyByUser))
			r.Get("/dashboard/home", dashboardHandler.GetHome)
			r.Get("/dashboard/readiness", dashboardHandler.GetReadinessValidation)
			r.Get("/internal/readiness/compare", dashboardHandler.GetInternalReadinessCompare)
			r.Get("/journey", questionProxy.ServeHTTP)
			r.Handle("/journey/*", questionProxy)
			r.Get("/skills/readiness", progressProxy.ServeHTTP)
			r.Get("/companies/readiness", progressProxy.ServeHTTP)
			r.Handle("/skills/*", questionProxy)
			r.Get("/skills", questionProxy.ServeHTTP)
			r.Post("/auth/logout", userProxy.ServeHTTP)
			r.Handle("/users/*", userProxy)
			r.Handle("/questions/*", questionProxy)
			r.Handle("/streaks/*", streakProxy)
			r.Handle("/progress/*", progressProxy)
			r.Get("/readiness/dashboard", progressProxy.ServeHTTP)
		})
	})

	port := envOrDefault("GATEWAY_PORT", "8080")
	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           r,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		log.Printf("gateway listening on :%s", port)
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
