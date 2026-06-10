package domain

import "time"

type User struct {
	UserID   string
	Email    string
	Role     string
	Status   string
	TenantID string
}

type Device struct {
	DeviceID      string
	Alias         string
	OS            string
	Arch          string
	IPAddress     string
	Status        string
	AgentVersion  string
	GroupID       string
	Tags          []string
	LastHeartbeat time.Time
}

type DeviceGroup struct {
	GroupID     string
	Name        string
	Description string
}

type DiscoveryCandidate struct {
	CandidateID          string
	IPAddress            string
	MACAddress           string
	DiscoveredByDeviceID string
}

type AgentRelease struct {
	Version        string
	SHA256Checksum string
}

type MetricsSnapshot struct {
	CPUUsagePercent     float64
	RAMUsageBytes       int64
	RAMTotalBytes       int64
	DiskIOReadBytesSec  int64
	DiskIOWriteBytesSec int64
}

type Process struct {
	PID           int
	Name          string
	CPUUsage      float64
	RAMUsageBytes int64
}

type CommandResult struct {
	ExitCode int
	Stdout   string
	Stderr   string
}

type ScriptTemplate struct {
	ScriptID    string
	Name        string
	Description string
	Content     string
	Interpreter string
}

type AlertRule struct {
	RuleID     string
	Name       string
	Metric     string
	Operator   string
	Threshold  float64
	Duration   string
	PlaybookID string
}

type PlaybookStep struct {
	Type     string
	TargetID string
	Delay    string
}

type Playbook struct {
	PlaybookID string
	Name       string
	Steps      []PlaybookStep
	CreatedAt  time.Time
}

type CronJob struct {
	CronID        string
	Expression    string
	ScriptID      string
	DeviceGroupID string
	IsActive      bool
	CreatedAt     time.Time
}

type NotificationChannel struct {
	ChannelID   string
	Name        string
	Type        string
	Destination string
	IsActive    bool
}

type Incident struct {
	IncidentID string
	DeviceID   string
	Type       string
	Severity   string
	Status     string
	CreatedAt  time.Time
}

type AuditLog struct {
	LogID      string
	Timestamp  time.Time
	UserID     string
	ActionType string
	DeviceID   string
	Payload    string
	Status     string
}
