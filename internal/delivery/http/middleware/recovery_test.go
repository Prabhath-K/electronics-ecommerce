package middleware

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRecovery(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var logBuffer bytes.Buffer

	logger := slog.New(
		slog.NewJSONHandler(
			&logBuffer,
			&slog.HandlerOptions{
				Level: slog.LevelError,
			},
		),
	)

	router := gin.New()

	router.Use(RequestID())
	router.Use(Recovery(logger))

	router.GET("/panic", func(c *gin.Context) {
		panic("test panic")
	})

	req := httptest.NewRequest("GET", "/panic", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != 500 {
		t.Fatalf(
			"expected status 500, got %d",
			rec.Code,
		)
	}

	var response map[string]string

	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf(
			"failed to decode response: %v",
			err,
		)
	}

	if response["error"] != "internal server error" {
		t.Fatalf(
			"expected error %q, got %q",
			"internal server error",
			response["error"],
		)
	}

	requestID := rec.Header().Get("X-Request-ID")

	if requestID == "" {
		t.Fatal("expected X-Request-ID response header")
	}

	if logBuffer.Len() == 0 {
		t.Fatal("expected panic to be logged")
	}
}
