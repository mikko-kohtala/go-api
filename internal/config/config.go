package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Environment       string
	Port              string
	LogLevel          string
	AllowedOrigins    []string
	RateLimitRequests int
	RateLimitDuration time.Duration
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	ShutdownTimeout   time.Duration
}

func Load() *Config {
	return &Config{
		Environment:       getEnv("ENVIRONMENT", "development"),
		Port:              getEnv("PORT", "8080"),
		LogLevel:          getEnv("LOG_LEVEL", "info"),
		AllowedOrigins:    getEnvAsSlice("ALLOWED_ORIGINS", []string{"*"}),
		RateLimitRequests: getEnvAsInt("RATE_LIMIT_REQUESTS", 100),
		RateLimitDuration: getEnvAsDuration("RATE_LIMIT_DURATION", time.Minute),
		ReadTimeout:       getEnvAsDuration("READ_TIMEOUT", 15*time.Second),
		WriteTimeout:      getEnvAsDuration("WRITE_TIMEOUT", 15*time.Second),
		IdleTimeout:       getEnvAsDuration("IDLE_TIMEOUT", 60*time.Second),
		ShutdownTimeout:   getEnvAsDuration("SHUTDOWN_TIMEOUT", 30*time.Second),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := getEnv(key, "")
	if value, err := time.ParseDuration(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsSlice(key string, defaultValue []string) []string {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	return strings.Split(valueStr, ",")
}