package delivery

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/prometheus/client_golang/prometheus"
	"websockets/internal/delivery/routers"
	"websockets/pkg/log"
)

func Start(logger *log.Logs, db *sqlx.DB, messagesCountMetric *prometheus.CounterVec) {
	r := gin.Default()

	routers.RegisterChatRouter(r, db)
	routers.RegisterWebSocketRouter(r, db, logger, messagesCountMetric)

	if err := r.Run("0.0.0.0:3002"); err != nil {
		panic(fmt.Sprintf("error running client: %v", err.Error()))
	}
}
