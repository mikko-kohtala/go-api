package httpserver

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequestID_FromHeaderAndEchoed(t *testing.T) {
	h := RequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// ensure context has ID and response header is set
		rid := GetRequestID(r.Context())
		if rid == "" {
			t.Fatalf("expected request id in context")
		}
		if got := w.Header().Get("X-Request-ID"); got != rid {
			t.Fatalf("expected response header to equal context id; got %q want %q", got, rid)
		}
		_, _ = io.WriteString(w, rid)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("X-Request-ID", "frontend-123")
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	if rr.Body.String() != "frontend-123" {
		t.Fatalf("unexpected body: %s", rr.Body.String())
	}
}
