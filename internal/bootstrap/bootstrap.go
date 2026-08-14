package bootstrap

import (
	"fmt"
	"log/slog"
	"os"

	"ecommerce-backend/configs"
	httpdelivery "ecommerce-backend/internal/delivery/http"
	"ecommerce-backend/internal/infrastructure/database"
	"ecommerce-backend/internal/infrastructure/logger"
	redisinfra "ecommerce-backend/internal/infrastructure/redis"

	redisclient "github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Application struct {
	Config *configs.Config
	Logger *slog.Logger
	DB     *gorm.DB
	Redis  *redisclient.Client
	HTTP   *httpdelivery.Server
}

func New(cfg *configs.Config) (*Application, error) {
	appLogger := logger.New(os.Stdout)

	db, err := database.NewPostgres(cfg.Database)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to initialize PostgreSQL: %w",
			err,
		)
	}

	redisClient, err := redisinfra.NewRedis(cfg.Redis)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to initialize Redis: %w",
			err,
		)
	}

	httpServer := httpdelivery.NewServer(
		cfg.Server.Port,
		appLogger,
	)

	return &Application{
		Config: cfg,
		Logger: appLogger,
		DB:     db,
		Redis:  redisClient,
		HTTP:   httpServer,
	}, nil
}

func (app *Application) Close() error {
	var closeErr error

	if app.Redis != nil {
		if err := app.Redis.Close(); err != nil {
			closeErr = fmt.Errorf(
				"failed to close Redis: %w",
				err,
			)
		}
	}

	if app.DB != nil {
		sqlDB, err := app.DB.DB()
		if err != nil {
			if closeErr != nil {
				return fmt.Errorf(
					"%w; failed to get SQL database for closing: %v",
					closeErr,
					err,
				)
			}

			return fmt.Errorf(
				"failed to get SQL database for closing: %w",
				err,
			)
		}

		if err := sqlDB.Close(); err != nil {
			if closeErr != nil {
				return fmt.Errorf(
					"%w; failed to close PostgreSQL: %v",
					closeErr,
					err,
				)
			}

			closeErr = fmt.Errorf(
				"failed to close PostgreSQL: %w",
				err,
			)
		}
	}

	return closeErr
}
