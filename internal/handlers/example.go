package handlers

import (
    "net/http"
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

// Ping godoc
// @Summary      Health check ping
// @Description  Returns a simple pong response.
// @Tags         example
// @Success      200 {object} map[string]string
// @Router       /api/v1/ping [get]
func Ping(w http.ResponseWriter, r *http.Request) {
    response.JSON(w, r, http.StatusOK, map[string]string{"pong": "ok"})
}

// NewEchoHandler creates an Echo handler with dependencies
func NewEchoHandler(svc services.ExampleService) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var req EchoRequest
        errs, err := validate.BindAndValidate(r, &req)
        if err != nil {
            response.Error(w, r, http.StatusBadRequest, "invalid_request", "invalid JSON", nil)
            return
        }
        if errs != nil {
            response.Error(w, r, http.StatusBadRequest, "validation_error", "validation failed", errs)
            return
        }

        // Use service layer with context
        result, err := svc.Echo(r.Context(), req.Message)
        if err != nil {
            response.Error(w, r, http.StatusInternalServerError, "internal_error", "failed to process request", nil)
            return
        }

        response.JSON(w, r, http.StatusOK, EchoResponse{Message: result})
    }
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
func Echo(w http.ResponseWriter, r *http.Request) {
    // Legacy handler for backward compatibility
    svc := services.NewExampleService()
    NewEchoHandler(svc)(w, r)
}
