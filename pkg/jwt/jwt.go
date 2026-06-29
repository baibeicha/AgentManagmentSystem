package jwt

import (
	"AgentManagmentSystem/pkg/config"
	"crypto/ecdsa"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenProvider struct {
	log        *slog.Logger
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
	repo       TokenRepository
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func (tp *TokenProvider) GetClaims(tokenString string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("unexpected signature algorithm: %v", t.Header["alg"])
		}
		return tp.publicKey, nil
	})
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return token.Claims.(*TokenClaims), nil
}

func NewTokenProvider(cfg *config.Config, repo TokenRepository, accessTTL, refreshTTL time.Duration) (*TokenProvider, error) {
	privBytes, err := os.ReadFile(cfg.JWT.PrivateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key: %w", err)
	}
	privateKey, err := jwt.ParseECPrivateKeyFromPEM(privBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	pubBytes, err := os.ReadFile(cfg.JWT.PublicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read public key: %w", err)
	}
	publicKey, err := jwt.ParseECPublicKeyFromPEM(pubBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to read public key: %w", err)
	}

	return &TokenProvider{
		privateKey: privateKey,
		publicKey:  publicKey,
		repo:       repo,
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
		log:        slog.Default(),
	}, nil
}
