package ws

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"strconv"
	scheduler2 "websockets/internal/delivery/ws/scheduler"
)

func InitHubHandler(scheduler *scheduler2.HubScheduler, tracer trace.Tracer) *HubHandler {
	if scheduler == nil {
		panic("cant be nil scheduler")
	}

	return &HubHandler{
		scheduler: scheduler,
		tracer:    tracer,
	}
}

type HubHandler struct {
	scheduler *scheduler2.HubScheduler
	tracer    trace.Tracer
}

func (h *HubHandler) JoinChat(c *gin.Context) {
	_, span := h.tracer.Start(c.Request.Context(), "Create chat")
	defer span.End()

	chatIDRaw := c.Query("id")

	chatID, err := strconv.Atoi(chatIDRaw)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String("Chat id not provided", err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())

		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
	}

	err = h.scheduler.JoinChat(chatID, c.Writer, c.Request)
	if err != nil {
		if errors.Is(err, scheduler2.RoomNotFound) {
			span.RecordError(err, trace.WithAttributes(
				attribute.String("Chat not found", err.Error())),
			)
			span.SetStatus(codes.Error, err.Error())

			c.JSON(http.StatusBadRequest, scheduler2.RoomNotFound)
			return
		}
		span.RecordError(err, trace.WithAttributes(
			attribute.String("Internal server error", err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())

		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	span.SetStatus(codes.Ok, "Successfully")
}
