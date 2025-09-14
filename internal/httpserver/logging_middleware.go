package httpserver

import (
    "fmt"
    "log/slog"
    "net/http"
    "os"
    "time"

    "github.com/go-chi/chi/v5/middleware"
    "github.com/mikko-kohtala/go-api/internal/logging"
)

// LoggingMiddleware logs basic request/response details using slog JSON.
func LoggingMiddleware(logger *slog.Logger) func(next http.Handler) http.Handler {
    // Add HTTP component to logger
    logger = logger.With(slog.String("component", "HTTP"))
    return func(next http.Handler) http.Handler {
        fn := func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()
            ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
            // Create request-scoped logger with request_id if available
            rid := GetRequestID(r.Context())
            reqLogger := logger
            if rid != "" {
                reqLogger = logger.With(slog.String("request_id", rid))
            }

            // Check if pretty logging is enabled
            prettyLogs := os.Getenv("PRETTY_LOGS") == "true"

            // Log incoming request (with arrow indicator in pretty handler)
            if prettyLogs {
                reqLogger.Info(fmt.Sprintf("%s %s", r.Method, r.URL.Path))
            }

            // store logger in context for handlers to use if desired
            ctx := logging.IntoContext(r.Context(), reqLogger)
            next.ServeHTTP(ww, r.WithContext(ctx))
            duration := time.Since(start)

            if prettyLogs {
                // Log the completed request with status and latency
                reqLogger.Info(fmt.Sprintf("%s %s", r.Method, r.URL.Path),
                    slog.Int("status", ww.Status()),
                    slog.Duration("latency", duration),
                )
            } else {
                // Full logging for production/JSON logs
                reqLogger.Info("request",
                    slog.String("remote_ip", r.RemoteAddr),
                    slog.String("method", r.Method),
                    slog.String("path", r.URL.Path),
                    slog.Int("status", ww.Status()),
                    slog.Int("bytes", ww.BytesWritten()),
                    slog.String("duration", duration.String()),
                    slog.String("user_agent", r.UserAgent()),
                )
            }
        }
        return http.HandlerFunc(fn)
    }
}
