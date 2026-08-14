package middleware

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRequestLogger(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var logBuffer bytes.Buffer

	logger := slog.New(
		slog.NewJSONHandler(
			&logBuffer,
			&slog.HandlerOptions{
				Level: slog.LevelInfo,
			},
		),
	)

	router := gin.New()

	router.Use(RequestID())
	router.Use(RequestLogger(logger))

	router.GET("/test", func(c *gin.Context) {
		c.Status(200)
	})

	req := httptest.NewRequest(
		"GET",
		"/test",
		nil,
	)

	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != 200 {
		t.Fatalf(
			"expected status 200, got %d",
			rec.Code,
		)
	}

	var record map[string]any

	if err := json.Unmarshal(
		logBuffer.Bytes(),
		&record,
	); err != nil {
		t.Fatalf(
			"failed to decode log output: %v",
			err,
		)
	}

	if record["msg"] != "http request" {
		t.Fatalf(
			"expected message %q, got %v",
			"http request",
			record["msg"],
		)
	}

	if record["method"] != "GET" {
		t.Fatalf(
			"expected method GET, got %v",
			record["method"],
		)
	}

	if record["path"] != "/test" {
		t.Fatalf(
			"expected path /test, got %v",
			record["path"],
		)
	}

	if record["status"] != float64(200) {
		t.Fatalf(
			"expected status 200, got %v",
			record["status"],
		)
	}

	if record["request_id"] == "" {
		t.Fatal("expected request_id in log")
	}

	if record["duration"] == "" {
		t.Fatal("expected duration in log")
	}
}
func TestRequestLoggerPreservesRequestID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var logBuffer bytes.Buffer

	logger := slog.New(
		slog.NewJSONHandler(
			&logBuffer,
			&slog.HandlerOptions{
				Level: slog.LevelInfo,
			},
		),
	)

	router := gin.New()

	router.Use(RequestID())
	router.Use(RequestLogger(logger))

	router.GET("/test", func(c *gin.Context) {
		c.Status(200)
	})

	req := httptest.NewRequest(
		"GET",
		"/test",
		nil,
	)

	req.Header.Set(
		"X-Request-ID",
		"review-test-123",
	)

	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != 200 {
		t.Fatalf(
			"expected status 200, got %d",
			rec.Code,
		)
	}

	var record map[string]any

	if err := json.Unmarshal(
		logBuffer.Bytes(),
		&record,
	); err != nil {
		t.Fatalf(
			"failed to decode log output: %v",
			err,
		)
	}

	if record["request_id"] != "review-test-123" {
		t.Fatalf(
			"expected request_id %q, got %v",
			"review-test-123",
			record["request_id"],
		)
	}
}
