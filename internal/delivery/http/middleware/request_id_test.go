package middleware

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRequestIDGeneratesID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(RequestID())

	router.GET("/test", func(c *gin.Context) {
		requestID, exists := c.Get(RequestIDKey)

		if !exists {
			t.Fatal("expected request ID in context")
		}

		if requestID == "" {
			t.Fatal("expected non-empty request ID")
		}

		c.Status(200)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != 200 {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	responseRequestID := rec.Header().Get("X-Request-ID")

	if responseRequestID == "" {
		t.Fatal("expected X-Request-ID response header")
	}
}

func TestRequestIDPreservesExistingID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(RequestID())

	router.GET("/test", func(c *gin.Context) {
		requestID, exists := c.Get(RequestIDKey)

		if !exists {
			t.Fatal("expected request ID in context")
		}

		if requestID != "test-request-id" {
			t.Fatalf(
				"expected request ID %q, got %q",
				"test-request-id",
				requestID,
			)
		}

		c.Status(200)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-Request-ID", "test-request-id")

	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != 200 {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	responseRequestID := rec.Header().Get("X-Request-ID")

	if responseRequestID != "test-request-id" {
		t.Fatalf(
			"expected response request ID %q, got %q",
			"test-request-id",
			responseRequestID,
		)
	}
}
