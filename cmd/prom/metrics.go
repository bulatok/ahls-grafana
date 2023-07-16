package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	HttpRequestNum      *prometheus.CounterVec
	HttpResponseStatus  *prometheus.CounterVec
	HttpRequestDuration *prometheus.HistogramVec
}

func NewMetrics() *Metrics {
	return &Metrics{
		HttpRequestNum: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "app",
				Name:      "http_requests_total",
				Help:      "Number of requests",
			},
			[]string{"path"},
		),
		HttpRequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: "app",
				Name:      "http_response_time_seconds",
				Help:      "Duration of HTTP requests.",
				Buckets:   []float64{0.1, 0.15, 0.2, 0.25, 0.3, 1, 3, 5},
			}, []string{"path", "status_code"}),
		HttpResponseStatus: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "app",
				Name:      "http_response_status",
				Help:      "Status code of HTTP response",
			},
			[]string{"status_code"},
		),
	}
}

func (m *Metrics) Register() error {
	if err := prometheus.Register(m.HttpResponseStatus); err != nil {
		return err
	}

	if err := prometheus.Register(m.HttpRequestNum); err != nil {
		return err
	}

	if err := prometheus.Register(m.HttpRequestDuration); err != nil {
		return err
	}
	return nil
}

type StatusCode int

func (s StatusCode) String() string { return strconv.Itoa(int(s)) }

type responseWriter struct {
	http.ResponseWriter
	statusCode StatusCode
}

func (wr *responseWriter) WriteHeader(statusCode int) {
	wr.statusCode = StatusCode(statusCode)
	wr.ResponseWriter.WriteHeader(statusCode)
}

func NewResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (m *Metrics) ApplyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		tNow := time.Now()

		rw := NewResponseWriter(w)
		next.ServeHTTP(rw, r)

		statusCode := rw.statusCode

		m.HttpRequestDuration.WithLabelValues(path, statusCode.String()).Observe(time.Since(tNow).Seconds())
		m.HttpResponseStatus.WithLabelValues(statusCode.String()).Inc()
		m.HttpRequestNum.WithLabelValues(path).Inc()
	})
}
