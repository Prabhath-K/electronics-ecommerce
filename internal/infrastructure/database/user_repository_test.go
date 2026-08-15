package database

import (
	"context"
	"fmt"
	"testing"
	"time"

	"ecommerce-backend/configs"
	"ecommerce-backend/internal/domain/user"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

func TestUserRepository_CreateAndFind(t *testing.T) {
	if err := godotenv.Load("../../../.env"); err != nil {
		t.Fatalf("failed to load environment: %v", err)
	}

	cfg, err := configs.LoadConfig()
	require.NoError(t, err)

	db, err := NewPostgres(cfg.Database)
	require.NoError(t, err)

	sqlDB, err := db.DB()
	require.NoError(t, err)

	t.Cleanup(func() {
		sqlDB.Close()
	})

	repo := NewUserRepository(db)
	ctx := context.Background()

	roleName := fmt.Sprintf(
		"repository-test-role-%d",
		time.Now().UnixNano(),
	)

	var roleID uint64

	err = db.Raw(
		`INSERT INTO role (role_name) VALUES (?) RETURNING role_id`,
		roleName,
	).Scan(&roleID).Error
	require.NoError(t, err)
	require.NotZero(t, roleID)

	t.Cleanup(func() {
		db.Exec(
			`DELETE FROM "user" WHERE email = ?`,
			"repository-test@example.com",
		)

		db.Exec(
			`DELETE FROM role WHERE role_id = ?`,
			roleID,
		)
	})

	phone := "9876543210"

	newUser := &user.User{
		RoleID:          roleID,
		Email:           "repository-test@example.com",
		PasswordHash:    "hashed-password",
		FirstName:       "Repository",
		LastName:        "Test",
		Phone:           &phone,
		IsEmailVerified: false,
		IsPhoneVerified: false,
		Status:          "active",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	err = repo.Create(ctx, newUser)
	require.NoError(t, err)
	require.NotZero(t, newUser.ID)

	foundByEmail, err := repo.FindByEmail(ctx, newUser.Email)
	require.NoError(t, err)
	require.NotNil(t, foundByEmail)

	require.Equal(t, newUser.ID, foundByEmail.ID)
	require.Equal(t, newUser.Email, foundByEmail.Email)
	require.Equal(t, newUser.FirstName, foundByEmail.FirstName)
	require.Equal(t, newUser.LastName, foundByEmail.LastName)

	foundByPhone, err := repo.FindByPhone(ctx, phone)
	require.NoError(t, err)
	require.NotNil(t, foundByPhone)

	require.Equal(t, newUser.ID, foundByPhone.ID)
	require.Equal(t, newUser.Email, foundByPhone.Email)
}

func TestUserRepository_FindByEmail_NotFound(t *testing.T) {
	if err := godotenv.Load("../../../.env"); err != nil {
		t.Fatalf("failed to load environment: %v", err)
	}

	cfg, err := configs.LoadConfig()
	require.NoError(t, err)

	db, err := NewPostgres(cfg.Database)
	require.NoError(t, err)

	sqlDB, err := db.DB()
	require.NoError(t, err)

	t.Cleanup(func() {
		sqlDB.Close()
	})

	repo := NewUserRepository(db)

	result, err := repo.FindByEmail(
		context.Background(),
		"does-not-exist@example.com",
	)

	require.NoError(t, err)
	require.Nil(t, result)
}

func TestUserRepository_FindByPhone_NotFound(t *testing.T) {
	if err := godotenv.Load("../../../.env"); err != nil {
		t.Fatalf("failed to load environment: %v", err)
	}

	cfg, err := configs.LoadConfig()
	require.NoError(t, err)

	db, err := NewPostgres(cfg.Database)
	require.NoError(t, err)

	sqlDB, err := db.DB()
	require.NoError(t, err)

	t.Cleanup(func() {
		sqlDB.Close()
	})

	repo := NewUserRepository(db)

	result, err := repo.FindByPhone(
		context.Background(),
		"0000000000",
	)

	require.NoError(t, err)
	require.Nil(t, result)
}

func TestUserRepository_Create_DuplicateEmail(t *testing.T) {
	if err := godotenv.Load("../../../.env"); err != nil {
		t.Fatalf("failed to load environment: %v", err)
	}

	cfg, err := configs.LoadConfig()
	require.NoError(t, err)

	db, err := NewPostgres(cfg.Database)
	require.NoError(t, err)

	repo := NewUserRepository(db)
	ctx := context.Background()

	roleName := fmt.Sprintf(
		"duplicate-email-test-role-%d",
		time.Now().UnixNano(),
	)

	var roleID uint64

	err = db.Raw(
		`INSERT INTO role (role_name) VALUES (?) RETURNING role_id`,
		roleName,
	).Scan(&roleID).Error
	require.NoError(t, err)
	require.NotZero(t, roleID)

	email := fmt.Sprintf(
		"duplicate-email-%d@example.com",
		time.Now().UnixNano(),
	)

	phone1 := fmt.Sprintf(
		"987%07d",
		time.Now().UnixNano()%10000000,
	)

	phone2 := fmt.Sprintf(
		"986%07d",
		time.Now().UnixNano()%10000000,
	)

	firstUser := &user.User{
		RoleID:          roleID,
		Email:           email,
		PasswordHash:    "hashed-password",
		FirstName:       "First",
		LastName:        "User",
		Phone:           &phone1,
		IsEmailVerified: false,
		IsPhoneVerified: false,
		Status:          "active",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	secondUser := &user.User{
		RoleID:          roleID,
		Email:           email,
		PasswordHash:    "another-hashed-password",
		FirstName:       "Second",
		LastName:        "User",
		Phone:           &phone2,
		IsEmailVerified: false,
		IsPhoneVerified: false,
		Status:          "active",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	t.Cleanup(func() {
		db.Exec(`DELETE FROM "user" WHERE email = ?`, email)
		db.Exec(`DELETE FROM role WHERE role_id = ?`, roleID)
	})

	require.NoError(t, repo.Create(ctx, firstUser))

	err = repo.Create(ctx, secondUser)
	require.Error(t, err)
}
func TestUserRepository_Create_DuplicatePhone(t *testing.T) {
	if err := godotenv.Load("../../../.env"); err != nil {
		t.Fatalf("failed to load environment: %v", err)
	}

	cfg, err := configs.LoadConfig()
	require.NoError(t, err)

	db, err := NewPostgres(cfg.Database)
	require.NoError(t, err)

	repo := NewUserRepository(db)
	ctx := context.Background()

	roleName := fmt.Sprintf(
		"duplicate-phone-test-role-%d",
		time.Now().UnixNano(),
	)

	var roleID uint64

	err = db.Raw(
		`INSERT INTO role (role_name) VALUES (?) RETURNING role_id`,
		roleName,
	).Scan(&roleID).Error
	require.NoError(t, err)
	require.NotZero(t, roleID)

	email1 := fmt.Sprintf(
		"duplicate-phone-a-%d@example.com",
		time.Now().UnixNano(),
	)

	email2 := fmt.Sprintf(
		"duplicate-phone-b-%d@example.com",
		time.Now().UnixNano(),
	)

	phone := fmt.Sprintf(
		"985%07d",
		time.Now().UnixNano()%10000000,
	)

	firstUser := &user.User{
		RoleID:          roleID,
		Email:           email1,
		PasswordHash:    "hashed-password",
		FirstName:       "First",
		LastName:        "User",
		Phone:           &phone,
		IsEmailVerified: false,
		IsPhoneVerified: false,
		Status:          "active",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	secondUser := &user.User{
		RoleID:          roleID,
		Email:           email2,
		PasswordHash:    "another-hashed-password",
		FirstName:       "Second",
		LastName:        "User",
		Phone:           &phone,
		IsEmailVerified: false,
		IsPhoneVerified: false,
		Status:          "active",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	t.Cleanup(func() {
		db.Exec(`DELETE FROM "user" WHERE email IN (?, ?)`, email1, email2)
		db.Exec(`DELETE FROM role WHERE role_id = ?`, roleID)
	})

	require.NoError(t, repo.Create(ctx, firstUser))

	err = repo.Create(ctx, secondUser)
	require.Error(t, err)
}
