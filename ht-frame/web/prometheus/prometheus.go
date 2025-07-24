package prometheus

import "github.com/prometheus/client_golang/prometheus"

var (
	// Requests 创建一个请求计数器指标
	Requests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "gin_http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "route"},
	)

	// Duration 创建一个请求持续时间指标
	Duration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "gin_http_duration_seconds",
			Help:    "Histogram of response durations for HTTP requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "route"},
	)
)

func init() {
	// 注册指标
	prometheus.MustRegister(Requests)
	prometheus.MustRegister(Duration)
}
