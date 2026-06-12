package handler

import (
	"AgentManagmentSystem/internal/gateway/delivery/http/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *GatewayHandlers) ListActiveIncidents(c *gin.Context) {
	incidents, _ := h.incident.ListActiveIncidents(c.Request.Context())
	res := make([]dto.IncidentResponse, len(incidents))
	for i, in := range incidents {
		res[i] = dto.IncidentResponse{IncidentID: in.IncidentID, DeviceID: in.DeviceID, Type: in.Type, Severity: in.Severity, Status: in.Status, CreatedAt: in.CreatedAt}
	}
	c.JSON(http.StatusOK, res)
}

func (h *GatewayHandlers) AcknowledgeIncident(c *gin.Context) {
	_ = h.incident.AcknowledgeIncident(c.Request.Context(), c.Param("id"))
	c.JSON(http.StatusOK, dto.StatusResponse{Status: "ok"})
}

func (h *GatewayHandlers) ResolveIncident(c *gin.Context) {
	_ = h.incident.ResolveIncident(c.Request.Context(), c.Param("id"))
	c.JSON(http.StatusOK, dto.StatusResponse{Status: "ok"})
}
