package service

import (
	server "AgentManagmentSystem/pkg/api/grpc/auth/v1"
	"AgentManagmentSystem/pkg/config"
	"AgentManagmentSystem/pkg/jwt"
	"context"
	"log/slog"

	"github.com/redis/go-redis/v9"
)

type AuthService struct {
	log           *slog.Logger
	cfg           *config.Config
	tokenProvider *jwt.TokenProvider
	redis         *redis.Client
}

func (a AuthService) Login(ctx context.Context, request *server.LoginRequest) (*server.LoginResponse, error) {
	panic("implement me")
}

func (a AuthService) CheckPermission(ctx context.Context, request *server.CheckPermissionRequest) (*server.CheckPermissionResponse, error) {
	panic("implement me")
}
