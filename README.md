# Go API Template

A modern, production-ready Go API template built with the latest technologies and best practices.

## Features

- 🚀 **Modern Go**: Built with Go 1.21+ and latest best practices
- 📝 **JSON Logging**: Structured logging with Zap logger
- 📚 **API Documentation**: Auto-generated Swagger/OpenAPI documentation
- 🛡️ **Rate Limiting**: Built-in rate limiting middleware
- 🔧 **Configuration**: Environment-based configuration management
- 🐳 **Docker Ready**: Docker and Docker Compose setup included
- 🏗️ **Clean Architecture**: Well-organized project structure
- 🔄 **CORS Support**: Configurable CORS middleware
- 📊 **Health Checks**: Built-in health check endpoints
- 🧪 **Testing Ready**: Test structure and examples included

## Tech Stack

- **Framework**: [Gin](https://github.com/gin-gonic/gin) - HTTP web framework
- **Logging**: [Zap](https://github.com/uber-go/zap) - Structured logging
- **Documentation**: [Swagger](https://github.com/swaggo/swag) - API documentation
- **Rate Limiting**: Custom middleware implementation
- **Configuration**: Environment variables with sensible defaults

## Quick Start

### Prerequisites

- Go 1.21 or higher
- Docker (optional)

### Installation

1. Clone the repository:
```bash
git clone <your-repo-url>
cd go-api-template
```

2. Install dependencies:
```bash
make deps
```

3. Run the application:
```bash
make run
```

The API will be available at `http://localhost:8080`

### Using Docker

1. Build and run with Docker Compose:
```bash
make docker-run
```

2. Or build the Docker image manually:
```bash
make docker-build
docker run -p 8080:8080 go-api-template
```

## API Endpoints

### Health Check
- `GET /health` - Health check endpoint

### Users
- `GET /api/v1/users` - Get all users
- `GET /api/v1/users/{id}` - Get user by ID
- `POST /api/v1/users` - Create a new user
- `PUT /api/v1/users/{id}` - Update user
- `DELETE /api/v1/users/{id}` - Delete user

### Examples
- `GET /api/v1/examples` - Get all examples
- `POST /api/v1/examples` - Create a new example

### Documentation
- `GET /swagger/index.html` - Swagger UI documentation

## Configuration

The application can be configured using environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `ENVIRONMENT` | `development` | Application environment |
| `PORT` | `8080` | Server port |
| `LOG_LEVEL` | `info` | Log level (debug, info, warn, error) |
| `RATE_LIMIT_REQUESTS` | `100` | Rate limit requests per window |
| `RATE_LIMIT_WINDOW` | `1m` | Rate limit time window |

See `.env.example` for all available configuration options.

## Development

### Available Commands

```bash
make help          # Show all available commands
make build         # Build the application
make run           # Run the application
make test          # Run tests
make clean         # Clean build artifacts
make deps          # Install dependencies
make swagger       # Generate Swagger documentation
make fmt           # Format code
make lint          # Lint code
make docker-build  # Build Docker image
make docker-run    # Run with Docker Compose
```

### Project Structure

```
.
├── cmd/                    # Application entrypoints
├── internal/              # Private application code
│   ├── config/           # Configuration management
│   ├── handlers/         # HTTP handlers
│   ├── logger/           # Logging setup
│   ├── middleware/       # HTTP middleware
│   └── models/           # Data models
├── docs/                 # Swagger documentation
├── Dockerfile           # Docker configuration
├── docker-compose.yml   # Docker Compose setup
├── Makefile            # Development commands
└── README.md           # This file
```

### Adding New Endpoints

1. Create a new handler in `internal/handlers/`
2. Add the route in `main.go`
3. Add Swagger documentation comments
4. Update the models if needed

Example handler:
```go
// GetExample godoc
// @Summary Get example
// @Description Get example by ID
// @Tags examples
// @Accept json
// @Produce json
// @Param id path int true "Example ID"
// @Success 200 {object} models.Example
// @Router /api/v1/examples/{id} [get]
func GetExample(c *gin.Context) {
    // Implementation
}
```

## Logging

The application uses structured JSON logging with Zap. Logs include:

- Request/response information
- Error details with stack traces
- Performance metrics
- Custom application events

Example log entry:
```json
{
  "timestamp": "2023-01-01T00:00:00Z",
  "level": "info",
  "message": "HTTP Request",
  "method": "GET",
  "path": "/api/v1/users",
  "status": 200,
  "latency": "1.234ms",
  "client_ip": "127.0.0.1"
}
```

## Rate Limiting

The API includes built-in rate limiting to prevent abuse:

- Configurable requests per time window
- Per-client IP limiting
- Automatic cleanup of expired entries
- HTTP 429 response when limit exceeded

## Testing

Run tests with:
```bash
make test
```

The project includes example tests and is ready for comprehensive testing setup.

## Production Deployment

### Environment Variables

Set the following environment variables for production:

```bash
ENVIRONMENT=production
PORT=8080
LOG_LEVEL=info
RATE_LIMIT_REQUESTS=1000
RATE_LIMIT_WINDOW=1m
```

### Docker Deployment

1. Build the production image:
```bash
docker build -t your-registry/go-api-template:latest .
```

2. Push to registry:
```bash
docker push your-registry/go-api-template:latest
```

3. Deploy with your orchestration tool (Kubernetes, Docker Swarm, etc.)

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Run `make fmt` and `make lint`
6. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

For support and questions, please open an issue in the repository.