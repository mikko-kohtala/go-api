package telemetry

import (
	"github.com/prometheus/client_golang/prometheus"
)

// Metrics holds all application metrics
type Metrics struct {
	RequestsTotal   *prometheus.CounterVec
	RequestDuration *prometheus.HistogramVec
	ActiveRequests  prometheus.Gauge
	ErrorsTotal     *prometheus.CounterVec
	Registry        *prometheus.Registry
}

// NewMetrics creates and registers all metrics
func NewMetrics() *Metrics {
	reg := prometheus.NewRegistry()

	requestsTotal := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request latencies in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path", "status"},
	)

	activeRequests := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_requests_active",
			Help: "Number of active HTTP requests",
		},
	)

	errorsTotal := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "errors_total",
			Help: "Total number of errors",
		},
		[]string{"type"},
	)

	// Register metrics with custom registry
	reg.MustRegister(requestsTotal, requestDuration, activeRequests, errorsTotal)

	return &Metrics{
		RequestsTotal:   requestsTotal,
		RequestDuration: requestDuration,
		ActiveRequests:  activeRequests,
		ErrorsTotal:     errorsTotal,
		Registry:        reg,
	}
}