package httpserver

import (
	"context"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/mikko-kohtala/go-api/internal/config"
)

func TestGracefulShutdownCompletesInFlightRequests(t *testing.T) {
	cfg := &config.Config{
		Env:                "development",
		Port:               0,
		RequestTimeout:     2 * time.Second,
		BodyLimitBytes:     1048576,
		CORSAllowedOrigins: []string{"*"},
		CORSAllowedMethods: []string{"GET"},
		CORSAllowedHeaders: []string{"*"},
		RateLimitEnabled:   false,
		RateLimit:          1,
		RateLimitPeriod:    "1m",
		CompressionLevel:   5,
	}

	handler := NewRouter(cfg, testLogger())

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}

	srv := &http.Server{
		Handler: handler,
	}

	go func() {
		if err := srv.Serve(ln); err != nil && err != http.ErrServerClosed {
			t.Errorf("serve returned error: %v", err)
		}
	}()

	baseURL := "http://" + ln.Addr().String()
	done := make(chan struct{})

	go func() {
		resp, err := http.Get(baseURL + "/test/sleep?duration_ms=200")
		if err != nil {
			t.Errorf("request failed: %v", err)
		} else {
			resp.Body.Close()
		}
		close(done)
	}()

	time.Sleep(50 * time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		t.Fatalf("shutdown failed: %v", err)
	}

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("in-flight request did not finish before timeout")
	}

	if _, err := http.Get(baseURL + "/healthz"); err == nil {
		t.Fatalf("expected server to be stopped")
	}
}
