package requestid

import (
    "context"
)

// Header names used for request ID propagation.
const (
    HeaderRequestID     = "X-Request-ID"
    HeaderCorrelationID = "X-Correlation-ID"
)

// ctxKey is an unexported type for keys defined in this package.
type ctxKey struct{}

// IntoContext stores a request ID in context.
func IntoContext(ctx context.Context, id string) context.Context {
    if id == "" {
        return ctx
    }
    return context.WithValue(ctx, ctxKey{}, id)
}

// FromContext returns the request ID from context, if present.
func FromContext(ctx context.Context) string {
    if v := ctx.Value(ctxKey{}); v != nil {
        if s, ok := v.(string); ok {
            return s
        }
    }
    return ""
}

