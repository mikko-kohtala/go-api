package services

import (
	"context"
	"runtime"
	"time"
)

type SystemStats struct {
	Uptime       time.Duration `json:"uptime"`
	MemoryUsage  uint64        `json:"memory_usage_mb"`
	NumGoroutine int           `json:"goroutines"`
	NumCPU       int           `json:"cpus"`
	GoVersion    string        `json:"go_version"`
}

type StatsService interface {
	GetSystemStats(ctx context.Context) (*SystemStats, error)
	GetAPIStats(ctx context.Context) (map[string]interface{}, error)
}

type statsService struct {
	startTime time.Time
}

func NewStatsService() StatsService {
	return &statsService{
		startTime: time.Now(),
	}
}

func (s *statsService) GetSystemStats(ctx context.Context) (*SystemStats, error) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	stats := &SystemStats{
		Uptime:       time.Since(s.startTime),
		MemoryUsage:  m.Alloc / 1024 / 1024, // Convert to MB
		NumGoroutine: runtime.NumGoroutine(),
		NumCPU:       runtime.NumCPU(),
		GoVersion:    runtime.Version(),
	}

	return stats, nil
}

func (s *statsService) GetAPIStats(ctx context.Context) (map[string]interface{}, error) {
	// In a real application, this would track API metrics
	// For now, return mock data
	activeConnections := runtime.NumGoroutine() - 2
	if activeConnections < 0 {
		activeConnections = 0
	}
	return map[string]interface{}{
		"total_requests":     1234,
		"requests_per_min":   42,
		"average_latency_ms": 15,
		"error_rate":         0.01,
		"active_connections": activeConnections,
	}, nil
}
