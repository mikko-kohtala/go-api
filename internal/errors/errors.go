package errors

import (
	"errors"
	"fmt"
)

// Domain error types
var (
	ErrNotFound      = errors.New("resource not found")
	ErrInvalidInput  = errors.New("invalid input")
	ErrUnauthorized  = errors.New("unauthorized")
	ErrForbidden     = errors.New("forbidden")
	ErrInternal      = errors.New("internal server error")
	ErrTimeout       = errors.New("operation timeout")
	ErrRateLimited   = errors.New("rate limit exceeded")
)

// APIError represents a structured API error
type APIError struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Details map[string]string `json:"details,omitempty"`
	Err     error             `json:"-"`
}

func (e *APIError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *APIError) Unwrap() error {
	return e.Err
}

// New creates a new APIError
func New(code, message string, err error) *APIError {
	return &APIError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// WithDetails adds field-level details to the error
func (e *APIError) WithDetails(details map[string]string) *APIError {
	e.Details = details
	return e
}

// Wrap wraps an error with additional context
func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", message, err)
}