package httpserver

import (
    "log/slog"
    "net/http"
    "time"

    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
    "github.com/mikko-kohtala/go-api/internal/logging"
    "github.com/mikko-kohtala/go-api/internal/requestid"
)

// LoggingMiddleware logs basic request/response details using slog JSON.
func LoggingMiddleware(logger *slog.Logger) func(next http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        fn := func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()
            ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
            // Create request-scoped logger with request_id if available
            rid := requestid.FromContext(r.Context())
            reqLogger := logger
            if rid != "" {
                reqLogger = logger.With(slog.String("request_id", rid))
            }
            // store logger in context for handlers to use if desired
            ctx := logging.IntoContext(r.Context(), reqLogger)
            next.ServeHTTP(ww, r.WithContext(ctx))
            duration := time.Since(start)
            // capture route pattern to reduce cardinality (e.g., /api/v1/echo)
            route := ""
            if rctx := chi.RouteContext(r.Context()); rctx != nil {
                route = rctx.RoutePattern()
            }
            reqLogger.Info("request",
                slog.String("remote_ip", r.RemoteAddr),
                slog.String("method", r.Method),
                slog.String("path", r.URL.Path),
                slog.String("route", route),
                slog.Int("status", ww.Status()),
                slog.Int("bytes", ww.BytesWritten()),
                slog.String("duration", duration.String()),
                slog.String("user_agent", r.UserAgent()),
            )
        }
        return http.HandlerFunc(fn)
    }
}
