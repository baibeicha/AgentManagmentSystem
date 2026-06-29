package service

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"AgentManagmentSystem/internal/auth/domain"
	"AgentManagmentSystem/internal/auth/repository"
	server "AgentManagmentSystem/pkg/api/grpc/auth/v1"
	"AgentManagmentSystem/pkg/config"
	"AgentManagmentSystem/pkg/encoder"
	"AgentManagmentSystem/pkg/jwt"

	"github.com/pquerna/otp/totp"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

var (
	StatusSuspended = "SUSPENDED"
)

var (
	RoleUser      = "USER"
	RoleTeamAdmin = "TEAM_ADMIN"
	RoleAdmin     = "ADMIN"
)

type AuthService struct {
	log           *slog.Logger
	cfg           *config.Config
	tokenProvider *jwt.TokenProvider
	userRepo      *repository.UserRepository
	policyRepo    *repository.ResourcePolicyRepository
	sessionRepo   *repository.SessionRepository
}

func NewAuthService(
	log *slog.Logger,
	cfg *config.Config,
	tokenProvider *jwt.TokenProvider,
	userRepo *repository.UserRepository,
	policyRepo *repository.ResourcePolicyRepository,
	sessionRepo *repository.SessionRepository,
) *AuthService {
	return &AuthService{
		log:           log,
		cfg:           cfg,
		tokenProvider: tokenProvider,
		userRepo:      userRepo,
		policyRepo:    policyRepo,
		sessionRepo:   sessionRepo,
	}
}

type userDetailsAdapter struct {
	user *domain.User
}

func (a userDetailsAdapter) GetID() string       { return a.user.ID }
func (a userDetailsAdapter) GetUsername() string { return a.user.Email }
func (a userDetailsAdapter) GetTenantID() string { return a.user.TenantID }
func (a userDetailsAdapter) GetRole() string     { return a.user.Role }

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
	user, err := a.userRepo.GetByEmail(ctx, request.Login)
	if err != nil {
		a.log.Warn("failed login attempt: user not found", "email", request.Login)
		return nil, ErrInvalidCredentials
	}

	if user.Status == StatusSuspended {
		return nil, ErrUserSuspended
	}

	if !encoder.CheckPassword(request.Password, user.PasswordHash) {
		a.log.Warn("failed login attempt: wrong password", "email", request.Login)
		return nil, ErrInvalidCredentials
	}

	clientIP, userAgent := extractClientMeta(ctx)

	adapter := userDetailsAdapter{user: user}
	tokens, err := a.tokenProvider.GenerateTokens(ctx, adapter, request.DeviceId, clientIP, userAgent)
	if err != nil {
		a.log.Error("failed to generate tokens", "error", err)
		return nil, ErrInternalError
	}

	refreshTtl := a.cfg.GetDuration("jwt.ttl.refresh")
	unit := config.GetTimeUnit(a.cfg.GetString("jwt.ttl.unit"))

	session := &domain.Session{
		UserID:       user.ID,
		RefreshToken: tokens.RefreshToken,
		DeviceID:     request.DeviceId,
		ClientIP:     clientIP,
		UserAgent:    userAgent,
		IsRevoked:    false,
		ExpiresAt:    time.Now().Add(refreshTtl * unit),
	}

	if err := a.sessionRepo.Create(ctx, session); err != nil {
		a.log.Error("failed to save session to db", "error", err)
		return nil, ErrInternalError
	}

	a.log.Info("user authenticated successfully", "user_id", user.ID, "ip", clientIP)

	return &server.LoginResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

func (a *AuthService) ValidateToken(ctx context.Context, request *server.ValidateTokenRequest) (*server.ValidateTokenResponse, error) {
	claims, err := a.tokenProvider.ParseAndVerify(request.AccessToken)
	if err != nil {
		return nil, ErrInvalidToken
	}

	userID, err := claims.GetIssuer()
	if err != nil {
		return nil, ErrInvalidToken
	}

	user, err := a.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if user.Status == StatusSuspended {
		return nil, ErrUserSuspended
	}

	return &server.ValidateTokenResponse{
		UserId:   user.ID,
		Email:    user.Email,
		TenantId: user.TenantID,
		Role:     user.Role,
		Status:   user.Status,
	}, nil
}

func (a *AuthService) RefreshToken(ctx context.Context, request *server.RefreshTokenRequest) (*server.RefreshTokenResponse, error) {
	session, err := a.sessionRepo.GetByRefreshToken(ctx, request.RefreshToken)
	if err != nil || session.IsRevoked || session.ExpiresAt.Before(time.Now()) {
		a.log.Warn("attempt to use invalid or revoked refresh token")
		return nil, ErrInvalidToken
	}

	user, err := a.userRepo.GetByID(ctx, session.UserID)
	if err != nil || user.Status == StatusSuspended {
		return nil, ErrUserNotFound
	}

	clientIP, userAgent := extractClientMeta(ctx)
	adapter := userDetailsAdapter{user: user}

	tokens, err := a.tokenProvider.RefreshTokens(ctx, request.RefreshToken, adapter, clientIP, userAgent)
	if err != nil {
		a.log.Warn("token rotation failed in provider", "error", err, "user_id", user.ID)
		_ = a.sessionRepo.RevokeSession(ctx, session.ID)
		return nil, ErrTokenReused
	}

	_ = a.sessionRepo.RevokeSession(ctx, session.ID)

	newSession := &domain.Session{
		UserID:       user.ID,
		RefreshToken: tokens.RefreshToken,
		DeviceID:     session.DeviceID,
		ClientIP:     clientIP,
		UserAgent:    userAgent,
		IsRevoked:    false,
		ExpiresAt:    time.Now().Add(30 * 24 * time.Hour),
	}
	_ = a.sessionRepo.Create(ctx, newSession)

	return &server.RefreshTokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

func (a *AuthService) Logout(ctx context.Context, request *server.LogoutRequest) (*server.LogoutResponse, error) {
	err := a.tokenProvider.DeleteToken(ctx, request.RefreshToken)
	if err != nil {
		a.log.Error("failed to delete token from redis", "error", err)
	}

	session, err := a.sessionRepo.GetByRefreshToken(ctx, request.RefreshToken)
	if err == nil {
		_ = a.sessionRepo.RevokeSession(ctx, session.ID)
	}

	return &server.LogoutResponse{Success: true}, nil
}

func (a *AuthService) GetMe(ctx context.Context, request *server.GetMeRequest) (*server.UserResponse, error) {
	user, err := a.userRepo.GetByID(ctx, request.UserId)
	if err != nil {
		return nil, ErrUserNotFound
	}

	return &server.UserResponse{
		UserId:   user.ID,
		Email:    user.Email,
		TenantId: user.TenantID,
		Role:     user.Role,
		Status:   user.Status,
	}, nil
}

func (a *AuthService) GetPermissions(ctx context.Context, request *server.GetPermissionsRequest) (*server.GetPermissionsResponse, error) {
	policies, err := a.policyRepo.GetPoliciesByUserID(ctx, request.UserId)
	if err != nil {
		a.log.Error("failed to get user policies", "error", err)
		return nil, ErrInternalError
	}

	var perms []string
	for _, p := range policies {
		perms = append(perms, p.ResourceID+":"+p.Action)
	}

	return &server.GetPermissionsResponse{Permissions: perms}, nil
}

func (a *AuthService) CheckPermission(ctx context.Context, request *server.CheckPermissionRequest) (*server.CheckPermissionResponse, error) {
	user, err := a.userRepo.GetByID(ctx, request.UserId)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if user.Role == RoleAdmin {
		return &server.CheckPermissionResponse{Allowed: true}, nil
	}

	if user.Role == RoleTeamAdmin && user.TenantID == request.TenantId {
		return &server.CheckPermissionResponse{Allowed: true}, nil
	}

	hasAccess, err := a.policyRepo.CheckPermission(ctx, user.ID, request.ResourceId, request.Action)
	if err != nil {
		a.log.Error("failed to check granular permissions", "error", err)
		return nil, ErrInternalError
	}

	return &server.CheckPermissionResponse{Allowed: hasAccess}, nil
}

func (a *AuthService) Setup2FA(ctx context.Context, request *server.Setup2FARequest) (*server.Setup2FAResponse, error) {
	user, err := a.userRepo.GetByID(ctx, request.UserId)
	if err != nil {
		return nil, ErrUserNotFound
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      a.cfg.GetString("app.name"),
		AccountName: user.Email,
	})
	if err != nil {
		a.log.Error("failed to generate TOTP secret", "error", err)
		return nil, ErrInternalError
	}

	user.TOTPSecret = key.Secret()
	if err := a.userRepo.Update(ctx, user); err != nil {
		a.log.Error("failed to save TOTP secret to db", "error", err)
		return nil, ErrInternalError
	}

	return &server.Setup2FAResponse{
		Secret:     key.Secret(),
		OtpauthUrl: key.URL(),
	}, nil
}

func (a *AuthService) Verify2FA(ctx context.Context, request *server.Verify2FARequest) (*server.Verify2FAResponse, error) {
	user, err := a.userRepo.GetByID(ctx, request.UserId)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if user.TOTPSecret == "" {
		return nil, errors.New("2FA is not initialized for this user")
	}

	valid := totp.Validate(request.Code, user.TOTPSecret)
	if !valid {
		a.log.Warn("invalid 2FA code provided", "user_id", user.ID)
		return nil, errors.New("invalid TOTP code")
	}

	if !user.Is2FAEnabled {
		user.Is2FAEnabled = true
		if err := a.userRepo.Update(ctx, user); err != nil {
			a.log.Error("failed to enable 2FA in db", "error", err)
			return nil, ErrInternalError
		}
		a.log.Info("2FA successfully enabled for user", "user_id", user.ID)
	}

	// В будущем здесь можно добавить логику проверки request.ActionId для подтверждения опасных операций (например, удаление БД).

	return &server.Verify2FAResponse{Status: "approved"}, nil
}

func (a *AuthService) Disable2FA(ctx context.Context, request *server.Disable2FARequest) (*server.Disable2FAResponse, error) {
	user, err := a.userRepo.GetByID(ctx, request.UserId)
	if err != nil {
		return nil, ErrUserNotFound
	}

	valid := totp.Validate(request.Code, user.TOTPSecret)
	if !valid {
		a.log.Warn("attempt to disable 2FA with invalid code", "user_id", user.ID)
		return nil, errors.New("invalid TOTP code")
	}

	user.Is2FAEnabled = false
	user.TOTPSecret = ""

	if err := a.userRepo.Update(ctx, user); err != nil {
		a.log.Error("failed to disable 2FA in db", "error", err)
		return nil, ErrInternalError
	}

	a.log.Info("2FA disabled for user", "user_id", user.ID)

	return &server.Disable2FAResponse{Success: true}, nil
}

func (a *AuthService) Register(ctx context.Context, request *server.RegisterRequest) (*server.RegisterResponse, error) {
	passwordHash, err := encoder.HashPassword(request.Password)
	if err != nil {
		a.log.Error("failed to hash password", "error", err)
		return nil, ErrInternalError
	}

	err = a.userRepo.Create(ctx, &domain.User{
		Email:        request.Email,
		PasswordHash: passwordHash,
		Role:         RoleUser,
	})

	if err != nil {
		a.log.Error("failed to create user", "error", err)
		return nil, ErrInvalidArgument
	}

	return &server.RegisterResponse{
		Success: true,
	}, nil
}
