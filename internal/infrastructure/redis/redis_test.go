package redis

import (
	"os"
	"testing"

	"ecommerce-backend/configs"

	"github.com/joho/godotenv"
)

func TestNewRedis(t *testing.T) {
	if err := godotenv.Load("../../../.env"); err != nil {
		t.Fatalf("failed to load environment: %v", err)
	}

	cfg := configs.RedisConfig{
		Host:     os.Getenv("REDIS_HOST"),
		Port:     os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       os.Getenv("REDIS_DB"),
	}

	client, err := NewRedis(cfg)
	if err != nil {
		t.Fatalf("failed to connect to Redis: %v", err)
	}

	if client == nil {
		t.Fatal("expected Redis client, got nil")
	}

	defer client.Close()
}
