package httpserver

import (
	"net/http"

	"github.com/mikko-kohtala/go-api/internal/services"
)

// ConnectionsMiddleware tracks active HTTP connections
func ConnectionsMiddleware(statsService services.StatsService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Increment active connections on request start
			statsService.IncrementActiveConnections()

			// Decrement when request completes
			defer statsService.DecrementActiveConnections()

			// Process request
			next.ServeHTTP(w, r)
		})
	}
}