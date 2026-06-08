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
	"github.com/prepio/prepio/services/user/internal/handler"
	"github.com/prepio/prepio/services/user/internal/service"
	"github.com/prepio/prepio/services/user/internal/store"
	"github.com/prepio/prepio/shared/jwt"
	"github.com/prepio/prepio/shared/middleware"
	"github.com/prepio/prepio/shared/postgres"
	redisclient "github.com/prepio/prepio/shared/redis"
)

func main() {
	ctx := context.Background()

	dsn := envOrDefault("DATABASE_URL", "postgres://prepio:prepio@localhost:5432/prepio?sslmode=disable")
	pool, err := postgres.New(ctx, dsn)
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

	signer, err := jwt.NewSigner(envOrDefault("JWT_SECRET", "dev-secret-change-in-production"))
	if err != nil {
		log.Fatalf("jwt: %v", err)
	}

	userStore := store.NewUserStore(pool)
	refreshStore := store.NewRefreshTokenStore(pool)
	deviceStore := store.NewUserDeviceStore(pool)

	authService := service.NewAuthService(userStore, refreshStore, signer, redisClient)
	userService := service.NewUserService(userStore, deviceStore)

	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)

	r := chi.NewRouter()
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/auth/register", authHandler.Register)
		r.Post("/auth/login", authHandler.Login)
		r.Post("/auth/refresh", authHandler.Refresh)

		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(signer, redisClient))
			r.Post("/auth/logout", authHandler.Logout)
			r.Get("/users/me", userHandler.GetMe)
			r.Patch("/users/me", userHandler.UpdateMe)
			r.Post("/users/me/devices", userHandler.RegisterDevice)
			r.Delete("/users/me/devices/{deviceID}", userHandler.DeleteDevice)
		})
	})

	port := envOrDefault("USER_SERVICE_PORT", "8081")
	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           r,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		log.Printf("user service listening on :%s", port)
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
