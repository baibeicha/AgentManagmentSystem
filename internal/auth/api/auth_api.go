package api

import (
	"context"
	"errors"

	"AgentManagmentSystem/internal/auth/service"
	server "AgentManagmentSystem/pkg/api/grpc/auth/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func mapError(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, service.ErrInvalidCredentials):
		return status.Error(codes.Unauthenticated, err.Error())
	case errors.Is(err, service.ErrUserNotFound):
		return status.Error(codes.Unauthenticated, err.Error())
	case errors.Is(err, service.ErrInvalidToken):
		return status.Error(codes.Unauthenticated, err.Error())
	case errors.Is(err, service.ErrTokenReused):
		return status.Error(codes.Unauthenticated, err.Error())
	case errors.Is(err, service.ErrResourceNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, service.ErrInvalidArgument):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, service.ErrInternalError):
		return status.Error(codes.Internal, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}

func (a *AuthServiceApi) Login(ctx context.Context, request *server.LoginRequest) (*server.LoginResponse, error) {
	resp, err := a.authService.Login(ctx, request)
	return resp, mapError(err)
}

func (a *AuthServiceApi) RefreshToken(ctx context.Context, request *server.RefreshTokenRequest) (*server.RefreshTokenResponse, error) {
	resp, err := a.authService.RefreshToken(ctx, request)
	return resp, mapError(err)
}

func (a *AuthServiceApi) Logout(ctx context.Context, request *server.LogoutRequest) (*server.LogoutResponse, error) {
	resp, err := a.authService.Logout(ctx, request)
	return resp, mapError(err)
}

func (a *AuthServiceApi) CheckPermission(ctx context.Context, request *server.CheckPermissionRequest) (*server.CheckPermissionResponse, error) {
	resp, err := a.authService.CheckPermission(ctx, request)
	return resp, mapError(err)
}
