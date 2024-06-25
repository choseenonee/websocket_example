package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel/trace"
	"websockets/internal/delivery/handlers"
	"websockets/internal/repository"
)

func RegisterChatRouter(r *gin.Engine, db *sqlx.DB, tracer trace.Tracer) *gin.RouterGroup {
	chatRouter := r.Group("/chat")

	chatRepo := repository.InitChatRepo(db)
	hubHandler := handlers.InitChatHandler(chatRepo, tracer)

	chatRouter.POST("/", hubHandler.CreateChat)
	chatRouter.GET("/by_page", hubHandler.GetChatsByPage)
	chatRouter.GET("/messages", hubHandler.GetChatMessagesByPage)

	return chatRouter
}
