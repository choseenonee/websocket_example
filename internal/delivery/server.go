package delivery

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"websockets/internal/delivery/routers"
	"websockets/pkg/log"
)

func Start(logger *log.Logs, db *sqlx.DB) {
	r := gin.Default()
	
	routers.RegisterChatRouter(r, db)
	routers.RegisterWebSocketRouter(r, db, logger)

	if err := r.Run("0.0.0.0:3002"); err != nil {
		panic(fmt.Sprintf("error running client: %v", err.Error()))
	}
}
