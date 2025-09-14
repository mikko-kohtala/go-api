package httpserver

import (
    "crypto/rand"
    "encoding/hex"
    "net/http"
    "regexp"
    "strings"

    "github.com/mikko-kohtala/go-api/internal/requestid"
)

// RequestID middleware trusts an incoming X-Request-ID or X-Correlation-ID header
// from the client. If absent, it generates a secure random ID. The chosen ID is
// set on the response header and stored in the request context.
func RequestID(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        rid := pickRequestID(r)
        w.Header().Set(requestid.HeaderRequestID, rid)
        ctx := requestid.IntoContext(r.Context(), rid)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

var requestIDPattern = regexp.MustCompile(`^[A-Za-z0-9_.:-]{1,128}$`)

func pickRequestID(r *http.Request) string {
    rid := strings.TrimSpace(r.Header.Get(requestid.HeaderRequestID))
    if rid == "" {
        rid = strings.TrimSpace(r.Header.Get(requestid.HeaderCorrelationID))
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
