package database

import (
	"os"
	"testing"

	"ecommerce-backend/configs"

	"github.com/joho/godotenv"
)

func TestNewPostgres(t *testing.T) {
	if err := godotenv.Load("../../../.env"); err != nil {
		t.Fatalf("failed to load environment: %v", err)
	}

	cfg := configs.DatabaseConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}

	db, err := NewPostgres(cfg)
	if err != nil {
		t.Fatalf("failed to connect to PostgreSQL: %v", err)
	}

	if db == nil {
		t.Fatal("expected database connection, got nil")
	}

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("failed to get SQL database: %v", err)
	}

	defer sqlDB.Close()
}
