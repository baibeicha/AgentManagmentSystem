package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenRepository interface {
	IsTokenValid(tokenString string) bool
	SaveToWhiteList(tokenString, userID string) error
	RevokeToken(tokenString string) error
}

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type TokenClaims struct {
	Username   string `json:"username"`
	TenantID   string `json:"tid"`
	GlobalRole string `json:"rol"`
	jwt.RegisteredClaims
}

type UserDetails interface {
	GetID() string
	GetUsername() string
	GetTenantID() string
	GetRole() string
}

func (tp *TokenProvider) GenerateAccess(user UserDetails) (string, error) {
	claims := TokenClaims{
		Username:   user.GetUsername(),
		TenantID:   user.GetTenantID(),
		GlobalRole: user.GetRole(),
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    user.GetID(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tp.accessTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(tp.privateKey)
}

func (tp *TokenProvider) GenerateRefresh(user UserDetails) (string, error) {
	claims := TokenClaims{
		Username:   user.GetUsername(),
		TenantID:   user.GetTenantID(),
		GlobalRole: user.GetRole(),
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    user.GetID(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tp.refreshTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(tp.privateKey)
}

func (tp *TokenProvider) GenerateTokens(user UserDetails) (*Tokens, error) {
	accessToken, err := tp.GenerateAccess(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := tp.GenerateRefresh(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	err = tp.repo.SaveToWhiteList(refreshToken, user.GetID())
	if err != nil {
		return nil, fmt.Errorf("failed to save refresh token to whitelist: %w", err)
	}

	return &Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (tp *TokenProvider) VerifyToken(tokenString string) (bool, error) {
	if tokenString == "" {
		return false, nil
	}

	claims := &TokenClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signature algorithm: %v", t.Header["alg"])
		}
		return tp.publicKey, nil
	})

	if err != nil {
		return false, fmt.Errorf("failed to parse token: %w", err)
	}

	if !tp.repo.IsTokenValid(tokenString) {
		return false, errors.New("token is not valid or has been revoked")
	}

	return token.Valid, nil
}

func (tp *TokenProvider) RefreshToken(refreshToken string, user UserDetails) (*Tokens, error) {
	isValid, err := tp.VerifyToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}
	if !isValid {
		return nil, errors.New("refresh token is not valid")
	}

	var claims TokenClaims
	_, err = jwt.ParseWithClaims(refreshToken, &claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signature algorithm: %v", t.Header["alg"])
		}
		return tp.publicKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	err = tp.repo.RevokeToken(refreshToken)

	if err != nil {
		tp.log.Warn("failed to revoke refresh token: ", "err", err)
	}

	return tp.GenerateTokens(user)
}

func (tp *TokenProvider) DeleteToken(refreshToken string) error {
	token, err := jwt.ParseWithClaims(refreshToken, &TokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signature algorithm: %v", t.Header["alg"])
		}
		return tp.publicKey, nil
	})

	if err != nil {
		return fmt.Errorf("failed to parse token for deletion: %w", err)
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return errors.New("invalid token claims")
	}

	userID, err := claims.GetIssuer()
	if err != nil || userID == "" {
		return errors.New("failed to get user ID from token")
	}

	err = tp.repo.RevokeToken(refreshToken)
	if err != nil {
		return fmt.Errorf("error revoking token: %w", err)
	}

	return nil
}
