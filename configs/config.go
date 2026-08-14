package configs

import "fmt"

// Config represents the application's configuration.
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
}

type ServerConfig struct {
	Environment string
	Port        string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       string
}

type JWTConfig struct {
	AccessSecret  string
	RefreshSecret string
	AccessExpiry  string
	RefreshExpiry string
}

func LoadConfig() (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Environment: getEnv("APP_ENV"),
			Port:        getEnv("APP_PORT"),
		},

		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST"),
			Port:     getEnv("DB_PORT"),
			User:     getEnv("DB_USER"),
			Password: getEnv("DB_PASSWORD"),
			Name:     getEnv("DB_NAME"),
			SSLMode:  getEnv("DB_SSLMODE"),
		},

		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST"),
			Port:     getEnv("REDIS_PORT"),
			Password: getEnv("REDIS_PASSWORD"),
			DB:       getEnv("REDIS_DB"),
		},

		JWT: JWTConfig{
			AccessSecret:  getEnv("JWT_ACCESS_SECRET"),
			RefreshSecret: getEnv("JWT_REFRESH_SECRET"),
			AccessExpiry:  getEnv("JWT_ACCESS_EXPIRY"),
			RefreshExpiry: getEnv("JWT_REFRESH_EXPIRY"),
		},
	}

	if err := validateConfig(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func validateConfig(cfg *Config) error {
	if cfg.Server.Environment == "" {
		return fmt.Errorf("APP_ENV is required")
	}

	if cfg.Server.Port == "" {
		return fmt.Errorf("APP_PORT is required")
	}

	if cfg.Database.Host == "" {
		return fmt.Errorf("DB_HOST is required")
	}

	if cfg.Database.Port == "" {
		return fmt.Errorf("DB_PORT is required")
	}

	if cfg.Database.User == "" {
		return fmt.Errorf("DB_USER is required")
	}

	if cfg.Database.Password == "" {
		return fmt.Errorf("DB_PASSWORD is required")
	}

	if cfg.Database.Name == "" {
		return fmt.Errorf("DB_NAME is required")
	}

	if cfg.Database.SSLMode == "" {
		return fmt.Errorf("DB_SSLMODE is required")
	}

	if cfg.Redis.Host == "" {
		return fmt.Errorf("REDIS_HOST is required")
	}

	if cfg.Redis.Port == "" {
		return fmt.Errorf("REDIS_PORT is required")
	}

	if cfg.JWT.AccessSecret == "" {
		return fmt.Errorf("JWT_ACCESS_SECRET is required")
	}

	if cfg.JWT.RefreshSecret == "" {
		return fmt.Errorf("JWT_REFRESH_SECRET is required")
	}

	if cfg.JWT.AccessExpiry == "" {
		return fmt.Errorf("JWT_ACCESS_EXPIRY is required")
	}

	if cfg.JWT.RefreshExpiry == "" {
		return fmt.Errorf("JWT_REFRESH_EXPIRY is required")
	}

	return nil
}
