package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"websockets/internal/models"
	"websockets/internal/repository"
)

func InitChatHandler(chatRepository repository.ChatRepo) *ChatHandler {
	return &ChatHandler{
		chatRepository: chatRepository,
	}
}

type ChatHandler struct {
	chatRepository repository.ChatRepo
}

func (ch *ChatHandler) CreateChat(c *gin.Context) {
	chatName := c.Query("name")

	chatID, err := ch.chatRepository.Create(c.Request.Context(), models.ChatCreate{Name: chatName})
	if err != nil {
		if strings.Contains(err.Error(), "unique") {
			c.JSON(http.StatusBadRequest, gin.H{"err": "chat with given name already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"chat_id": chatID})
}
