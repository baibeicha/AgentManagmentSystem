package mock

import (
	"AgentManagmentSystem/internal/gateway/domain"
	"context"
)

type TeamMock struct{}

func NewTeamMock() *TeamMock {
	return &TeamMock{}
}

func (m *TeamMock) ListUsers(ctx context.Context) ([]domain.User, error) {
	return []domain.User{
		{
			UserID: "u-1111-2222-3333-4444",
			Email:  "admin@omniwatch.local",
			Role:   "TEAM_ADMIN",
			Status: "ACTIVE",
		},
		{
			UserID: "u-5555-6666-7777-8888",
			Email:  "operator@omniwatch.local",
			Role:   "OPERATOR",
			Status: "ACTIVE",
		},
		{
			UserID: "u-9999-0000-1111-2222",
			Email:  "new_dev@omniwatch.local",
			Role:   "VIEWER",
			Status: "PENDING_INVITE",
		},
	}, nil
}

func (m *TeamMock) InviteUser(ctx context.Context, email, role string) error {
	return nil
}

func (m *TeamMock) UpdateUserRoles(ctx context.Context, userID string, role string, allowedGroupIDs []string) error {
	return nil
}

func (m *TeamMock) DeleteUser(ctx context.Context, userID string) error {
	return nil
}
