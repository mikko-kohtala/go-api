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

Quick start
-----------

- Install Go 1.23+
- Generate docs, then run:

```
make docs
make run
```

Open http://localhost:8080/swagger/index.html for interactive API docs.

Configuration
-------------

Environment variables (see `.env.example`):

- `APP_ENV` (development|production)
- `PORT` (default 8080)
- `REQUEST_TIMEOUT` (e.g. 15s)
- `BODY_LIMIT_BYTES` (default 10485760 = 10MiB)
- `CORS_ALLOWED_ORIGINS`, `CORS_ALLOWED_METHODS`, `CORS_ALLOWED_HEADERS`
- `RATE_LIMIT_ENABLED` (true|false)
- `RATE_LIMIT_PERIOD` (e.g. 1m)
- `RATE_LIMIT` (requests per period per IP)

Endpoints
---------

- `GET /` — basic info
- `GET /healthz` — liveness probe
- `GET /readyz` — readiness probe
- `GET /api/v1/ping` — returns `{ "pong": "ok" }`
- `POST /api/v1/echo` — `{ "message": "..." }` → echoes back
- `GET /swagger/index.html` — docs UI

Docker
------

```
docker build -t init-codex:local .
docker run --rm -p 8080:8080 init-codex:local
```

Notes
-----

- Logs are structured JSON using Go’s `slog`.
- Rate limiting uses `github.com/go-chi/httprate` and is configurable.
- The Swagger docs are generated from comments (`swag init`).
- Request ID propagation: the server trusts `X-Request-ID` (or `X-Correlation-ID`) from the client, echoes it back on responses, and includes it in every log line.
- Validation: JSON bodies are decoded with `DisallowUnknownFields` and validated via struct tags (e.g. `validate:"required,min=1"`).
 - Rate limiting is applied to `/api/*` routes, not to health endpoints.
