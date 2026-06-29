package grpc

import (
	"AgentManagmentSystem/internal/gateway/delivery/http/middleware"
	"AgentManagmentSystem/internal/gateway/domain"
	auth "AgentManagmentSystem/pkg/api/grpc/auth/v1"
	"context"
	"errors"
	"fmt"
	"log/slog"

	"google.golang.org/grpc"
)

var (
	ErrLogoutUnsuccessful     = errors.New("unsuccessful logout")
	ErrNoAuthorization        = errors.New("no authorization")
	ErrDisable2FAUnsuccessful = errors.New("disable 2FA unsuccessful")
)

type AuthServiceClient struct {
	conn       *grpc.ClientConn
	authClient auth.AuthServiceClient
	log        *slog.Logger
}

func NewAuthServiceClient(conn *grpc.ClientConn) *AuthServiceClient {
	authClient := auth.NewAuthServiceClient(conn)
	return &AuthServiceClient{
		conn:       conn,
		authClient: authClient,
		log:        slog.Default(),
	}
}

func (a AuthServiceClient) Register(ctx context.Context, email, password string) (bool, error) {
	resp, err := a.authClient.Register(ctx, &auth.RegisterRequest{
		Email:    email,
		Password: password,
	})

	if err != nil {
		return false, fmt.Errorf("register failed: %w", err)
	}

	return resp.Success, nil
}

func (a AuthServiceClient) Login(ctx context.Context, login, password, deviceID string) (accessToken, refreshToken string, err error) {
	resp, err := a.authClient.Login(ctx, &auth.LoginRequest{
		Login:    login,
		Password: password,
		DeviceId: deviceID,
	})

	if err != nil {
		return "", "", fmt.Errorf("login error: %w", err)
	}

	return resp.AccessToken, resp.RefreshToken, nil
}

func (a AuthServiceClient) ValidateToken(ctx context.Context, accessToken string) (*domain.User, error) {
	resp, err := a.authClient.ValidateToken(ctx, &auth.ValidateTokenRequest{
		AccessToken: accessToken,
	})

	if err != nil {
		return nil, fmt.Errorf("validate token error: %w", err)
	}

	return &domain.User{
		UserID:   resp.UserId,
		Email:    resp.Email,
		Role:     resp.Role,
		Status:   resp.Status,
		TenantID: resp.TenantId,
	}, nil
}

func (a AuthServiceClient) RefreshToken(ctx context.Context, refreshToken string) (newAccess, newRefresh string, err error) {
	resp, err := a.authClient.RefreshToken(ctx, &auth.RefreshTokenRequest{
		RefreshToken: refreshToken,
	})

	if err != nil {
		return "", "", fmt.Errorf("refresh token error: %w", err)
	}

	return resp.AccessToken, resp.RefreshToken, nil
}

func (a AuthServiceClient) Logout(ctx context.Context, refreshToken string) error {
	resp, err := a.authClient.Logout(ctx, &auth.LogoutRequest{
		RefreshToken: refreshToken,
	})

	if err != nil {
		return fmt.Errorf("logout error: %w", err)
	}

	if resp.Success {
		return nil
	}

	return ErrLogoutUnsuccessful
}

func (a AuthServiceClient) GetMe(ctx context.Context) (*domain.User, error) {
	userId, ok := middleware.GetUserID(ctx)
	if !ok {
		return nil, ErrNoAuthorization
	}

	resp, err := a.authClient.GetMe(ctx, &auth.GetMeRequest{
		UserId: userId,
	})

	if err != nil {
		return nil, fmt.Errorf("get me error: %w", err)
	}

	return &domain.User{
		UserID:   resp.UserId,
		Email:    resp.Email,
		Role:     resp.Role,
		Status:   resp.Status,
		TenantID: resp.TenantId,
	}, nil
}

func (a AuthServiceClient) GetPermissions(ctx context.Context) ([]string, error) {
	userId, ok := middleware.GetUserID(ctx)
	if !ok {
		return nil, ErrNoAuthorization
	}

	tenantId, ok := middleware.GetTenantID(ctx)
	if !ok {
		return nil, ErrNoAuthorization
	}

	resp, err := a.authClient.GetPermissions(ctx, &auth.GetPermissionsRequest{
		UserId:   userId,
		TenantId: tenantId,
	})

	if err != nil {
		return nil, fmt.Errorf("get permissions error: %w", err)
	}

	return resp.Permissions, nil
}

func (a AuthServiceClient) Setup2FA(ctx context.Context) (secret, otpauthURL string, err error) {
	userId, ok := middleware.GetUserID(ctx)
	if !ok {
		return "", "", ErrNoAuthorization
	}

	resp, err := a.authClient.Setup2FA(ctx, &auth.Setup2FARequest{
		UserId: userId,
	})

	if err != nil {
		return "", "", fmt.Errorf("setup 2fa error: %w", err)
	}

	return resp.Secret, resp.OtpauthUrl, nil
}

func (a AuthServiceClient) Verify2FA(ctx context.Context, code, actionID string) (status string, err error) {
	userId, ok := middleware.GetUserID(ctx)
	if !ok {
		return "", ErrNoAuthorization
	}

	resp, err := a.authClient.Verify2FA(ctx, &auth.Verify2FARequest{
		UserId:   userId,
		Code:     code,
		ActionId: actionID,
	})

	if err != nil {
		return "", fmt.Errorf("verify 2fa error: %w", err)
	}

	return resp.Status, nil
}

func (a AuthServiceClient) Disable2FA(ctx context.Context, code string) error {
	userId, ok := middleware.GetUserID(ctx)
	if !ok {
		return ErrNoAuthorization
	}

	resp, err := a.authClient.Disable2FA(ctx, &auth.Disable2FARequest{
		UserId: userId,
		Code:   code,
	})

	if err != nil {
		return fmt.Errorf("disable 2fa error: %w", err)
	}

	if resp.Success {
		return nil
	}

	return ErrDisable2FAUnsuccessful
}
