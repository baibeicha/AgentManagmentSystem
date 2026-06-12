package handler

import (
	"AgentManagmentSystem/internal/gateway/delivery/http/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *GatewayHandlers) ListNotificationChannels(c *gin.Context) {
	channels, _ := h.notify.ListNotificationChannels(c.Request.Context())
	res := make([]dto.NotificationChannelResponse, len(channels))
	for i, ch := range channels {
		res[i] = dto.NotificationChannelResponse{ChannelID: ch.ChannelID, Name: ch.Name, Type: ch.Type, Destination: ch.Destination, IsActive: ch.IsActive}
	}
	c.JSON(http.StatusOK, res)
}

func (h *GatewayHandlers) CreateNotificationChannel(c *gin.Context) { c.Status(http.StatusCreated) }
func (h *GatewayHandlers) UpdateNotificationChannel(c *gin.Context) {
	c.JSON(http.StatusOK, dto.StatusResponse{Status: "ok"})
}
func (h *GatewayHandlers) DeleteNotificationChannel(c *gin.Context) {
	c.JSON(http.StatusOK, dto.SuccessResponse{Success: true})
}
