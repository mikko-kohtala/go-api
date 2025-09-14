package main

import (
	"log/slog"
	"time"

	"github.com/mikko-kohtala/go-api/pkg/logger"
)

func main() {
	// Create logger with pretty format
	log := logger.New(
		logger.WithFormat("pretty"),
		logger.WithLevel(slog.LevelDebug),
	)

	// Different log levels
	log.Debug("Starting application", slog.String("version", "1.0.0"))
	log.Info("Server started", slog.Int("port", 8080))
	log.Warn("Cache miss", slog.String("key", "user:123"))
	log.Error("Database connection failed",
		slog.String("error", "connection timeout"),
		slog.Duration("retry_after", 5*time.Second),
	)

	// Structured logging with groups
	log.Info("User action",
		slog.Group("user",
			slog.String("id", "usr_123"),
			slog.String("email", "user@example.com"),
		),
		slog.Group("action",
			slog.String("type", "login"),
			slog.Time("timestamp", time.Now()),
		),
	)

	// JSON format
	jsonLog := logger.New(
		logger.WithFormat("json"),
	)

	jsonLog.Info("JSON formatted log",
		slog.String("format", "json"),
		slog.Bool("structured", true),
	)
}