package router

import (
	"AgentManagmentSystem/internal/gateway/delivery/http/router/handler"
	"log/slog"

	"github.com/gin-gonic/gin"

	"AgentManagmentSystem/internal/gateway/delivery/http/middleware"
)

func SetupRouter(h *handler.GatewayHandlers) *gin.Engine {
	log := slog.Default()

	r := gin.New()

	r.Use(middleware.RequestLogger(log))
	r.Use(middleware.Recovery(log))
	r.Use(middleware.CorsMiddleware())

	r.GET("/healthz", h.Liveness)
	r.GET("/readyz", h.Readiness)

	api := r.Group("/api/v1")

	api.POST("/auth/register", h.Register)
	api.POST("/auth/login", h.Login)
	api.POST("/auth/refresh", h.Refresh)
	api.POST("/agent/enroll", h.EnrollAgent)

	api.GET("/downloads/agent/:os/:arch", h.DownloadAgent)
	api.GET("/downloads/agent/latest-version", h.LatestAgentVersion)

	protected := api.Group("", middleware.RequireAuth(h.GetAuthService()))
	{
		protected.POST("/auth/logout", h.Logout)
		protected.GET("/auth/me", h.GetMe)
		protected.GET("/auth/permissions", h.GetPermissions)
		protected.POST("/auth/2fa/setup", h.Setup2FA)
		protected.POST("/auth/2fa/verify", h.Verify2FA)
		protected.POST("/auth/2fa/disable", h.Disable2FA)

		protected.GET("/devices", h.ListDevices)
		protected.GET("/devices/:id", h.GetDevice)
		protected.GET("/device-groups", h.ListDeviceGroups)

		protected.GET("/metrics/:id/current", h.GetCurrentMetrics)
		protected.GET("/metrics/:id/history", h.GetHistoricalMetrics)
		protected.GET("/metrics/:id/processes", h.GetActiveProcesses)

		protected.GET("/alerts/rules", h.ListAlertRules)
		protected.GET("/playbooks", h.ListPlaybooks)
		protected.GET("/playbooks/:id", h.GetPlaybook)
		protected.GET("/cron", h.ListCronJobs)

		protected.GET("/incidents/active", h.ListActiveIncidents)
	}

	operator := protected.Group("", middleware.RequireRole("TEAM_ADMIN", "OPERATOR"))
	{
		operator.POST("/commands/execute", h.ExecuteCommand)
		operator.GET("/scripts", h.ListScripts)
		operator.GET("/scripts/:id", h.GetScript)

		operator.GET("/discovery/candidates", h.ListDiscoveryCandidates)

		operator.PUT("/incidents/:id/acknowledge", h.AcknowledgeIncident)
		operator.PUT("/incidents/:id/resolve", h.ResolveIncident)
	}

	admin := protected.Group("", middleware.RequireRole("TEAM_ADMIN"))
	{
		admin.DELETE("/devices/:id", h.DeleteDevice)
		admin.PUT("/devices/:id/metadata", h.UpdateDeviceMetadata)
		admin.POST("/devices/bootstrap", h.CreateBootstrapToken)
		admin.POST("/devices/update", h.TriggerDeviceUpdate)
		admin.POST("/device-groups", h.CreateDeviceGroup)

		admin.POST("/discovery/candidates/:id/approve", h.ApproveDiscoveryCandidate)
		admin.DELETE("/discovery/candidates/:id/ignore", h.IgnoreDiscoveryCandidate)

		admin.POST("/system/agent-releases", h.UploadAgentRelease)

		admin.POST("/scripts", h.CreateScript)
		admin.PUT("/scripts/:id", h.UpdateScript)
		admin.DELETE("/scripts/:id", h.DeleteScript)

		admin.POST("/alerts/rules", h.CreateAlertRule)
		admin.PUT("/alerts/rules/:id", h.UpdateAlertRule)
		admin.DELETE("/alerts/rules/:id", h.DeleteAlertRule)

		admin.POST("/playbooks", h.CreatePlaybook)
		admin.PUT("/playbooks/:id", h.UpdatePlaybook)
		admin.DELETE("/playbooks/:id", h.DeletePlaybook)

		admin.POST("/cron", h.CreateCronJob)
		admin.PUT("/cron/:id", h.UpdateCronJob)
		admin.DELETE("/cron/:id", h.DeleteCronJob)

		admin.GET("/notifications/channels", h.ListNotificationChannels)
		admin.POST("/notifications/channels", h.CreateNotificationChannel)
		admin.PUT("/notifications/channels/:id", h.UpdateNotificationChannel)
		admin.DELETE("/notifications/channels/:id", h.DeleteNotificationChannel)

		admin.GET("/users", h.ListUsers)
		admin.POST("/users", h.InviteUser)
		admin.PUT("/users/:id/roles", h.UpdateUserRoles)
		admin.DELETE("/users/:id", h.DeleteUser)

		admin.GET("/audit", h.ListAuditLogs)
		admin.POST("/audit/verify", h.VerifyAuditIntegrity)
		admin.GET("/export", h.ExportConfig)
		admin.POST("/import", h.ImportConfig)
	}

	ws := r.Group("/ws/v1", middleware.RequireAuth(h.GetAuthService()))
	{
		ws.GET("/stream", h.StreamAllMetrics)
		ws.GET("/terminal/:id", middleware.RequireRole("TEAM_ADMIN", "OPERATOR"), h.StreamTerminalShell)
	}

	return r
}
