package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// 请求总数（按 method、path、status 三维度）
	RequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "currencyflow",
			Subsystem: "http",
			Name:      "requests_total",
			Help:      "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	// 请求延迟分布（毫秒级桶）
	RequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "currencyflow",
			Subsystem: "http",
			Name:      "request_duration_seconds",
			Help:      "HTTP request duration in seconds",
			Buckets:   []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
		},
		[]string{"method", "path", "status"},
	)

	// 当前并发处理数
	InflightRequests = promauto.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "currencyflow",
			Subsystem: "http",
			Name:      "inflight_requests",
			Help:      "Number of HTTP requests currently being processed",
		},
	)
)
