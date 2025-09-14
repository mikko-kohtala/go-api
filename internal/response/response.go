package response

import (
    "encoding/json"
    "log/slog"
    "net/http"

    "github.com/mikko-kohtala/go-api/internal/logging"
    "github.com/mikko-kohtala/go-api/internal/requestid"
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
        if l := logging.FromContext(r.Context()); l != nil {
            l.Error("encode json response failed", slog.String("error", err.Error()))
        }
    }
}

// Error writes a standardized error response.
func Error(w http.ResponseWriter, r *http.Request, status int, code, message string, fields map[string]string) {
    // Prefer ID stored in context by middleware; fall back to trusted headers.
    rid := requestid.FromContext(r.Context())
    if rid == "" {
        rid = r.Header.Get(requestid.HeaderRequestID)
        if rid == "" {
            rid = r.Header.Get(requestid.HeaderCorrelationID)
        }
    }
    JSON(w, r, status, ErrorResponse{
        Error:     code,
        Message:   message,
        Fields:    fields,
        RequestID: rid,
    })
}
