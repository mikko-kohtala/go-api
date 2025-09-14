package httpserver

import (
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"github.com/mikko-kohtala/go-api/internal/metrics"
	"github.com/mikko-kohtala/go-api/internal/tracing"
)

// MetricsMiddleware records Prometheus metrics for HTTP requests
func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		// Increment active connections
		metrics.IncrementActiveConnections()
		defer metrics.DecrementActiveConnections()

		// Wrap response writer to capture status and size
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		
		// Get request size
		var requestSize int64
		if r.ContentLength > 0 {
			requestSize = r.ContentLength
		}

		// Process request
		next.ServeHTTP(ww, r)

		// Record metrics
		duration := time.Since(start)
		method := r.Method
		path := sanitizePath(r.URL.Path)

		metrics.RecordHTTPRequest(method, path, ww.Status(), duration, requestSize, int64(ww.BytesWritten()))
	})
}

// TracingMiddleware adds OpenTelemetry tracing to HTTP requests
func TracingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract trace context from headers
		propagator := otel.GetTextMapPropagator()
		ctx := propagator.Extract(r.Context(), propagation.HeaderCarrier(r.Header))
		
		// Start span
		spanName := r.Method + " " + r.URL.Path
		ctx, span := tracing.StartSpan(ctx, spanName,
			trace.WithSpanKind(trace.SpanKindServer),
			trace.WithAttributes(
				attribute.String("http.method", r.Method),
				attribute.String("http.url", r.URL.String()),
				attribute.String("http.user_agent", r.UserAgent()),
				attribute.String("http.remote_addr", r.RemoteAddr),
			),
		)
		defer span.End()

		// Add request ID to span
		if requestID := GetRequestID(ctx); requestID != "" {
			span.SetAttributes(attribute.String("request.id", requestID))
		}

		// Wrap response writer to capture status
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		
		// Process request with traced context
		next.ServeHTTP(ww, r.WithContext(ctx))

		// Record span attributes
		span.SetAttributes(
			attribute.Int("http.status_code", ww.Status()),
			attribute.Int64("http.response.size", int64(ww.BytesWritten())),
		)

		// Set trace context in response headers
		propagator.Inject(ctx, propagation.HeaderCarrier(ww.Header()))
	})
}

// sanitizePath removes dynamic path segments for better metrics grouping
func sanitizePath(path string) string {
	// Replace common dynamic segments
	path = strings.ReplaceAll(path, "/api/v1/", "/api/v1/")
	
	// Remove UUIDs, IDs, etc. (basic pattern matching)
	// This is a simple implementation - you might want more sophisticated path sanitization
	if strings.Contains(path, "/") {
		parts := strings.Split(path, "/")
		for i, part := range parts {
			// Replace UUID-like patterns
			if len(part) == 36 && strings.Contains(part, "-") {
				parts[i] = "{id}"
			}
			// Replace numeric IDs
			if len(part) > 0 && isNumeric(part) {
				parts[i] = "{id}"
			}
		}
		path = strings.Join(parts, "/")
	}
	
	return path
}

// isNumeric checks if a string contains only digits
func isNumeric(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return len(s) > 0
}