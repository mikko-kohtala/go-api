package errors

import (
	"fmt"
	"net/http"
)

// ErrorCode represents a machine-readable error code
type ErrorCode string

const (
	// Validation errors
	ErrCodeValidation     ErrorCode = "validation_error"
	ErrCodeInvalidRequest ErrorCode = "invalid_request"
	
	// Authentication/Authorization errors
	ErrCodeUnauthorized   ErrorCode = "unauthorized"
	ErrCodeForbidden      ErrorCode = "forbidden"
	
	// Resource errors
	ErrCodeNotFound       ErrorCode = "not_found"
	ErrCodeConflict       ErrorCode = "conflict"
	
	// Server errors
	ErrCodeInternal       ErrorCode = "internal_error"
	ErrCodeServiceUnavailable ErrorCode = "service_unavailable"
)

// APIError represents a structured API error
type APIError struct {
	Code      ErrorCode         `json:"code"`
	Message   string            `json:"message"`
	Fields    map[string]string `json:"fields,omitempty"`
	RequestID string            `json:"request_id,omitempty"`
	Err       error             `json:"-"`
}

// Error implements the error interface
func (e *APIError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s", e.Message, e.Err.Error())
	}
	return e.Message
}

// Unwrap returns the underlying error
func (e *APIError) Unwrap() error {
	return e.Err
}

// HTTPStatus returns the appropriate HTTP status code for the error
func (e *APIError) HTTPStatus() int {
	switch e.Code {
	case ErrCodeValidation, ErrCodeInvalidRequest:
		return http.StatusBadRequest
	case ErrCodeUnauthorized:
		return http.StatusUnauthorized
	case ErrCodeForbidden:
		return http.StatusForbidden
	case ErrCodeNotFound:
		return http.StatusNotFound
	case ErrCodeConflict:
		return http.StatusConflict
	case ErrCodeServiceUnavailable:
		return http.StatusServiceUnavailable
	default:
		return http.StatusInternalServerError
	}
}

// New creates a new APIError
func New(code ErrorCode, message string) *APIError {
	return &APIError{
		Code:    code,
		Message: message,
	}
}

// Wrap wraps an existing error with additional context
func Wrap(err error, code ErrorCode, message string) *APIError {
	return &APIError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// WithFields adds field-specific error details
func (e *APIError) WithFields(fields map[string]string) *APIError {
	e.Fields = fields
	return e
}

// WithRequestID adds request ID to the error
func (e *APIError) WithRequestID(requestID string) *APIError {
	e.RequestID = requestID
	return e
}