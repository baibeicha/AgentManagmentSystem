package jwt

import (
	"AgentManagmentSystem/pkg/config"
	"crypto/rsa"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenProvider struct {
	log        *slog.Logger
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	repo       TokenRepository
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func (tp *TokenProvider) GetClaims(tokenString string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		return tp.publicKey, nil
	})
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return token.Claims.(*TokenClaims), nil
}

func NewTokenProvider(cfg *config.Config, log *slog.Logger, repo TokenRepository, accessTTL, refreshTTL time.Duration) (*TokenProvider, error) {
	privBytes, err := os.ReadFile(cfg.GetString("jwt.private_key_path"))
	if err != nil {
		return nil, fmt.Errorf("failed to read private key: %w", err)
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	pubBytes, err := os.ReadFile(cfg.GetString("jwt.public_key_path"))
	if err != nil {
		return nil, fmt.Errorf("failed to read public key: %w", err)
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(pubBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to read public key: %w", err)
	}

	return &TokenProvider{
		privateKey: privateKey,
		publicKey:  publicKey,
		repo:       repo,
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
		log:        log,
	}, nil
}

func GetTimeUnit(unit string) time.Duration {
	switch unit {
	case "s":
		return time.Second
	case "m":
		return time.Minute
	case "h":
		return time.Hour
	case "d":
		return time.Hour * 24
	default:
		return time.Minute
	}
}
