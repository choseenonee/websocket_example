package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"log"
)

func InitPrometheusMetrics() *prometheus.CounterVec {
	messagesSent := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "chat_messages_sent_total",
			Help: "Total number of chat messages sent",
		},
		[]string{"status"},
	)

	// Регистрация метрики в Prometheus
	err := prometheus.Register(messagesSent)
	if err != nil {
		log.Fatalf("PROMETHEUS ERR: %v", err)
	}

	return messagesSent
}
