package httpserver

import "net/http"

// SecurityHeaders adds a few simple, safe security headers.
// Kept minimal to avoid breaking Swagger UI or other tooling.
func SecurityHeaders() func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            w.Header().Set("X-Content-Type-Options", "nosniff")
            w.Header().Set("X-Frame-Options", "DENY")
            w.Header().Set("Referrer-Policy", "no-referrer")
            next.ServeHTTP(w, r)
        })
    }
}

