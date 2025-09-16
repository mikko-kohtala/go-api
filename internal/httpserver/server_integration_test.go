package httpserver

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/mikko-kohtala/go-api/internal/config"
	"log/slog"
)

// minimal logger for tests
func testLogger() *slog.Logger { return slog.New(slog.NewTextHandler(testDiscard{}, nil)) }

type testDiscard struct{}

func (testDiscard) Write(p []byte) (int, error) { return len(p), nil }

func TestBodyLimit_EchoTooLarge(t *testing.T) {
	cfg := &config.Config{
		Env:                "test",
		Port:               0,
		RequestTimeout:     time.Second,
		BodyLimitBytes:     10,
		CORSAllowedOrigins: []string{"*"},
		CORSAllowedMethods: []string{"GET", "POST"},
		CORSAllowedHeaders: []string{"Content-Type"},
		RateLimitEnabled:   false,
		RateLimit:          1,
		RateLimitPeriod:    "1m",
		CompressionLevel:   5,
	}
	h := NewRouter(cfg, testLogger())

	// Body > 10 bytes triggers MaxBytesReader error during JSON decode
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/echo", bytes.NewBufferString(`{"message":"0123456789ABC"}`))
	req.Header.Set("Content-Type", "application/json")
	h.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for too large body, got %d", rr.Code)
	}
}

func TestHealth_NotRateLimited(t *testing.T) {
	cfg := &config.Config{
		Env:                "test",
		Port:               0,
		RequestTimeout:     time.Second,
		BodyLimitBytes:     1048576,
		CORSAllowedOrigins: []string{"*"},
		CORSAllowedMethods: []string{"GET"},
		CORSAllowedHeaders: []string{"*"},
		RateLimitEnabled:   true,
		RateLimit:          1,
		RateLimitPeriod:    "10s",
		CompressionLevel:   5,
	}
	h := NewRouter(cfg, testLogger())
	// Call /healthz twice quickly; should not be limited
	for i := 0; i < 2; i++ {
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
		h.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d on iteration %d", rr.Code, i)
		}
	}
}

func TestTestRoutesDisabledInProduction(t *testing.T) {
	cfg := &config.Config{
		Env:                "production",
		Port:               0,
		RequestTimeout:     time.Second,
		BodyLimitBytes:     1048576,
		CORSAllowedOrigins: []string{"*"},
		CORSAllowedMethods: []string{"GET"},
		CORSAllowedHeaders: []string{"*"},
		RateLimitEnabled:   false,
		RateLimit:          1,
		RateLimitPeriod:    "1m",
		CompressionLevel:   5,
	}

	h := NewRouter(cfg, testLogger())
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test/logs", nil)
	h.ServeHTTP(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected 404 when test routes disabled, got %d", rr.Code)
	}
}

func TestMetricsEndpointAvailable(t *testing.T) {
	cfg := &config.Config{
		Env:                "test",
		Port:               0,
		RequestTimeout:     time.Second,
		BodyLimitBytes:     1048576,
		CORSAllowedOrigins: []string{"*"},
		CORSAllowedMethods: []string{"GET"},
		CORSAllowedHeaders: []string{"*"},
		RateLimitEnabled:   false,
		RateLimit:          1,
		RateLimitPeriod:    "1m",
		CompressionLevel:   5,
	}

	h := NewRouter(cfg, testLogger())
	server := httptest.NewServer(h)
	defer server.Close()

	resp, err := http.Get(server.URL + "/metrics")
	if err != nil {
		t.Fatalf("GET /metrics failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 from /metrics, got %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed reading metrics body: %v", err)
	}
	if !bytes.Contains(body, []byte("api_requests_total")) {
		t.Fatalf("expected metrics output to contain api_requests_total, got %s", string(body))
	}
}
