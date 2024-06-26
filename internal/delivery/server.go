package delivery

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"net/http"
	"websockets/internal/delivery/docs"
	"websockets/internal/delivery/middleware"
	"websockets/internal/delivery/routers"
	"websockets/internal/metrics"
	"websockets/pkg/log"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Start(logger *log.Logs, db *sqlx.DB, prometheusMetrics *metrics.PrometheusMetrics) {
	r := gin.Default()

	r.Static("/static", "../internal/delivery/docs/asyncapi/")
	r.LoadHTMLGlob("../internal/delivery/docs/asyncapi/*.html")
	r.GET("/asyncapi", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	docs.SwaggerInfo.BasePath = "/"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	mdw := middleware.InitMiddleware(logger)
	r.Use(mdw.CORSMiddleware())

	routers.RegisterChatRouter(r, db)
	routers.RegisterWebSocketRouter(r, db, logger, prometheusMetrics)

	if err := r.Run("0.0.0.0:3002"); err != nil {
		panic(fmt.Sprintf("error running client: %v", err.Error()))
	}
}
