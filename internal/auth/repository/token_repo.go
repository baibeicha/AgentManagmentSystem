package repository

import (
	"AgentManagmentSystem/pkg/jwt"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"AgentManagmentSystem/pkg/config"

	"github.com/redis/go-redis/v9"
)

type RedisTokenRepository struct {
	client *redis.Client
	cfg    *config.Config
	ttl    time.Duration
}

type SessionMetadata struct {
	TokenHash string    `json:"token_hash"`
	IP        string    `json:"ip"`
	UserAgent string    `json:"user_agent"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

func NewRedisTokenRepo(cfg *config.Config, client *redis.Client) *RedisTokenRepository {
	ttlUnit := jwt.GetTimeUnit(cfg.GetString("jwt.ttl.unit"))
	refreshTTL := cfg.GetDuration("jwt.ttl.refresh") * ttlUnit

	return &RedisTokenRepository{
		cfg:    cfg,
		client: client,
		ttl:    refreshTTL,
	}
}

func (r *RedisTokenRepository) hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func (r *RedisTokenRepository) makeKey(tenantId, userId, deviceId string) string {
	return fmt.Sprintf("auth:tenant:%s:user:%s:session:%s", tenantId, userId, deviceId)
}

func (r *RedisTokenRepository) SaveToWhiteList(ctx context.Context,
	tokenString, tenantId, userId, deviceId, ip, userAgent string) error {
	now := time.Now()

	session := SessionMetadata{
		TokenHash: r.hashToken(tokenString),
		IP:        ip,
		UserAgent: userAgent,
		CreatedAt: now,
		ExpiresAt: now.Add(r.ttl),
	}

	data, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("failed to marshal session metadata: %w", err)
	}

	key := r.makeKey(tenantId, userId, deviceId)

	return r.client.Set(ctx, key, data, r.ttl).Err()
}

func (r *RedisTokenRepository) VerifyAndRotate(ctx context.Context,
	oldToken, newToken, tenantId, userId, deviceId, ip, userAgent string) error {
	key := r.makeKey(tenantId, userId, deviceId)

	val, err := r.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return fmt.Errorf("session not found or expired")
	} else if err != nil {
		return err
	}

	var session SessionMetadata
	if err := json.Unmarshal([]byte(val), &session); err != nil {
		return err
	}

	if session.TokenHash != r.hashToken(oldToken) {
		err := r.RevokeAllUserSessions(ctx, tenantId, userId)
		if err != nil {
			return fmt.Errorf("SECURITY ALERT: token reuse detected, error while revoking sessions: %w", err)
		}
		return fmt.Errorf("SECURITY ALERT: token reuse detected, all sessions revoked")
	}

	return r.SaveToWhiteList(ctx, newToken, tenantId, userId, deviceId, ip, userAgent)
}

func (r *RedisTokenRepository) RevokeSession(ctx context.Context, tenantId, userId, deviceId string) error {
	key := r.makeKey(tenantId, userId, deviceId)
	return r.client.Del(ctx, key).Err()
}

func (r *RedisTokenRepository) GetActiveSessions(ctx context.Context, tenantId, userId string) (map[string]SessionMetadata, error) {
	pattern := fmt.Sprintf("auth:tenant:%s:user:%s:session:*", tenantId, userId)
	sessions := make(map[string]SessionMetadata)

	iter := r.client.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		key := iter.Val()
		val, err := r.client.Get(ctx, key).Result()
		if err == nil {
			var meta SessionMetadata
			if json.Unmarshal([]byte(val), &meta) == nil {
				parts := strings.Split(key, ":")
				if len(parts) >= 7 {
					deviceId := parts[6]
					sessions[deviceId] = meta
				}
			}
		}
	}

	if err := iter.Err(); err != nil {
		return nil, err
	}
	return sessions, nil
}

func (r *RedisTokenRepository) RevokeAllUserSessions(ctx context.Context, tenantId, userId string) error {
	pattern := fmt.Sprintf("auth:tenant:%s:user:%s:session:*", tenantId, userId)
	return r.deleteByPattern(ctx, pattern)
}

func (r *RedisTokenRepository) RevokeTenantSessions(ctx context.Context, tenantId string) error {
	pattern := fmt.Sprintf("auth:tenant:%s:user:*:session:*", tenantId)
	return r.deleteByPattern(ctx, pattern)
}

func (r *RedisTokenRepository) deleteByPattern(ctx context.Context, pattern string) error {
	iter := r.client.Scan(ctx, 0, pattern, 0).Iterator()
	var keys []string

	for iter.Next(ctx) {
		keys = append(keys, iter.Val())

		if len(keys) >= 100 {
			if err := r.client.Del(ctx, keys...).Err(); err != nil {
				return err
			}
			keys = keys[:0]
		}
	}

	if err := iter.Err(); err != nil {
		return err
	}

	if len(keys) > 0 {
		if err := r.client.Del(ctx, keys...).Err(); err != nil {
			return err
		}
	}
	return nil
}
