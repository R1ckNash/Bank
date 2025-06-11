package httpdelivery

import (
	"account/internal/metrics"
	"github.com/go-chi/chi/v5"
	"net/http"
	"time"
)

type statusWriter struct {
	http.ResponseWriter
	status int
}

func PrometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := &statusWriter{ResponseWriter: w, status: 200}
		next.ServeHTTP(ww, r)
		duration := time.Since(start).Seconds()
		routePattern := chi.RouteContext(r.Context()).RoutePattern()
		metrics.HTTPDuration.WithLabelValues(routePattern, r.Method, http.StatusText(ww.status)).Observe(duration)
	})
}
