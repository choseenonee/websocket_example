package repository

import (
	"context"
	"websockets/internal/models"
)

type ChatRepo interface {
	Create(ctx context.Context, chatCreate models.ChatCreate) (int, error)
	GetChatByID(ctx context.Context, chatID int) (models.Chat, error)
	GetChatMessagesByPage(ctx context.Context, chatID, page int) ([]models.MessageWithStringContent, error)
	GetChatsByName(ctx context.Context, name string, page int) ([]models.Chat, error)
	GetChatsByPage(ctx context.Context, page int) ([]models.Chat, error)
	CreateMessage(ctx context.Context, message models.MessageCreate) (int, error)
}

type ChatGetterRepo interface {
	IsChatExists(chatID int) (bool, error)
}
