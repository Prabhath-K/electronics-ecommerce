package http

import (
	"context"
	"log/slog"
	stdhttp "net/http"
	"time"

	"ecommerce-backend/internal/delivery/http/middleware"

	"github.com/gin-gonic/gin"
)

type Server struct {
	HTTP *stdhttp.Server
}

func NewServer(port string, logger *slog.Logger) *Server {
	router := gin.New()

	router.Use(middleware.RequestID())
	router.Use(middleware.RequestLogger(logger))
	router.Use(middleware.Recovery(logger))

	router.GET("/health", func(c *gin.Context) {
		c.JSON(stdhttp.StatusOK, gin.H{
			"status": "ok",
		})
	})

	httpServer := &stdhttp.Server{
		Addr:              ":" + port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	return &Server{
		HTTP: httpServer,
	}
}

func (s *Server) Start() error {
	return s.HTTP.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.HTTP.Shutdown(ctx)
}
