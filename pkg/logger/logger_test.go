package logger_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"testing"
	"time"

	"github.com/mikko-kohtala/go-api/pkg/logger"
)

func TestLogger_SimpleMessages(t *testing.T) {
	var buf bytes.Buffer
	log := logger.New(
		logger.WithOutput(&buf),
		logger.WithFormat("json"),
	)

	// Simple info message
	log.Info("Server started successfully")

	// Warning message
	log.Warn("Cache miss rate is high", slog.Float64("rate", 0.85))

	// Debug message
	log.Debug("Processing request", slog.String("endpoint", "/api/users"))

	output := buf.String()
	if !strings.Contains(output, "Server started successfully") {
		t.Errorf("Expected info message in output")
	}
	if !strings.Contains(output, "Cache miss rate is high") {
		t.Errorf("Expected warning message in output")
	}
}

func TestLogger_ErrorLogging(t *testing.T) {
	var buf bytes.Buffer
	log := logger.New(
		logger.WithOutput(&buf),
		logger.WithFormat("json"),
	)

	// Common error scenarios
	err := errors.New("connection refused")
	log.Error("Database connection failed",
		slog.String("error", err.Error()),
		slog.String("host", "localhost:5432"),
		slog.Int("retry_count", 3),
	)

	// Validation error with field details
	validationErrors := map[string]string{
		"email":    "invalid format",
		"password": "too short",
	}
	log.Error("Validation failed",
		slog.Any("fields", validationErrors),
		slog.String("user_id", "usr_123"),
	)

	// Critical error with stack trace simulation
	log.Error("Panic recovered",
		slog.String("error", "nil pointer dereference"),
		slog.String("stack", "main.HandleRequest:45\nmain.ProcessData:23"),
	)

	output := buf.String()
	if !strings.Contains(output, "Database connection failed") {
		t.Errorf("Expected database error in output")
	}
	if !strings.Contains(output, "Validation failed") {
		t.Errorf("Expected validation error in output")
	}
}

func TestLogger_ObjectLogging(t *testing.T) {
	var buf bytes.Buffer
	log := logger.New(
		logger.WithOutput(&buf),
		logger.WithFormat("json"),
	)

	// User object
	type User struct {
		ID        string    `json:"id"`
		Email     string    `json:"email"`
		Name      string    `json:"name"`
		CreatedAt time.Time `json:"created_at"`
	}

	user := User{
		ID:        "usr_abc123",
		Email:     "john@example.com",
		Name:      "John Doe",
		CreatedAt: time.Now(),
	}

	log.Info("User registered",
		slog.Any("user", user),
		slog.String("ip_address", "192.168.1.1"),
	)

	// Request/Response objects
	type APIRequest struct {
		Method  string            `json:"method"`
		Path    string            `json:"path"`
		Headers map[string]string `json:"headers"`
		Body    json.RawMessage   `json:"body"`
	}

	req := APIRequest{
		Method: "POST",
		Path:   "/api/orders",
		Headers: map[string]string{
			"Content-Type":  "application/json",
			"Authorization": "Bearer [redacted]",
		},
		Body: json.RawMessage(`{"product_id":"prod_123","quantity":2}`),
	}

	log.Info("Incoming API request",
		slog.Any("request", req),
		slog.Duration("timeout", 30*time.Second),
	)

	// Metrics object
	type Metrics struct {
		RequestCount   int64   `json:"request_count"`
		ErrorRate      float64 `json:"error_rate"`
		AvgLatency     float64 `json:"avg_latency_ms"`
		ActiveSessions int     `json:"active_sessions"`
	}

	metrics := Metrics{
		RequestCount:   15234,
		ErrorRate:      0.02,
		AvgLatency:     45.3,
		ActiveSessions: 127,
	}

	log.Info("System metrics",
		slog.Any("metrics", metrics),
		slog.Time("timestamp", time.Now()),
	)

	output := buf.String()
	if !strings.Contains(output, "usr_abc123") {
		t.Errorf("Expected user ID in output")
	}
	if !strings.Contains(output, "/api/orders") {
		t.Errorf("Expected API path in output")
	}
	if !strings.Contains(output, "request_count") {
		t.Errorf("Expected metrics in output")
	}
}

func TestLogger_PrettyFormat(t *testing.T) {
	var buf bytes.Buffer
	log := logger.New(
		logger.WithOutput(&buf),
		logger.WithFormat("pretty"),
	)

	// Test that pretty format produces colored output
	log.Info("Testing pretty format")
	log.Warn("Warning in pretty format")
	log.Error("Error in pretty format")

	output := buf.String()
	// Pretty format includes ANSI color codes
	if !strings.Contains(output, "[32m") && !strings.Contains(output, "[33m") && !strings.Contains(output, "[31m") {
		t.Errorf("Expected ANSI color codes in pretty format output")
	}
}

func TestLogger_Context(t *testing.T) {
	var buf bytes.Buffer
	baseLogger := logger.New(
		logger.WithOutput(&buf),
		logger.WithFormat("json"),
	)

	// Simulate request context with logger
	ctx := context.Background()

	// Add request-scoped logger with request ID
	requestLogger := baseLogger.With(
		slog.String("request_id", "req_xyz789"),
		slog.String("user_id", "usr_456"),
	)
	ctx = logger.IntoContext(ctx, requestLogger)

	// Retrieve and use logger from context
	ctxLogger := logger.FromContext(ctx)
	ctxLogger.Info("Processing payment",
		slog.Float64("amount", 99.99),
		slog.String("currency", "USD"),
	)

	output := buf.String()
	if !strings.Contains(output, "req_xyz789") {
		t.Errorf("Expected request_id in context logger output")
	}
	if !strings.Contains(output, "usr_456") {
		t.Errorf("Expected user_id in context logger output")
	}
}

func TestLogger_Groups(t *testing.T) {
	var buf bytes.Buffer
	log := logger.New(
		logger.WithOutput(&buf),
		logger.WithFormat("json"),
	)

	// Using groups to organize related attributes
	log.Info("Order processed",
		slog.Group("order",
			slog.String("id", "ord_789"),
			slog.Float64("total", 249.99),
			slog.Int("items", 3),
		),
		slog.Group("customer",
			slog.String("id", "cust_123"),
			slog.String("email", "jane@example.com"),
			slog.String("tier", "premium"),
		),
		slog.Group("shipping",
			slog.String("method", "express"),
			slog.String("tracking", "1Z999AA1012345678"),
			slog.Time("estimated_delivery", time.Now().Add(48*time.Hour)),
		),
	)

	output := buf.String()
	var logEntry map[string]interface{}
	if err := json.Unmarshal([]byte(output), &logEntry); err != nil {
		t.Fatalf("Failed to parse JSON log output: %v", err)
	}

	// Check that groups are properly structured
	if order, ok := logEntry["order"].(map[string]interface{}); ok {
		if order["id"] != "ord_789" {
			t.Errorf("Expected order.id to be ord_789")
		}
	} else {
		t.Errorf("Expected order group in log output")
	}
}

func TestLogger_RealWorldScenarios(t *testing.T) {
	var buf bytes.Buffer
	log := logger.New(
		logger.WithOutput(&buf),
		logger.WithFormat("json"),
		logger.WithLevel(slog.LevelDebug),
	)

	// Also create a logger that outputs to stdout for visibility
	visibleLog := logger.New(
		logger.WithFormat("pretty"), // Use pretty format for better readability
		logger.WithLevel(slog.LevelDebug),
	)

	// Scenario 1: API endpoint hit
	startTime := time.Now()
	attrs := []any{
		slog.String("method", "GET"),
		slog.String("path", "/api/products"),
		slog.String("client_ip", "203.0.113.45"),
	}
	log.Info("API request received", attrs...)
	visibleLog.Info("API request received", attrs...)

	// Scenario 2: Database query
	attrs = []any{
		slog.String("query", "SELECT * FROM products WHERE category = $1"),
		slog.String("params", "[electronics]"),
		slog.Duration("duration", 23*time.Millisecond),
	}
	log.Debug("Executing database query", attrs...)
	visibleLog.Debug("Executing database query", attrs...)

	// Scenario 3: Cache operation
	attrs = []any{
		slog.String("key", "products:electronics"),
		slog.Bool("hit", false),
		slog.String("action", "fetch_from_db"),
	}
	log.Debug("Cache lookup", attrs...)
	visibleLog.Debug("Cache lookup", attrs...)

	// Scenario 4: Business logic
	attrs = []any{
		slog.String("order_id", "ord_555"),
		slog.String("code", "SUMMER20"),
		slog.Float64("discount_amount", 49.98),
		slog.Float64("original_total", 249.90),
		slog.Float64("final_total", 199.92),
	}
	log.Info("Discount applied", attrs...)
	visibleLog.Info("Discount applied", attrs...)

	// Scenario 5: External service call
	attrs = []any{
		slog.String("gateway", "stripe"),
		slog.Duration("response_time", 5*time.Second),
		slog.String("transaction_id", "txn_abc123"),
		slog.Bool("succeeded", true),
	}
	log.Warn("Payment gateway slow response", attrs...)
	visibleLog.Warn("Payment gateway slow response", attrs...)

	// Scenario 6: Background job
	attrs = []any{
		slog.String("job_type", "email_notification"),
		slog.Int("recipients", 1523),
		slog.Duration("processing_time", 2*time.Minute+30*time.Second),
		slog.Int("failed", 3),
		slog.Float64("success_rate", 99.8),
	}
	log.Info("Background job completed", attrs...)
	visibleLog.Info("Background job completed", attrs...)

	// Scenario 7: Security event
	attrs = []any{
		slog.String("username", "admin@example.com"),
		slog.String("ip_address", "198.51.100.42"),
		slog.Int("attempt_number", 3),
		slog.String("action", "account_locked"),
	}
	log.Warn("Failed login attempt", attrs...)
	visibleLog.Warn("Failed login attempt", attrs...)

	// Scenario 8: Request completed
	attrs = []any{
		slog.String("method", "GET"),
		slog.String("path", "/api/products"),
		slog.Int("status", 200),
		slog.Duration("latency", time.Since(startTime)),
		slog.Int("response_size", 4523),
	}
	log.Info("API request completed", attrs...)
	visibleLog.Info("API request completed", attrs...)

	output := buf.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")

	// Should have logged all scenarios
	if len(lines) < 8 {
		t.Errorf("Expected at least 8 log entries, got %d", len(lines))
	}

	// Verify some key scenarios are present
	if !strings.Contains(output, "API request received") {
		t.Errorf("Missing API request log")
	}
	if !strings.Contains(output, "Payment gateway slow response") {
		t.Errorf("Missing payment gateway log")
	}
	if !strings.Contains(output, "Failed login attempt") {
		t.Errorf("Missing security event log")
	}
}

func ExampleLogger_usage() {
	// Create a logger for your application
	log := logger.New(
		logger.WithFormat("json"),
		logger.WithLevel(slog.LevelInfo),
	)

	// Simple logging
	log.Info("Application started")

	// Log with attributes
	log.Info("User signed in",
		slog.String("user_id", "usr_123"),
		slog.String("method", "oauth"),
	)

	// Error logging
	err := fmt.Errorf("database connection timeout")
	log.Error("Failed to fetch user data",
		slog.String("error", err.Error()),
		slog.String("user_id", "usr_123"),
	)

	// Using groups for structured data
	log.Info("Order placed",
		slog.Group("order",
			slog.String("id", "ord_456"),
			slog.Float64("total", 99.99),
		),
		slog.Group("customer",
			slog.String("id", "cust_789"),
			slog.String("email", "customer@example.com"),
		),
	)
}
