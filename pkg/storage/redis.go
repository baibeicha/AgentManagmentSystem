package storage

import (
	"AgentManagmentSystem/pkg/config"
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(ctx context.Context, cfg *config.Config) (*redis.Client, error) {
	db := redis.NewClient(&redis.Options{
		Addr:         cfg.Redis.Addr,
		Password:     cfg.Redis.Password,
		DB:           cfg.Redis.DB,
		Username:     cfg.Redis.User,
		MaxRetries:   cfg.Redis.MaxRetries,
		DialTimeout:  cfg.Redis.DialTimeout,
		ReadTimeout:  cfg.Redis.Timeout,
		WriteTimeout: cfg.Redis.Timeout,
	})

	if err := db.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis server: %w", err)
	}

	return db, nil
}
