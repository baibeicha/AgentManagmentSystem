package jwt

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenRepository interface {
	SaveToWhiteList(ctx context.Context, tokenString, tenantId, userId, deviceId, ip, userAgent string) error
	VerifyAndRotate(ctx context.Context, oldToken, newToken, tenantId, userId, deviceId, ip, userAgent string) error
	RevokeSession(ctx context.Context, tenantId, userId, deviceId string) error
}

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type TokenClaims struct {
	Username   string `json:"username"`
	TenantID   string `json:"tid"`
	GlobalRole string `json:"rol"`
	DeviceID   string `json:"did"`
	jwt.RegisteredClaims
}

type UserDetails interface {
	GetID() string
	GetUsername() string
	GetTenantID() string
	GetRole() string
}

func (tp *TokenProvider) GenerateAccess(user UserDetails, deviceID string) (string, error) {
	claims := TokenClaims{
		Username:   user.GetUsername(),
		TenantID:   user.GetTenantID(),
		GlobalRole: user.GetRole(),
		DeviceID:   deviceID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    user.GetID(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tp.accessTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	return token.SignedString(tp.privateKey)
}

func (tp *TokenProvider) GenerateRefresh(user UserDetails, deviceID string) (string, error) {
	claims := TokenClaims{
		Username:   user.GetUsername(),
		TenantID:   user.GetTenantID(),
		GlobalRole: user.GetRole(),
		DeviceID:   deviceID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    user.GetID(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tp.refreshTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	return token.SignedString(tp.privateKey)
}

func (tp *TokenProvider) GenerateTokens(ctx context.Context, user UserDetails, deviceID, ip, userAgent string) (*Tokens, error) {
	accessToken, err := tp.GenerateAccess(user, deviceID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := tp.GenerateRefresh(user, deviceID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	err = tp.repo.SaveToWhiteList(ctx, refreshToken, user.GetTenantID(), user.GetID(), deviceID, ip, userAgent)
	if err != nil {
		return nil, fmt.Errorf("failed to save session to whitelist: %w", err)
	}

	return &Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (tp *TokenProvider) ParseAndVerify(tokenString string) (*TokenClaims, error) {
	if tokenString == "" {
		return nil, errors.New("empty token string")
	}

	claims := &TokenClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signature algorithm: %v", t.Header["alg"])
		}
		return tp.publicKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, errors.New("token is not valid")
	}

	return claims, nil
}

func (tp *TokenProvider) RefreshTokens(ctx context.Context, oldRefreshToken string, user UserDetails, ip, userAgent string) (*Tokens, error) {
	claims, err := tp.ParseAndVerify(oldRefreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid old refresh token: %w", err)
	}

	deviceID := claims.DeviceID
	if deviceID == "" {
		return nil, errors.New("token claims do not contain device ID")
	}

	newAccessToken, err := tp.GenerateAccess(user, deviceID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate new access token: %w", err)
	}

	newRefreshToken, err := tp.GenerateRefresh(user, deviceID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate new refresh token: %w", err)
	}

	err = tp.repo.VerifyAndRotate(ctx, oldRefreshToken, newRefreshToken, user.GetTenantID(), user.GetID(), deviceID, ip, userAgent)
	if err != nil {
		tp.log.Warn("token rotation failed (possible theft attempt)",
			"user_id", user.GetID(),
			"device_id", deviceID,
			"error", err,
		)
		return nil, fmt.Errorf("session rotation failed: %w", err)
	}

	return &Tokens{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (tp *TokenProvider) DeleteToken(ctx context.Context, refreshToken string) error {
	claims := &TokenClaims{}
	_, err := jwt.ParseWithClaims(refreshToken, claims, func(t *jwt.Token) (interface{}, error) {
		return tp.publicKey, nil
	})

	if err != nil && !errors.Is(err, jwt.ErrTokenExpired) {
		return fmt.Errorf("failed to parse token for deletion: %w", err)
	}

	userID, _ := claims.GetIssuer()
	if userID == "" || claims.TenantID == "" || claims.DeviceID == "" {
		return errors.New("token misses required session identifiers (uid, tid, did)")
	}

	err = tp.repo.RevokeSession(ctx, claims.TenantID, userID, claims.DeviceID)
	if err != nil {
		return fmt.Errorf("error revoking token session: %w", err)
	}

	return nil
}
