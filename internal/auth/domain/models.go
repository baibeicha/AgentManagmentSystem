package domain

import (
	"time"
)

// Tenant represents an organization or isolated workspace
type Tenant struct {
	ID        string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// User represents an identity within a Tenant
type User struct {
	ID           string
	TenantID     string
	Email        string
	PasswordHash string
	Role         string
	Status       string

	Is2FAEnabled bool
	TOTPSecret   string

	CreatedAt time.Time
	UpdatedAt time.Time
}

// Session represents an active Refresh Token lifecycle.
type Session struct {
	ID           string
	UserID       string
	RefreshToken string
	DeviceID     string
	ClientIP     string
	UserAgent    string
	IsRevoked    bool
	ExpiresAt    time.Time
	CreatedAt    time.Time
}

// ResourcePolicy replaces GroupPolicy.
type ResourcePolicy struct {
	ID         string
	UserID     string
	ResourceID string
	Action     string
	CreatedAt  time.Time
}
