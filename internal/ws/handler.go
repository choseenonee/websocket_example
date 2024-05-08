package ws

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitHubHandler(scheduler *HubScheduler) *HubHandler {
	if scheduler == nil {
		panic("cant be nil scheduler")
	}

	return &HubHandler{
		scheduler: scheduler,
	}
}

type HubHandler struct {
	scheduler *HubScheduler
}

func (h *HubHandler) CreateRoom(c *gin.Context) {
	roomName := c.Query("name")
	err := h.scheduler.CreateRoom(roomName, c.Writer, c.Request)
	if err != nil {
		if errors.Is(err, RoomAlreadyExists) {
			c.JSON(http.StatusBadRequest, RoomAlreadyExists)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	return
}

func (h *HubHandler) JoinRoom(c *gin.Context) {
	roomName := c.Query("name")
	err := h.scheduler.JoinRoom(roomName, c.Writer, c.Request)
	if err != nil {
		if errors.Is(err, RoomNotFound) {
			c.JSON(http.StatusBadRequest, RoomNotFound)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	return
}

func (h *HubHandler) LeaveRoom(c *gin.Context) {
	// TODO: implement me, just call the removeGarbageConn func
}

func (h *HubHandler) GetRooms(c *gin.Context) {
	// TODO: implement me
}
