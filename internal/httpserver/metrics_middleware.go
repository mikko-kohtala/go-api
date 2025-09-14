package httpserver

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/mikko-kohtala/go-api/internal/telemetry"
)

// MetricsMiddleware records HTTP metrics
func MetricsMiddleware(metrics *telemetry.Metrics) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Track active requests
			metrics.ActiveRequests.Inc()
			defer metrics.ActiveRequests.Dec()

			// Wrap response writer to capture status code
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			// Process request
			next.ServeHTTP(ww, r)

			// Record metrics
			duration := time.Since(start).Seconds()
			status := strconv.Itoa(ww.Status())
			path := r.URL.Path
			method := r.Method

			metrics.RequestsTotal.WithLabelValues(method, path, status).Inc()
			metrics.RequestDuration.WithLabelValues(method, path, status).Observe(duration)
		})
	}
}