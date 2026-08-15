package otp

import "time"

type OTP struct {
	ID           uint64
	UserID       uint64
	CodeHash     string
	Purpose      string
	ExpiresAt    time.Time
	IsUsed       bool
	AttemptCount int16
	CreatedAt    time.Time
}
