package ws

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	scheduler2 "websockets/internal/delivery/ws/scheduler"
)

func InitHubHandler(scheduler *scheduler2.HubScheduler) *HubHandler {
	if scheduler == nil {
		panic("cant be nil scheduler")
	}

	return &HubHandler{
		scheduler: scheduler,
	}
}

type HubHandler struct {
	scheduler *scheduler2.HubScheduler
}

func (h *HubHandler) JoinChat(c *gin.Context) {
	chatIDRaw := c.Query("id")

	chatID, err := strconv.Atoi(chatIDRaw)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
	}

	err = h.scheduler.JoinChat(chatID, c.Writer, c.Request)
	if err != nil {
		if errors.Is(err, scheduler2.RoomNotFound) {
			c.JSON(http.StatusBadRequest, scheduler2.RoomNotFound)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
}
