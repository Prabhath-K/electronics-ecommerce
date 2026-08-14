package redis

import (
	"context"
	"fmt"
	"strconv"

	"ecommerce-backend/configs"

	"github.com/redis/go-redis/v9"
)

func NewRedis(cfg configs.RedisConfig) (*redis.Client, error) {
	db, err := strconv.Atoi(cfg.DB)
	if err != nil {
		return nil, fmt.Errorf("invalid Redis database number: %w", err)
	}

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       db,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		_ = client.Close()

		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return client, nil
}
