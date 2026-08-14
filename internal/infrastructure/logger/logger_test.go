package logger

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestNew(t *testing.T) {
	var buffer bytes.Buffer

	log := New(&buffer)

	if log == nil {
		t.Fatal("expected logger, got nil")
	}

	log.Info(
		"test message",
		"user_id", 42,
	)

	var record map[string]any

	if err := json.Unmarshal(buffer.Bytes(), &record); err != nil {
		t.Fatalf("failed to decode log output: %v", err)
	}

	if record["level"] != "INFO" {
		t.Errorf("expected level INFO, got %v", record["level"])
	}

	if record["msg"] != "test message" {
		t.Errorf("expected message 'test message', got %v", record["msg"])
	}

	if record["service"] != "ecommerce-backend" {
		t.Errorf("expected service ecommerce-backend, got %v", record["service"])
	}

	if record["user_id"] != float64(42) {
		t.Errorf("expected user_id 42, got %v", record["user_id"])
	}
}

func TestLogLevelFiltering(t *testing.T) {
	var buffer bytes.Buffer

	log := New(&buffer)

	log.Debug("debug message")

	if buffer.Len() != 0 {
		t.Fatal("expected DEBUG log to be filtered out")
	}

	log.Info("info message")

	if buffer.Len() == 0 {
		t.Fatal("expected INFO log to be written")
	}
}
