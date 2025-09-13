package validation

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Tag     string `json:"tag"`
	Value   string `json:"value"`
	Message string `json:"message"`
}

// ValidationErrors represents multiple validation errors
type ValidationErrors struct {
	Errors []ValidationError `json:"errors"`
	Count  int              `json:"count"`
}

// CustomValidator wraps the validator with custom validation rules
type CustomValidator struct {
	validator *validator.Validate
	logger    *zap.Logger
}

// NewValidator creates a new custom validator instance
func NewValidator(logger *zap.Logger) *CustomValidator {
	v := validator.New()
	
	// Register custom validators
	v.RegisterValidation("notempty", notEmptyValidator)
	v.RegisterValidation("alphanumeric", alphanumericValidator)
	v.RegisterValidation("email_format", emailFormatValidator)
	v.RegisterValidation("phone_format", phoneFormatValidator)
	
	return &CustomValidator{
		validator: v,
		logger:    logger,
	}
}

// ValidateStruct validates a struct and returns detailed error information
func (cv *CustomValidator) ValidateStruct(obj interface{}) *ValidationErrors {
	err := cv.validator.Struct(obj)
	if err == nil {
		return nil
	}

	var validationErrors []ValidationError
	
	if validationErr, ok := err.(validator.ValidationErrors); ok {
		for _, err := range validationErr {
			validationErrors = append(validationErrors, ValidationError{
				Field:   err.Field(),
				Tag:     err.Tag(),
				Value:   fmt.Sprintf("%v", err.Value()),
				Message: getErrorMessage(err),
			})
		}
	}

	return &ValidationErrors{
		Errors: validationErrors,
		Count:  len(validationErrors),
	}
}

// ValidateAndBind validates and binds JSON request body
func (cv *CustomValidator) ValidateAndBind(c *gin.Context, obj interface{}) *ValidationErrors {
	// First bind the JSON
	if err := c.ShouldBindJSON(obj); err != nil {
		cv.logger.Error("Failed to bind JSON", zap.Error(err))
		return &ValidationErrors{
			Errors: []ValidationError{{
				Field:   "request_body",
				Tag:     "json_binding",
				Value:   "",
				Message: "Invalid JSON format",
			}},
			Count: 1,
		}
	}

	// Then validate the struct
	return cv.ValidateStruct(obj)
}

// HandleValidationError handles validation errors and returns appropriate HTTP response
func (cv *CustomValidator) HandleValidationError(c *gin.Context, validationErrors *ValidationErrors) {
	if validationErrors == nil {
		return
	}

	cv.logger.Warn("Validation failed",
		zap.String("request_id", getRequestID(c)),
		zap.Int("error_count", validationErrors.Count),
		zap.Any("errors", validationErrors.Errors),
	)

	c.JSON(http.StatusBadRequest, gin.H{
		"error":   "Validation failed",
		"details": validationErrors,
	})
	c.Abort()
}

// Custom validation functions

// notEmptyValidator validates that a string is not empty (including whitespace)
func notEmptyValidator(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	return strings.TrimSpace(value) != ""
}

// alphanumericValidator validates that a string contains only alphanumeric characters
func alphanumericValidator(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	matched, _ := regexp.MatchString("^[a-zA-Z0-9]+$", value)
	return matched
}

// emailFormatValidator validates email format
func emailFormatValidator(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(value)
}

// phoneFormatValidator validates phone number format
func phoneFormatValidator(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	phoneRegex := regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)
	return phoneRegex.MatchString(value)
}

// getErrorMessage returns a human-readable error message for validation errors
func getErrorMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", err.Field())
	case "min":
		return fmt.Sprintf("%s must be at least %s characters long", err.Field(), err.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters long", err.Field(), err.Param())
	case "email":
		return fmt.Sprintf("%s must be a valid email address", err.Field())
	case "notempty":
		return fmt.Sprintf("%s cannot be empty", err.Field())
	case "alphanumeric":
		return fmt.Sprintf("%s must contain only alphanumeric characters", err.Field())
	case "email_format":
		return fmt.Sprintf("%s must be a valid email address", err.Field())
	case "phone_format":
		return fmt.Sprintf("%s must be a valid phone number", err.Field())
	default:
		return fmt.Sprintf("%s is invalid", err.Field())
	}
}

// getRequestID extracts request ID from context
func getRequestID(c *gin.Context) string {
	if requestID, exists := c.Get("request_id"); exists {
		if id, ok := requestID.(string); ok {
			return id
		}
	}
	return ""
}