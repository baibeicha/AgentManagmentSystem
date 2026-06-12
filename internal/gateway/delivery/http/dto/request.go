package dto

type RegistrationRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginRequest struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
	DeviceID string `json:"device_id,omitempty"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type Verify2FARequest struct {
	Code     string `json:"code" binding:"required,len=6"`
	ActionID string `json:"action_id" binding:"required"`
}

type Disable2FARequest struct {
	Code string `json:"code" binding:"required,len=6"`
}

type UpdateDeviceMetadataRequest struct {
	Alias   string   `json:"alias,omitempty"`
	GroupID string   `json:"group_id,omitempty" binding:"omitempty"`
	Tags    []string `json:"tags,omitempty"`
}

type AgentEnrollRequest struct {
	ProvisioningKey string `json:"provisioning_key" binding:"required"`
	CSRPem          string `json:"csr_pem" binding:"required"`
}

type UpdateDeviceRequest struct {
	DeviceID      string `json:"device_id" binding:"required"`
	TargetVersion string `json:"target_version" binding:"required"`
}

type CreateDeviceGroupRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description,omitempty"`
}

type ExecuteCommandRequest struct {
	DeviceID string `json:"device_id" binding:"required"`
	Payload  string `json:"payload" binding:"required"`
}

type ScriptTemplateRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description,omitempty"`
	Content     string `json:"content" binding:"required"`
	Interpreter string `json:"interpreter" binding:"required,oneof=bash sh powershell cmd"`
}

type AlertRuleRequest struct {
	Name       string  `json:"name" binding:"required"`
	Metric     string  `json:"metric" binding:"required,oneof=cpu ram disk network"`
	Operator   string  `json:"operator" binding:"required,oneof=> < == >="`
	Threshold  float64 `json:"threshold" binding:"required"`
	Duration   string  `json:"duration" binding:"required"`
	PlaybookID string  `json:"playbook_id,omitempty" binding:"omitempty"`
}

type PlaybookRequest struct {
	Name  string         `json:"name" binding:"required"`
	Steps []PlaybookStep `json:"steps" binding:"required"`
}

type PlaybookStep struct {
	Type     string `json:"type" binding:"required"`
	TargetID string `json:"target_id,omitempty"`
	Delay    string `json:"delay,omitempty"`
}

type CronRequest struct {
	Expression    string `json:"expression" binding:"required"`
	ScriptID      string `json:"script_id" binding:"required"`
	DeviceGroupID string `json:"device_group_id" binding:"required"`
}

type NotificationChannelRequest struct {
	Name        string `json:"name" binding:"required"`
	Type        string `json:"type" binding:"required,oneof=telegram discord webhook"`
	Destination string `json:"destination" binding:"required"`
	IsActive    bool   `json:"is_active"`
}

type InviteUserRequest struct {
	Email string `json:"email" binding:"required,email"`
	Role  string `json:"role" binding:"required,oneof=TEAM_ADMIN OPERATOR VIEWER"`
}

type UpdateUserRolesRequest struct {
	Role            string   `json:"role" binding:"required,oneof=TEAM_ADMIN OPERATOR VIEWER"`
	AllowedGroupIDs []string `json:"allowed_group_ids,omitempty" binding:"omitempty,dive"`
}

// UploadAgentReleaseRequest использует теги 'form' вместо 'json', так как
// файлы передаются через multipart/form-data
type UploadAgentReleaseRequest struct {
	Version string `form:"version" binding:"required"`
	OS      string `form:"os" binding:"required,oneof=linux windows darwin"`
	Arch    string `form:"arch" binding:"required,oneof=amd64 arm64"`
}
