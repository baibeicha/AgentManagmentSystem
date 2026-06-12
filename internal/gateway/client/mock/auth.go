package mock

import (
	"AgentManagmentSystem/internal/gateway/domain"
	"context"
	"errors"
)

type AuthMock struct{}

func NewAuthMock() *AuthMock {
	return &AuthMock{}
}

func (m *AuthMock) Register(ctx context.Context, email, password string) (bool, error) {
	if email == "" || password == "" {
		return false, errors.New("invalid credentials")
	}
	return true, nil
}

func (m *AuthMock) Login(ctx context.Context, login, password, deviceID string) (string, string, error) {
	if login == "" || password == "" {
		return "", "", errors.New("invalid credentials")
	}
	return "mock_access_jwt_token_12345", "mock_refresh_opaque_token_67890", nil
}

func (m *AuthMock) ValidateToken(ctx context.Context, accessToken string) (*domain.User, error) {
	if accessToken != "mock_access_jwt_token_12345" {
		return nil, errors.New("unauthorized")
	}

	return &domain.User{
		UserID:   "u-1111-2222-3333-4444",
		Email:    "admin@ams.local",
		TenantID: "tenant-001",
		Role:     "TEAM_ADMIN",
		Status:   "ACTIVE",
	}, nil
}

func (m *AuthMock) RefreshToken(ctx context.Context, refreshToken string) (string, string, error) {
	if refreshToken != "mock_refresh_opaque_token_67890" {
		return "", "", errors.New("invalid refresh token")
	}
	return "new_mock_access_jwt_token", "new_mock_refresh_opaque_token", nil
}

func (m *AuthMock) Logout(ctx context.Context, refreshToken string) error {
	return nil
}

func (m *AuthMock) GetMe(ctx context.Context) (*domain.User, error) {
	return &domain.User{
		UserID: "u-1111-2222-3333-4444",
		Email:  "admin@omniwatch.local",
		Role:   "TEAM_ADMIN",
		Status: "ACTIVE",
	}, nil
}

func (m *AuthMock) GetPermissions(ctx context.Context) ([]string, error) {
	return []string{"Metrics:View", "Terminal:Open", "Command:Execute", "Device:Delete"}, nil
}

func (m *AuthMock) Setup2FA(ctx context.Context) (string, string, error) {
	return "JBSWY3DPEHPK3PXP", "otpauth://totp/OmniWatch:admin?secret=JBSWY3DPEHPK3PXP&issuer=OmniWatch", nil
}

func (m *AuthMock) Verify2FA(ctx context.Context, code, actionID string) (string, error) {
	if code == "123456" {
		return "approved", nil
	}
	return "", errors.New("invalid TOTP code")
}

func (m *AuthMock) Disable2FA(ctx context.Context, code string) error {
	if code != "123456" {
		return errors.New("invalid TOTP code")
	}
	return nil
}
