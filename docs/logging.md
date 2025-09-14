# Logging Documentation

## Quick Start

```bash
# Run with pretty logs (recommended)
make run

# Or manually
PRETTY_LOGS=true go run ./cmd/api
```

## Log Formats

### Pretty Format (Development)
```
11:16:32.110 INFO Started server {port:3000}
11:16:41.142 INFO ▶ GET / {id:"test123"}
11:16:41.142 INFO ◀ GET / {id:"test123", status:200, latency:99.71µs}
```
- Timestamps: `HH:MM:SS.mmm`
- Colors: INFO=green, WARN=yellow, ERROR=red, DEBUG=gray
- Arrows: `▶` incoming, `◀` outgoing

### JSON Format (Production)
```json
{"time":"2024-01-15T10:22:20.742Z","level":"INFO","msg":"request","request_id":"test123","method":"GET","path":"/","status":200,"duration":"140.92µs"}
```

## Test Endpoint

Generate example logs at `/test/logs`:

```bash
# All log types (default)
curl "http://localhost:3000/test/logs"

# Only errors
curl "http://localhost:3000/test/logs?debug=false&info=false&warn=false"

# Multiple iterations
curl "http://localhost:3000/test/logs?count=3"

# Skip grouped examples
curl "http://localhost:3000/test/logs?groups=false"
```

### Parameters
- `debug`, `info`, `warn`, `error` - Include level (default: true)
- `groups` - Include grouped attributes (default: true)
- `count` - Number of iterations, max 10 (default: 1)

### Example Output
```
12:42:11.883 DEBUG Debug message {id:"...", detail:"This is for debugging", iteration:1}
12:42:11.883 INFO Processing request {id:"...", user_id:usr_101, action:view_dashboard}
12:42:11.883 WARN Cache miss {id:"...", key:user:session:abc1, fallback:database}
12:42:11.883 ERROR External API timeout {id:"...", service:payment-gateway, timeout:5s}
12:42:11.883 INFO Order processed {order:[id=ord_789 total=299.99], customer:[email=customer@example.com]}
```

## Code Usage

```go
import "github.com/mikko-kohtala/go-api/pkg/logger"

// Create logger
log := logger.New(
    logger.WithFormat("pretty"),
    logger.WithLevel(slog.LevelDebug),
)

// Simple logging
log.Info("Server started", slog.Int("port", 3000))

// From request context (includes request_id)
l := logger.FromContext(r.Context())
l.Info("Processing payment", slog.Float64("amount", 99.99))

// Grouped attributes
log.Info("Order completed",
    slog.Group("order", slog.String("id", "ord_123"), slog.Float64("total", 299.99)),
    slog.Group("customer", slog.String("email", "user@example.com")),
)
```

## Testing

```bash
# Show logs during tests
go test ./internal/httpserver -v -run TestLogger -show-logs

# Format code
make format
```

## Best Practices

- Use structured logging with key-value pairs
- Include request IDs from context
- Use appropriate log levels (DEBUG, INFO, WARN, ERROR)
- Group related data with `slog.Group`
- Never log passwords or sensitive data