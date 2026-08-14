package middleware

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Recovery(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if recovered := recover(); recovered != nil {
				requestID, _ := c.Get(RequestIDKey)

				logger.Error(
					"panic recovered",
					"request_id", requestID,
					"panic", recovered,
				)

				c.AbortWithStatusJSON(
					http.StatusInternalServerError,
					gin.H{
						"error": "internal server error",
					},
				)
			}
		}()

		c.Next()
	}
}
