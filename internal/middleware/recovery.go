package middleware

import (
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Recovery returns a gin.HandlerFunc for recovering from panics
func Recovery(logger *zap.Logger) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			logger.Error("Panic recovered",
				zap.String("request_id", GetRequestID(c)),
				zap.String("error", err),
				zap.String("path", c.Request.URL.Path),
				zap.String("method", c.Request.Method),
				zap.String("client_ip", c.ClientIP()),
				zap.String("user_agent", c.Request.UserAgent()),
				zap.String("stack", string(debug.Stack())),
			)
		}

		// Check for a broken connection, as it is not really a
		// condition that warrants a panic stack trace.
		var brokenPipe bool
		if ne, ok := recovered.(*net.OpError); ok {
			if se, ok := ne.Err.(*os.SyscallError); ok {
				if strings.Contains(strings.ToLower(se.Error()), "broken pipe") ||
					strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
					brokenPipe = true
				}
			}
		}

		if brokenPipe {
			logger.Error("Broken pipe",
				zap.String("request_id", GetRequestID(c)),
				zap.String("path", c.Request.URL.Path),
				zap.String("method", c.Request.Method),
				zap.String("client_ip", c.ClientIP()),
			)
			// If the connection is dead, we can't write a status to it.
			c.Error(recovered.(error)) // nolint: errcheck
			c.Abort()
			return
		}

		httpRequest, _ := httputil.DumpRequest(c.Request, false)
		logger.Error("Panic recovered",
			zap.String("request_id", GetRequestID(c)),
			zap.String("request", string(httpRequest)),
			zap.String("stack", string(debug.Stack())),
		)

		c.AbortWithStatus(http.StatusInternalServerError)
	})
}