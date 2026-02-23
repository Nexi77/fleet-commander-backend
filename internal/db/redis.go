package db

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Nexi77/fleetcommander-backend/internal/config"
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	Client *redis.Client
}

func NewRedisConnection(ctx context.Context, cfg *config.Config) (*RedisClient, error) {
	slog.Debug("Initializing Redis connection", "host", cfg.RedisUrl)

	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisUrl,
		Password: cfg.RedisPassword,
		DB:       0, // Use default database instance (0)
		// Optional but good practice: connection pooling is handled automatically,
		// but you can tweak settings like PoolSize if needed for high load.
		// PoolSize: 100,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		client.Close()
		return nil, fmt.Errorf("Unable to ping redis: %w", err)
	}

	slog.Info("Redis connection established successfully")

	return &RedisClient{
		Client: client,
	}, nil
}

func (r *RedisClient) Close() {
	if r.Client != nil {
		slog.Info("Closing Redis connection...")
		r.Client.Close()
	}
}
