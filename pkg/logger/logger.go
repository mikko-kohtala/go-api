// Package logger provides a structured logging solution built on slog
// with support for both JSON and pretty-printed human-readable formats.
package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
)

// Config represents the configuration for the logger
type Config struct {
	// Level sets the minimum log level (Debug, Info, Warn, Error)
	Level slog.Level
	// Format specifies the output format ("json" or "pretty")
	Format string
	// AddSource adds source code location to logs
	AddSource bool
	// Output specifies where to write logs (defaults to os.Stdout)
	Output io.Writer
}

// Option is a functional option for configuring the logger
type Option func(*Config)

// WithLevel sets the log level
func WithLevel(level slog.Level) Option {
	return func(c *Config) {
		c.Level = level
	}
}

// WithFormat sets the output format ("json" or "pretty")
func WithFormat(format string) Option {
	return func(c *Config) {
		c.Format = format
	}
}

// WithSource enables source code location in logs
func WithSource(addSource bool) Option {
	return func(c *Config) {
		c.AddSource = addSource
	}
}

// WithOutput sets the output writer
func WithOutput(w io.Writer) Option {
	return func(c *Config) {
		c.Output = w
	}
}

// New creates a new slog.Logger with the specified options
func New(opts ...Option) *slog.Logger {
	cfg := &Config{
		Level:     slog.LevelInfo,
		Format:    "json",
		AddSource: false,
		Output:    os.Stdout,
	}

	// Apply options
	for _, opt := range opts {
		opt(cfg)
	}

	// Check environment variable for format override
	if os.Getenv("PRETTY_LOGS") == "true" {
		cfg.Format = "pretty"
	}

	var handler slog.Handler
	handlerOpts := &slog.HandlerOptions{
		Level:     cfg.Level,
		AddSource: cfg.AddSource,
	}

	switch cfg.Format {
	case "pretty":
		handler = NewPrettyHandler(cfg.Output, handlerOpts)
	default:
		handler = slog.NewJSONHandler(cfg.Output, handlerOpts)
	}

	return slog.New(handler)
}

// NewForEnvironment creates a logger configured for the specified environment
func NewForEnvironment(env string) *slog.Logger {
	switch env {
	case "development", "dev":
		return New(
			WithLevel(slog.LevelDebug),
			WithFormat("pretty"),
			WithSource(true),
		)
	case "production", "prod":
		return New(
			WithLevel(slog.LevelInfo),
			WithFormat("json"),
			WithSource(false),
		)
	default:
		return New()
	}
}

// FromContext retrieves a logger from the context
func FromContext(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(contextKey).(*slog.Logger); ok {
		return logger
	}
	return slog.Default()
}

// IntoContext stores a logger in the context
func IntoContext(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, contextKey, logger)
}

// contextKey is used to store the logger in context
type ctxKey struct{}

var contextKey = ctxKey{}