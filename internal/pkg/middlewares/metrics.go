package middlewares

import (
	"Redioteka/internal/pkg/utils/log"
	"net/http"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	hits = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "hits",
		Help: "requests count to api method",
	}, []string{"status", "path"})
	errors = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "errors",
		Help: "error requests count",
	}, []string{"status", "path"})
	duration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "duration",
		Help:    "histogram of requests duration in seconds",
		Buckets: []float64{1e-6, 1e-4, 1e-3, 5e-3, 0.01, 0.025, 0.1, 0.5, 1, 2, 5, 10},
	}, []string{"path"})
)

func RegisterMetrics() {
	prometheus.MustRegister(hits, errors, duration)
}

func (m *GoMiddleware) MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uri := r.RequestURI
		if queryPos := strings.Index(uri, "?"); queryPos >= 0 {
			uri = uri[:queryPos+1] + "queryParams"
		}
		url := strings.Split(uri, "/")
		for i := range url {
			if _, err := strconv.Atoi(url[i]); err == nil {
				url[i] = "id"
			}
		}
		newURL := strings.Join(url, "/")

		sw := NewResponseWriter(w)
		timer := prometheus.NewTimer(duration.WithLabelValues(newURL))
		next.ServeHTTP(sw, r)

		if s := timer.ObserveDuration().Seconds(); s < 0 {
			log.Log.Debug("negative request duration")
		}
		hits.WithLabelValues(strconv.Itoa(sw.Status), newURL).Inc()
		if sw.Status >= 400 {
			errors.WithLabelValues(strconv.Itoa(sw.Status), newURL).Inc()
		}
	})
}

type StatusRecorder struct {
	http.ResponseWriter
	Status int
}

func (r *StatusRecorder) WriteHeader(status int) {
	r.Status = status
	r.ResponseWriter.WriteHeader(status)
}

func NewResponseWriter(w http.ResponseWriter) *StatusRecorder {
	return &StatusRecorder{w, http.StatusOK}
}
