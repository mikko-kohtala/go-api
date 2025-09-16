package httpserver

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/mikko-kohtala/go-api/internal/metrics"
)

// MetricsMiddleware collects HTTP metrics for Prometheus
func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap the ResponseWriter to capture status code and response size
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		// Record request size
		requestSize := float64(r.ContentLength)
		if requestSize > 0 {
			metrics.HttpRequestSize.WithLabelValues(r.Method, r.URL.Path).Observe(requestSize)
		}

		// Process request
		next.ServeHTTP(ww, r)

		// Record metrics after request is processed
		duration := time.Since(start).Seconds()
		status := strconv.Itoa(ww.Status())

		// Normalize path for metrics to avoid high cardinality
		path := normalizePath(r.URL.Path)

		// Record request metrics
		metrics.HttpRequestsTotal.WithLabelValues(r.Method, path, status).Inc()
		metrics.HttpRequestDuration.WithLabelValues(r.Method, path).Observe(duration)

		// Record response size
		responseSize := float64(ww.BytesWritten())
		if responseSize > 0 {
			metrics.HttpResponseSize.WithLabelValues(r.Method, path).Observe(responseSize)
		}
	})
}

// normalizePath normalizes URL paths for metrics to avoid high cardinality
// It replaces dynamic segments with placeholders
func normalizePath(path string) string {
	// Common patterns to normalize
	normalizations := map[string]string{
		"/api/v1/users/": "/api/v1/users/{id}",
		"/api/v1/stats/": "/api/v1/stats/{type}",
	}

	for prefix, normalized := range normalizations {
		if len(path) > len(prefix) && path[:len(prefix)] == prefix {
			return normalized
		}
	}

	return path
}
