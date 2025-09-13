package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/user/go-api-template/internal/config"
	"github.com/user/go-api-template/internal/handlers"
	"github.com/user/go-api-template/internal/middleware"
	"github.com/user/go-api-template/pkg/logger"
)

func main() {
	cfg := config.Load()

	log := logger.New(cfg.LogLevel)
	slog.SetDefault(log)

	router := setupRouter(cfg, log)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		slog.Info("starting server", "port", cfg.Port, "env", cfg.Environment)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("failed to start server", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("server forced to shutdown", "error", err)
		os.Exit(1)
	}

	slog.Info("server stopped")
}

func setupRouter(cfg *config.Config, log *slog.Logger) http.Handler {
	mux := http.NewServeMux()

	h := handlers.New(log)

	mux.HandleFunc("GET /health", h.Health)
	mux.HandleFunc("GET /api/v1/users", h.GetUsers)
	mux.HandleFunc("GET /api/v1/users/{id}", h.GetUser)
	mux.HandleFunc("POST /api/v1/users", h.CreateUser)
	mux.HandleFunc("PUT /api/v1/users/{id}", h.UpdateUser)
	mux.HandleFunc("DELETE /api/v1/users/{id}", h.DeleteUser)

	handler := middleware.Chain(
		mux,
		middleware.RateLimit(cfg.RateLimitRequests, cfg.RateLimitDuration),
		middleware.CORS(cfg.AllowedOrigins),
		middleware.RequestID,
		middleware.Logger(log),
		middleware.Recover(log),
	)

	return handler
}