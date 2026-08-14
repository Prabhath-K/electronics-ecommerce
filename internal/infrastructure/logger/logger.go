package logger

import (
	"io"
	"log/slog"
)

func New(w io.Writer) *slog.Logger {
	options := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	handler := slog.NewJSONHandler(w, options)

	return slog.New(handler).With(
		"service", "ecommerce-backend",
	)
}
