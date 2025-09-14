package httpserver

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"regexp"
	"strings"
)

type ctxKey string

const requestIDKey ctxKey = "request_id"

// RequestID middleware trusts an incoming X-Request-ID or X-Correlation-ID header
// from the client. If absent, it generates a secure random ID.
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rid := pickRequestID(r)
		w.Header().Set("X-Request-ID", rid)
		ctx := context.WithValue(r.Context(), requestIDKey, rid)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

var requestIDPattern = regexp.MustCompile(`^[A-Za-z0-9_.:-]{1,128}$`)

func pickRequestID(r *http.Request) string {
	rid := strings.TrimSpace(r.Header.Get("X-Request-ID"))
	if rid == "" {
		rid = strings.TrimSpace(r.Header.Get("X-Correlation-ID"))
	}
	if rid != "" && requestIDPattern.MatchString(rid) {
		return rid
	}
	// fallback: random 16 bytes hex
	var b [16]byte
	if _, err := rand.Read(b[:]); err == nil {
		return hex.EncodeToString(b[:])
	}
	return "unknown"
}

// GetRequestID returns the request id from context, if set.
func GetRequestID(ctx context.Context) string {
	if v := ctx.Value(requestIDKey); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}
