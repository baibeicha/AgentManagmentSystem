package api

import (
	"context"

	"AgentManagmentSystem/internal/auth/service"
	server "AgentManagmentSystem/pkg/api/grpc/auth/v1"
)

type AuthServiceApi struct {
	server.UnimplementedAuthServiceServer
	authService *service.AuthService
}

func NewAuthServiceApi(authService *service.AuthService) *AuthServiceApi {
	return &AuthServiceApi{
		authService: authService,
	}
}

func (a *AuthServiceApi) Login(ctx context.Context, request *server.LoginRequest) (*server.LoginResponse, error) {
	return a.authService.Login(ctx, request)
}

func (a *AuthServiceApi) RefreshToken(ctx context.Context, request *server.RefreshTokenRequest) (*server.RefreshTokenResponse, error) {
	return a.authService.RefreshToken(ctx, request)
}

func (a *AuthServiceApi) Logout(ctx context.Context, request *server.LogoutRequest) (*server.LogoutResponse, error) {
	return a.authService.Logout(ctx, request)
}

func (a *AuthServiceApi) CheckPermission(ctx context.Context, request *server.CheckPermissionRequest) (*server.CheckPermissionResponse, error) {
	return a.authService.CheckPermission(ctx, request)
}
