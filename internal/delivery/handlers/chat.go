package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
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

// @Summary Create chat
// @Tags chat
// @Accept  json
// @Produce  json
// @Param name query string true "Chat name"
// @Success 200 {object} int "Successfully created chat with id"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /chat [post]
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

// @Summary Get chats by page
// @Tags chat
// @Accept  json
// @Produce  json
// @Param page query int true "Page"
// @Success 200 {object} int "Successfully returned chats"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /chat/by_page [get]
func (ch *ChatHandler) GetChatsByPage(c *gin.Context) {
	pageRaw := c.Query("page")

	page, err := strconv.Atoi(pageRaw)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page query"})
		return
	}

	chats, err := ch.chatRepository.GetChatsByPage(c.Request.Context(), page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"chats": chats})
}

// @Summary Get chat messages
// @Tags chat
// @Accept  json
// @Produce  json
// @Param chat_id query string true "Chat id"
// @Param page query string true "Page"
// @Success 200 {object} int "Successfully returned messages"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /chat/messages [get]
func (ch *ChatHandler) GetChatMessagesByPage(c *gin.Context) {
	chatIDRaw := c.Query("chat_id")
	chatID, err := strconv.Atoi(chatIDRaw)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page query"})
		return
	}

	pageRaw := c.Query("page")
	page, err := strconv.Atoi(pageRaw)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page query"})
		return
	}

	messages, err := ch.chatRepository.GetChatMessagesByPage(c.Request.Context(), chatID, page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"messages": messages})
}
