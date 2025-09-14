package httpserver

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    "log/slog"
    "github.com/mikko-kohtala/go-api/internal/config"
)

// minimal logger for tests
func testLogger() *slog.Logger { return slog.New(slog.NewTextHandler(testDiscard{}, nil)) }

type testDiscard struct{}
func (testDiscard) Write(p []byte) (int, error) { return len(p), nil }

func TestBodyLimit_EchoTooLarge(t *testing.T) {
    cfg := &config.Config{
        Env:              "test",
        Port:             0,
        RequestTimeout:   0, // not used in test server
        BodyLimitBytes:   10,
        CORSAllowedOrigins: []string{"*"},
        CORSAllowedMethods: []string{"GET","POST"},
        CORSAllowedHeaders: []string{"Content-Type"},
        RateLimitEnabled: false,
        RateLimit:        1,
        RateLimitPeriod:  "1m",
        CompressionLevel: 5,
    }
    // Avoid zero timeout by setting small positive duration
    if cfg.RequestTimeout <= 0 { cfg.RequestTimeout = 1 }
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
        Env:              "test",
        Port:             0,
        RequestTimeout:   1,
        BodyLimitBytes:   1048576,
        CORSAllowedOrigins: []string{"*"},
        CORSAllowedMethods: []string{"GET"},
        CORSAllowedHeaders: []string{"*"},
        RateLimitEnabled: true,
        RateLimit:        1,
        RateLimitPeriod:  "10s",
        CompressionLevel: 5,
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

func TestRequestID_PropagatesToErrorBody(t *testing.T) {
    cfg := &config.Config{
        Env:                "test",
        Port:               0,
        RequestTimeout:     1,
        BodyLimitBytes:     1048576,
        CORSAllowedOrigins: []string{"*"},
        CORSAllowedMethods: []string{"GET","POST"},
        CORSAllowedHeaders: []string{"*"},
        RateLimitEnabled:   false,
        CompressionLevel:   5,
    }
    h := NewRouter(cfg, testLogger())

    rr := httptest.NewRecorder()
    // invalid JSON to trigger decode error path and response.Error usage
    req := httptest.NewRequest(http.MethodPost, "/api/v1/echo", bytes.NewBufferString(`{"message":`))
    req.Header.Set("Content-Type", "application/json")
    h.ServeHTTP(rr, req)
    if rr.Code != http.StatusBadRequest {
        t.Fatalf("expected 400, got %d", rr.Code)
    }
    rid := rr.Header().Get("X-Request-ID")
    if rid == "" {
        t.Fatalf("expected X-Request-ID header to be set")
    }
    var body struct {
        RequestID string `json:"request_id"`
    }
    if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
        t.Fatalf("invalid json body: %v", err)
    }
    if body.RequestID != rid {
        t.Fatalf("request_id mismatch: body=%q header=%q", body.RequestID, rid)
    }
}

func TestSecurityHeadersPresent(t *testing.T) {
    cfg := &config.Config{
        Env:                "test",
        Port:               0,
        RequestTimeout:     1,
        BodyLimitBytes:     1048576,
        CORSAllowedOrigins: []string{"*"},
        CORSAllowedMethods: []string{"GET"},
        CORSAllowedHeaders: []string{"*"},
        RateLimitEnabled:   false,
        CompressionLevel:   5,
    }
    h := NewRouter(cfg, testLogger())
    rr := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodGet, "/api/v1/ping", nil)
    h.ServeHTTP(rr, req)
    if rr.Header().Get("X-Content-Type-Options") != "nosniff" {
        t.Fatalf("missing X-Content-Type-Options header")
    }
    if rr.Header().Get("X-Frame-Options") != "DENY" {
        t.Fatalf("missing X-Frame-Options header")
    }
    if rr.Header().Get("Referrer-Policy") != "no-referrer" {
        t.Fatalf("missing Referrer-Policy header")
    }
}
