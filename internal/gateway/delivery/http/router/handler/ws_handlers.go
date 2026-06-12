package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *GatewayHandlers) StreamAllMetrics(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "WebSocket upgrade mock. Real channels requires gorilla/websocket implementation."})
}

func (h *GatewayHandlers) StreamTerminalShell(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "WebSocket reverse-shell upgrade mock."})
}
