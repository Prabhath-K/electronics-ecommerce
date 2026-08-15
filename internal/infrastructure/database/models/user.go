package models

import "time"

type User struct {
	ID              uint64     `gorm:"column:user_id;primaryKey;autoIncrement"`
	RoleID          uint64     `gorm:"column:role_id"`
	Email           string     `gorm:"column:email"`
	PasswordHash    string     `gorm:"column:password_hash"`
	FirstName       string     `gorm:"column:first_name"`
	LastName        string     `gorm:"column:last_name"`
	Phone           *string    `gorm:"column:phone"`
	IsEmailVerified bool       `gorm:"column:is_email_verified"`
	IsPhoneVerified bool       `gorm:"column:is_phone_verified"`
	Status          string     `gorm:"column:status"`
	CreatedAt       time.Time  `gorm:"column:created_at"`
	UpdatedAt       time.Time  `gorm:"column:updated_at"`
	LastLogin       *time.Time `gorm:"column:last_login"`
	DeletedAt       *time.Time `gorm:"column:deleted_at"`
}

func (User) TableName() string {
	return "user"
}
