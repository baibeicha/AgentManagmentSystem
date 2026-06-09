package service

import (
	"context"
	"errors"
	"log/slog"
	"strconv"

	"AgentManagmentSystem/internal/auth/domain"
	"AgentManagmentSystem/internal/auth/repository"
	server "AgentManagmentSystem/pkg/api/grpc/auth/v1"
	"AgentManagmentSystem/pkg/config"
	"AgentManagmentSystem/pkg/encoder"
	"AgentManagmentSystem/pkg/jwt"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
	ErrResourceNotFound   = errors.New("resource not found")
	ErrInternalError      = errors.New("internal server error")
	ErrInvalidToken       = errors.New("invalid refresh token")
	ErrTokenReused        = errors.New("session rotation failed or token reused")
	ErrInvalidArgument    = errors.New("invalid argument format")
)

type AuthService struct {
	log           *slog.Logger
	cfg           *config.Config
	tokenProvider *jwt.TokenProvider
	userRepo      repository.UserRepository
	hostRepo      repository.HostRepository
	policyRepo    repository.GroupPolicyRepository
}

func NewAuthService(
	log *slog.Logger,
	cfg *config.Config,
	tokenProvider *jwt.TokenProvider,
	userRepo repository.UserRepository,
	hostRepo repository.HostRepository,
	policyRepo repository.GroupPolicyRepository,
) *AuthService {
	return &AuthService{
		log:           log,
		cfg:           cfg,
		tokenProvider: tokenProvider,
		userRepo:      userRepo,
		hostRepo:      hostRepo,
		policyRepo:    policyRepo,
	}
}

type userDetailsAdapter struct {
	user *domain.User
}

func (a userDetailsAdapter) GetID() string       { return strconv.FormatInt(a.user.ID, 10) }
func (a userDetailsAdapter) GetUsername() string { return a.user.Login }
func (a userDetailsAdapter) GetTenantID() string { return strconv.FormatInt(a.user.TenantID, 10) }
func (a userDetailsAdapter) GetRole() string     { return a.user.GlobalRole }

func extractClientMeta(ctx context.Context) (ip string, userAgent string) {
	ip = "unknown"
	if p, ok := peer.FromContext(ctx); ok {
		ip = p.Addr.String()
	}

	userAgent = "unknown"
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if ua := md.Get("user-agent"); len(ua) > 0 {
			userAgent = ua[0]
		}
	}
	return ip, userAgent
}

func (a *AuthService) Login(ctx context.Context, request *server.LoginRequest) (*server.LoginResponse, error) {
	user, err := a.userRepo.GetByLogin(ctx, request.Login)
	if err != nil {
		a.log.Warn("failed login attempt: user not found", "login", request.Login)
		return nil, ErrInvalidCredentials
	}

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if !encoder.CheckPassword(request.Password, user.PasswordHash) {
		a.log.Warn("failed login attempt: wrong password", "login", request.Login)
		return nil, ErrInvalidCredentials
	}

	clientIP, userAgent := extractClientMeta(ctx)

	a.log.Info("user successfully authenticated",
		"user_id", user.ID,
		"tenant_id", user.TenantID,
		"ip", clientIP,
	)

	adapter := userDetailsAdapter{user: user}
	tokens, err := a.tokenProvider.GenerateTokens(ctx, adapter, request.DeviceId, clientIP, userAgent)
	if err != nil {
		a.log.Error("failed to generate tokens", "error", err)
		return nil, ErrInternalError
	}

	return &server.LoginResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

func (a *AuthService) RefreshToken(ctx context.Context, request *server.RefreshTokenRequest) (*server.RefreshTokenResponse, error) {
	claims, err := a.tokenProvider.ParseAndVerify(request.RefreshToken)
	if err != nil {
		return nil, ErrInvalidToken
	}

	issuer, err := claims.GetIssuer()
	if err != nil {
		return nil, ErrInvalidToken
	}

	userID, _ := strconv.ParseInt(issuer, 10, 64)
	user, err := a.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	clientIP, userAgent := extractClientMeta(ctx)

	adapter := userDetailsAdapter{user: user}
	tokens, err := a.tokenProvider.RefreshTokens(ctx, request.RefreshToken, adapter, clientIP, userAgent)
	if err != nil {
		a.log.Warn("token rotation failed", "error", err, "user_id", userID)
		return nil, ErrTokenReused
	}

	return &server.RefreshTokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

func (a *AuthService) Logout(ctx context.Context, request *server.LogoutRequest) (*server.LogoutResponse, error) {
	err := a.tokenProvider.DeleteToken(ctx, request.RefreshToken)
	if err != nil {
		a.log.Error("failed to delete token during logout", "error", err)
		return nil, ErrInternalError
	}

	return &server.LogoutResponse{Success: true}, nil
}

func (a *AuthService) CheckPermission(ctx context.Context, request *server.CheckPermissionRequest) (*server.CheckPermissionResponse, error) {
	userID, err := strconv.ParseInt(request.UserId, 10, 64)
	if err != nil {
		return nil, ErrInvalidArgument
	}

	tenantID, err := strconv.ParseInt(request.TenantId, 10, 64)
	if err != nil {
		return nil, ErrInvalidArgument
	}

	resourceID, err := strconv.ParseInt(request.ResourceId, 10, 64)
	if err != nil {
		return nil, ErrInvalidArgument
	}

	user, err := a.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if user.TenantID != tenantID {
		a.log.Warn("cross-tenant access attempt detected!", "user_id", userID, "attempted_tenant", tenantID)
		return &server.CheckPermissionResponse{Allowed: false}, nil
	}

	if user.GlobalRole == "global_admin" {
		return &server.CheckPermissionResponse{Allowed: true}, nil
	}

	if user.GlobalRole == "tenant_admin" {
		host, err := a.hostRepo.GetByID(ctx, resourceID)
		if err != nil {
			return nil, ErrResourceNotFound
		}
		if host.TenantID == user.TenantID {
			return &server.CheckPermissionResponse{Allowed: true}, nil
		}
		return &server.CheckPermissionResponse{Allowed: false}, nil
	}

	host, err := a.hostRepo.GetByID(ctx, resourceID)
	if err != nil || host.GroupID == nil {
		return &server.CheckPermissionResponse{Allowed: false}, nil
	}

	hasAccess, err := a.policyRepo.CheckPermission(ctx, user.ID, *host.GroupID, request.Action)
	if err != nil {
		a.log.Error("failed to check granular permissions", "error", err)
		return nil, ErrInternalError
	}

	return &server.CheckPermissionResponse{Allowed: hasAccess}, nil
}
