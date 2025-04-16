package metrics

type MetricsCollector interface {
	RecordRequest(namespace, user, model string, duration float64, status string)
	RecordTokenUsage(namespace, user, model string, inputTokens, outputTokens int64)
	RecordRateLimitHit(namespace, user, model, ruleType string)
	UpdateRateLimitTokens(namespace, user, model, ruleType string, tokens float64)
	UpdateQuotaUsage(namespace, model, quotaName, quotaType string, usage float64)
	UpdateQuotaLimit(namespace, model, quotaName, quotaType string, limit float64)
	RecordError(namespace, model, errorType string)
}

type DefaultMetricsCollector struct{}

func NewMetricsCollector() MetricsCollector {
	return &DefaultMetricsCollector{}
}

func (m *DefaultMetricsCollector) RecordRequest(namespace, user, model string, duration float64, status string) {
	requestTotal.WithLabelValues(namespace, user, model, status).Inc()
	requestDuration.WithLabelValues(namespace, user, model).Observe(duration)
}

func (m *DefaultMetricsCollector) RecordTokenUsage(namespace, user, model string, inputTokens, outputTokens int64) {
	tokenUsage.WithLabelValues(namespace, user, model, "input").Add(float64(inputTokens))
	tokenUsage.WithLabelValues(namespace, user, model, "output").Add(float64(outputTokens))
	tokenDistribution.WithLabelValues(namespace, user, model, "input").Observe(float64(inputTokens))
	tokenDistribution.WithLabelValues(namespace, user, model, "output").Observe(float64(outputTokens))
}

func (m *DefaultMetricsCollector) RecordRateLimitHit(namespace, user, model, ruleType string) {
	rateLimitHits.WithLabelValues(namespace, user, model, ruleType).Inc()
}

// TODO: 
func (m *DefaultMetricsCollector) UpdateRateLimitTokens(namespace, user, model, ruleType string, tokens float64) {
	rateLimitTokens.WithLabelValues(namespace, user, model, ruleType).Set(tokens)
}

// TODO:
func (m *DefaultMetricsCollector) UpdateQuotaUsage(namespace, model, quotaName, quotaType string, usage float64) {
	quotaUsage.WithLabelValues(namespace, model, quotaName, quotaType).Set(usage)
}

// TODO:
func (m *DefaultMetricsCollector) UpdateQuotaLimit(namespace, model, quotaName, quotaType string, limit float64) {
	quotaLimit.WithLabelValues(namespace, model, quotaName, quotaType).Set(limit)
}

// TODO: 
func (m *DefaultMetricsCollector) RecordError(namespace, model, errorType string) {
	errorTotal.WithLabelValues(namespace, model, errorType).Inc()
}
