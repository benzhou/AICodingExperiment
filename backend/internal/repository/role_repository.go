package repository

import (
	"backend/internal/db"
	"backend/internal/models"
	"database/sql"
	"errors"
	"log"
)

var (
	ErrRoleNotFound = errors.New("role not found")
	ErrRoleExists   = errors.New("role already exists for user")
)

// RoleRepository defines the interface for role-related operations
type RoleRepository interface {
	AssignRoleToUser(userID string, role models.Role) error
	RemoveRoleFromUser(userID string, role models.Role) error
	GetUserRoles(userID string) ([]models.Role, error)
	HasRole(userID string, role models.Role) (bool, error)
}

// PostgresRoleRepository implements RoleRepository for PostgreSQL
type PostgresRoleRepository struct {
	db *sql.DB
}

// NewRoleRepository creates a new role repository
func NewRoleRepository() RoleRepository {
	if db.DB == nil {
		// Return a mock repository for development
		return &MockRoleRepository{
			userRoles: make(map[string][]models.Role),
		}
	}
	return &PostgresRoleRepository{
		db: db.DB,
	}
}

// AssignRoleToUser assigns a role to a user
func (r *PostgresRoleRepository) AssignRoleToUser(userID string, role models.Role) error {
	// Check if the role already exists for this user
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM user_roles WHERE user_id = $1 AND role = $2)", userID, role).Scan(&exists)
	if err != nil {
		log.Printf("Error checking if role exists: %v", err)
		return err
	}

	if exists {
		return ErrRoleExists
	}

	// Add the role
	query := `
		INSERT INTO user_roles (user_id, role)
		VALUES ($1, $2)
	`
	_, err = r.db.Exec(query, userID, role)
	return err
}

// RemoveRoleFromUser removes a role from a user
func (r *PostgresRoleRepository) RemoveRoleFromUser(userID string, role models.Role) error {
	query := `
		DELETE FROM user_roles
		WHERE user_id = $1 AND role = $2
	`
	result, err := r.db.Exec(query, userID, role)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRoleNotFound
	}

	return nil
}

// GetUserRoles returns all roles assigned to a user
func (r *PostgresRoleRepository) GetUserRoles(userID string) ([]models.Role, error) {
	query := `
		SELECT role
		FROM user_roles
		WHERE user_id = $1
	`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []models.Role
	for rows.Next() {
		var role models.Role
		if err := rows.Scan(&role); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return roles, nil
}

// HasRole checks if a user has a specific role
func (r *PostgresRoleRepository) HasRole(userID string, role models.Role) (bool, error) {
	var exists bool
	query := `
		SELECT EXISTS(
			SELECT 1
			FROM user_roles
			WHERE user_id = $1 AND role = $2
		)
	`
	err := r.db.QueryRow(query, userID, role).Scan(&exists)
	return exists, err
}

// MockRoleRepository is a mock implementation for development
type MockRoleRepository struct {
	userRoles map[string][]models.Role
}

// AssignRoleToUser assigns a role to a user in the mock repository
func (r *MockRoleRepository) AssignRoleToUser(userID string, role models.Role) error {
	roles, exists := r.userRoles[userID]
	if !exists {
		r.userRoles[userID] = []models.Role{role}
		return nil
	}

	// Check if role already exists
	for _, existingRole := range roles {
		if existingRole == role {
			return ErrRoleExists
		}
	}

	// Add the role
	r.userRoles[userID] = append(r.userRoles[userID], role)
	return nil
}

// RemoveRoleFromUser removes a role from a user in the mock repository
func (r *MockRoleRepository) RemoveRoleFromUser(userID string, role models.Role) error {
	roles, exists := r.userRoles[userID]
	if !exists {
		return ErrRoleNotFound
	}

	// Find and remove the role
	for i, existingRole := range roles {
		if existingRole == role {
			r.userRoles[userID] = append(roles[:i], roles[i+1:]...)
			return nil
		}
	}

	return ErrRoleNotFound
}

// GetUserRoles returns all roles assigned to a user in the mock repository
func (r *MockRoleRepository) GetUserRoles(userID string) ([]models.Role, error) {
	roles, exists := r.userRoles[userID]
	if !exists {
		return []models.Role{}, nil
	}
	return roles, nil
}

// HasRole checks if a user has a specific role in the mock repository
func (r *MockRoleRepository) HasRole(userID string, role models.Role) (bool, error) {
	roles, exists := r.userRoles[userID]
	if !exists {
		return false, nil
	}

	for _, existingRole := range roles {
		if existingRole == role {
			return true, nil
		}
	}

	return false, nil
}
