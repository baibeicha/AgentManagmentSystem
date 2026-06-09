package service

import (
	"AgentManagmentSystem/internal/gateway/domain"
	"context"
	"io"
	"time"
)

// AuthService handles user sessions, token issuance, and 2FA management.
type AuthService interface {
	Login(ctx context.Context, login, password, deviceID string) (accessToken, refreshToken string, err error)
	RefreshToken(ctx context.Context, refreshToken string) (newAccess, newRefresh string, err error)
	Logout(ctx context.Context, refreshToken string) error
	GetMe(ctx context.Context) (*domain.User, error)
	GetPermissions(ctx context.Context) ([]string, error)
	Setup2FA(ctx context.Context) (secret, otpauthURL string, err error)
	Verify2FA(ctx context.Context, code, actionID string) (status string, err error)
	Disable2FA(ctx context.Context, code string) error
}

// DeviceService manages host inventory, logical grouping, and infrastructure bootstrap.
type DeviceService interface {
	ListDevices(ctx context.Context) ([]domain.Device, error)
	GetDevice(ctx context.Context, deviceID string) (*domain.Device, error)
	DeleteDevice(ctx context.Context, deviceID string) error
	UpdateDeviceMetadata(ctx context.Context, deviceID string, alias, groupID string, tags []string) error
	CreateBootstrapToken(ctx context.Context) (token string, expiresAt time.Time, err error)
	EnrollAgent(ctx context.Context, provisioningKey, csrPEM string) (clientCertPEM, caCertPEM string, err error)
	TriggerDeviceUpdate(ctx context.Context, deviceID, targetVersion string) error
	ListDeviceGroups(ctx context.Context) ([]domain.DeviceGroup, error)
	CreateDeviceGroup(ctx context.Context, name, description string) error
}

// DiscoveryService tracks unmanaged assets found via active local network sweeps.
type DiscoveryService interface {
	ListDiscoveryCandidates(ctx context.Context) ([]domain.DiscoveryCandidate, error)
	ApproveDiscoveryCandidate(ctx context.Context, candidateID string) error
	IgnoreDiscoveryCandidate(ctx context.Context, candidateID string) error
}

// ReleaseService controls agent lifecycle binaries distribution.
type ReleaseService interface {
	GetAgentReleaseStream(ctx context.Context, os, arch string) (io.ReadCloser, error)
	GetLatestReleaseVersion(ctx context.Context) (*domain.AgentRelease, error)
	UploadAgentRelease(ctx context.Context, version, os, arch string, fileStream io.Reader) error
}

// MetricsService acts as a read-layer for timeseries telemetry engine (TimescaleDB).
type MetricsService interface {
	GetCurrentMetrics(ctx context.Context, deviceID string) (*domain.MetricsSnapshot, error)
	GetHistoricalMetrics(ctx context.Context, deviceID, metricType string, from, to time.Time, step string) ([][]string, error)
	GetActiveProcesses(ctx context.Context, deviceID string) ([]domain.Process, error)
}

// CommandService handles remote code execution and custom script deployment safely.
type CommandService interface {
	ExecuteCommand(ctx context.Context, deviceID, payload string) (*domain.CommandResult, error)
	ListScripts(ctx context.Context) ([]domain.ScriptTemplate, error)
	CreateScript(ctx context.Context, name, description, content, interpreter string) error
	GetScript(ctx context.Context, scriptID string) (*domain.ScriptTemplate, error)
	UpdateScript(ctx context.Context, scriptID string, name, description, content, interpreter string) error
	DeleteScript(ctx context.Context, scriptID string) error
}

// AutomationService controls alerting thresholds, playbooks state machines, and cron tasks.
type AutomationService interface {
	ListAlertRules(ctx context.Context) ([]domain.AlertRule, error)
	CreateAlertRule(ctx context.Context, rule domain.AlertRule) error
	UpdateAlertRule(ctx context.Context, ruleID string, rule domain.AlertRule) error
	DeleteAlertRule(ctx context.Context, ruleID string) error
	ListPlaybooks(ctx context.Context) ([]domain.Playbook, error)
	CreatePlaybook(ctx context.Context, name string, steps []domain.PlaybookStep) error
	GetPlaybook(ctx context.Context, playbookID string) (*domain.Playbook, error)
	UpdatePlaybook(ctx context.Context, playbookID string, name string, steps []domain.PlaybookStep) error
	DeletePlaybook(ctx context.Context, playbookID string) error
	ListCronJobs(ctx context.Context) ([]domain.CronJob, error)
	CreateCronJob(ctx context.Context, expression, scriptID, deviceGroupID string) error
	UpdateCronJob(ctx context.Context, cronID string, expression, scriptID, deviceGroupID string, isActive bool) error
	DeleteCronJob(ctx context.Context, cronID string) error
}

// NotificationService provisions dispatch pipelines like webhooks or messenger bots.
type NotificationService interface {
	ListNotificationChannels(ctx context.Context) ([]domain.NotificationChannel, error)
	CreateNotificationChannel(ctx context.Context, channel domain.NotificationChannel) error
	UpdateNotificationChannel(ctx context.Context, channelID string, channel domain.NotificationChannel) error
	DeleteNotificationChannel(ctx context.Context, channelID string) error
}

// IncidentService operates over the lifecycle of triggered infrastructure alerts (e.g. Dead Man's Switch).
type IncidentService interface {
	ListActiveIncidents(ctx context.Context) ([]domain.Incident, error)
	AcknowledgeIncident(ctx context.Context, incidentID string) error
	ResolveIncident(ctx context.Context, incidentID string) error
}

// TeamService manages multi-tenant administrative groups and RBAC assignments.
type TeamService interface {
	ListUsers(ctx context.Context) ([]domain.User, error)
	InviteUser(ctx context.Context, email, role string) error
	UpdateUserRoles(ctx context.Context, userID string, role string, allowedGroupIDs []string) error
	DeleteUser(ctx context.Context, userID string) error
}

// AuditService ensures historical traceability and structural system-wide configurations migration.
type AuditService interface {
	ListAuditLogs(ctx context.Context) ([]domain.AuditLog, error)
	VerifyAuditIntegrity(ctx context.Context) (status string, corruptedRowID string, err error)
	ExportConfig(ctx context.Context) (yamlManifest string, err error)
	ImportConfig(ctx context.Context, yamlManifest string) error
}
