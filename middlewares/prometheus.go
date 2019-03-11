package middlewares

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	dflBuckets = []float64{300, 1200, 5000}
	labelNames = []string{"code", "method", "path"}
)

const (
	reqsName    = "http_requests_total"
	latencyName = "http_request_duration_milliseconds"
)

type PrometheusMiddleware struct {
	reqs    *prometheus.CounterVec
	latency *prometheus.HistogramVec
}

func NewPrometheusMiddleware(name string, buckets ...float64) *PrometheusMiddleware {
	labels := prometheus.Labels{"service": name}

	var m PrometheusMiddleware
	m.reqs = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name:        reqsName,
			Help:        "How many HTTP requests processed, partitioned by status code, method and HTTP path.",
			ConstLabels: labels,
		},
		labelNames,
	)
	prometheus.MustRegister(m.reqs)

	if len(buckets) == 0 {
		buckets = dflBuckets
	}
	m.latency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:        latencyName,
			Help:        "How long it took to process the request, partitioned by status code, method and HTTP path.",
			ConstLabels: labels,
			Buckets:     buckets,
		},
		labelNames,
	)
	prometheus.MustRegister(m.latency)

	return &m
}

func (m *PrometheusMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		m.reqs.WithLabelValues(strconv.Itoa(c.Writer.Status()), c.Request.Method, c.Request.URL.Path).Inc()
		m.latency.WithLabelValues(strconv.Itoa(c.Writer.Status()), c.Request.Method, c.Request.URL.Path).Observe(float64(time.Since(start).Nanoseconds()) / 1000000)
	}
}
