package ws

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"websockets/internal/ws/scheduler"
)

func InitHubHandler(scheduler *scheduler.HubScheduler) *HubHandler {
	if scheduler == nil {
		panic("cant be nil scheduler")
	}

	return &HubHandler{
		scheduler: scheduler,
	}
}

type HubHandler struct {
	scheduler *scheduler.HubScheduler
}

func (h *HubHandler) JoinChat(c *gin.Context) {
	chatIDRaw := c.Query("id")

	chatID, err := strconv.Atoi(chatIDRaw)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
	}

	err = h.scheduler.JoinChat(chatID, c.Writer, c.Request)
	if err != nil {
		if errors.Is(err, scheduler.RoomNotFound) {
			c.JSON(http.StatusBadRequest, scheduler.RoomNotFound)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
}
