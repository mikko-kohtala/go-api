# Logging Documentation

## Overview

This API uses a custom logging system built on Go's `slog` package, providing structured logging with beautiful formatting for development and JSON output for production.

## Features

- **Pretty formatting** for local development with colors and symbols
- **JSON formatting** for production environments
- **Request ID tracking** throughout request lifecycle
- **Structured logging** with groups and attributes
- **Multiple log levels** (DEBUG, INFO, WARN, ERROR)
- **Test endpoint** for demonstrating logging capabilities

## Configuration

### Environment Variables

- `PRETTY_LOGS=true` - Enable pretty formatted logs (automatically set by `make run`)
- `LOG_LEVEL` - Set minimum log level (debug, info, warn, error)

### Running with Pretty Logs

```bash
# Using make (recommended - sets PRETTY_LOGS automatically)
make run

# Or manually with environment variable
PRETTY_LOGS=true go run ./cmd/api

# Or set the environment variable
export PRETTY_LOGS=true
go run ./cmd/api
```

## Log Format Examples

### Pretty Format (Development)

```
11:16:32.110 INFO Started server {port:3000}
11:16:41.142 INFO ▶ GET / {id:"test123"}
11:16:41.142 INFO Root endpoint accessed {id:"test123"}
11:16:41.142 INFO ◀ GET / {id:"test123", status:200, latency:99.71µs}
```

**Format details:**
- Gray timestamps (HH:MM:SS.mmm)
- Colored log levels (INFO=green, WARN=yellow, ERROR=red, DEBUG=gray)
- Arrow symbols for HTTP traffic:
  - `▶` = incoming request
  - `◀` = outgoing response
- Gray metadata in curly braces

### JSON Format (Production)

```json
{"time":"2024-01-15T10:22:20.742Z","level":"INFO","msg":"Started server","port":3000}
{"time":"2024-01-15T10:23:41.773Z","level":"INFO","msg":"request","request_id":"test123","method":"GET","path":"/","status":200,"bytes":83,"duration":"140.92µs"}
```

## Test Endpoint

The `/test/logs` endpoint allows you to generate example logs to see the logging system in action.

### Basic Usage

```bash
# Generate all types of logs (default)
curl "http://localhost:3000/test/logs"

# The endpoint returns a JSON response
{
  "message": "Log examples generated",
  "parameters": {
    "debug": true,
    "info": true,
    "warn": true,
    "error": true,
    "groups": true,
    "count": 1
  },
  "usage": "Use query parameters to control output: debug=false, info=false, warn=false, error=false, groups=false, count=3"
}
```

### Query Parameters

| Parameter | Default | Description |
|-----------|---------|-------------|
| `debug` | `true` | Include DEBUG level logs |
| `info` | `true` | Include INFO level logs |
| `warn` | `true` | Include WARN level logs |
| `error` | `true` | Include ERROR level logs |
| `groups` | `true` | Include examples with grouped attributes |
| `count` | `1` | Number of iterations (max: 10) |

### Example Requests

```bash
# Only show error logs
curl "http://localhost:3000/test/logs?debug=false&info=false&warn=false"

# Generate 3 iterations of all log types
curl "http://localhost:3000/test/logs?count=3"

# Show only INFO and WARN logs
curl "http://localhost:3000/test/logs?debug=false&error=false"

# All logs without grouped examples
curl "http://localhost:3000/test/logs?groups=false"

# Multiple iterations of errors only
curl "http://localhost:3000/test/logs?debug=false&info=false&warn=false&count=5"
```

### Generated Log Examples

When you call the test endpoint, it generates various realistic log scenarios:

#### DEBUG Level
```
12:42:11.883 DEBUG Debug message {id:"...", detail:"This is for debugging", iteration:1, environment:development}
```

#### INFO Level
```
12:42:11.883 INFO Processing request {id:"...", user_id:usr_101, action:view_dashboard, iteration:1}
12:42:11.883 INFO Database query executed {id:"...", query:"SELECT * FROM users WHERE active = true", duration:21.00ms, rows:160}
12:42:11.883 INFO Cache hit {id:"...", key:user:session:1, latency:1.00µs, hit_rate:0.95}
```

#### WARN Level
```
12:42:11.883 WARN Cache miss {id:"...", key:user:session:abc1, fallback:database, latency:150.00ms}
12:42:11.883 WARN Rate limit approaching {id:"...", client_ip:192.168.1.1, requests:91, limit:100, reset_in:55s}
```

#### ERROR Level
```
12:42:11.883 ERROR External API timeout {id:"...", service:payment-gateway, endpoint:https://api.payment.com/charge/1, timeout:5s, retrying:true}
12:42:11.883 ERROR Database connection lost {id:"...", host:db.example.com:5432, error:"connection reset by peer", pool_size:10, active_connections:0}
```

#### Grouped Attributes
```
12:42:11.883 INFO Order processed {id:"...", order:[id=ord_789 total=299.99 currency=USD created_at=2024-01-15T12:42:11Z], customer:[id=cust_456 email=customer@example.com tier=premium], shipping:[method=express carrier=FedEx cost=12.99]}
12:42:11.883 INFO Analytics event {id:"...", event:[type=page_view page=/products time_on_page=45s], user:[id=usr_789 segment=power_user authenticated=true]}
```

## Using the Logger in Code

### Basic Usage

```go
import "github.com/mikko-kohtala/go-api/pkg/logger"

// Create a logger
log := logger.New(
    logger.WithFormat("pretty"),
    logger.WithLevel(slog.LevelDebug),
)

// Simple logging
log.Info("Server started", slog.Int("port", 3000))
log.Error("Connection failed", slog.String("error", err.Error()))
```

### With Request Context

```go
// Get logger from request context (includes request_id)
l := logger.FromContext(r.Context())
if l != nil {
    l.Info("Processing payment",
        slog.Float64("amount", 99.99),
        slog.String("currency", "USD"),
    )
}
```

### Structured Logging with Groups

```go
log.Info("Order completed",
    slog.Group("order",
        slog.String("id", "ord_123"),
        slog.Float64("total", 299.99),
    ),
    slog.Group("customer",
        slog.String("id", "cust_456"),
        slog.String("email", "user@example.com"),
    ),
)
```

## Testing

### Run Tests with Visible Logs

```bash
# Run specific logging tests with output
go test ./internal/httpserver -v -run TestLogger_APIRequestResponse -show-logs

# Run all tests (logs hidden by default)
go test ./...
```

### Format Code

```bash
# Format all Go code including tests
make format
```

## Performance Considerations

- Logs are written synchronously - consider async logging for high-throughput scenarios
- Pretty formatting has minimal overhead suitable for development
- JSON formatting is optimized for production use
- Request ID tracking adds negligible overhead

## Best Practices

1. **Use structured logging** - Add context with key-value pairs rather than formatting strings
2. **Include request IDs** - Always use the logger from request context when available
3. **Choose appropriate levels** - DEBUG for development details, INFO for normal operations, WARN for issues that don't require immediate action, ERROR for failures
4. **Group related data** - Use `slog.Group` for related attributes (order details, user info, etc.)
5. **Keep sensitive data out** - Never log passwords, tokens, or full credit card numbers
6. **Use consistent keys** - Establish naming conventions for common fields (user_id, order_id, etc.)

## Troubleshooting

### Logs Not Showing Colors

Ensure `PRETTY_LOGS=true` is set:
```bash
export PRETTY_LOGS=true
make run
# or
PRETTY_LOGS=true go run ./cmd/api
```

### Debug Logs Not Appearing

Check the log level configuration:
```go
logger.New(
    logger.WithLevel(slog.LevelDebug), // Enable debug logs
)
```

### Request ID Not Appearing

Ensure the request passes through the middleware chain that sets the request ID. The logger should be retrieved from the request context:
```go
l := logger.FromContext(r.Context())
```