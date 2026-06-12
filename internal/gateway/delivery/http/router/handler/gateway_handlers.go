package handler

import (
	"AgentManagmentSystem/internal/gateway/delivery/http/dto"
	"AgentManagmentSystem/internal/gateway/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GatewayHandlers struct {
	auth       service.AuthService
	device     service.DeviceService
	discovery  service.DiscoveryService
	metrics    service.MetricsService
	command    service.CommandService
	automation service.AutomationService
	notify     service.NotificationService
	incident   service.IncidentService
	team       service.TeamService
	audit      service.AuditService
}

func (h *GatewayHandlers) GetAuthService() service.AuthService {
	return h.auth
}

func NewGatewayHandlers(
	auth service.AuthService, device service.DeviceService, discovery service.DiscoveryService, metrics service.MetricsService,
	command service.CommandService, automation service.AutomationService, notify service.NotificationService,
	incident service.IncidentService, team service.TeamService, audit service.AuditService,
) *GatewayHandlers {
	return &GatewayHandlers{
		auth:       auth,
		device:     device,
		discovery:  discovery,
		metrics:    metrics,
		command:    command,
		automation: automation,
		notify:     notify,
		incident:   incident,
		team:       team,
		audit:      audit,
	}
}

func (h *GatewayHandlers) Liveness(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "alive"})
}

func (h *GatewayHandlers) Readiness(c *gin.Context) {
	c.JSON(http.StatusOK, dto.ReadyResponse{
		Status: "ready",
		Components: map[string]string{
			"database":  "connected_mock",
			"redis":     "connected_mock",
			"kafka":     "connected_mock",
			"grpc_auth": "connected_mock",
		},
	})
}
