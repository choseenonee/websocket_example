package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"websockets/internal/delivery/handlers"
	"websockets/internal/repository"
)

func RegisterChatRouter(r *gin.Engine, db *sqlx.DB) *gin.RouterGroup {
	chatRouter := r.Group("/chat")

	chatRepo := repository.InitChatRepo(db)
	hubHandler := handlers.InitChatHandler(chatRepo)

	chatRouter.POST("/create_chat", hubHandler.CreateChat)

	return chatRouter
}
