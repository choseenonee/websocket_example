package scheduler

import (
	"context"
	"encoding/json"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"sync/atomic"
	"time"
	"websockets/internal/models"
	"websockets/internal/repository"
	"websockets/pkg/config"
	"websockets/pkg/log"
)

type RepoMessageCreator interface {
	CreateMessage(messageCreate *models.MessageCreate)
}

type ChatRepoScheduler struct {
	messages     chan *models.MessageCreate
	chatRepo     repository.ChatRepo
	logger       *log.Logs
	workersCount atomic.Int32
	promMetrics  prometheus.Gauge
	tracer       trace.Tracer
}

func InitChatRepoScheduler(chatRepo repository.ChatRepo, logger *log.Logs,
	promMetrics prometheus.Gauge, tracer trace.Tracer) RepoMessageCreator {
	chatRepoScheduler := ChatRepoScheduler{
		messages:    make(chan *models.MessageCreate, 100),
		chatRepo:    chatRepo,
		logger:      logger,
		promMetrics: promMetrics,
		tracer:      tracer,
	}

	go chatRepoScheduler.run()

	return &chatRepoScheduler
}

func writeMessageToFile(message *models.MessageCreate, logger *log.Logs) {
	jsonMessage, err := json.Marshal(*message)
	if err != nil {
		panic(err.Error())
	}
	logger.Info(string(jsonMessage))
}

func (c *ChatRepoScheduler) run() {
	for {
		select {
		case message := <-c.messages:
			ctx, cancel := context.WithTimeout(context.Background(),
				time.Duration(viper.GetInt(config.DBTimeout))*time.Millisecond)
			ctx, span := c.tracer.Start(ctx, "Create message")

			_, err := c.chatRepo.CreateMessage(ctx, *message)
			if err != nil {
				span.RecordError(err, trace.WithAttributes(
					attribute.String("Error creating message, writing its to file", err.Error())),
				)
				span.SetStatus(codes.Error, err.Error())

				c.logger.Error(err.Error())
				writeMessageToFile(message, c.logger)
			}

			span.End()
			cancel()

			if len(c.messages) >= viper.GetInt(config.ChatRepoMessagesNewWorkerOn) && c.workersCount.Load() < 100 {
				c.promMetrics.Inc()
				c.workersCount.Add(1)
				go c.newRunWorker()
			}
		}
	}
}

func (c *ChatRepoScheduler) CreateMessage(messageCreate *models.MessageCreate) {
	c.messages <- messageCreate
}

func (c *ChatRepoScheduler) newRunWorker() {
	for {
		select {
		case message := <-c.messages:
			ctx, cancel := context.WithTimeout(context.Background(),
				time.Duration(viper.GetInt(config.DBTimeout))*time.Millisecond)
			ctx, span := c.tracer.Start(ctx, "Create message")

			_, err := c.chatRepo.CreateMessage(ctx, *message)
			if err != nil {
				span.RecordError(err, trace.WithAttributes(
					attribute.String("Error creating message, writing its to file", err.Error())),
				)
				span.SetStatus(codes.Error, err.Error())

				c.logger.Error(err.Error())
				writeMessageToFile(message, c.logger)
			}

			span.End()
			cancel()

			if len(c.messages) < viper.GetInt(config.ChatRepoMessagesNewWorkerOn) {
				c.promMetrics.Dec()
				c.workersCount.Add(-1)
				return
			}
		}
	}
}
