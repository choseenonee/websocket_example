package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel/trace"
	"websockets/internal/delivery/ws"
	scheduler2 "websockets/internal/delivery/ws/scheduler"
	"websockets/internal/metrics"
	"websockets/internal/repository"
	"websockets/pkg/log"
)

func RegisterWebSocketRouter(r *gin.Engine, db *sqlx.DB, logger *log.Logs, prometheusMetrics *metrics.PrometheusMetrics,
	tracer trace.Tracer) *gin.RouterGroup {
	wsRouter := r.Group("/ws")

	chatRepo := repository.InitChatRepo(db)
	chatScheduler := scheduler2.InitChatRepoScheduler(chatRepo, logger)
	chatGetterRepo := repository.InitChatGetterRepo(chatRepo)
	hubScheduler := scheduler2.InitHubScheduler(logger, chatScheduler, chatGetterRepo, *prometheusMetrics)
	hubHandler := ws.InitHubHandler(hubScheduler, tracer)

	wsRouter.GET("/join_chat", hubHandler.JoinChat)

	return wsRouter
}
