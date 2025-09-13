package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiterEntry represents a rate limiter entry for a client
type RateLimiterEntry struct {
	requests int
	window   time.Time
}

// RateLimiter returns a gin.HandlerFunc for rate limiting
func RateLimiter(requests int, window time.Duration) gin.HandlerFunc {
	// Create a simple in-memory rate limiter
	store := make(map[string]*RateLimiterEntry)
	mutex := sync.RWMutex{}

	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		now := time.Now()
		
		mutex.Lock()
		entry, exists := store[clientIP]
		if !exists {
			entry = &RateLimiterEntry{
				requests: 1,
				window:   now.Add(window),
			}
			store[clientIP] = entry
		} else {
			// Reset if window has expired
			if now.After(entry.window) {
				entry.requests = 1
				entry.window = now.Add(window)
			} else {
				entry.requests++
			}
		}
		mutex.Unlock()

		if entry.requests > requests {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":       "Rate limit exceeded",
				"retry_after": entry.window.Sub(now).Seconds(),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}