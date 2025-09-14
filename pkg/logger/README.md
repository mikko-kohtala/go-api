# Logger Package

A structured logging package for Go applications built on top of `slog` with support for both JSON and pretty-printed human-readable formats.

## Features

- 🎨 **Pretty Logging**: Human-readable colored output for local development
- 📊 **JSON Logging**: Structured JSON output for production
- 🔧 **Flexible Configuration**: Configure via code or environment variables
- 🎯 **Context Support**: Store and retrieve loggers from context
- 🚦 **Request Flow Visualization**: Visual indicators for incoming/outgoing HTTP requests
- 🏷️ **Request ID Tracking**: Automatic correlation of logs within a request lifecycle

## Installation

```go
import "github.com/mikko-kohtala/go-api/pkg/logger"
```

## Quick Start

### Basic Usage

```go
// Create a logger with default settings (JSON format, Info level)
log := logger.New()
log.Info("Application started")

// Create a logger with pretty format
log := logger.New(logger.WithFormat("pretty"))
log.Info("Application started")

// Create a logger for specific environment
log := logger.NewForEnvironment("development") // Pretty format, Debug level
log := logger.NewForEnvironment("production")  // JSON format, Info level
```

### Configuration Options

```go
log := logger.New(
    logger.WithLevel(slog.LevelDebug),      // Set log level
    logger.WithFormat("pretty"),            // Set format (json/pretty)
    logger.WithSource(true),                // Include source location
    logger.WithOutput(os.Stderr),           // Set output writer
)
```

### Environment Variables

The logger respects the `PRETTY_LOGS` environment variable:

```bash
PRETTY_LOGS=true go run main.go  # Forces pretty format regardless of config
```

## HTTP Middleware Integration

The package provides excellent integration with HTTP middleware for request tracking:

```go
func LoggingMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Create request-scoped logger
            reqLogger := logger.With(slog.String("request_id", generateRequestID()))

            // Log incoming request with direction indicator
            if prettyLogsEnabled {
                incomingLogger := reqLogger.With(slog.String("direction", "incoming"))
                incomingLogger.Info(fmt.Sprintf("%s %s", r.Method, r.URL.Path))
            }

            // Store logger in context for handlers
            ctx := logger.IntoContext(r.Context(), reqLogger)

            // Process request
            next.ServeHTTP(w, r.WithContext(ctx))

            // Log outgoing response
            if prettyLogsEnabled {
                outgoingLogger := reqLogger.With(slog.String("direction", "outgoing"))
                outgoingLogger.Info(fmt.Sprintf("%s %s", r.Method, r.URL.Path),
                    slog.Int("status", statusCode),
                    slog.Duration("latency", duration),
                )
            }
        })
    }
}
```

### Using Logger in Handlers

```go
func MyHandler(w http.ResponseWriter, r *http.Request) {
    // Retrieve logger from context
    log := logger.FromContext(r.Context())

    // All logs will automatically include the request_id
    log.Info("Processing request")
    log.Debug("Request details", slog.String("user", userID))

    if err != nil {
        log.Error("Failed to process", slog.String("error", err.Error()))
    }
}
```

## Pretty Format Output

When using pretty format, logs are displayed with:
- Color-coded log levels (INFO=green, WARN=yellow, ERROR=red, DEBUG=gray)
- Clean timestamps (HH:MM:SS.mmm)
- Request flow indicators (→ for incoming, ← for outgoing)
- Minimal, relevant information for local development

Example output:
```
10:23:41.773 INFO Started server {port:3000}
10:23:41.773 INFO → GET / {id:"56ada389..."}
10:23:41.773 INFO Processing request {id:"56ada389..."}
10:23:41.774 INFO ← GET / {id:"56ada389...", status:200, latency:121.54µs}
```

## JSON Format Output

In production, logs are output as structured JSON:
```json
{
  "time": "2025-09-14T10:23:41.773Z",
  "level": "INFO",
  "msg": "request",
  "request_id": "56ada389-4d0f-862e-93d7-5e07c09dc35a",
  "method": "GET",
  "path": "/",
  "status": 200,
  "duration": "121.54µs"
}
```

## Context Integration

The package provides helpers for storing and retrieving loggers from context:

```go
// Store logger in context
ctx := logger.IntoContext(ctx, myLogger)

// Retrieve logger from context (returns default logger if not found)
log := logger.FromContext(ctx)
```

## Advanced Usage

### Custom Handler Options

```go
opts := &slog.HandlerOptions{
    Level:     slog.LevelDebug,
    AddSource: true,
    ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
        // Custom attribute replacement
        return a
    },
}

handler := logger.NewPrettyHandler(os.Stdout, opts)
log := slog.New(handler)
```

### Multiple Loggers

```go
// Application logger
appLogger := logger.New(logger.WithFormat("json"))

// Audit logger with different configuration
auditLogger := logger.New(
    logger.WithFormat("json"),
    logger.WithOutput(auditFile),
    logger.WithLevel(slog.LevelInfo),
)

// Development logger for debugging
debugLogger := logger.New(
    logger.WithFormat("pretty"),
    logger.WithLevel(slog.LevelDebug),
    logger.WithSource(true),
)
```

## Best Practices

1. **Use Request-Scoped Loggers**: Always create request-scoped loggers in middleware to maintain request correlation
2. **Leverage Context**: Pass loggers through context rather than global variables
3. **Environment-Specific Configuration**: Use `NewForEnvironment()` for automatic environment-based configuration
4. **Structured Logging**: Use slog attributes for structured data rather than string interpolation
5. **Error Handling**: Always log errors with appropriate context

## License

This package is part of the go-api project and follows the same license terms.