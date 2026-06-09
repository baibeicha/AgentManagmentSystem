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

func (a AuthServiceApi) Login(ctx context.Context, request *server.LoginRequest) (*server.LoginResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a AuthServiceApi) RefreshToken(ctx context.Context, request *server.RefreshTokenRequest) (*server.RefreshTokenResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a AuthServiceApi) Logout(ctx context.Context, request *server.LogoutRequest) (*server.LogoutResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a AuthServiceApi) CheckPermission(ctx context.Context, request *server.CheckPermissionRequest) (*server.CheckPermissionResponse, error) {
	//TODO implement me
	panic("implement me")
}
