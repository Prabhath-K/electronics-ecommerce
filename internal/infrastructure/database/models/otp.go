package models

import "time"

type OTP struct {
	ID           uint64    `gorm:"column:otp_id;primaryKey;autoIncrement"`
	UserID       uint64    `gorm:"column:user_id"`
	CodeHash     string    `gorm:"column:otp_code"`
	Purpose      string    `gorm:"column:purpose"`
	ExpiresAt    time.Time `gorm:"column:expires_at"`
	IsUsed       bool      `gorm:"column:is_used"`
	AttemptCount int16     `gorm:"column:attempt_count"`
	CreatedAt    time.Time `gorm:"column:created_at"`
}

func (OTP) TableName() string {
	return "otp"
}
