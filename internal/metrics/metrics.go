package metrics

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	registerOnce     sync.Once
	requestLatency   *prometheus.HistogramVec
	requestTotal     *prometheus.CounterVec
	requestsInFlight prometheus.Gauge
)

func ensureMetrics() {
	registerOnce.Do(func() {
		requestLatency = prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: "api",
				Name:      "request_duration_seconds",
				Help:      "Duration of HTTP requests.",
				Buckets:   prometheus.DefBuckets,
			},
			[]string{"method", "route", "status"},
		)

		requestTotal = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "api",
				Name:      "requests_total",
				Help:      "Total number of HTTP requests processed.",
			},
			[]string{"method", "route", "status"},
		)

		requestsInFlight = prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: "api",
				Name:      "requests_in_flight",
				Help:      "Current number of in-flight requests.",
			},
		)

		prometheus.MustRegister(requestLatency, requestTotal, requestsInFlight)
	})
}

// Middleware instruments HTTP handlers with Prometheus metrics.
func Middleware(next http.Handler) http.Handler {
	ensureMetrics()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		recorder := &statusRecorder{ResponseWriter: w, status: http.StatusOK}

		requestsInFlight.Inc()
		defer requestsInFlight.Dec()

		next.ServeHTTP(recorder, r)

		route := chi.RouteContext(r.Context())
		pattern := r.URL.Path
		if route != nil {
			if rp := route.RoutePattern(); rp != "" {
				pattern = rp
			}
		}

		labels := []string{r.Method, pattern, strconv.Itoa(recorder.status)}

		duration := time.Since(start).Seconds()
		requestLatency.WithLabelValues(labels...).Observe(duration)
		requestTotal.WithLabelValues(labels...).Inc()
	})
}

// Handler exposes the Prometheus metrics endpoint.
func Handler() http.Handler {
	ensureMetrics()
	return promhttp.Handler()
}

// statusRecorder captures the status code written by a handler.
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (s *statusRecorder) WriteHeader(code int) {
	s.status = code
	s.ResponseWriter.WriteHeader(code)
}
