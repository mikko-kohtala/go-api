package logger

import (
	"context"
	"log/slog"
)

type contextKey string

const (
	RequestIDKey contextKey = "request_id"
	TraceIDKey   contextKey = "trace_id"
	UserIDKey    contextKey = "user_id"
)

func FromContext(ctx context.Context) *slog.Logger {
	logger := slog.Default()

	var args []any

	if requestID, ok := ctx.Value(RequestIDKey).(string); ok && requestID != "" {
		args = append(args, "request_id", requestID)
	}

	if traceID, ok := ctx.Value(TraceIDKey).(string); ok && traceID != "" {
		args = append(args, "trace_id", traceID)
	}

	if userID, ok := ctx.Value(UserIDKey).(string); ok && userID != "" {
		args = append(args, "user_id", userID)
	}

	if len(args) > 0 {
		logger = logger.With(args...)
	}

	return logger
}

func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, RequestIDKey, requestID)
}

func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, TraceIDKey, traceID)
}

func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}

func GetRequestID(ctx context.Context) string {
	if id, ok := ctx.Value(RequestIDKey).(string); ok {
		return id
	}
	return ""
}

func GetTraceID(ctx context.Context) string {
	if id, ok := ctx.Value(TraceIDKey).(string); ok {
		return id
	}
	return ""
}

func GetUserID(ctx context.Context) string {
	if id, ok := ctx.Value(UserIDKey).(string); ok {
		return id
	}
	return ""
}