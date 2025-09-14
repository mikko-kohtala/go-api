package logging

import (
    "context"
    "log/slog"
)

type ctxKey struct{}

// IntoContext stores logger in context.
func IntoContext(ctx context.Context, l *slog.Logger) context.Context {
    return context.WithValue(ctx, ctxKey{}, l)
}

// FromContext returns logger from context, or nil if not set.
func FromContext(ctx context.Context) *slog.Logger {
    if v := ctx.Value(ctxKey{}); v != nil {
        if l, ok := v.(*slog.Logger); ok {
            return l
        }
    }
    return nil
}

