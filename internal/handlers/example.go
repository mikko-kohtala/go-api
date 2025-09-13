package handlers

import (
    "net/http"
)

import "init-codex/internal/validate"

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
func Ping(w http.ResponseWriter, _ *http.Request) {
    respondJSON(w, http.StatusOK, map[string]string{"pong": "ok"})
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
    var req EchoRequest
    errs, err := validate.BindAndValidate(r, &req)
    if err != nil {
        respondJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid JSON", "detail": err.Error()})
        return
    }
    if errs != nil {
        respondJSON(w, http.StatusBadRequest, map[string]any{"error": "validation failed", "fields": errs})
        return
    }
    respondJSON(w, http.StatusOK, EchoResponse{Message: req.Message})
}
