package handlers

import (
	"github.com/mikko-kohtala/go-api/internal/response"
	"github.com/mikko-kohtala/go-api/internal/validate"
	"net/http"
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
		response.Error(w, r, http.StatusBadRequest, "invalid_request", "invalid JSON", nil)
		return
	}
	if errs != nil {
		response.Error(w, r, http.StatusBadRequest, "validation_error", "validation failed", errs)
		return
	}
	response.JSON(w, r, http.StatusOK, EchoResponse{Message: req.Message})
}
