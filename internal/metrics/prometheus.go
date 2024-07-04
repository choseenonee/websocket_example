package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"log"
)

type PrometheusMetrics struct {
	MessagesSent       prometheus.Counter
	MessagesLatency    prometheus.Histogram
	ChatsOnline        prometheus.Gauge
	UsersOnline        prometheus.Gauge
	MessageToDbWorkers prometheus.Gauge
}

func InitPrometheusMetrics() *PrometheusMetrics {
	messagesSent := prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "chat_messages_sent_total",
			Help: "Total number of chat messages sent",
		},
	)
	err := prometheus.Register(messagesSent)
	if err != nil {
		log.Fatalf("PROMETHEUS ERR: %v", err)
	}

	messagesLatency := prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "chat_messages_latency_seconds",
			Help:    "Latency between message receiving and sending it to last chat member",
			Buckets: []float64{.005, .01, .025, .05, .075, .1, .25, .5, .75, 1.0},
		},
	)
	err = prometheus.Register(messagesLatency)
	if err != nil {
		log.Fatalf("PROMETHEUS ERR: %v", err)
	}

	chatsOnline := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "chats_alive_total",
			Help: "Total number of chats with online members",
		},
	)
	err = prometheus.Register(chatsOnline)
	if err != nil {
		log.Fatalf("PROMETHEUS ERR: %v", err)
	}

	usersOnline := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "users_online_total",
			Help: "Total number of chats with online members",
		},
	)
	err = prometheus.Register(usersOnline)
	if err != nil {
		log.Fatalf("PROMETHEUS ERR: %v", err)
	}

	messageToDbWorkers := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "message_to_db_workers_total",
			Help: "Total number of goroutines that sends new messages to db",
		},
	)
	err = prometheus.Register(messageToDbWorkers)
	if err != nil {
		log.Fatalf("PROMETHEUS ERR: %v", err)
	}

	return &PrometheusMetrics{
		MessagesSent:       messagesSent,
		MessagesLatency:    messagesLatency,
		ChatsOnline:        chatsOnline,
		UsersOnline:        usersOnline,
		MessageToDbWorkers: messageToDbWorkers,
	}
}
