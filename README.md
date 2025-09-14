Init Codex — Minimal Go API Template
===================================

A clean Go 1.23 HTTP API starter using:

- chi router and modern middleware
- slog JSON logging
- Swagger/OpenAPI docs via swag
- Optional per‑IP rate limiting
- Graceful shutdown and sane defaults
- JSON request validation (go-playground/validator) with unknown-field rejection
- Request body size limit via `BODY_LIMIT_BYTES` (default 10 MiB)
- Configurable gzip compression level (`COMPRESSION_LEVEL`, default 5)
- **NEW**: Prometheus metrics collection
- **NEW**: OpenTelemetry distributed tracing
- **NEW**: Enhanced error handling with structured error types
- **NEW**: Service layer with dependency injection
- **NEW**: Comprehensive test coverage

Quick start
-----------

- Install Go 1.23+
- Generate docs, then run:

```
make docs
make run
```

Open http://localhost:3000/swagger/index.html for interactive API docs.

Configuration
-------------

Environment variables (see `.env.example`):

- `APP_ENV` (development|production)
- `PORT` (default 3000)
- `REQUEST_TIMEOUT` (e.g. 15s)
- `BODY_LIMIT_BYTES` (default 10485760 = 10MiB)
- `COMPRESSION_LEVEL` (1–9, default 5)
- `CORS_ALLOWED_ORIGINS`, `CORS_ALLOWED_METHODS`, `CORS_ALLOWED_HEADERS`
- `RATE_LIMIT_ENABLED` (true|false)
- `RATE_LIMIT_PERIOD` (e.g. 1m)
- `RATE_LIMIT` (requests per period per IP)
- `METRICS_ENABLED` (true|false, default true)
- `METRICS_PATH` (default /metrics)
- `TRACING_ENABLED` (true|false, default false)

Endpoints
---------

- `GET /` — basic info
- `GET /healthz` — liveness probe
- `GET /readyz` — readiness probe
- `GET /api/v1/ping` — returns `{ "pong": "ok" }`
- `POST /api/v1/echo` — `{ "message": "..." }` → echoes back
- `GET /metrics` — Prometheus metrics (if enabled)
- `GET /swagger/index.html` — docs UI
- `GET /api-docs` — docs UI (alias for Swagger)

Docker
------

```
docker build -t init-codex:local .
docker run --rm -p 3000:3000 init-codex:local
```

Notes
-----

- Logs are structured JSON using Go’s `slog`.
- Rate limiting uses `github.com/go-chi/httprate` and is configurable.
- The Swagger docs are generated from comments (`swag init`).
- Request ID propagation: the server trusts `X-Request-ID` (or `X-Correlation-ID`) from the client, echoes it back on responses, and includes it in every log line.
- Validation: JSON bodies are decoded with `DisallowUnknownFields` and validated via struct tags (e.g. `validate:"required,min=1"`).
- Rate limiting is applied to `/api/*` routes, not to health endpoints.
- In-memory rate limiting is per-instance; for multi-instance deployments, use sticky sessions or replace with a distributed limiter.
- CORS strict mode: set `CORS_STRICT=true` to fail startup if `*` is used in production.
- **NEW**: Prometheus metrics are collected automatically for all HTTP requests.
- **NEW**: OpenTelemetry tracing is available (disabled by default).
- **NEW**: Structured error handling with proper HTTP status codes.
- **NEW**: Service layer architecture with dependency injection for better testability.
