package http

import (
	"bytes"
	"encoding/json"
	"log/slog"
	stdhttp "net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func newTestLogger(buffer *bytes.Buffer) *slog.Logger {
	return slog.New(
		slog.NewJSONHandler(
			buffer,
			&slog.HandlerOptions{
				Level: slog.LevelInfo,
			},
		),
	)
}

func TestHealth(t *testing.T) {
	var logBuffer bytes.Buffer

	logger := newTestLogger(&logBuffer)

	server := NewServer("8080", logger)

	req := httptest.NewRequest(
		stdhttp.MethodGet,
		"/health",
		nil,
	)

	rec := httptest.NewRecorder()

	server.HTTP.Handler.ServeHTTP(rec, req)

	if rec.Code != stdhttp.StatusOK {
		t.Fatalf(
			"expected status %d, got %d",
			stdhttp.StatusOK,
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

	if response["status"] != "ok" {
		t.Fatalf(
			"expected status %q, got %q",
			"ok",
			response["status"],
		)
	}

	requestID := rec.Header().Get("X-Request-ID")

	if requestID == "" {
		t.Fatal("expected X-Request-ID response header")
	}

	if logBuffer.Len() == 0 {
		t.Fatal("expected HTTP request to be logged")
	}

	var record map[string]any

	if err := json.Unmarshal(logBuffer.Bytes(), &record); err != nil {
		t.Fatalf(
			"failed to decode HTTP log: %v",
			err,
		)
	}

	if record["msg"] != "http request" {
		t.Fatalf(
			"expected log message %q, got %v",
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

	if record["path"] != "/health" {
		t.Fatalf(
			"expected path /health, got %v",
			record["path"],
		)
	}

	if record["status"] != float64(200) {
		t.Fatalf(
			"expected status 200, got %v",
			record["status"],
		)
	}

	if record["request_id"] != requestID {
		t.Fatalf(
			"expected request ID %q in log, got %v",
			requestID,
			record["request_id"],
		)
	}
}

func TestRecoveryIntegration(t *testing.T) {
	var logBuffer bytes.Buffer

	logger := newTestLogger(&logBuffer)

	server := NewServer("8080", logger)

	router := server.HTTP.Handler.(*gin.Engine)

	router.GET("/panic-test", func(c *gin.Context) {
		panic("integration test panic")
	})

	req := httptest.NewRequest(
		stdhttp.MethodGet,
		"/panic-test",
		nil,
	)

	rec := httptest.NewRecorder()

	server.HTTP.Handler.ServeHTTP(rec, req)

	if rec.Code != stdhttp.StatusInternalServerError {
		t.Fatalf(
			"expected status %d, got %d",
			stdhttp.StatusInternalServerError,
			rec.Code,
		)
	}

	requestID := rec.Header().Get("X-Request-ID")

	if requestID == "" {
		t.Fatal("expected X-Request-ID response header")
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

	if logBuffer.Len() == 0 {
		t.Fatal("expected panic and request logs")
	}

	// The request logger writes the final access log after recovery.
	var requestLog map[string]any

	lines := bytes.Split(
		bytes.TrimSpace(logBuffer.Bytes()),
		[]byte("\n"),
	)

	if len(lines) == 0 {
		t.Fatal("expected at least one log entry")
	}

	if err := json.Unmarshal(lines[len(lines)-1], &requestLog); err != nil {
		t.Fatalf(
			"failed to decode request log: %v",
			err,
		)
	}

	if requestLog["msg"] != "http request" {
		t.Fatalf(
			"expected final log message %q, got %v",
			"http request",
			requestLog["msg"],
		)
	}

	if requestLog["method"] != "GET" {
		t.Fatalf(
			"expected method GET, got %v",
			requestLog["method"],
		)
	}

	if requestLog["path"] != "/panic-test" {
		t.Fatalf(
			"expected path /panic-test, got %v",
			requestLog["path"],
		)
	}

	if requestLog["status"] != float64(500) {
		t.Fatalf(
			"expected status 500, got %v",
			requestLog["status"],
		)
	}

	if requestLog["request_id"] != requestID {
		t.Fatalf(
			"expected request ID %q in log, got %v",
			requestID,
			requestLog["request_id"],
		)
	}
}
