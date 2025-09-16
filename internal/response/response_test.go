package response

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	pkglogger "github.com/mikko-kohtala/go-api/pkg/logger"
)

func TestErrorUsesContextRequestID(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req = req.WithContext(pkglogger.WithRequestID(req.Context(), "generated-123"))

	Error(rr, req, http.StatusBadRequest, "invalid_request", "Invalid", nil)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}

	var resp ErrorResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if resp.RequestID != "generated-123" {
		t.Fatalf("expected request_id to be generated-123, got %q", resp.RequestID)
	}
}

func TestErrorPrefersHeaderRequestID(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Request-ID", "client-abc")
	req = req.WithContext(pkglogger.WithRequestID(req.Context(), "generated-123"))

	Error(rr, req, http.StatusBadRequest, "invalid_request", "Invalid", nil)

	var resp ErrorResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if resp.RequestID != "client-abc" {
		t.Fatalf("expected request_id to prefer header value, got %q", resp.RequestID)
	}
}

func TestJSONSkipsWhenContextCanceled(t *testing.T) {
	rr := &recordingResponseWriter{ResponseWriter: httptest.NewRecorder()}
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	ctx, cancel := context.WithCancel(req.Context())
	cancel()
	req = req.WithContext(ctx)

	JSON(rr, req, http.StatusOK, map[string]string{"status": "ok"})

	if rr.writeHeaderCalls != 0 {
		t.Fatalf("expected WriteHeader not to be called, got %d", rr.writeHeaderCalls)
	}
}

type recordingResponseWriter struct {
	http.ResponseWriter
	writeHeaderCalls int
}

func (rw *recordingResponseWriter) WriteHeader(statusCode int) {
	rw.writeHeaderCalls++
	rw.ResponseWriter.WriteHeader(statusCode)
}
