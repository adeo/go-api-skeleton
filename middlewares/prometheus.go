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

	ContextKeyPrometheusURI = "prometheus_route"
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

		route := c.Request.URL.Path
		if v, ok := c.Get(ContextKeyPrometheusURI); ok {
			route = v.(string)
		}

		m.latency.WithLabelValues(strconv.Itoa(c.Writer.Status()), c.Request.Method, route).Observe(time.Since(start).Seconds())
	}
}
