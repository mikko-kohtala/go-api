package handlers

import (
    "bytes"
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestEcho(t *testing.T) {
    rr := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodPost, "/api/v1/echo", bytes.NewBufferString(`{"message":"x"}`))
    req.Header.Set("Content-Type", "application/json")
    Echo(rr, req)
    if rr.Code != http.StatusOK {
        t.Fatalf("expected 200, got %d", rr.Code)
    }
    if got := rr.Body.String(); got != "{\"message\":\"x\"}\n" {
        t.Fatalf("unexpected body: %q", got)
    }
}

