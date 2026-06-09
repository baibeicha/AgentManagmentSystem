package api

import (
	"AgentManagmentSystem/internal/auth/service"
	server "AgentManagmentSystem/pkg/api/grpc/auth/v1"
	"context"
)

type AuthServiceApi struct {
	server.UnimplementedAuthServiceServer
	authService service.AuthService
}

func (a AuthServiceApi) Login(ctx context.Context, request *server.LoginRequest) (*server.LoginResponse, error) {
	return a.authService.Login(ctx, request)
}

func (a AuthServiceApi) CheckPermission(ctx context.Context, request *server.CheckPermissionRequest) (*server.CheckPermissionResponse, error) {
	return a.authService.CheckPermission(ctx, request)
}
