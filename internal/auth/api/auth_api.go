package api

import (
	"AgentManagmentSystem/internal/auth/service"
	server "AgentManagmentSystem/pkg/api/grpc/auth/v1"
	"context"
)

type AuthServiceApi struct {
	server.UnimplementedAuthServiceServer
	authService *service.AuthService
}

func (a *AuthServiceApi) Login(ctx context.Context, request *server.LoginRequest) (*server.LoginResponse, error) {
	resp, err := a.authService.Login(ctx, request)
	return resp, service.MapError(err)
}

func (a *AuthServiceApi) RefreshToken(ctx context.Context, request *server.RefreshTokenRequest) (*server.RefreshTokenResponse, error) {
	resp, err := a.authService.RefreshToken(ctx, request)
	return resp, service.MapError(err)
}

func (a *AuthServiceApi) Logout(ctx context.Context, request *server.LogoutRequest) (*server.LogoutResponse, error) {
	resp, err := a.authService.Logout(ctx, request)
	return resp, service.MapError(err)
}

func (a *AuthServiceApi) ValidateToken(ctx context.Context, request *server.ValidateTokenRequest) (*server.ValidateTokenResponse, error) {
	resp, err := a.authService.ValidateToken(ctx, request)
	return resp, service.MapError(err)
}

func (a *AuthServiceApi) GetMe(ctx context.Context, request *server.GetMeRequest) (*server.UserResponse, error) {
	resp, err := a.authService.GetMe(ctx, request)
	return resp, service.MapError(err)
}

func (a *AuthServiceApi) CheckPermission(ctx context.Context, request *server.CheckPermissionRequest) (*server.CheckPermissionResponse, error) {
	resp, err := a.authService.CheckPermission(ctx, request)
	return resp, service.MapError(err)
}

func (a *AuthServiceApi) GetPermissions(ctx context.Context, request *server.GetPermissionsRequest) (*server.GetPermissionsResponse, error) {
	resp, err := a.authService.GetPermissions(ctx, request)
	return resp, service.MapError(err)
}

func (a *AuthServiceApi) Setup2FA(ctx context.Context, request *server.Setup2FARequest) (*server.Setup2FAResponse, error) {
	resp, err := a.authService.Setup2FA(ctx, request)
	return resp, service.MapError(err)
}

func (a *AuthServiceApi) Verify2FA(ctx context.Context, request *server.Verify2FARequest) (*server.Verify2FAResponse, error) {
	resp, err := a.authService.Verify2FA(ctx, request)
	return resp, service.MapError(err)
}

func (a *AuthServiceApi) Disable2FA(ctx context.Context, request *server.Disable2FARequest) (*server.Disable2FAResponse, error) {
	resp, err := a.authService.Disable2FA(ctx, request)
	return resp, service.MapError(err)
}

func (a *AuthServiceApi) Register(ctx context.Context, request *server.RegisterRequest) (*server.RegisterResponse, error) {
	resp, err := a.authService.Register(ctx, request)
	return resp, service.MapError(err)
}

func NewAuthServiceApi(authService *service.AuthService) *AuthServiceApi {
	return &AuthServiceApi{
		authService: authService,
	}
}
