package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"websockets/internal/models"
	"websockets/pkg/config"
)

type chatRepo struct {
	db *sqlx.DB
}

func InitChatRepo(db *sqlx.DB) ChatRepo {
	return chatRepo{
		db: db,
	}
}

func (c chatRepo) parseMessages(rows *sql.Rows) ([]models.Message, error) {
	var messages []models.Message

	for rows.Next() {
		var message models.Message
		err := rows.Scan(&message.ID, &message.Sender, &message.Content, &message.SendTimeStamp)
		if err != nil {
			return nil, err
		}

		messages = append(messages, message)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}

func (c chatRepo) parseChats(rows *sql.Rows) ([]models.Chat, error) {
	var chats []models.Chat

	for rows.Next() {
		var chat models.Chat
		err := rows.Scan(&chat.ID, &chat.Name)
		if err != nil {
			return nil, err
		}

		chats = append(chats, chat)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return chats, nil
}

func (c chatRepo) Create(ctx context.Context, chatCreate models.ChatCreate) (int, error) {
	tx, err := c.db.Beginx()
	if err != nil {
		return 0, err
	}

	query := `INSERT INTO chats (name) VALUES ($1) RETURNING id;`

	res, err := tx.ExecContext(ctx, query, chatCreate.Name)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return 0, fmt.Errorf("err: %v, also rbErr: %v", err, rbErr)
		}
		return 0, err
	}

	chatID, err := res.LastInsertId()
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return 0, err
		}
		return 0, err
	}

	if err = tx.Commit(); err != nil {
		return 0, err
	}

	return int(chatID), nil
}

func (c chatRepo) GetChatMessagesByPage(ctx context.Context, chatID, page int) ([]models.Message, error) {
	query := `SELECT id, sender, content, send_timestamp FROM messages WHERE messages.chat_id = $1 OFFSET $2 LIMIT $3`

	paginationPageLength := viper.GetInt(config.PaginationPageLength)

	rows, err := c.db.QueryContext(ctx, query, chatID, (page-1)*paginationPageLength, page)
	if err != nil {
		return nil, err
	}

	return c.parseMessages(rows)
}

func (c chatRepo) GetChatsByName(ctx context.Context, name string, page int) ([]models.Chat, error) {
	query := `SELECT id, name FROM chats WHERE name LIKE $1 OFFSET $2 LIMIT $3`

	paginationPageLength := viper.GetInt(config.PaginationPageLength)

	rows, err := c.db.QueryContext(ctx, query, fmt.Sprintf("%v%", name), (page-1)*paginationPageLength, page)
	if err != nil {
		return nil, err
	}

	return c.parseChats(rows)
}

func (c chatRepo) GetChatsByPage(ctx context.Context, page int) ([]models.Chat, error) {
	query := `SELECT id, name FROM chats OFFSET $1 LIMIT $2`

	paginationPageLength := viper.GetInt(config.PaginationPageLength)

	rows, err := c.db.QueryContext(ctx, query, (page-1)*paginationPageLength, page)
	if err != nil {
		return nil, err
	}

	return c.parseChats(rows)
}
