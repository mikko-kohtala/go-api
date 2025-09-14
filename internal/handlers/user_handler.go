package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mikko-kohtala/go-api/internal/response"
	"github.com/mikko-kohtala/go-api/internal/services"
	"github.com/mikko-kohtala/go-api/internal/validate"
)

type UserHandler struct {
	userService services.UserService
	logger      *slog.Logger
}

func NewUserHandler(userService services.UserService, logger *slog.Logger) *UserHandler {
	return &UserHandler{
		userService: userService,
		logger:      logger,
	}
}

type CreateUserRequest struct {
	Email string `json:"email" validate:"required,email"`
	Name  string `json:"name" validate:"required,min=1,max=100"`
}

type UpdateUserRequest struct {
	Email string `json:"email,omitempty" validate:"omitempty,email"`
	Name  string `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
	Role  string `json:"role,omitempty" validate:"omitempty,oneof=admin user moderator"`
}

// GetAllUsers godoc
// @Summary      Get all users
// @Description  Returns a list of all users
// @Tags         users
// @Produce      json
// @Success      200 {array} services.User
// @Failure      500 {object} map[string]interface{}
// @Router       /api/v1/users [get]
func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.userService.GetAllUsers(r.Context())
	if err != nil {
		h.logger.Error("failed to get users", slog.String("error", err.Error()))
		response.Error(w, r, http.StatusInternalServerError, "internal_error", "Failed to retrieve users", nil)
		return
	}

	response.JSON(w, r, http.StatusOK, map[string]interface{}{
		"users": users,
		"count": len(users),
	})
}

// GetUserByID godoc
// @Summary      Get user by ID
// @Description  Returns a single user by ID
// @Tags         users
// @Produce      json
// @Param        userID path string true "User ID"
// @Success      200 {object} services.User
// @Failure      404 {object} map[string]interface{}
// @Failure      500 {object} map[string]interface{}
// @Router       /api/v1/users/{userID} [get]
func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	if userID == "" {
		response.Error(w, r, http.StatusBadRequest, "invalid_request", "User ID is required", nil)
		return
	}

	user, err := h.userService.GetUserByID(r.Context(), userID)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			h.logger.Debug("user not found", slog.String("user_id", userID))
			response.Error(w, r, http.StatusNotFound, "not_found", "User not found", nil)
			return
		}
		h.logger.Error("failed to get user", slog.String("error", err.Error()))
		response.Error(w, r, http.StatusInternalServerError, "internal_error", "Failed to retrieve user", nil)
		return
	}

	response.JSON(w, r, http.StatusOK, user)
}

// CreateUser godoc
// @Summary      Create a new user
// @Description  Creates a new user with the provided information
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user body CreateUserRequest true "User information"
// @Success      201 {object} services.User
// @Failure      400 {object} map[string]interface{}
// @Failure      409 {object} map[string]interface{}
// @Failure      500 {object} map[string]interface{}
// @Router       /api/v1/users [post]
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	errs, err := validate.BindAndValidate(r, &req)
	if err != nil {
		response.Error(w, r, http.StatusBadRequest, "invalid_request", "Invalid JSON", nil)
		return
	}
	if errs != nil {
		response.Error(w, r, http.StatusBadRequest, "validation_error", "Validation failed", errs)
		return
	}

	user, err := h.userService.CreateUser(r.Context(), req.Email, req.Name)
	if err != nil {
		if errors.Is(err, services.ErrEmailAlreadyExists) {
			response.Error(w, r, http.StatusConflict, "duplicate_email", "Email already exists", nil)
			return
		}
		if errors.Is(err, services.ErrInvalidEmail) {
			response.Error(w, r, http.StatusBadRequest, "invalid_email", "Invalid email address", nil)
			return
		}
		h.logger.Error("failed to create user", slog.String("error", err.Error()))
		response.Error(w, r, http.StatusInternalServerError, "internal_error", "Failed to create user", nil)
		return
	}

	h.logger.Info("user created", slog.String("user_id", user.ID), slog.String("email", user.Email))
	response.JSON(w, r, http.StatusCreated, user)
}

// UpdateUser godoc
// @Summary      Update a user
// @Description  Updates user information
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        userID path string true "User ID"
// @Param        user body UpdateUserRequest true "User update information"
// @Success      200 {object} services.User
// @Failure      400 {object} map[string]interface{}
// @Failure      404 {object} map[string]interface{}
// @Failure      500 {object} map[string]interface{}
// @Router       /api/v1/users/{userID} [put]
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	if userID == "" {
		response.Error(w, r, http.StatusBadRequest, "invalid_request", "User ID is required", nil)
		return
	}

	var req UpdateUserRequest
	errs, err := validate.BindAndValidate(r, &req)
	if err != nil {
		response.Error(w, r, http.StatusBadRequest, "invalid_request", "Invalid JSON", nil)
		return
	}
	if errs != nil {
		response.Error(w, r, http.StatusBadRequest, "validation_error", "Validation failed", errs)
		return
	}

	// Convert request to map for updates
	updates := make(map[string]interface{})
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Role != "" {
		updates["role"] = req.Role
	}

	user, err := h.userService.UpdateUser(r.Context(), userID, updates)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			response.Error(w, r, http.StatusNotFound, "not_found", "User not found", nil)
			return
		}
		if errors.Is(err, services.ErrEmailAlreadyExists) {
			response.Error(w, r, http.StatusConflict, "duplicate_email", "Email already exists", nil)
			return
		}
		h.logger.Error("failed to update user", slog.String("error", err.Error()))
		response.Error(w, r, http.StatusInternalServerError, "internal_error", "Failed to update user", nil)
		return
	}

	h.logger.Info("user updated", slog.String("user_id", user.ID))
	response.JSON(w, r, http.StatusOK, user)
}

// DeleteUser godoc
// @Summary      Delete a user
// @Description  Deletes a user by ID
// @Tags         users
// @Param        userID path string true "User ID"
// @Success      204 "No Content"
// @Failure      404 {object} map[string]interface{}
// @Failure      500 {object} map[string]interface{}
// @Router       /api/v1/users/{userID} [delete]
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	if userID == "" {
		response.Error(w, r, http.StatusBadRequest, "invalid_request", "User ID is required", nil)
		return
	}

	err := h.userService.DeleteUser(r.Context(), userID)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			response.Error(w, r, http.StatusNotFound, "not_found", "User not found", nil)
			return
		}
		h.logger.Error("failed to delete user", slog.String("error", err.Error()))
		response.Error(w, r, http.StatusInternalServerError, "internal_error", "Failed to delete user", nil)
		return
	}

	h.logger.Info("user deleted", slog.String("user_id", userID))
	w.WriteHeader(http.StatusNoContent)
}