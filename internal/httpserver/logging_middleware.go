package httpserver

import (
    "log/slog"
    "net/http"
    "time"

    "github.com/go-chi/chi/v5/middleware"
)

// LoggingMiddleware logs basic request/response details using slog JSON.
func LoggingMiddleware(logger *slog.Logger) func(next http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        fn := func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()
            ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
            next.ServeHTTP(ww, r)
            duration := time.Since(start)

            logger.Info("request",
                slog.String("request_id", middleware.GetReqID(r.Context())),
                slog.String("remote_ip", r.RemoteAddr),
                slog.String("method", r.Method),
                slog.String("path", r.URL.Path),
                slog.Int("status", ww.Status()),
                slog.Int("bytes", ww.BytesWritten()),
                slog.String("duration", duration.String()),
                slog.String("user_agent", r.UserAgent()),
            )
        }
        return http.HandlerFunc(fn)
    }
}

