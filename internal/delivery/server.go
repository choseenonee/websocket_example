package delivery

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/prometheus/client_golang/prometheus"
	"websockets/internal/delivery/docs"
	"websockets/internal/delivery/middleware"
	"websockets/internal/delivery/routers"
	"websockets/pkg/log"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Start(logger *log.Logs, db *sqlx.DB, messagesCountMetric *prometheus.CounterVec) {
	r := gin.Default()

	docs.SwaggerInfo.BasePath = "/"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	mdw := middleware.InitMiddleware(logger)
	r.Use(mdw.CORSMiddleware())

	routers.RegisterChatRouter(r, db)
	routers.RegisterWebSocketRouter(r, db, logger, messagesCountMetric)

	if err := r.Run("0.0.0.0:3002"); err != nil {
		panic(fmt.Sprintf("error running client: %v", err.Error()))
	}
}
