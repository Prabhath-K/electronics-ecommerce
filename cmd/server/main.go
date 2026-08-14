package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"ecommerce-backend/configs"
	"ecommerce-backend/internal/bootstrap"
)

func main() {
	if err := configs.LoadEnv(); err != nil {
		log.Fatal(err)
	}

	cfg, err := configs.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	app, err := bootstrap.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	app.Logger.Info(
		"application initialized",
		"environment", app.Config.Server.Environment,
		"port", app.Config.Server.Port,
	)

	serverErr := make(chan error, 1)

	go func() {
		serverErr <- app.HTTP.Start()
	}()

	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	select {
	case err := <-serverErr:
		if !errors.Is(err, http.ErrServerClosed) {
			app.Logger.Error(
				"HTTP server failed",
				"error", err,
			)
		}

	case <-ctx.Done():
		app.Logger.Info("shutdown signal received")
	}

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	if err := app.HTTP.Shutdown(shutdownCtx); err != nil {
		app.Logger.Error(
			"HTTP server shutdown failed",
			"error", err,
		)
	}

	if err := app.Close(); err != nil {
		app.Logger.Error(
			"application resources shutdown failed",
			"error", err,
		)
	}

	app.Logger.Info("application shutdown completed")
}
