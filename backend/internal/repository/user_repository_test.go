package repository

import (
	"backend/internal/models"
	"backend/internal/testutil"
	"testing"
	"time"
)

func TestUserRepository_Create(t *testing.T) {
	db := testutil.SetupTestDB(t)
	repo := &PostgresUserRepository{db: db}

	tests := []struct {
		name    string
		user    *models.User
		wantErr bool
	}{
		{
			name: "valid user",
			user: &models.User{
				Email:        "test@example.com",
				Name:         "Test User",
				PasswordHash: "hashed_password",
				AuthProvider: "local",
			},
			wantErr: false,
		},
		{
			name: "duplicate email",
			user: &models.User{
				Email:        "test@example.com",
				Name:         "Another User",
				PasswordHash: "different_password",
				AuthProvider: "local",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Create(tt.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				if tt.user.ID == "" {
					t.Error("Create() didn't set user ID")
				}
				if tt.user.CreatedAt.IsZero() {
					t.Error("Create() didn't set CreatedAt")
				}
			}
		})
	}
}

func TestUserRepository_FindByEmail(t *testing.T) {
	db := testutil.SetupTestDB(t)
	repo := &PostgresUserRepository{db: db}

	// Create test user
	testUser := &models.User{
		Email:        "test@example.com",
		Name:         "Test User",
		PasswordHash: "hashed_password",
		AuthProvider: "local",
	}
	if err := repo.Create(testUser); err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	tests := []struct {
		name    string
		email   string
		want    *models.User
		wantErr error
	}{
		{
			name:    "existing user",
			email:   "test@example.com",
			want:    testUser,
			wantErr: nil,
		},
		{
			name:    "non-existent user",
			email:   "nonexistent@example.com",
			want:    nil,
			wantErr: ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.FindByEmail(tt.email)
			if err != tt.wantErr {
				t.Errorf("FindByEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr == nil {
				if got.Email != tt.want.Email {
					t.Errorf("FindByEmail() got = %v, want %v", got.Email, tt.want.Email)
				}
			}
		})
	}
}
