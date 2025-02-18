package repository

import (
	"backend/internal/db"
	"backend/internal/models"
	"database/sql"
	"errors"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrUserExists   = errors.New("user already exists")
)

type UserRepository interface {
	Create(user *models.User) error
	FindByID(id string) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	Update(user *models.User) error
}

type PostgresUserRepository struct {
	db *sql.DB
}

func NewUserRepository() UserRepository {
	if db.DB == nil {
		// Return a mock repository for development
		return &MockUserRepository{
			users: make(map[string]*models.User),
		}
	}
	return &PostgresUserRepository{
		db: db.DB,
	}
}

// Add a mock repository for development
type MockUserRepository struct {
	users map[string]*models.User
}

func (r *MockUserRepository) Create(user *models.User) error {
	if _, exists := r.users[user.Email]; exists {
		return ErrUserExists
	}
	r.users[user.Email] = user
	return nil
}

func (r *MockUserRepository) FindByID(id string) (*models.User, error) {
	for _, user := range r.users {
		if user.ID == id {
			return user, nil
		}
	}
	return nil, ErrUserNotFound
}

func (r *MockUserRepository) FindByEmail(email string) (*models.User, error) {
	if user, exists := r.users[email]; exists {
		return user, nil
	}
	return nil, ErrUserNotFound
}

func (r *MockUserRepository) Update(user *models.User) error {
	if _, exists := r.users[user.Email]; !exists {
		return ErrUserNotFound
	}
	r.users[user.Email] = user
	return nil
}

func (r *PostgresUserRepository) Create(user *models.User) error {
	query := `
		INSERT INTO users (email, name, password_hash, auth_provider)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at`

	return r.db.QueryRow(
		query,
		user.Email,
		user.Name,
		user.PasswordHash,
		user.AuthProvider,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

func (r *PostgresUserRepository) FindByID(id string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, email, name, password_hash, auth_provider, created_at, updated_at
		FROM users
		WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.PasswordHash,
		&user.AuthProvider,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}
	return user, err
}

func (r *PostgresUserRepository) FindByEmail(email string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, email, name, password_hash, auth_provider, created_at, updated_at
		FROM users
		WHERE email = $1`

	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.PasswordHash,
		&user.AuthProvider,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}
	return user, err
}

func (r *PostgresUserRepository) Update(user *models.User) error {
	query := `
		UPDATE users
		SET email = $1, name = $2, password_hash = $3, auth_provider = $4, updated_at = NOW()
		WHERE id = $5
		RETURNING updated_at`

	err := r.db.QueryRow(
		query,
		user.Email,
		user.Name,
		user.PasswordHash,
		user.AuthProvider,
		user.ID,
	).Scan(&user.UpdatedAt)

	if err == sql.ErrNoRows {
		return ErrUserNotFound
	}
	return err
}
