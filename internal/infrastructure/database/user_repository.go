package database

import (
	"context"
	"errors"
	"fmt"

	"ecommerce-backend/internal/domain/user"
	"ecommerce-backend/internal/infrastructure/database/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

var _ user.Repository = (*UserRepository)(nil)

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Create(
	ctx context.Context,
	domainUser *user.User,
) error {
	if domainUser == nil {
		return fmt.Errorf("user cannot be nil")
	}

	modelUser := toUserModel(domainUser)

	if err := r.db.WithContext(ctx).
		Create(&modelUser).
		Error; err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	*domainUser = *toUserDomain(&modelUser)

	return nil
}

func (r *UserRepository) FindByEmail(
	ctx context.Context,
	email string,
) (*user.User, error) {
	var modelUser models.User

	err := r.db.WithContext(ctx).
		Where("email = ?", email).
		First(&modelUser).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf(
			"failed to find user by email: %w",
			err,
		)
	}

	return toUserDomain(&modelUser), nil
}

func (r *UserRepository) FindByPhone(
	ctx context.Context,
	phone string,
) (*user.User, error) {
	var modelUser models.User

	err := r.db.WithContext(ctx).
		Where("phone = ?", phone).
		First(&modelUser).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf(
			"failed to find user by phone: %w",
			err,
		)
	}

	return toUserDomain(&modelUser), nil
}

func toUserModel(domainUser *user.User) models.User {
	return models.User{
		ID:              domainUser.ID,
		RoleID:          domainUser.RoleID,
		Email:           domainUser.Email,
		PasswordHash:    domainUser.PasswordHash,
		FirstName:       domainUser.FirstName,
		LastName:        domainUser.LastName,
		Phone:           domainUser.Phone,
		IsEmailVerified: domainUser.IsEmailVerified,
		IsPhoneVerified: domainUser.IsPhoneVerified,
		Status:          domainUser.Status,
		CreatedAt:       domainUser.CreatedAt,
		UpdatedAt:       domainUser.UpdatedAt,
		LastLogin:       domainUser.LastLogin,
		DeletedAt:       domainUser.DeletedAt,
	}
}

func toUserDomain(modelUser *models.User) *user.User {
	return &user.User{
		ID:              modelUser.ID,
		RoleID:          modelUser.RoleID,
		Email:           modelUser.Email,
		PasswordHash:    modelUser.PasswordHash,
		FirstName:       modelUser.FirstName,
		LastName:        modelUser.LastName,
		Phone:           modelUser.Phone,
		IsEmailVerified: modelUser.IsEmailVerified,
		IsPhoneVerified: modelUser.IsPhoneVerified,
		Status:          modelUser.Status,
		CreatedAt:       modelUser.CreatedAt,
		UpdatedAt:       modelUser.UpdatedAt,
		LastLogin:       modelUser.LastLogin,
		DeletedAt:       modelUser.DeletedAt,
	}
}
