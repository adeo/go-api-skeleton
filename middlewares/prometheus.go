package middlewares

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	labelNames = []string{"status", "method", "uri"}
)

const (
	latencyName = "http_server_requests_seconds"
)

type PrometheusMiddleware struct {
	latency *prometheus.HistogramVec
}

func NewPrometheusMiddleware(name string) *PrometheusMiddleware {
	labels := prometheus.Labels{"service": name}

	m := &PrometheusMiddleware{
		latency: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:        latencyName,
				Help:        "How long it took to process the request, partitioned by status code, method and HTTP path.",
				ConstLabels: labels,
			},
			labelNames,
		),
	}

	prometheus.MustRegister(m.latency)

	return m
}

func (m *PrometheusMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		m.latency.WithLabelValues(strconv.Itoa(c.Writer.Status()), c.Request.Method, c.Request.URL.Path).Observe(time.Since(start).Seconds())
	}
}
