package metrics

func InitMetrics() *PrometheusMetrics {
	InitMetricsHandler()
	return InitPrometheusMetrics()
}
