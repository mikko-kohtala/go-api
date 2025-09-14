package handlers

import (
    "net/http"
    "github.com/mikko-kohtala/go-api/internal/errors"
    "github.com/mikko-kohtala/go-api/internal/response"
    "github.com/mikko-kohtala/go-api/internal/services"
    "github.com/mikko-kohtala/go-api/internal/validate"
)

type EchoRequest struct {
    Message string `json:"message" validate:"required,min=1"`
}

type EchoResponse struct {
    Message string `json:"message"`
}

// Handler holds dependencies for handlers
type Handler struct {
    services *services.ServiceContainer
}

// NewHandler creates a new handler with dependencies
func NewHandler(services *services.ServiceContainer) *Handler {
    return &Handler{
        services: services,
    }
}

// Ping godoc
// @Summary      Health check ping
// @Description  Returns a simple pong response.
// @Tags         example
// @Success      200 {object} map[string]string
// @Router       /api/v1/ping [get]
func Ping(w http.ResponseWriter, r *http.Request) {
    response.JSON(w, r, http.StatusOK, map[string]string{"pong": "ok"})
}

// Echo godoc
// @Summary      Echo a JSON payload
// @Description  Returns a JSON payload with the same message.
// @Tags         example
// @Accept       json
// @Produce      json
// @Param        request  body      EchoRequest  true  "Echo input"
// @Success      200      {object}  EchoResponse
// @Failure      400      {object}  map[string]string
// @Router       /api/v1/echo [post]
func (h *Handler) Echo(w http.ResponseWriter, r *http.Request) {
    var req EchoRequest
    errs, err := validate.BindAndValidate(r, &req)
    if err != nil {
        apiErr := errors.Wrap(err, errors.ErrCodeInvalidRequest, "invalid JSON request")
        response.APIError(w, r, apiErr)
        return
    }
    if errs != nil {
        apiErr := errors.New(errors.ErrCodeValidation, "validation failed").WithFields(errs)
        response.APIError(w, r, apiErr)
        return
    }
    
    // Use service layer for business logic
    result, err := h.services.Echo.Echo(r.Context(), req.Message)
    if err != nil {
        apiErr := errors.Wrap(err, errors.ErrCodeInternal, "failed to process echo request")
        response.APIError(w, r, apiErr)
        return
    }
    
    response.JSON(w, r, http.StatusOK, EchoResponse{Message: result})
}

// Echo is a standalone function for backward compatibility
func Echo(w http.ResponseWriter, r *http.Request) {
    // Create a temporary handler for backward compatibility
    // In a real application, you'd want to refactor this to use dependency injection
    handler := &Handler{
        services: &services.ServiceContainer{
            Echo: services.NewEchoService(nil), // Use nil logger for tests
        },
    }
    handler.Echo(w, r)
}
