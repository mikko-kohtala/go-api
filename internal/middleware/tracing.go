package middleware

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	RequestIDKey = "request_id"
	StartTimeKey = "start_time"
)

// TracingMiddleware creates a middleware for request tracing
func TracingMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate or extract request ID
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Set request ID in context
		c.Set(RequestIDKey, requestID)
		c.Set(StartTimeKey, time.Now())

		// Add request ID to response headers
		c.Header("X-Request-ID", requestID)

		// Create a new context with request ID for downstream services
		ctx := context.WithValue(c.Request.Context(), RequestIDKey, requestID)
		c.Request = c.Request.WithContext(ctx)

		// Log request start
		logger.Info("Request started",
			zap.String("request_id", requestID),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("client_ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
			zap.String("referer", c.Request.Referer()),
		)

		// Process request
		c.Next()

		// Log request completion
		startTime, _ := c.Get(StartTimeKey)
		duration := time.Since(startTime.(time.Time))

		logger.Info("Request completed",
			zap.String("request_id", requestID),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("duration", duration),
			zap.Int("response_size", c.Writer.Size()),
		)
	}
}

// GetRequestID extracts request ID from gin context
func GetRequestID(c *gin.Context) string {
	if requestID, exists := c.Get(RequestIDKey); exists {
		if id, ok := requestID.(string); ok {
			return id
		}
	}
	return ""
}

// GetRequestDuration calculates request duration from start time
func GetRequestDuration(c *gin.Context) time.Duration {
	if startTime, exists := c.Get(StartTimeKey); exists {
		if start, ok := startTime.(time.Time); ok {
			return time.Since(start)
		}
	}
	return 0
}