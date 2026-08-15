package user

import "time"

type User struct {
	ID              uint64
	RoleID          uint64
	Email           string
	PasswordHash    string
	FirstName       string
	LastName        string
    Phone           *string	
	IsEmailVerified bool
	IsPhoneVerified bool
	Status          string
	LastLogin       *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       *time.Time
}
