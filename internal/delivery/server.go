package delivery

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"websockets/internal/delivery/docs"
	"websockets/internal/delivery/middleware"
	"websockets/internal/delivery/routers"
	"websockets/internal/metrics"
	"websockets/pkg/log"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Start(logger *log.Logs, db *sqlx.DB, prometheusMetrics *metrics.PrometheusMetrics, tracer trace.Tracer) {
	r := gin.Default()

	mdw := middleware.InitMiddleware(logger)
	r.Use(mdw.CORSMiddleware())

	r.Static("/static", "../internal/delivery/docs/asyncapi/")
	r.LoadHTMLGlob("../internal/delivery/docs/asyncapi/*.html")
	r.GET("/asyncapi", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	docs.SwaggerInfo.BasePath = "/"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	routers.RegisterChatRouter(r, db, tracer)
	routers.RegisterWebSocketRouter(r, db, logger, prometheusMetrics, tracer)

	if err := r.Run("0.0.0.0:8080"); err != nil {
		panic(fmt.Sprintf("error running client: %v", err.Error()))
	}
}
