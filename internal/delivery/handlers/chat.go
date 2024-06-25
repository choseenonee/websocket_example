package handlers

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"strconv"
	"strings"
	"websockets/internal/models"
	"websockets/internal/repository"
)

func InitChatHandler(chatRepository repository.ChatRepo, tracer trace.Tracer) *ChatHandler {
	return &ChatHandler{
		chatRepository: chatRepository,
		tracer:         tracer,
	}
}

type ChatHandler struct {
	chatRepository repository.ChatRepo
	tracer         trace.Tracer
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
	ctx, span := ch.tracer.Start(c.Request.Context(), "Create chat")
	defer span.End()

	chatName := c.Query("name")

	chatID, err := ch.chatRepository.Create(ctx, models.ChatCreate{Name: chatName})
	if err != nil {
		if strings.Contains(err.Error(), "unique") {
			span.RecordError(err, trace.WithAttributes(
				attribute.String("Chat already exists", err.Error())),
			)
			span.SetStatus(codes.Error, err.Error())

			c.JSON(http.StatusBadRequest, gin.H{"err": "chat with given name already exists"})
			return
		}

		span.RecordError(err, trace.WithAttributes(
			attribute.String("Internal server error", err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())

		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	span.SetStatus(codes.Ok, "Successfully")

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
	ctx, span := ch.tracer.Start(c.Request.Context(), "Create chat")
	defer span.End()

	pageRaw := c.Query("page")

	page, err := strconv.Atoi(pageRaw)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String("Page not provided", err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())

		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page query"})
		return
	}

	chats, err := ch.chatRepository.GetChatsByPage(ctx, page)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String("Internal server error", err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())

		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	span.SetStatus(codes.Ok, "Successfully")

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
	ctx, span := ch.tracer.Start(c.Request.Context(), "Create chat")
	defer span.End()

	chatIDRaw := c.Query("chat_id")
	chatID, err := strconv.Atoi(chatIDRaw)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page query"})
		return
	}

	pageRaw := c.Query("page")
	page, err := strconv.Atoi(pageRaw)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String("Page not provided", err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())

		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page query"})
		return
	}

	messages, err := ch.chatRepository.GetChatMessagesByPage(ctx, chatID, page)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String("Internal server error", err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())

		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	span.SetStatus(codes.Ok, "Successfully")

	c.JSON(http.StatusOK, gin.H{"messages": messages})
}
