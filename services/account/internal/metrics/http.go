// internal/metrics/http.go
package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	HTTPDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_request_duration_seconds",
			Help: "Duration of HTTP requests.",
		},
		[]string{"handler", "method", "status"},
	)
)

func Register() {
	prometheus.MustRegister(HTTPDuration)
}
