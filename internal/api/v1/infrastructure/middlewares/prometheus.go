package middlewares

import (
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Metrics struct {
	requestsTotal    *prometheus.CounterVec
	requestDuration  *prometheus.HistogramVec
	responseSize     *prometheus.SummaryVec
	requestsInFlight prometheus.Gauge
	errorsTotal      *prometheus.CounterVec
}

func NewPrometheusMetrics() *Metrics {
	return &Metrics{
		requestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests",
			},
			[]string{"method", "path", "status"},
		),
		requestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "Duration of HTTP requests",
				Buckets: []float64{0.1, 0.3, 0.5, 1, 3, 5, 10},
			},
			[]string{"method", "path", "status"},
		),
		responseSize: promauto.NewSummaryVec(
			prometheus.SummaryOpts{
				Name: "http_response_size_bytes",
				Help: "Size of HTTP responses",
				Objectives: map[float64]float64{
					0.5:  0.05,
					0.9:  0.01,
					0.99: 0.001,
				},
			},
			[]string{"method", "path", "status"},
		),
		requestsInFlight: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "http_requests_in_flight",
				Help: "Current number of incoming HTTP requests",
			},
		),
		errorsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_errors_total",
				Help: "Total number of HTTP errors",
			},
			[]string{"method", "path", "status"},
		),
	}
}

func normalizePath(path string) string {
	re := regexp.MustCompile(`/(\d+)`)
	return re.ReplaceAllString(path, "/:id")
}

func (m *Metrics) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/metrics" {
			c.Next()
			return
		}

		start := time.Now()
		path := normalizePath(c.FullPath())

		if path == "" {
			path = "not_found"
		}

		m.requestsInFlight.Inc()
		defer m.requestsInFlight.Dec()

		c.Next()

		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Writer.Status())
		method := c.Request.Method

		labels := prometheus.Labels{
			"method": method,
			"path":   path,
			"status": status,
		}

		m.requestsTotal.With(labels).Inc()
		m.requestDuration.With(labels).Observe(duration)
		m.responseSize.With(labels).Observe(float64(c.Writer.Size()))

		if c.Writer.Status() >= http.StatusBadRequest {
			m.errorsTotal.With(labels).Inc()
		}
	}
}
