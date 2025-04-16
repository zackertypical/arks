// pkg/gateway/metrics/metrics.go

package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	requestTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "gateway_requests_total",
			Help: "Total number of requests processed by the gateway",
		},
		[]string{"namespace", "user", "model", "status"},
	)

	requestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "gateway_request_duration_seconds",
			Help:    "Request duration in seconds",
			Buckets: []float64{0.1, 1, 2, 5, 10, 20, 30, 45, 60},
		},
		[]string{"namespace", "user", "model"},
	)

	tokenUsage = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "gateway_token_usage",
			Help: "Token usage statistics",
		},
		[]string{"namespace", "user", "token", "type"}, // type: input/output
	)

	tokenDistribution = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "gateway_token_distribution",
			Help:    "Distribution of token usage per request",
			Buckets: prometheus.ExponentialBuckets(1, 2, 17), // [1, 2, 4, 8, 16, 32, 64, 128, 256, 512, 1024, 2048, 4096, 8192, 16384, 32768, 65536]
		},
		[]string{"namespace", "user", "token", "type"}, // type: input/output
	)

	rateLimitHits = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "gateway_rate_limit_hits_total",
			Help: "Total number of rate limit hits",
		},
		[]string{"namespace", "user", "model", "rule_type"},
	)

	rateLimitTokens = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gateway_rate_limit_tokens",
			Help: "Current number of available tokens",
		},
		[]string{"namespace", "user", "model", "rule_type"},
	)

	quotaUsage = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gateway_quota_usage",
			Help: "Current quota usage",
		},
		[]string{"namespace", "model", "quota_name", "quota_type"},
	)

	quotaLimit = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gateway_quota_limit",
			Help: "Current quota limit",
		},
		[]string{"namespace", "model", "quota_name", "quota_type"},
	)

	// cacheHits = promauto.NewCounterVec(
	// 	prometheus.CounterOpts{
	// 		Name: "gateway_cache_hits_total",
	// 		Help: "Total number of cache hits",
	// 	},
	// 	[]string{"cache_type"},
	// )

	// cacheMisses = promauto.NewCounterVec(
	// 	prometheus.CounterOpts{
	// 		Name: "gateway_cache_misses_total",
	// 		Help: "Total number of cache misses",
	// 	},
	// 	[]string{"cache_type"},
	// )

	// cacheLatency = promauto.NewHistogramVec(
	// 	prometheus.HistogramOpts{
	// 		Name:    "gateway_cache_latency_seconds",
	// 		Help:    "Cache operation latency in seconds",
	// 		Buckets: prometheus.DefBuckets,
	// 	},
	// 	[]string{"cache_type", "operation"},
	// )

	errorTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "gateway_errors_total",
			Help: "Total number of errors",
		},
		[]string{"namespace", "model", "error_type"},
	)
)
