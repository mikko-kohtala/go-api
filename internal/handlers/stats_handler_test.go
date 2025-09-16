package handlers

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mikko-kohtala/go-api/internal/services"
)

func TestStatsHandler_GetSystemStats(t *testing.T) {
	handler := NewStatsHandler(services.NewStatsService(), slog.New(slog.NewTextHandler(io.Discard, nil)))

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/stats/system", nil)
	handler.GetSystemStats(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}

func TestStatsHandler_GetAPIStats(t *testing.T) {
	handler := NewStatsHandler(services.NewStatsService(), slog.New(slog.NewTextHandler(io.Discard, nil)))

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/stats/api", nil)
	handler.GetAPIStats(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}
