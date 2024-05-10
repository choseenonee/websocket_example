package scheduler

import (
	"context"
	"encoding/json"
	"github.com/spf13/viper"
	"time"
	"websockets/internal/models"
	"websockets/internal/repository"
	"websockets/pkg/config"
	"websockets/pkg/log"
)

type RepoMessageCreator interface {
	CreateMessage(messageCreate models.MessageCreate)
}

type ChatRepoScheduler struct {
	messages chan models.MessageCreate
	chatRepo repository.ChatRepo
	logger   *log.Logs
}

func InitChatRepoScheduler(chatRepo repository.ChatRepo, logger *log.Logs) RepoMessageCreator {
	chatRepoScheduler := ChatRepoScheduler{
		messages: make(chan models.MessageCreate, 100),
		chatRepo: chatRepo,
		logger:   logger,
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

			_, err := c.chatRepo.CreateMessage(ctx, message)
			if err != nil {
				c.logger.Error(err.Error())
				writeMessageToFile(&message, c.logger)
			}
			cancel()
		}
	}
}

func (c *ChatRepoScheduler) CreateMessage(messageCreate models.MessageCreate) {
	c.messages <- messageCreate
}
