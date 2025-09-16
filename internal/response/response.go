package response

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/mikko-kohtala/go-api/pkg/logger"
)

// ErrorResponse is a consistent error envelope for API responses.
//
// Fields semantics:
// - Error: stable machine‑readable error code (e.g., "invalid_request", "validation_error").
// - Message: human‑readable message safe to show to clients.
// - Fields: optional field‑level messages for validation errors.
// - RequestID: echoes client request id when present.
type ErrorResponse struct {
	Error     string            `json:"error"`
	Message   string            `json:"message,omitempty"`
	Fields    map[string]string `json:"fields,omitempty"`
	RequestID string            `json:"request_id,omitempty"`
}

// JSON writes a JSON response with a status code and logs encoding failures.
func JSON(w http.ResponseWriter, r *http.Request, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		if l := logger.FromContext(r.Context()); l != nil {
			l.Error("encode json response failed", slog.String("error", err.Error()))
		}
	}
}

// ctxKey is a custom type for context keys (matching httpserver package)
type ctxKey string

const requestIDKey ctxKey = "request_id"

// GetRequestID retrieves the request ID from context
func GetRequestID(ctx context.Context) string {
	if v := ctx.Value(requestIDKey); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// Error writes a standardized error response.
func Error(w http.ResponseWriter, r *http.Request, status int, code, message string, fields map[string]string) {
	// Get request ID from context (includes server-generated IDs)
	rid := GetRequestID(r.Context())

	// If no ID in context, fall back to headers (shouldn't happen with middleware)
	if rid == "" {
		rid = r.Header.Get("X-Request-ID")
		if rid == "" {
			rid = r.Header.Get("X-Correlation-ID")
		}
	}

	JSON(w, r, status, ErrorResponse{
		Error:     code,
		Message:   message,
		Fields:    fields,
		RequestID: rid,
	})
}
