package database

import (
	"fmt"
	"net/url"
	"time"

	"ecommerce-backend/configs"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgres(cfg configs.DatabaseConfig) (*gorm.DB, error) {
	dsnURL := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(cfg.User, cfg.Password),
		Host:   fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Path:   "/" + cfg.Name,
	}

	query := dsnURL.Query()
	query.Set("sslmode", cfg.SSLMode)
	dsnURL.RawQuery = query.Encode()

	db, err := gorm.Open(
		postgres.Open(dsnURL.String()),
		&gorm.Config{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize PostgreSQL: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get SQL database: %w", err)
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(time.Hour)
	sqlDB.SetConnMaxIdleTime(10 * time.Minute)

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping PostgreSQL: %w", err)
	}

	return db, nil
}
