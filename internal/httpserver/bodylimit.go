package httpserver

import (
    "net/http"
)

// BodyLimit returns middleware that limits request body size using http.MaxBytesReader.
func BodyLimit(maxBytes int64) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            if maxBytes > 0 && r.Body != nil {
                r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
            }
            next.ServeHTTP(w, r)
        })
    }
}

