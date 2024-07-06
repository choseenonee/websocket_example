package repository

import (
	"context"
	"database/sql"
	"errors"
)

type chatGetterRepo struct {
	chatRepo ChatRepo
}

func InitChatGetterRepo(chatRepo ChatRepo) ChatGetterRepo {
	return chatGetterRepo{
		chatRepo: chatRepo,
	}
}

func (cg chatGetterRepo) IsChatExists(ctx context.Context, chatID int) (bool, error) {
	_, err := cg.chatRepo.GetChatByID(ctx, chatID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
