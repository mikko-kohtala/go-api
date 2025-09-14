# Swagger API Documentation

This project uses [Swaggo](https://github.com/swaggo/swag) to auto-generate OpenAPI documentation from code annotations.

## Why Generated Docs?

The `internal/docs/docs.go` file is auto-generated and embedded in the Go binary, providing:
- Type-safe documentation that matches actual code
- No separate files to deploy
- Build-time validation of annotations
- Single source of truth (prevents drift)

## Usage

### Add Annotations

Add Swaggo comments above handlers:

```go
// @Summary Get user by ID
// @Tags users
// @Param userID path string true "User ID"
// @Success 200 {object} services.User
// @Router /api/v1/users/{userID} [get]
func GetUser(w http.ResponseWriter, r *http.Request) {
    // implementation
}
```

### Regenerate Docs

After modifying annotations:

```bash
swag init -g cmd/api/main.go -o internal/docs
```

### View Documentation

Access Swagger UI at: `http://localhost:8080/swagger/`

## Common Annotations

- `@Summary` - Short description
- `@Tags` - Group endpoints
- `@Param` - Document parameters
- `@Success` - Success response
- `@Failure` - Error response
- `@Router` - Path and method