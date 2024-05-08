package internal

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"websockets/internal/ws"
	"websockets/pkg/log"
)

func Start(logger *log.Logs) {
	r := gin.Default()

	wsRouter := r.Group("/ws")

	hubScheduler := ws.InitHubScheduler(logger)
	hubHandler := ws.InitHubHandler(hubScheduler)

	wsRouter.GET("/create_room", hubHandler.CreateRoom)
	wsRouter.GET("/join_room", hubHandler.JoinRoom)

	if err := r.Run("0.0.0.0:3002"); err != nil {
		panic(fmt.Sprintf("error running client: %v", err.Error()))
	}
}
