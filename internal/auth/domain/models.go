package domain

import "time"

type Tenant struct {
	ID        int64
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type User struct {
	ID           int64
	TenantID     int64
	Login        string
	PasswordHash string
	GlobalRole   string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type HostGroup struct {
	ID        int64
	TenantID  int64
	Name      string
	CreatedAt time.Time
}

type Host struct {
	ID        int64
	TenantID  int64
	GroupID   *int64
	Hostname  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type GroupPolicy struct {
	UserID     int64
	GroupID    int64
	Permission string
	CreatedAt  time.Time
}
