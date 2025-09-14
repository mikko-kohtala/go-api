package logger_test

import (
	"errors"
	"log/slog"
	"os"
	"time"

	"github.com/mikko-kohtala/go-api/pkg/logger"
)

// Example_basicUsage demonstrates simple logging scenarios
func Example_basicUsage() {
	// Create logger with pretty format for development
	log := logger.New(
		logger.WithFormat("pretty"),
		logger.WithOutput(os.Stdout),
	)

	// Simple messages
	log.Info("Server started", slog.Int("port", 8080))
	log.Warn("Cache miss", slog.String("key", "user:123"))
	log.Debug("Query executed", slog.Duration("time", 23*time.Millisecond))
}

// Example_errorHandling shows error logging patterns
func Example_errorHandling() {
	log := logger.New(logger.WithFormat("pretty"))

	// Database error
	err := errors.New("connection timeout")
	log.Error("Database query failed",
		slog.String("error", err.Error()),
		slog.String("query", "SELECT * FROM users"),
		slog.Int("retry", 3),
	)

	// Validation errors
	log.Error("Request validation failed",
		slog.Any("errors", map[string]string{
			"email": "invalid format",
			"age":   "must be positive",
		}),
	)
}

// Example_structuredData demonstrates logging complex objects
func Example_structuredData() {
	log := logger.New(logger.WithFormat("json"))

	// User registration
	type User struct {
		ID    string `json:"id"`
		Email string `json:"email"`
		Plan  string `json:"plan"`
	}

	user := User{
		ID:    "usr_abc123",
		Email: "alice@example.com",
		Plan:  "premium",
	}

	log.Info("New user registered",
		slog.Any("user", user),
		slog.String("referral", "google"),
	)

	// API metrics
	log.Info("API metrics",
		slog.Group("performance",
			slog.Int64("requests", 10523),
			slog.Float64("avg_latency_ms", 42.3),
			slog.Float64("p99_latency_ms", 125.5),
		),
		slog.Group("errors",
			slog.Int("4xx", 23),
			slog.Int("5xx", 2),
			slog.Float64("error_rate", 0.24),
		),
	)
}

// Example_requestContext shows request-scoped logging
func Example_requestContext() {
	log := logger.New(logger.WithFormat("pretty"))

	// Create request-scoped logger
	requestLog := log.With(
		slog.String("request_id", "req_xyz789"),
		slog.String("user_id", "usr_456"),
		slog.String("method", "POST"),
		slog.String("path", "/api/orders"),
	)

	// All logs include request context
	requestLog.Info("Processing order")
	requestLog.Info("Payment authorized", slog.Float64("amount", 99.99))
	requestLog.Info("Order completed", slog.Duration("total_time", 1250*time.Millisecond))
}