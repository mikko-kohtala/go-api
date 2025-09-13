package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Logger returns a gin.HandlerFunc for logging requests with request ID
func Logger(logger *zap.Logger) gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// Extract request ID from context
		requestID := ""
		if param.Keys != nil {
			if id, exists := param.Keys[RequestIDKey]; exists {
				if idStr, ok := id.(string); ok {
					requestID = idStr
				}
			}
		}

		logger.Info("HTTP Request",
			zap.String("request_id", requestID),
			zap.String("method", param.Method),
			zap.String("path", param.Path),
			zap.Int("status", param.StatusCode),
			zap.Duration("latency", param.Latency),
			zap.String("client_ip", param.ClientIP),
			zap.String("user_agent", param.Request.UserAgent()),
			zap.Time("timestamp", param.TimeStamp),
		)
		return ""
	})
}