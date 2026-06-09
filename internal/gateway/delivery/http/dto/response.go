package dto

import "time"

type SuccessResponse struct {
	Success bool `json:"success"`
}

type StatusResponse struct {
	Status string `json:"status"`
}

type ErrorResponse struct {
	Error    string `json:"error"`
	ActionID string `json:"action_id,omitempty"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type UserProfileResponse struct {
	UserID   string `json:"user_id"`
	Login    string `json:"login"`
	TenantID string `json:"tenant_id"`
	Role     string `json:"role"`
}

type PermissionsResponse struct {
	Permissions []string `json:"permissions"`
}

type Setup2FAResponse struct {
	Secret     string `json:"secret"`
	OTPAuthURL string `json:"otpauth_url"`
}

type DeviceResponse struct {
	DeviceID      string    `json:"device_id"`
	Alias         string    `json:"alias"`
	OS            string    `json:"os"`
	Arch          string    `json:"arch"`
	IPAddress     string    `json:"ip_address"`
	Status        string    `json:"status"` // online, offline
	AgentVersion  string    `json:"agent_version"`
	GroupID       string    `json:"group_id,omitempty"`
	Tags          []string  `json:"tags,omitempty"`
	LastHeartbeat time.Time `json:"last_heartbeat"`
}

type BootstrapResponse struct {
	BootstrapToken string    `json:"bootstrap_token"`
	ExpiresAt      time.Time `json:"expires_at"`
}

type AgentEnrollResponse struct {
	ClientCertificatePEM string `json:"client_certificate_pem"`
	CACertificatePEM     string `json:"ca_certificate_pem"`
}

type DeviceGroupResponse struct {
	GroupID     string `json:"group_id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type DiscoveryCandidateResponse struct {
	CandidateID          string `json:"candidate_id"`
	IPAddress            string `json:"ip_address"`
	MACAddress           string `json:"mac_address"`
	DiscoveredByDeviceID string `json:"discovered_by_device_id"`
}

type MetricsCurrentResponse struct {
	CPUUsagePercent     float64 `json:"cpu_usage_percent"`
	RAMUsageBytes       int64   `json:"ram_usage_bytes"`
	RAMTotalBytes       int64   `json:"ram_total_bytes"`
	DiskIOReadBytesSec  int64   `json:"disk_io_read_bytes_sec"`
	DiskIOWriteBytesSec int64   `json:"disk_io_write_bytes_sec"`
}

type ProcessResponse struct {
	PID           int     `json:"pid"`
	Name          string  `json:"name"`
	CPUUsage      float64 `json:"cpu_usage"`
	RAMUsageBytes int64   `json:"ram_usage_bytes"`
}

type CommandExecuteResponse struct {
	ExitCode int    `json:"exit_code"`
	Stdout   string `json:"stdout"`
	Stderr   string `json:"stderr"`
}

type ScriptTemplateResponse struct {
	ScriptID    string `json:"script_id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Content     string `json:"content"`
	Interpreter string `json:"interpreter"`
}

type AlertRuleResponse struct {
	RuleID     string  `json:"rule_id"`
	Name       string  `json:"name"`
	Metric     string  `json:"metric"`
	Operator   string  `json:"operator"`
	Threshold  float64 `json:"threshold"`
	Duration   string  `json:"duration"`
	PlaybookID string  `json:"playbook_id,omitempty"`
}

type IncidentResponse struct {
	IncidentID string    `json:"incident_id"`
	DeviceID   string    `json:"device_id"`
	Type       string    `json:"type"`
	Severity   string    `json:"severity"` // info, warning, critical
	Status     string    `json:"status"`   // OPEN, ACKNOWLEDGED, RESOLVED
	CreatedAt  time.Time `json:"created_at"`
}

type AuditVerifyResponse struct {
	IntegrityStatus string `json:"integrity_status"` // e.g., VALID, CORRUPTED
	CorruptedRowID  string `json:"corrupted_row_id,omitempty"`
}

type LatestReleaseResponse struct {
	Version        string `json:"version"`
	Sha256Checksum string `json:"sha256_checksum"`
}

type ReadyResponse struct {
	Status     string            `json:"status"`
	Components map[string]string `json:"components"` // e.g., {"database": "connected", "redis": "connected"}
}

type UserResponse struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	Status string `json:"status"`
}

type PlaybookResponse struct {
	PlaybookID string         `json:"playbook_id"`
	Name       string         `json:"name"`
	Steps      []PlaybookStep `json:"steps"`
	CreatedAt  time.Time      `json:"created_at"`
}

type CronResponse struct {
	CronID        string    `json:"cron_id"`
	Expression    string    `json:"expression"`
	ScriptID      string    `json:"script_id"`
	DeviceGroupID string    `json:"device_group_id"`
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
}

type NotificationChannelResponse struct {
	ChannelID   string `json:"channel_id"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Destination string `json:"destination"`
	IsActive    bool   `json:"is_active"`
}

type AuditLogResponse struct {
	LogID      string    `json:"log_id"`
	Timestamp  time.Time `json:"timestamp"`
	UserID     string    `json:"user_id"`
	ActionType string    `json:"action_type"`
	DeviceID   string    `json:"device_id,omitempty"`
	Payload    string    `json:"payload"`
	Status     string    `json:"status"`
}

type TimeseriesResponse struct {
	Metric string     `json:"metric"`
	Values [][]string `json:"values"`
}
