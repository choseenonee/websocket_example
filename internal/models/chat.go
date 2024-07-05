package models

import "time"

type ChatCreate struct {
	Name string `json:"name"`
}

type Chat struct {
	ChatCreate
	ID int `json:"id"`
}

type MessageCreate struct {
	Sender        string    `json:"sender"`
	Content       []byte    `json:"content,omitempty"`
	SendTimeStamp time.Time `json:"timestamp"`
	ChatID        int       `json:"chat_id"`
}

type Message struct {
	MessageCreate
	ID int `json:"id"`
}

type MessageWithStringContent struct {
	Message
	Content string `json:"content"`
}

func InitMessageCreate(sender string, content []byte, timestamp time.Time, chatID int) *MessageCreate {
	return &MessageCreate{
		Sender:        sender,
		Content:       content,
		SendTimeStamp: timestamp,
		ChatID:        chatID,
	}
}
