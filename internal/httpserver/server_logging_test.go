package httpserver

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/mikko-kohtala/go-api/internal/config"
	"github.com/mikko-kohtala/go-api/pkg/logger"
	"log/slog"
)

// Test flag to control log visibility
var showLogs = flag.Bool("show-logs", false, "Show logs during tests")

// TestLogger_APIRequestResponse tests logging of real API requests and responses
func TestLogger_APIRequestResponse(t *testing.T) {
	// Parse flags if not already parsed
	if !flag.Parsed() {
		flag.Parse()
	}

	// Create logger based on flag
	var log *slog.Logger
	if *showLogs {
		// Use pretty format to stdout for visibility
		log = logger.New(
			logger.WithFormat("pretty"),
			logger.WithLevel(slog.LevelDebug),
		)
		t.Log("Logs enabled - use -show-logs flag to control visibility")
	} else {
		// Use buffer for testing without output
		var buf bytes.Buffer
		log = logger.New(
			logger.WithOutput(&buf),
			logger.WithFormat("json"),
			logger.WithLevel(slog.LevelDebug),
		)
	}

	cfg := &config.Config{
		Env:                "test",
		Port:               8080,
		RequestTimeout:     30,
		BodyLimitBytes:     1048576,
		CORSAllowedOrigins: []string{"*"},
		CORSAllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		CORSAllowedHeaders: []string{"Content-Type", "Authorization"},
		RateLimitEnabled:   true,
		RateLimit:          100,
		RateLimitPeriod:    "1m",
		CompressionLevel:   5,
	}

	// Create router with our logger
	handler := NewRouter(cfg, log)
	server := httptest.NewServer(handler)
	defer server.Close()

	// Run test scenarios
	t.Run("GET_HealthCheck", func(t *testing.T) {
		testAPICall(t, log, server.URL, "GET", "/healthz", nil, http.StatusOK)
	})

	t.Run("POST_Echo", func(t *testing.T) {
		payload := map[string]string{"message": "Hello, World!"}
		testAPICall(t, log, server.URL, "POST", "/api/v1/echo", payload, http.StatusOK)
	})

	t.Run("GET_Example", func(t *testing.T) {
		testAPICall(t, log, server.URL, "GET", "/api/v1/example", nil, http.StatusNotFound) // No example endpoint exists
	})

	t.Run("POST_LargePayload", func(t *testing.T) {
		// Create a moderately large payload (within body limit)
		largeData := make([]byte, 1000)
		for i := range largeData {
			largeData[i] = 'A'
		}
		payload := map[string]string{"message": string(largeData)}
		testAPICall(t, log, server.URL, "POST", "/api/v1/echo", payload, http.StatusOK)
	})

	t.Run("GET_NotFound", func(t *testing.T) {
		testAPICall(t, log, server.URL, "GET", "/api/v1/nonexistent", nil, http.StatusNotFound)
	})

	t.Run("POST_InvalidJSON", func(t *testing.T) {
		// Send invalid JSON
		req, _ := http.NewRequest("POST", server.URL+"/api/v1/echo", bytes.NewBufferString("{invalid json"))
		req.Header.Set("Content-Type", "application/json")

		startTime := time.Now()
		log.Info("API request initiated",
			slog.String("method", req.Method),
			slog.String("url", req.URL.String()),
			slog.String("test", "POST_InvalidJSON"),
		)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Error("Request failed",
				slog.String("error", err.Error()),
				slog.Duration("duration", time.Since(startTime)),
			)
			t.Fatalf("Request failed: %v", err)
		}
		defer resp.Body.Close()

		log.Info("API response received",
			slog.Int("status", resp.StatusCode),
			slog.String("status_text", resp.Status),
			slog.Duration("latency", time.Since(startTime)),
		)

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("Concurrent_Requests", func(t *testing.T) {
		// Test multiple concurrent requests
		concurrentCount := 5
		done := make(chan bool, concurrentCount)
		errors := make(chan error, concurrentCount)

		for i := 0; i < concurrentCount; i++ {
			go func(requestNum int) {
				defer func() { done <- true }()

				payload := map[string]interface{}{
					"message": fmt.Sprintf("Concurrent request %d", requestNum),
				}

				// Create a sub-test to avoid race conditions with t.Helper()
				body, _ := json.Marshal(payload)
				req, err := http.NewRequest("POST", server.URL+"/api/v1/echo", bytes.NewBuffer(body))
				if err != nil {
					errors <- err
					return
				}

				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("User-Agent", "TestClient/1.0")

				resp, err := http.DefaultClient.Do(req)
				if err != nil {
					errors <- err
					return
				}
				defer resp.Body.Close()

				if resp.StatusCode != http.StatusOK {
					respBody, _ := io.ReadAll(resp.Body)
					errors <- fmt.Errorf("request %d: expected status 200, got %d. Body: %s", requestNum, resp.StatusCode, string(respBody))
				}
			}(i)
		}

		// Wait for all requests to complete
		for i := 0; i < concurrentCount; i++ {
			<-done
		}

		// Check for errors
		close(errors)
		for err := range errors {
			if err != nil {
				t.Error(err)
			}
		}
	})

	// Log summary
	if *showLogs {
		log.Info("Test scenarios completed",
			slog.Int("total_scenarios", 7),
			slog.String("server_url", server.URL),
		)
	}
}

// Helper function to test API calls with logging
func testAPICall(t *testing.T, log *slog.Logger, baseURL, method, path string, payload interface{}, expectedStatus int) {
	t.Helper()

	var body io.Reader
	var requestBody []byte
	if payload != nil {
		var err error
		requestBody, err = json.Marshal(payload)
		if err != nil {
			t.Fatalf("Failed to marshal payload: %v", err)
		}
		body = bytes.NewBuffer(requestBody)
	}

	req, err := http.NewRequest(method, baseURL+path, body)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// Add common headers
	req.Header.Set("User-Agent", "TestClient/1.0")
	req.Header.Set("X-Request-ID", "test-"+time.Now().Format("20060102-150405"))

	// Log request
	startTime := time.Now()
	log.Info("API request",
		slog.Group("request",
			slog.String("method", method),
			slog.String("path", path),
			slog.String("url", req.URL.String()),
			slog.Any("headers", req.Header),
			slog.Int("body_size", len(requestBody)),
		),
	)

	if len(requestBody) > 0 && len(requestBody) < 1000 { // Only log small bodies
		log.Debug("Request body",
			slog.String("body", string(requestBody)),
		)
	}

	// Execute request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Error("Request failed",
			slog.String("error", err.Error()),
			slog.Duration("duration", time.Since(startTime)),
		)
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	// Read response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("Failed to read response body",
			slog.String("error", err.Error()),
		)
		t.Fatalf("Failed to read response: %v", err)
	}

	// Log response
	log.Info("API response",
		slog.Group("response",
			slog.Int("status_code", resp.StatusCode),
			slog.String("status", resp.Status),
			slog.Any("headers", resp.Header),
			slog.Int("body_size", len(responseBody)),
			slog.Duration("latency", time.Since(startTime)),
		),
	)

	if len(responseBody) > 0 && len(responseBody) < 1000 { // Only log small bodies
		log.Debug("Response body",
			slog.String("body", string(responseBody)),
		)
	}

	// Performance metrics
	if time.Since(startTime) > 100*time.Millisecond {
		log.Warn("Slow API response",
			slog.String("path", path),
			slog.Duration("latency", time.Since(startTime)),
		)
	}

	// Check status code
	if resp.StatusCode != expectedStatus {
		log.Error("Unexpected status code",
			slog.Int("expected", expectedStatus),
			slog.Int("actual", resp.StatusCode),
			slog.String("response", string(responseBody)),
		)
		t.Errorf("Expected status %d, got %d", expectedStatus, resp.StatusCode)
	}
}

// TestLogger_RealWorldScenarios tests various real-world logging scenarios
func TestLogger_RealWorldScenarios(t *testing.T) {
	// Parse flags if not already parsed
	if !flag.Parsed() {
		flag.Parse()
	}

	// Setup logger with optional output
	var log *slog.Logger
	var buf bytes.Buffer

	if *showLogs {
		log = logger.New(
			logger.WithFormat("pretty"),
			logger.WithLevel(slog.LevelDebug),
		)
	} else {
		log = logger.New(
			logger.WithOutput(&buf),
			logger.WithFormat("json"),
			logger.WithLevel(slog.LevelDebug),
		)
	}

	// Scenario 1: API endpoint hit with full request/response cycle
	startTime := time.Now()
	requestID := "req_" + time.Now().Format("20060102150405")

	log.Info("Incoming HTTP request",
		slog.Group("http",
			slog.String("method", "POST"),
			slog.String("path", "/api/v1/users"),
		),
		slog.String("request_id", requestID),
	)

	// Scenario 2: Authentication & Authorization
	log.Info("Authentication successful",
		slog.String("request_id", requestID),
		slog.String("user_id", "usr_abc123"),
		slog.String("auth_method", "jwt"),
		slog.String("token_exp", "2024-12-31T23:59:59Z"),
	)

	// Scenario 3: Request validation
	log.Debug("Request validation",
		slog.String("request_id", requestID),
		slog.Group("validation",
			slog.Bool("passed", true),
			slog.Duration("duration", 2*time.Millisecond),
		),
	)

	// Scenario 4: Database operations
	log.Debug("Database query",
		slog.String("request_id", requestID),
		slog.Group("db",
			slog.String("operation", "INSERT"),
			slog.String("table", "users"),
			slog.Duration("duration", 15*time.Millisecond),
			slog.Int("rows_affected", 1),
		),
	)

	// Scenario 5: Cache interaction
	log.Debug("Cache operation",
		slog.String("request_id", requestID),
		slog.Group("cache",
			slog.String("operation", "SET"),
			slog.String("key", "user:abc123"),
			slog.Duration("ttl", 1*time.Hour),
			slog.Bool("success", true),
		),
	)

	// Scenario 6: External API call
	log.Info("External API call",
		slog.String("request_id", requestID),
		slog.Group("external",
			slog.String("service", "email-service"),
			slog.String("endpoint", "https://api.email.com/send"),
			slog.Int("status_code", 200),
			slog.Duration("latency", 250*time.Millisecond),
		),
	)

	// Scenario 7: Business logic events
	log.Info("User created",
		slog.String("request_id", requestID),
		slog.Group("user",
			slog.String("id", "usr_abc123"),
			slog.String("email", "user@example.com"),
			slog.String("plan", "premium"),
		),
		slog.Group("metadata",
			slog.String("referrer", "google"),
			slog.String("campaign", "summer-2024"),
		),
	)

	// Scenario 8: Response sent
	responseTime := time.Since(startTime)
	log.Info("HTTP response sent",
		slog.String("request_id", requestID),
		slog.Group("response",
			slog.Int("status_code", 201),
			slog.String("content_type", "application/json"),
			slog.Int("body_size", 256),
			slog.Duration("total_latency", responseTime),
		),
		slog.Group("performance",
			slog.Duration("db_time", 15*time.Millisecond),
			slog.Duration("cache_time", 3*time.Millisecond),
			slog.Duration("external_api_time", 250*time.Millisecond),
		),
	)

	// Scenario 9: Error scenarios
	log.Error("Payment processing failed",
		slog.String("request_id", "req_error123"),
		slog.Group("error",
			slog.String("type", "payment_declined"),
			slog.String("message", "Insufficient funds"),
			slog.String("card_last4", "4242"),
		),
		slog.Group("retry",
			slog.Bool("will_retry", true),
			slog.Int("attempt", 1),
			slog.Duration("next_retry_in", 5*time.Second),
		),
	)

	// Scenario 10: Rate limiting
	log.Warn("Rate limit exceeded",
		slog.String("client_ip", "203.0.113.42"),
		slog.Group("rate_limit",
			slog.Int("limit", 100),
			slog.String("period", "1m"),
			slog.Int("current_count", 101),
			slog.Duration("reset_in", 45*time.Second),
		),
	)

	// Verify logs were generated (when not showing)
	if !*showLogs {
		output := buf.String()
		if output == "" {
			t.Error("Expected log output but got none")
		}
		// Count log lines
		lines := bytes.Count([]byte(output), []byte("\n"))
		if lines < 10 {
			t.Errorf("Expected at least 10 log entries, got %d", lines)
		}
	}
}

// TestMain allows running tests with custom flags
func TestMain(m *testing.M) {
	flag.Parse()
	os.Exit(m.Run())
}
