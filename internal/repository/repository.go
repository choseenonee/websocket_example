package repository

import (
	"context"
	"websockets/internal/models"
)

type ChatRepo interface {
	Create(ctx context.Context, chatCreate models.ChatCreate) (int, error)
	GetChatMessagesByPage(ctx context.Context, chatID, page int) ([]models.Message, error)
	GetChatsByName(ctx context.Context, name string, page int) ([]models.Chat, error)
	GetChatsByPage(ctx context.Context, page int) ([]models.Chat, error)
}
