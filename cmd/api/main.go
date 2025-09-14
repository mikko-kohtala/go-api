package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	_ "go.uber.org/automaxprocs"

	"github.com/mikko-kohtala/go-api/internal/config"
	"github.com/mikko-kohtala/go-api/internal/httpserver"
	"github.com/mikko-kohtala/go-api/internal/logging"
)

//go:generate swag init -g cmd/api/main.go -o internal/docs --parseDependency --parseInternal
// @title           Init Codex API
// @version         1.0
// @description     A minimal, modern Go HTTP API template using chi, slog, Swagger, and optional rate limiting.
// @BasePath        /

func main() {
	// Load configuration from env with sane defaults
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Configure slog JSON logger
	logger := logging.New(cfg.Env)

	// CORS strict enforcement in production if enabled
	if (cfg.Env == "production" || cfg.Env == "prod") && cfg.CORSStrict {
		for _, o := range cfg.CORSAllowedOrigins {
			if strings.TrimSpace(o) == "*" {
				log.Fatalf("CORS_STRICT=true in production but CORS_ALLOWED_ORIGINS contains '*'")
			}
		}
	}

	// Build the HTTP server (router, middleware, handlers)
	mux := httpserver.NewRouter(cfg, logger)

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.Port),
		Handler:           mux,
		ReadTimeout:       15 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       120 * time.Second,
		MaxHeaderBytes:    1 << 20, // 1 MiB
	}

	// Start server in background
	go func() {
		logger.Info("http server starting", slog.Int("port", cfg.Port))
		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			logger.Error("http server failed", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("graceful shutdown failed", slog.String("error", err.Error()))
		_ = srv.Close()
	}
	logger.Info("server stopped")
}