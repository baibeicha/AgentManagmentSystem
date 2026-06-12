package handler

import (
	"AgentManagmentSystem/internal/gateway/delivery/http/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *GatewayHandlers) ListAuditLogs(c *gin.Context) {
	logs, _ := h.audit.ListAuditLogs(c.Request.Context())
	res := make([]dto.AuditLogResponse, len(logs))
	for i, l := range logs {
		res[i] = dto.AuditLogResponse{LogID: l.LogID, Timestamp: l.Timestamp, UserID: l.UserID, ActionType: l.ActionType, DeviceID: l.DeviceID, Payload: l.Payload, Status: l.Status}
	}
	c.JSON(http.StatusOK, res)
}

func (h *GatewayHandlers) VerifyAuditIntegrity(c *gin.Context) {
	status, corruptedID, _ := h.audit.VerifyAuditIntegrity(c.Request.Context())
	c.JSON(http.StatusOK, dto.AuditVerifyResponse{IntegrityStatus: status, CorruptedRowID: corruptedID})
}

func (h *GatewayHandlers) ExportConfig(c *gin.Context) {
	manifest, _ := h.audit.ExportConfig(c.Request.Context())
	c.Header("Content-Type", "application/x-yaml")
	c.String(http.StatusOK, manifest)
}

func (h *GatewayHandlers) ImportConfig(c *gin.Context) {
	c.Status(http.StatusOK)
}
