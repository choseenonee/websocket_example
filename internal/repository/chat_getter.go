package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/spf13/viper"
	"time"
	"websockets/pkg/config"
)

type chatGetterRepo struct {
	chatRepo ChatRepo
}

func InitChatGetterRepo(chatRepo ChatRepo) ChatGetterRepo {
	return chatGetterRepo{
		chatRepo: chatRepo,
	}
}

func (cg chatGetterRepo) IsChatExists(chatID int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(viper.GetInt(config.DBTimeout))*time.Millisecond)
	defer cancel()

	_, err := cg.chatRepo.GetChatByID(ctx, chatID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
