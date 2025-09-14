package handlers

import (
    "bytes"
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestEcho(t *testing.T) {
    tests := []struct {
        name           string
        requestBody    string
        expectedStatus int
        expectedBody   string
    }{
        {
            name:           "valid request",
            requestBody:    `{"message":"hello world"}`,
            expectedStatus: http.StatusOK,
            expectedBody:   `{"message":"hello world"}`,
        },
        {
            name:           "empty message",
            requestBody:    `{"message":""}`,
            expectedStatus: http.StatusBadRequest,
            expectedBody:   `{"code":"validation_error","message":"validation failed","fields":{"message":"is required"}}`,
        },
        {
            name:           "missing message field",
            requestBody:    `{}`,
            expectedStatus: http.StatusBadRequest,
            expectedBody:   `{"code":"validation_error","message":"validation failed","fields":{"message":"is required"}}`,
        },
        {
            name:           "invalid JSON",
            requestBody:    `{"message":"hello"`,
            expectedStatus: http.StatusBadRequest,
            expectedBody:   `{"code":"invalid_request","message":"invalid JSON request"}`,
        },
        {
            name:           "unknown field",
            requestBody:    `{"message":"hello","unknown":"field"}`,
            expectedStatus: http.StatusBadRequest,
            expectedBody:   `{"code":"invalid_request","message":"invalid JSON request"}`,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            rr := httptest.NewRecorder()
            req := httptest.NewRequest(http.MethodPost, "/api/v1/echo", bytes.NewBufferString(tt.requestBody))
            req.Header.Set("Content-Type", "application/json")
            
            Echo(rr, req)
            
            if rr.Code != tt.expectedStatus {
                t.Errorf("expected status %d, got %d", tt.expectedStatus, rr.Code)
            }
            
            body := rr.Body.String()
            if body != tt.expectedBody+"\n" {
                t.Errorf("expected body %q, got %q", tt.expectedBody+"\n", body)
            }
        })
    }
}

func TestPing(t *testing.T) {
    rr := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodGet, "/api/v1/ping", nil)
    
    Ping(rr, req)
    
    if rr.Code != http.StatusOK {
        t.Errorf("expected status 200, got %d", rr.Code)
    }
    
    expectedBody := `{"pong":"ok"}`
    if rr.Body.String() != expectedBody+"\n" {
        t.Errorf("expected body %q, got %q", expectedBody+"\n", rr.Body.String())
    }
}

