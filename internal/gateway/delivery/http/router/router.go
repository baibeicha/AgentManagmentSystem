package router

import (
	"AgentManagmentSystem/internal/gateway/domain"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"AgentManagmentSystem/internal/gateway/client/mock"
	"AgentManagmentSystem/internal/gateway/delivery/http/dto"
	"AgentManagmentSystem/internal/gateway/delivery/http/middleware"
	"AgentManagmentSystem/internal/gateway/service"
)

func SetupRouter(log *slog.Logger) *gin.Engine {
	r := gin.New()

	r.Use(middleware.RequestLogger(log))
	r.Use(middleware.Recovery(log))

	authMock := mock.NewAuthMock()
	deviceMock := mock.NewDeviceMock()
	metricsMock := mock.NewMetricsMock()
	commandMock := mock.NewCommandMock()
	automationMock := mock.NewAutomationMock()
	notificationMock := mock.NewNotificationMock()
	incidentMock := mock.NewIncidentMock()
	teamMock := mock.NewTeamMock()
	auditMock := mock.NewAuditMock()

	h := NewGatewayHandlers(
		authMock, deviceMock, deviceMock, metricsMock, commandMock,
		automationMock, notificationMock, incidentMock, teamMock, auditMock,
	)

	r.GET("/healthz", h.Liveness)
	r.GET("/readyz", h.Readiness)

	api := r.Group("/api/v1")

	api.POST("/auth/login", h.Login)
	api.POST("/auth/refresh", h.Refresh)
	api.POST("/agent/enroll", h.EnrollAgent)

	api.GET("/downloads/agent/:os/:arch", h.DownloadAgent)
	api.GET("/downloads/agent/latest-version", h.LatestAgentVersion)

	protected := api.Group("", middleware.RequireAuth(authMock))
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

	ws := r.Group("/ws/v1", middleware.RequireAuth(authMock))
	{
		ws.GET("/stream", h.StreamAllMetrics)
		ws.GET("/terminal/:id", middleware.RequireRole("TEAM_ADMIN", "OPERATOR"), h.StreamTerminalShell)
	}

	return r
}

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

func (h *GatewayHandlers) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	access, refresh, err := h.auth.Login(c.Request.Context(), req.Login, req.Password, req.DeviceID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, dto.TokenResponse{AccessToken: access, RefreshToken: refresh})
}

func (h *GatewayHandlers) Refresh(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	access, refresh, err := h.auth.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, dto.TokenResponse{AccessToken: access, RefreshToken: refresh})
}

func (h *GatewayHandlers) Logout(c *gin.Context) {
	var req dto.LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	_ = h.auth.Logout(c.Request.Context(), req.RefreshToken)
	c.JSON(http.StatusOK, dto.SuccessResponse{Success: true})
}

func (h *GatewayHandlers) GetMe(c *gin.Context) {
	u, err := h.auth.GetMe(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, dto.UserProfileResponse{
		UserID: u.UserID, Login: u.Email, TenantID: u.Status, Role: u.Role,
	})
}

func (h *GatewayHandlers) GetPermissions(c *gin.Context) {
	perms, _ := h.auth.GetPermissions(c.Request.Context())
	c.JSON(http.StatusOK, dto.PermissionsResponse{Permissions: perms})
}

func (h *GatewayHandlers) Setup2FA(c *gin.Context) {
	secret, url, _ := h.auth.Setup2FA(c.Request.Context())
	c.JSON(http.StatusOK, dto.Setup2FAResponse{Secret: secret, OTPAuthURL: url})
}

func (h *GatewayHandlers) Verify2FA(c *gin.Context) {
	var req dto.Verify2FARequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	status, err := h.auth.Verify2FA(c.Request.Context(), req.Code, req.ActionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, dto.StatusResponse{Status: status})
}

func (h *GatewayHandlers) Disable2FA(c *gin.Context) {
	var req dto.Disable2FARequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	if err := h.auth.Disable2FA(c.Request.Context(), req.Code); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, dto.StatusResponse{Status: "disabled"})
}

func (h *GatewayHandlers) ListDevices(c *gin.Context) {
	devices, _ := h.device.ListDevices(c.Request.Context())
	res := make([]dto.DeviceResponse, len(devices))
	for i, d := range devices {
		res[i] = dto.DeviceResponse{
			DeviceID: d.DeviceID, Alias: d.Alias, OS: d.OS, Arch: d.Arch,
			IPAddress: d.IPAddress, Status: d.Status, AgentVersion: d.AgentVersion,
			GroupID: d.GroupID, Tags: d.Tags, LastHeartbeat: d.LastHeartbeat,
		}
	}
	c.JSON(http.StatusOK, res)
}

func (h *GatewayHandlers) GetDevice(c *gin.Context) {
	d, _ := h.device.GetDevice(c.Request.Context(), c.Param("id"))
	c.JSON(http.StatusOK, dto.DeviceResponse{
		DeviceID: d.DeviceID, Alias: d.Alias, OS: d.OS, Arch: d.Arch,
		IPAddress: d.IPAddress, Status: d.Status, AgentVersion: d.AgentVersion,
		LastHeartbeat: d.LastHeartbeat,
	})
}

func (h *GatewayHandlers) DeleteDevice(c *gin.Context) {
	_ = h.device.DeleteDevice(c.Request.Context(), c.Param("id"))
	c.JSON(http.StatusOK, dto.SuccessResponse{Success: true})
}

func (h *GatewayHandlers) UpdateDeviceMetadata(c *gin.Context) {
	var req dto.UpdateDeviceMetadataRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	_ = h.device.UpdateDeviceMetadata(c.Request.Context(), c.Param("id"), req.Alias, req.GroupID, req.Tags)
	c.JSON(http.StatusOK, dto.StatusResponse{Status: "ok"})
}

func (h *GatewayHandlers) CreateBootstrapToken(c *gin.Context) {
	token, exp, _ := h.device.CreateBootstrapToken(c.Request.Context())
	c.JSON(http.StatusOK, dto.BootstrapResponse{BootstrapToken: token, ExpiresAt: exp})
}

func (h *GatewayHandlers) EnrollAgent(c *gin.Context) {
	var req dto.AgentEnrollRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	client, ca, err := h.device.EnrollAgent(c.Request.Context(), req.ProvisioningKey, req.CSRPem)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, dto.AgentEnrollResponse{ClientCertificatePEM: client, CACertificatePEM: ca})
}

func (h *GatewayHandlers) TriggerDeviceUpdate(c *gin.Context) {
	var req dto.UpdateDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	_ = h.device.TriggerDeviceUpdate(c.Request.Context(), req.DeviceID, req.TargetVersion)
	c.Status(http.StatusAccepted)
}

func (h *GatewayHandlers) ListDeviceGroups(c *gin.Context) {
	groups, _ := h.device.ListDeviceGroups(c.Request.Context())
	res := make([]dto.DeviceGroupResponse, len(groups))
	for i, g := range groups {
		res[i] = dto.DeviceGroupResponse{GroupID: g.GroupID, Name: g.Name, Description: g.Description}
	}
	c.JSON(http.StatusOK, res)
}

func (h *GatewayHandlers) CreateDeviceGroup(c *gin.Context) {
	var req dto.CreateDeviceGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	_ = h.device.CreateDeviceGroup(c.Request.Context(), req.Name, req.Description)
	c.Status(http.StatusCreated)
}

func (h *GatewayHandlers) ListDiscoveryCandidates(c *gin.Context) {
	candidates, _ := h.discovery.ListDiscoveryCandidates(c.Request.Context())
	res := make([]dto.DiscoveryCandidateResponse, len(candidates))
	for i, cd := range candidates {
		res[i] = dto.DiscoveryCandidateResponse{
			CandidateID: cd.CandidateID, IPAddress: cd.IPAddress,
			MACAddress: cd.MACAddress, DiscoveredByDeviceID: cd.DiscoveredByDeviceID,
		}
	}
	c.JSON(http.StatusOK, res)
}

func (h *GatewayHandlers) ApproveDiscoveryCandidate(c *gin.Context) {
	_ = h.discovery.ApproveDiscoveryCandidate(c.Request.Context(), c.Param("id"))
	c.JSON(http.StatusOK, dto.StatusResponse{Status: "approved"})
}

func (h *GatewayHandlers) IgnoreDiscoveryCandidate(c *gin.Context) {
	_ = h.discovery.IgnoreDiscoveryCandidate(c.Request.Context(), c.Param("id"))
	c.JSON(http.StatusOK, dto.StatusResponse{Status: "ignored"})
}

func (h *GatewayHandlers) DownloadAgent(c *gin.Context) {
	stream, err := h.device.(service.ReleaseService).GetAgentReleaseStream(c.Request.Context(), c.Param("os"), c.Param("arch"))
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "release not found"})
		return
	}
	defer stream.Close()
	c.Header("Content-Disposition", "attachment; filename=agent.bin")
	c.Header("Content-Type", "application/octet-stream")
	_, _ = io.Copy(c.Writer, stream)
}

func (h *GatewayHandlers) LatestAgentVersion(c *gin.Context) {
	rel, _ := h.device.(service.ReleaseService).GetLatestReleaseVersion(c.Request.Context())
	c.JSON(http.StatusOK, dto.LatestReleaseResponse{Version: rel.Version, Sha256Checksum: rel.SHA256Checksum})
}

func (h *GatewayHandlers) UploadAgentRelease(c *gin.Context) {
	file, _ := c.FormFile("file")
	opened, _ := file.Open()
	defer opened.Close()
	_ = h.device.(service.ReleaseService).UploadAgentRelease(c.Request.Context(), c.PostForm("version"), c.PostForm("os"), c.PostForm("arch"), opened)
	c.Status(http.StatusCreated)
}

func (h *GatewayHandlers) GetCurrentMetrics(c *gin.Context) {
	m, _ := h.metrics.GetCurrentMetrics(c.Request.Context(), c.Param("id"))
	c.JSON(http.StatusOK, dto.MetricsCurrentResponse{
		CPUUsagePercent: m.CPUUsagePercent, RAMUsageBytes: m.RAMUsageBytes, RAMTotalBytes: m.RAMTotalBytes,
		DiskIOReadBytesSec: m.DiskIOReadBytesSec, DiskIOWriteBytesSec: m.DiskIOWriteBytesSec,
	})
}

func (h *GatewayHandlers) GetHistoricalMetrics(c *gin.Context) {
	metrics, _ := h.metrics.GetHistoricalMetrics(c.Request.Context(), c.Param("id"), c.Query("metric_type"), time.Now().Add(-1*time.Hour), time.Now(), "1m")
	c.JSON(http.StatusOK, dto.TimeseriesResponse{Metric: c.Query("metric_type"), Values: metrics})
}

func (h *GatewayHandlers) GetActiveProcesses(c *gin.Context) {
	processes, _ := h.metrics.GetActiveProcesses(c.Request.Context(), c.Param("id"))
	res := make([]dto.ProcessResponse, len(processes))
	for i, p := range processes {
		res[i] = dto.ProcessResponse{PID: p.PID, Name: p.Name, CPUUsage: p.CPUUsage, RAMUsageBytes: p.RAMUsageBytes}
	}
	c.JSON(http.StatusOK, res)
}

func (h *GatewayHandlers) ExecuteCommand(c *gin.Context) {
	var req dto.ExecuteCommandRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	res, err := h.command.ExecuteCommand(c.Request.Context(), req.DeviceID, req.Payload)
	if err != nil {
		c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, dto.CommandExecuteResponse{ExitCode: res.ExitCode, Stdout: res.Stdout, Stderr: res.Stderr})
}

func (h *GatewayHandlers) ListScripts(c *gin.Context) {
	scripts, _ := h.command.ListScripts(c.Request.Context())
	res := make([]dto.ScriptTemplateResponse, len(scripts))
	for i, s := range scripts {
		res[i] = dto.ScriptTemplateResponse{ScriptID: s.ScriptID, Name: s.Name, Description: s.Description, Content: s.Content, Interpreter: s.Interpreter}
	}
	c.JSON(http.StatusOK, res)
}

func (h *GatewayHandlers) CreateScript(c *gin.Context) {
	var req dto.ScriptTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	_ = h.command.CreateScript(c.Request.Context(), req.Name, req.Description, req.Content, req.Interpreter)
	c.Status(http.StatusCreated)
}

func (h *GatewayHandlers) GetScript(c *gin.Context) {
	s, _ := h.command.GetScript(c.Request.Context(), c.Param("id"))
	c.JSON(http.StatusOK, dto.ScriptTemplateResponse{ScriptID: s.ScriptID, Name: s.Name, Content: s.Content, Interpreter: s.Interpreter})
}

func (h *GatewayHandlers) UpdateScript(c *gin.Context) {
	var req dto.ScriptTemplateRequest
	_ = c.ShouldBindJSON(&req)
	_ = h.command.UpdateScript(c.Request.Context(), c.Param("id"), req.Name, req.Description, req.Content, req.Interpreter)
	c.JSON(http.StatusOK, dto.StatusResponse{Status: "ok"})
}

func (h *GatewayHandlers) DeleteScript(c *gin.Context) {
	_ = h.command.DeleteScript(c.Request.Context(), c.Param("id"))
	c.JSON(http.StatusOK, dto.SuccessResponse{Success: true})
}

func (h *GatewayHandlers) ListAlertRules(c *gin.Context) {
	rules, _ := h.automation.ListAlertRules(c.Request.Context())
	res := make([]dto.AlertRuleResponse, len(rules))
	for i, r := range rules {
		res[i] = dto.AlertRuleResponse{RuleID: r.RuleID, Name: r.Name, Metric: r.Metric, Operator: r.Operator, Threshold: r.Threshold, Duration: r.Duration, PlaybookID: r.PlaybookID}
	}
	c.JSON(http.StatusOK, res)
}

func (h *GatewayHandlers) CreateAlertRule(c *gin.Context) {
	var req dto.AlertRuleRequest
	_ = c.ShouldBindJSON(&req)
	_ = h.automation.CreateAlertRule(c.Request.Context(), domain.AlertRule{Name: req.Name, Metric: req.Metric, Operator: req.Operator, Threshold: req.Threshold, Duration: req.Duration, PlaybookID: req.PlaybookID})
	c.Status(http.StatusCreated)
}

func (h *GatewayHandlers) UpdateAlertRule(c *gin.Context) {
	var req dto.AlertRuleRequest
	_ = c.ShouldBindJSON(&req)
	_ = h.automation.UpdateAlertRule(c.Request.Context(), c.Param("id"), domain.AlertRule{Name: req.Name, Metric: req.Metric, Operator: req.Operator, Threshold: req.Threshold, Duration: req.Duration, PlaybookID: req.PlaybookID})
	c.JSON(http.StatusOK, dto.StatusResponse{Status: "ok"})
}

func (h *GatewayHandlers) DeleteAlertRule(c *gin.Context) {
	_ = h.automation.DeleteAlertRule(c.Request.Context(), c.Param("id"))
	c.JSON(http.StatusOK, dto.SuccessResponse{Success: true})
}

func (h *GatewayHandlers) ListPlaybooks(c *gin.Context) {
	playbooks, _ := h.automation.ListPlaybooks(c.Request.Context())
	res := make([]dto.PlaybookResponse, len(playbooks))
	for i, p := range playbooks {
		steps := make([]dto.PlaybookStep, len(p.Steps))
		for j, s := range p.Steps {
			steps[j] = dto.PlaybookStep{Type: s.Type, TargetID: s.TargetID, Delay: s.Delay}
		}
		res[i] = dto.PlaybookResponse{PlaybookID: p.PlaybookID, Name: p.Name, Steps: steps, CreatedAt: p.CreatedAt}
	}
	c.JSON(http.StatusOK, res)
}

func (h *GatewayHandlers) CreatePlaybook(c *gin.Context) {
	var req dto.PlaybookRequest
	_ = c.ShouldBindJSON(&req)
	c.Status(http.StatusCreated)
}

func (h *GatewayHandlers) GetPlaybook(c *gin.Context) {
	p, _ := h.automation.GetPlaybook(c.Request.Context(), c.Param("id"))
	c.JSON(http.StatusOK, dto.PlaybookResponse{PlaybookID: p.PlaybookID, Name: p.Name, CreatedAt: p.CreatedAt})
}

func (h *GatewayHandlers) UpdatePlaybook(c *gin.Context) {
	c.JSON(http.StatusOK, dto.StatusResponse{Status: "ok"})
}

func (h *GatewayHandlers) DeletePlaybook(c *gin.Context) {
	_ = h.automation.DeletePlaybook(c.Request.Context(), c.Param("id"))
	c.JSON(http.StatusOK, dto.SuccessResponse{Success: true})
}

func (h *GatewayHandlers) ListCronJobs(c *gin.Context) {
	crons, _ := h.automation.ListCronJobs(c.Request.Context())
	res := make([]dto.CronResponse, len(crons))
	for i, cr := range crons {
		res[i] = dto.CronResponse{CronID: cr.CronID, Expression: cr.Expression, ScriptID: cr.ScriptID, DeviceGroupID: cr.DeviceGroupID, IsActive: cr.IsActive, CreatedAt: cr.CreatedAt}
	}
	c.JSON(http.StatusOK, res)
}

func (h *GatewayHandlers) CreateCronJob(c *gin.Context) { c.Status(http.StatusCreated) }
func (h *GatewayHandlers) UpdateCronJob(c *gin.Context) {
	c.JSON(http.StatusOK, dto.StatusResponse{Status: "ok"})
}
func (h *GatewayHandlers) DeleteCronJob(c *gin.Context) {
	c.JSON(http.StatusOK, dto.SuccessResponse{Success: true})
}

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

func (h *GatewayHandlers) ListUsers(c *gin.Context) {
	users, _ := h.team.ListUsers(c.Request.Context())
	res := make([]dto.UserResponse, len(users))
	for i, u := range users {
		res[i] = dto.UserResponse{UserID: u.UserID, Email: u.Email, Role: u.Role, Status: u.Status}
	}
	c.JSON(http.StatusOK, res)
}

func (h *GatewayHandlers) InviteUser(c *gin.Context) {
	c.JSON(http.StatusOK, dto.StatusResponse{Status: "invited"})
}
func (h *GatewayHandlers) UpdateUserRoles(c *gin.Context) {
	c.JSON(http.StatusOK, dto.StatusResponse{Status: "ok"})
}
func (h *GatewayHandlers) DeleteUser(c *gin.Context) {
	c.JSON(http.StatusOK, dto.SuccessResponse{Success: true})
}

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

func (h *GatewayHandlers) StreamAllMetrics(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "WebSocket upgrade mock. Real channels requires gorilla/websocket implementation."})
}

func (h *GatewayHandlers) StreamTerminalShell(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "WebSocket reverse-shell upgrade mock."})
}
