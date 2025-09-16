package config

import (
	"os"
	"testing"
	"time"
)

func TestLoadConfig(t *testing.T) {
	// Save and restore original env vars
	originalEnv := map[string]string{}
	envVars := []string{
		"APP_ENV", "PORT", "REQUEST_TIMEOUT", "BODY_LIMIT_BYTES",
		"CORS_ALLOWED_ORIGINS", "CORS_ALLOWED_METHODS", "CORS_ALLOWED_HEADERS",
		"RATE_LIMIT_ENABLED", "RATE_LIMIT_PERIOD", "RATE_LIMIT",
		"CORS_STRICT", "COMPRESSION_LEVEL",
	}
	for _, key := range envVars {
		originalEnv[key] = os.Getenv(key)
		os.Unsetenv(key)
	}
	defer func() {
		for key, value := range originalEnv {
			if value != "" {
				os.Setenv(key, value)
			}
		}
	}()

	tests := []struct {
		name    string
		envVars map[string]string
		wantErr bool
		verify  func(*Config, *testing.T)
	}{
		{
			name:    "default values",
			envVars: map[string]string{},
			wantErr: false,
			verify: func(cfg *Config, t *testing.T) {
				if cfg.Env != "development" {
					t.Errorf("Expected Env=development, got %s", cfg.Env)
				}
				if cfg.Port != 8080 {
					t.Errorf("Expected Port=8080, got %d", cfg.Port)
				}
				if cfg.RequestTimeout != 15*time.Second {
					t.Errorf("Expected RequestTimeout=15s, got %v", cfg.RequestTimeout)
				}
				if cfg.BodyLimitBytes != 10485760 {
					t.Errorf("Expected BodyLimitBytes=10485760, got %d", cfg.BodyLimitBytes)
				}
				if !cfg.RateLimitEnabled {
					t.Error("Expected RateLimitEnabled=true")
				}
				if cfg.RateLimit != 100 {
					t.Errorf("Expected RateLimit=100, got %d", cfg.RateLimit)
				}
				if cfg.CompressionLevel != 5 {
					t.Errorf("Expected CompressionLevel=5, got %d", cfg.CompressionLevel)
				}
			},
		},
		{
			name: "production environment",
			envVars: map[string]string{
				"APP_ENV":               "production",
				"PORT":                  "8080",
				"REQUEST_TIMEOUT":       "30s",
				"CORS_ALLOWED_ORIGINS":  "https://example.com,https://app.example.com",
				"CORS_STRICT":           "true",
				"RATE_LIMIT":            "200",
				"COMPRESSION_LEVEL":     "9",
			},
			wantErr: false,
			verify: func(cfg *Config, t *testing.T) {
				if cfg.Env != "production" {
					t.Errorf("Expected Env=production, got %s", cfg.Env)
				}
				if cfg.Port != 8080 {
					t.Errorf("Expected Port=8080, got %d", cfg.Port)
				}
				if cfg.RequestTimeout != 30*time.Second {
					t.Errorf("Expected RequestTimeout=30s, got %v", cfg.RequestTimeout)
				}
				if len(cfg.CORSAllowedOrigins) != 2 {
					t.Errorf("Expected 2 CORS origins, got %d", len(cfg.CORSAllowedOrigins))
				}
				if !cfg.CORSStrict {
					t.Error("Expected CORSStrict=true")
				}
				if cfg.RateLimit != 200 {
					t.Errorf("Expected RateLimit=200, got %d", cfg.RateLimit)
				}
				if cfg.CompressionLevel != 9 {
					t.Errorf("Expected CompressionLevel=9, got %d", cfg.CompressionLevel)
				}
			},
		},
		{
			name: "rate limiting disabled",
			envVars: map[string]string{
				"RATE_LIMIT_ENABLED": "false",
			},
			wantErr: false,
			verify: func(cfg *Config, t *testing.T) {
				if cfg.RateLimitEnabled {
					t.Error("Expected RateLimitEnabled=false")
				}
			},
		},
		{
			name: "invalid port - zero",
			envVars: map[string]string{
				"PORT": "0",
			},
			wantErr: true,
		},
		{
			name: "invalid port - too high",
			envVars: map[string]string{
				"PORT": "70000",
			},
			wantErr: true,
		},
		{
			name: "invalid port - not a number",
			envVars: map[string]string{
				"PORT": "abc",
			},
			wantErr: true,
		},
		{
			name: "invalid request timeout",
			envVars: map[string]string{
				"REQUEST_TIMEOUT": "0s",
			},
			wantErr: true,
		},
		{
			name: "invalid body limit - zero",
			envVars: map[string]string{
				"BODY_LIMIT_BYTES": "0",
			},
			wantErr: true,
		},
		{
			name: "invalid body limit - too large",
			envVars: map[string]string{
				"BODY_LIMIT_BYTES": "1073741825", // > 1 GiB
			},
			wantErr: true,
		},
		{
			name: "invalid rate limit with enabled",
			envVars: map[string]string{
				"RATE_LIMIT_ENABLED": "true",
				"RATE_LIMIT":         "0",
			},
			wantErr: true,
		},
		{
			name: "invalid compression level - too low",
			envVars: map[string]string{
				"COMPRESSION_LEVEL": "0",
			},
			wantErr: true,
		},
		{
			name: "invalid compression level - too high",
			envVars: map[string]string{
				"COMPRESSION_LEVEL": "10",
			},
			wantErr: true,
		},
		{
			name: "custom CORS headers",
			envVars: map[string]string{
				"CORS_ALLOWED_METHODS": "GET,POST",
				"CORS_ALLOWED_HEADERS": "Content-Type,Authorization",
			},
			wantErr: false,
			verify: func(cfg *Config, t *testing.T) {
				if len(cfg.CORSAllowedMethods) != 2 {
					t.Errorf("Expected 2 CORS methods, got %d", len(cfg.CORSAllowedMethods))
				}
				if len(cfg.CORSAllowedHeaders) != 2 {
					t.Errorf("Expected 2 CORS headers, got %d", len(cfg.CORSAllowedHeaders))
				}
			},
		},
		{
			name: "custom rate limit period",
			envVars: map[string]string{
				"RATE_LIMIT_PERIOD": "5m",
			},
			wantErr: false,
			verify: func(cfg *Config, t *testing.T) {
				if cfg.RateLimitPeriod != "5m" {
					t.Errorf("Expected RateLimitPeriod=5m, got %s", cfg.RateLimitPeriod)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up environment
			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}
			defer func() {
				for key := range tt.envVars {
					os.Unsetenv(key)
				}
			}()

			// Load config
			cfg, err := Load()

			// Check error expectation
			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// If we expected an error, we're done
			if tt.wantErr {
				return
			}

			// Verify the config
			if cfg == nil {
				t.Fatal("Expected non-nil config")
			}

			if tt.verify != nil {
				tt.verify(cfg, t)
			}
		})
	}
}

func TestConfigValidation(t *testing.T) {
	// Test that validation catches specific edge cases
	tests := []struct {
		name    string
		envVars map[string]string
		wantErr bool
		errMsg  string
	}{
		{
			name: "negative port",
			envVars: map[string]string{
				"PORT": "-1",
			},
			wantErr: true,
			errMsg:  "invalid PORT",
		},
		{
			name: "negative timeout",
			envVars: map[string]string{
				"REQUEST_TIMEOUT": "-1s",
			},
			wantErr: true,
			errMsg:  "REQUEST_TIMEOUT must be > 0",
		},
		{
			name: "rate limit enabled but zero limit",
			envVars: map[string]string{
				"RATE_LIMIT_ENABLED": "true",
				"RATE_LIMIT":         "0",
			},
			wantErr: true,
			errMsg:  "RATE_LIMIT must be > 0 when RATE_LIMIT_ENABLED=true",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear all env vars first
			envVars := []string{
				"APP_ENV", "PORT", "REQUEST_TIMEOUT", "BODY_LIMIT_BYTES",
				"CORS_ALLOWED_ORIGINS", "CORS_ALLOWED_METHODS", "CORS_ALLOWED_HEADERS",
				"RATE_LIMIT_ENABLED", "RATE_LIMIT_PERIOD", "RATE_LIMIT",
				"CORS_STRICT", "COMPRESSION_LEVEL",
			}
			for _, key := range envVars {
				os.Unsetenv(key)
			}

			// Set test env vars
			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}
			defer func() {
				for key := range tt.envVars {
					os.Unsetenv(key)
				}
			}()

			cfg, err := Load()

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error containing '%s', got nil", tt.errMsg)
				} else if tt.errMsg != "" && err.Error() != tt.errMsg {
					// Check if error message contains expected text
					if err.Error() != tt.errMsg {
						t.Errorf("Expected error '%s', got '%s'", tt.errMsg, err.Error())
					}
				}
				if cfg != nil {
					t.Error("Expected nil config on error")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if cfg == nil {
					t.Error("Expected non-nil config")
				}
			}
		})
	}
}