package repository

import (
	"backend/internal/db"
	"backend/internal/models"
	"database/sql"
	"errors"
	"log"
	"strings"
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
	GetUserByID(userID string) (*models.User, error)
	GetAllUsers() ([]models.User, error)
}

type PostgresUserRepository struct {
	db              *sql.DB
	findByEmailStmt *sql.Stmt
	// Add other prepared statements as needed
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
	// Make lookup case-insensitive
	email = strings.ToLower(strings.TrimSpace(email))

	// Check for exact match first
	if user, exists := r.users[email]; exists {
		return user, nil
	}

	// If no exact match, try case-insensitive matching
	for storedEmail, user := range r.users {
		if strings.EqualFold(storedEmail, email) {
			return user, nil
		}
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

// GetUserByID retrieves a user by their ID
func (r *MockUserRepository) GetUserByID(userID string) (*models.User, error) {
	for _, user := range r.users {
		if user.ID == userID {
			return user, nil
		}
	}
	return nil, ErrUserNotFound
}

// GetAllUsers retrieves all users from the mock store
func (r *MockUserRepository) GetAllUsers() ([]models.User, error) {
	var users []models.User
	for _, user := range r.users {
		users = append(users, *user)
	}
	return users, nil
}

// Add these validation functions
func validateEmail(email string) error {
	if len(email) > 255 { // Adjust max length as needed
		return errors.New("email too long")
	}
	// Add more email validation as needed
	return nil
}

func validateUserInput(user *models.User) error {
	if err := validateEmail(user.Email); err != nil {
		return err
	}
	if len(user.Name) > 100 { // Adjust max length as needed
		return errors.New("name too long")
	}
	return nil
}

// Update the Create method to include validation
func (r *PostgresUserRepository) Create(user *models.User) error {
	if err := validateUserInput(user); err != nil {
		return err
	}

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

	// Use sql.NullString for auth_provider to handle NULL values
	var authProvider sql.NullString

	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.PasswordHash,
		&authProvider,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}

	if err != nil {
		log.Printf("Database error when finding user by ID: %v", err)
		return nil, err
	}

	// Set auth_provider to "local" if NULL
	if authProvider.Valid {
		user.AuthProvider = authProvider.String
	} else {
		user.AuthProvider = "local" // Default value
		log.Printf("Auth provider is NULL for user: %s, setting default value 'local'", user.ID)
	}

	return user, nil
}

func (r *PostgresUserRepository) FindByEmail(email string) (*models.User, error) {
	user := &models.User{}

	// Normalize email - trim whitespace and convert to lowercase
	email = strings.TrimSpace(strings.ToLower(email))

	// Use LOWER() for case-insensitive comparison
	query := `
		SELECT id, email, name, password_hash, auth_provider, created_at, updated_at
		FROM users
		WHERE LOWER(email) = LOWER($1)`

	log.Printf("Looking up user with email: %s", email)

	// Use sql.NullString for auth_provider to handle NULL values
	var authProvider sql.NullString

	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.PasswordHash,
		&authProvider,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		log.Printf("User not found with email: %s", email)
		return nil, ErrUserNotFound
	}

	if err != nil {
		log.Printf("Database error when finding user by email: %v", err)
		return nil, err
	}

	// Set auth_provider to "local" if NULL
	if authProvider.Valid {
		user.AuthProvider = authProvider.String
	} else {
		user.AuthProvider = "local" // Default value
		log.Printf("Auth provider is NULL for user: %s, setting default value 'local'", user.ID)
	}

	log.Printf("Found user with ID: %s for email: %s", user.ID, email)
	return user, nil
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

// GetUserByID retrieves a user by their ID
func (r *PostgresUserRepository) GetUserByID(userID string) (*models.User, error) {
	var user models.User

	query := "SELECT id, name, email FROM users WHERE id = $1"
	err := r.db.QueryRow(query, userID).Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		log.Println("Error querying user:", err)
		return nil, err
	}

	return &user, nil
}

func NewPostgresUserRepository(db *sql.DB) (*PostgresUserRepository, error) {
	findByEmailStmt, err := db.Prepare(`
		SELECT id, email, name, password_hash, auth_provider, created_at, updated_at
		FROM users
		WHERE LOWER(email) = LOWER($1)
	`)
	if err != nil {
		return nil, err
	}

	return &PostgresUserRepository{
		db:              db,
		findByEmailStmt: findByEmailStmt,
	}, nil
}

// Don't forget to close prepared statements
func (r *PostgresUserRepository) Close() error {
	if r.findByEmailStmt != nil {
		return r.findByEmailStmt.Close()
	}
	return nil
}

// GetAllUsers retrieves all users from the database
func (r *PostgresUserRepository) GetAllUsers() ([]models.User, error) {
	query := `
		SELECT id, name, email
		FROM users
		ORDER BY name
	`
	rows, err := r.db.Query(query)
	if err != nil {
		log.Printf("Error querying all users: %v", err)
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			log.Printf("Error scanning user row: %v", err)
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating user rows: %v", err)
		return nil, err
	}

	return users, nil
}
