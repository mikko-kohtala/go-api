package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/mikko-kohtala/go-api/internal/response"
	pkglogger "github.com/mikko-kohtala/go-api/pkg/logger"
)

type TestLogResponse struct {
	Message    string                 `json:"message"`
	Parameters map[string]interface{} `json:"parameters"`
	Usage      string                 `json:"usage"`
}

// TestLogs godoc
// @Summary      Generate test log entries
// @Description  Generates various log entries for testing logging capabilities
// @Tags         test
// @Produce      json
// @Param        debug  query    bool     false  "Include debug logs (default: true)"
// @Param        info   query    bool     false  "Include info logs (default: true)"
// @Param        warn   query    bool     false  "Include warning logs (default: true)"
// @Param        error  query    bool     false  "Include error logs (default: true)"
// @Param        groups query    bool     false  "Include grouped logs (default: true)"
// @Param        count  query    int      false  "Number of iterations (1-10, default: 1)"
// @Success      200    {object} TestLogResponse
// @Router       /test/logs [get]
func TestLogs(w http.ResponseWriter, r *http.Request) {
	l := pkglogger.FromContext(r.Context())
	if l == nil {
		http.Error(w, "Logger not available", http.StatusInternalServerError)
		return
	}

	// Parse query parameters
	query := r.URL.Query()

	// Optional parameters with defaults
	includeDebug := query.Get("debug") != "false"   // default: true
	includeInfo := query.Get("info") != "false"     // default: true
	includeWarn := query.Get("warn") != "false"     // default: true
	includeError := query.Get("error") != "false"   // default: true
	includeGroups := query.Get("groups") != "false" // default: true
	count := 1                                      // default: 1 iteration
	if c := query.Get("count"); c != "" {
		if parsed, err := fmt.Sscanf(c, "%d", &count); err == nil && parsed == 1 && count > 0 && count <= 10 {
			// Use the parsed count (max 10 to prevent abuse)
		} else {
			count = 1
		}
	}

	// Generate log examples based on parameters
	for i := 0; i < count; i++ {
		iteration := i + 1

		if includeDebug {
			l.Debug("Debug message",
				slog.String("detail", "This is for debugging"),
				slog.Int("iteration", iteration),
				slog.String("environment", "development"))
		}

		if includeInfo {
			l.Info("Processing request",
				slog.String("user_id", fmt.Sprintf("usr_%d", 100+iteration)),
				slog.String("action", "view_dashboard"),
				slog.Int("iteration", iteration))

			l.Info("Database query executed",
				slog.String("query", "SELECT * FROM users WHERE active = true"),
				slog.Duration("duration", time.Duration(20+iteration)*time.Millisecond),
				slog.Int("rows", 150+iteration*10))

			l.Info("Cache hit",
				slog.String("key", fmt.Sprintf("user:session:%d", iteration)),
				slog.Duration("latency", time.Duration(iteration)*time.Microsecond),
				slog.Float64("hit_rate", 0.95))
		}

		if includeWarn {
			l.Warn("Cache miss",
				slog.String("key", fmt.Sprintf("user:session:abc%d", iteration)),
				slog.String("fallback", "database"),
				slog.Duration("latency", time.Duration(100+iteration*50)*time.Millisecond))

			l.Warn("Rate limit approaching",
				slog.String("client_ip", fmt.Sprintf("192.168.1.%d", iteration)),
				slog.Int("requests", 90+iteration),
				slog.Int("limit", 100),
				slog.Duration("reset_in", time.Duration(60-iteration*5)*time.Second))
		}

		if includeError {
			l.Error("External API timeout",
				slog.String("service", "payment-gateway"),
				slog.String("endpoint", fmt.Sprintf("https://api.payment.com/charge/%d", iteration)),
				slog.Duration("timeout", time.Duration(5)*time.Second),
				slog.Bool("retrying", iteration < 3))

			if iteration == 1 {
				l.Error("Database connection lost",
					slog.String("host", "db.example.com:5432"),
					slog.String("error", "connection reset by peer"),
					slog.Int("pool_size", 10),
					slog.Int("active_connections", 0))
			}
		}

		if includeGroups && iteration == 1 {
			l.Info("Order processed",
				slog.Group("order",
					slog.String("id", "ord_789"),
					slog.Float64("total", 299.99),
					slog.String("currency", "USD"),
					slog.Time("created_at", time.Now()),
				),
				slog.Group("customer",
					slog.String("id", "cust_456"),
					slog.String("email", "customer@example.com"),
					slog.String("tier", "premium"),
				),
				slog.Group("shipping",
					slog.String("method", "express"),
					slog.String("carrier", "FedEx"),
					slog.Float64("cost", 12.99),
				))

			l.Info("Analytics event",
				slog.Group("event",
					slog.String("type", "page_view"),
					slog.String("page", "/products"),
					slog.Duration("time_on_page", 45*time.Second),
				),
				slog.Group("user",
					slog.String("id", "usr_789"),
					slog.String("segment", "power_user"),
					slog.Bool("authenticated", true),
				))
		}
	}

	// Response with summary
	resp := TestLogResponse{
		Message: "Log examples generated",
		Parameters: map[string]interface{}{
			"debug":  includeDebug,
			"info":   includeInfo,
			"warn":   includeWarn,
			"error":  includeError,
			"groups": includeGroups,
			"count":  count,
		},
		Usage: "Use query parameters to control output: debug=false, info=false, warn=false, error=false, groups=false, count=3",
	}

	response.JSON(w, r, http.StatusOK, resp)
}