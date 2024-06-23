package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/prometheus/client_golang/prometheus"
	"websockets/internal/delivery/ws"
	scheduler2 "websockets/internal/delivery/ws/scheduler"
	"websockets/internal/repository"
	"websockets/pkg/log"
)

func RegisterWebSocketRouter(r *gin.Engine, db *sqlx.DB, logger *log.Logs, messagesCountMetric *prometheus.CounterVec) *gin.RouterGroup {
	wsRouter := r.Group("/ws")

	chatRepo := repository.InitChatRepo(db)
	chatScheduler := scheduler2.InitChatRepoScheduler(chatRepo, logger)
	chatGetterRepo := repository.InitChatGetterRepo(chatRepo)
	hubScheduler := scheduler2.InitHubScheduler(logger, chatScheduler, chatGetterRepo, messagesCountMetric)
	hubHandler := ws.InitHubHandler(hubScheduler)

	wsRouter.GET("/join_chat", hubHandler.JoinChat)

	return wsRouter
}
