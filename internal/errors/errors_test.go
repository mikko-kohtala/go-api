package errors

import (
	"errors"
	"testing"
)

func TestAPIError_Error(t *testing.T) {
	tests := []struct {
		name    string
		err     *APIError
		want    string
	}{
		{
			name: "with underlying error",
			err: &APIError{
				Code:    "test_error",
				Message: "test message",
				Err:     errors.New("underlying error"),
			},
			want: "test message: underlying error",
		},
		{
			name: "without underlying error",
			err: &APIError{
				Code:    "test_error",
				Message: "test message",
			},
			want: "test message",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Errorf("APIError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAPIError_Unwrap(t *testing.T) {
	underlying := errors.New("underlying error")
	err := &APIError{
		Code:    "test_error",
		Message: "test message",
		Err:     underlying,
	}

	if got := err.Unwrap(); got != underlying {
		t.Errorf("APIError.Unwrap() = %v, want %v", got, underlying)
	}
}

func TestNew(t *testing.T) {
	underlying := errors.New("underlying")
	err := New("code", "message", underlying)

	if err.Code != "code" {
		t.Errorf("New() Code = %v, want 'code'", err.Code)
	}
	if err.Message != "message" {
		t.Errorf("New() Message = %v, want 'message'", err.Message)
	}
	if err.Err != underlying {
		t.Errorf("New() Err = %v, want %v", err.Err, underlying)
	}
}

func TestAPIError_WithDetails(t *testing.T) {
	err := New("validation_error", "validation failed", nil)
	details := map[string]string{
		"field1": "required",
		"field2": "invalid format",
	}

	err = err.WithDetails(details)

	if len(err.Details) != 2 {
		t.Errorf("WithDetails() Details length = %v, want 2", len(err.Details))
	}
	if err.Details["field1"] != "required" {
		t.Errorf("WithDetails() Details[field1] = %v, want 'required'", err.Details["field1"])
	}
}

func TestWrap(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		message string
		wantNil bool
		wantMsg string
	}{
		{
			name:    "nil error",
			err:     nil,
			message: "wrapper",
			wantNil: true,
		},
		{
			name:    "non-nil error",
			err:     errors.New("original"),
			message: "wrapper",
			wantNil: false,
			wantMsg: "wrapper: original",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Wrap(tt.err, tt.message)
			if (got == nil) != tt.wantNil {
				t.Errorf("Wrap() = %v, wantNil %v", got, tt.wantNil)
			}
			if !tt.wantNil && got.Error() != tt.wantMsg {
				t.Errorf("Wrap() error = %v, want %v", got.Error(), tt.wantMsg)
			}
		})
	}
}

func TestErrorConstants(t *testing.T) {
	// Verify error constants are defined correctly
	if ErrNotFound.Error() != "resource not found" {
		t.Errorf("ErrNotFound = %v, want 'resource not found'", ErrNotFound)
	}
	if ErrInvalidInput.Error() != "invalid input" {
		t.Errorf("ErrInvalidInput = %v, want 'invalid input'", ErrInvalidInput)
	}
	if ErrUnauthorized.Error() != "unauthorized" {
		t.Errorf("ErrUnauthorized = %v, want 'unauthorized'", ErrUnauthorized)
	}
	if ErrForbidden.Error() != "forbidden" {
		t.Errorf("ErrForbidden = %v, want 'forbidden'", ErrForbidden)
	}
	if ErrInternal.Error() != "internal server error" {
		t.Errorf("ErrInternal = %v, want 'internal server error'", ErrInternal)
	}
	if ErrTimeout.Error() != "operation timeout" {
		t.Errorf("ErrTimeout = %v, want 'operation timeout'", ErrTimeout)
	}
	if ErrRateLimited.Error() != "rate limit exceeded" {
		t.Errorf("ErrRateLimited = %v, want 'rate limit exceeded'", ErrRateLimited)
	}
}