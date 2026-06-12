package repository

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
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
	ttlUnit := config.GetTimeUnit(cfg.GetString("jwt.ttl.unit"))
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

func (r *RedisTokenRepository) userSessionsKey(tenantId, userId string) string {
	return fmt.Sprintf("auth:tenant:%s:user:%s:sessions", tenantId, userId)
}

func (r *RedisTokenRepository) tenantUsersKey(tenantId string) string {
	return fmt.Sprintf("auth:tenant:%s:users", tenantId)
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

	pipe := r.client.Pipeline()
	pipe.Set(ctx, key, data, r.ttl)
	pipe.SAdd(ctx, r.userSessionsKey(tenantId, userId), deviceId)
	pipe.SAdd(ctx, r.tenantUsersKey(tenantId), userId)

	_, err = pipe.Exec(ctx)
	return err
}

var rotateScript = redis.NewScript(`
	local session = redis.call("GET", KEYS[1])
	if not session then return "ERR_NOT_FOUND" end

	local decoded = cjson.decode(session)
	if decoded.token_hash ~= ARGV[1] then
		return "ERR_REUSE"
	end

	redis.call("SET", KEYS[1], ARGV[2], "EX", ARGV[3])
	return "OK"
`)

func (r *RedisTokenRepository) VerifyAndRotate(ctx context.Context,
	oldToken, newToken, tenantId, userId, deviceId, ip, userAgent string) error {
	key := r.makeKey(tenantId, userId, deviceId)

	now := time.Now()
	newSession := SessionMetadata{
		TokenHash: r.hashToken(newToken),
		IP:        ip,
		UserAgent: userAgent,
		CreatedAt: now,
		ExpiresAt: now.Add(r.ttl),
	}

	newData, err := json.Marshal(newSession)
	if err != nil {
		return fmt.Errorf("failed to marshal new session metadata: %w", err)
	}

	oldHash := r.hashToken(oldToken)
	ttlSec := int(r.ttl.Seconds())

	result, err := rotateScript.Run(ctx, r.client, []string{key}, oldHash, string(newData), ttlSec).Result()
	if err != nil {
		return fmt.Errorf("failed to execute rotate script: %w", err)
	}

	resStr, ok := result.(string)
	if !ok {
		return fmt.Errorf("unexpected script result type")
	}

	if resStr == "ERR_NOT_FOUND" {
		return fmt.Errorf("session not found or expired")
	}

	if resStr == "ERR_REUSE" {
		err := r.RevokeAllUserSessions(ctx, tenantId, userId)
		if err != nil {
			return fmt.Errorf("SECURITY ALERT: token reuse detected, error while revoking sessions: %w", err)
		}
		return fmt.Errorf("SECURITY ALERT: token reuse detected, all sessions revoked")
	}

	return nil
}

func (r *RedisTokenRepository) RevokeSession(ctx context.Context, tenantId, userId, deviceId string) error {
	key := r.makeKey(tenantId, userId, deviceId)

	pipe := r.client.Pipeline()
	pipe.Del(ctx, key)
	pipe.SRem(ctx, r.userSessionsKey(tenantId, userId), deviceId)
	_, err := pipe.Exec(ctx)

	return err
}

func (r *RedisTokenRepository) GetActiveSessions(ctx context.Context, tenantId, userId string) (map[string]SessionMetadata, error) {
	deviceIds, err := r.client.SMembers(ctx, r.userSessionsKey(tenantId, userId)).Result()
	if err != nil {
		return nil, err
	}

	sessions := make(map[string]SessionMetadata)
	if len(deviceIds) == 0 {
		return sessions, nil
	}

	pipe := r.client.Pipeline()
	var cmds []*redis.StringCmd
	for _, deviceId := range deviceIds {
		key := r.makeKey(tenantId, userId, deviceId)
		cmds = append(cmds, pipe.Get(ctx, key))
	}

	_, _ = pipe.Exec(ctx)

	for i, cmd := range cmds {
		val, err := cmd.Result()
		if err == nil {
			var meta SessionMetadata
			if json.Unmarshal([]byte(val), &meta) == nil {
				sessions[deviceIds[i]] = meta
			}
		} else if errors.Is(err, redis.Nil) {
			// Clean up expired session from the set
			r.client.SRem(ctx, r.userSessionsKey(tenantId, userId), deviceIds[i])
		}
	}

	return sessions, nil
}

func (r *RedisTokenRepository) RevokeAllUserSessions(ctx context.Context, tenantId, userId string) error {
	deviceIds, err := r.client.SMembers(ctx, r.userSessionsKey(tenantId, userId)).Result()
	if err != nil {
		return err
	}

	if len(deviceIds) == 0 {
		return nil
	}

	pipe := r.client.Pipeline()
	for _, deviceId := range deviceIds {
		key := r.makeKey(tenantId, userId, deviceId)
		pipe.Del(ctx, key)
	}
	pipe.Del(ctx, r.userSessionsKey(tenantId, userId))

	_, err = pipe.Exec(ctx)
	return err
}

func (r *RedisTokenRepository) RevokeTenantSessions(ctx context.Context, tenantId string) error {
	userIds, err := r.client.SMembers(ctx, r.tenantUsersKey(tenantId)).Result()
	if err != nil {
		return err
	}

	for _, userId := range userIds {
		err := r.RevokeAllUserSessions(ctx, tenantId, userId)
		if err != nil {
			return err
		}
	}

	return r.client.Del(ctx, r.tenantUsersKey(tenantId)).Err()
}
