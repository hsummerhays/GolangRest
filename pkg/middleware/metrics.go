package middleware

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	httpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "http_request_duration_seconds",
		Help: "Duration of HTTP requests.",
	}, []string{"path", "method", "status"})

	httpRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests.",
	}, []string{"path", "method", "status"})
)

// MetricsMiddleware tracks the duration of HTTP requests and exposes it to Prometheus,
// while also logging the duration using slog.
func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap ResponseWriter to capture the status code
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		next.ServeHTTP(ww, r)

		duration := time.Since(start)

		// Record Prometheus metrics
		statusStr := strconv.Itoa(ww.Status())
		
		// For prometheus cardinality, we often group by route pattern rather than exact URL,
		// but for a simple project r.URL.Path is okay.
		path := r.URL.Path

		httpDuration.WithLabelValues(path, r.Method, statusStr).Observe(duration.Seconds())
		httpRequestsTotal.WithLabelValues(path, r.Method, statusStr).Inc()

		// Log using slog
		slog.Info("Request completed",
			"method", r.Method,
			"path", r.URL.Path,
			"status", ww.Status(),
			"duration_ms", duration.Milliseconds(),
			"duration_human", duration.String(),
			"bytes_written", ww.BytesWritten(),
			"ip", r.RemoteAddr,
		)
	})
}
