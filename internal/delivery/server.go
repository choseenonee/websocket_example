package delivery

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel/trace"
	"websockets/internal/delivery/middleware"
	"websockets/internal/delivery/routers"
	"websockets/internal/metrics"
	"websockets/pkg/log"
)

func Start(logger *log.Logs, db *sqlx.DB, prometheusMetrics *metrics.PrometheusMetrics, tracer trace.Tracer) {
	r := gin.Default()

	mdw := middleware.InitMiddleware(logger)
	r.Use(mdw.CORSMiddleware())

	routers.RegisterStatic(r)

	routers.RegisterChatRouter(r, db, tracer)
	routers.RegisterWebSocketRouter(r, db, logger, prometheusMetrics, tracer)

	if err := r.Run("0.0.0.0:3002"); err != nil {
		panic(fmt.Sprintf("error running client: %v", err.Error()))
	}
}
