package errors

import (
	"net/http"
	"testing"
)

func TestAPIError(t *testing.T) {
	tests := []struct {
		name           string
		code           ErrorCode
		message        string
		expectedStatus int
	}{
		{
			name:           "validation error",
			code:           ErrCodeValidation,
			message:        "validation failed",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "not found error",
			code:           ErrCodeNotFound,
			message:        "resource not found",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "internal error",
			code:           ErrCodeInternal,
			message:        "internal server error",
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := New(tt.code, tt.message)
			
			if err.Code != tt.code {
				t.Errorf("expected code %s, got %s", tt.code, err.Code)
			}
			
			if err.Message != tt.message {
				t.Errorf("expected message %s, got %s", tt.message, err.Message)
			}
			
			if err.HTTPStatus() != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, err.HTTPStatus())
			}
		})
	}
}

func TestAPIErrorWithFields(t *testing.T) {
	err := New(ErrCodeValidation, "validation failed")
	fields := map[string]string{
		"email": "invalid email format",
		"name":  "name is required",
	}
	
	err = err.WithFields(fields)
	
	if err.Fields == nil {
		t.Error("expected fields to be set")
	}
	
	if err.Fields["email"] != "invalid email format" {
		t.Errorf("expected email field error, got %s", err.Fields["email"])
	}
}

func TestAPIErrorWithRequestID(t *testing.T) {
	err := New(ErrCodeValidation, "validation failed")
	requestID := "req-123"
	
	err = err.WithRequestID(requestID)
	
	if err.RequestID != requestID {
		t.Errorf("expected request ID %s, got %s", requestID, err.RequestID)
	}
}

func TestWrapError(t *testing.T) {
	originalErr := &APIError{
		Code:    ErrCodeNotFound,
		Message: "original error",
	}
	
	wrappedErr := Wrap(originalErr, ErrCodeInternal, "wrapped error")
	
	if wrappedErr.Code != ErrCodeInternal {
		t.Errorf("expected code %s, got %s", ErrCodeInternal, wrappedErr.Code)
	}
	
	if wrappedErr.Message != "wrapped error" {
		t.Errorf("expected message 'wrapped error', got %s", wrappedErr.Message)
	}
	
	if wrappedErr.Err != originalErr {
		t.Error("expected wrapped error to contain original error")
	}
}