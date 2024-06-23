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
		err := rows.Scan(&message.ID, &message.Sender, &message.Content, &message.SendTimeStamp, &message.ChatID)
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

	row := tx.QueryRowContext(ctx, query, chatCreate.Name)

	var chatID int
	err = row.Scan(&chatID)
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			return 0, fmt.Errorf("err: %v, rbErr: %v", err.Error(), rbErr.Error())
		}
		return 0, err
	}

	if err = tx.Commit(); err != nil {
		return 0, err
	}

	return chatID, nil
}

func (c chatRepo) GetChatMessagesByPage(ctx context.Context, chatID, page int) ([]models.Message, error) {
	query := `SELECT id, sender, content, send_timestamp, chat_id FROM messages WHERE messages.chat_id = $1 ORDER BY send_timestamp DESC OFFSET $2 LIMIT $3`

	paginationPageLength := viper.GetInt(config.PaginationPageLength)

	rows, err := c.db.QueryContext(ctx, query, chatID, (page-1)*paginationPageLength, viper.GetInt(config.PaginationPageLength))
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

func (c chatRepo) CreateMessage(ctx context.Context, messageCreate models.MessageCreate) (int, error) {
	tx, err := c.db.Beginx()
	if err != nil {
		return 0, err
	}

	query := `INSERT INTO messages (sender, content, send_timestamp, chat_id) VALUES ($1, $2, $3, $4) RETURNING id;`

	row := tx.QueryRowContext(ctx, query, messageCreate.Sender, messageCreate.Content, messageCreate.SendTimeStamp,
		messageCreate.ChatID)

	var messageID int
	err = row.Scan(&messageID)
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			return 0, fmt.Errorf("err: %v, rbErr: %v", err.Error(), rbErr.Error())
		}
		return 0, err
	}

	if err = tx.Commit(); err != nil {
		return 0, err
	}

	return messageID, nil
}

func (c chatRepo) GetChatByID(ctx context.Context, chatID int) (models.Chat, error) {
	query := `SELECT name FROM chats WHERE id = $1`

	row := c.db.QueryRowContext(ctx, query, chatID)

	var chat models.Chat
	chat.ID = chatID
	err := row.Scan(&chat.Name)
	if err != nil {
		return models.Chat{}, err
	}

	return chat, nil
}
