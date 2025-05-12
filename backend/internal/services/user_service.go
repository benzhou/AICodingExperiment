package services

import (
	"backend/internal/models"
	"backend/internal/repository"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// UserService provides methods for managing users
type UserService struct {
	userRepo    repository.UserRepository
	roleService *RoleService
}

// NewUserService creates a new user service
func NewUserService(
	userRepo repository.UserRepository,
	roleService *RoleService,
) *UserService {
	return &UserService{
		userRepo:    userRepo,
		roleService: roleService,
	}
}

// GetAllUsers retrieves all users (This would need to be implemented in the repository)
func (s *UserService) GetAllUsers() ([]models.User, error) {
	// For now, return empty slice with a note that this needs to be implemented
	return []models.User{}, errors.New("GetAllUsers not implemented yet")
}

// GetUserByID retrieves a user by ID
func (s *UserService) GetUserByID(userID string) (*models.User, error) {
	return s.userRepo.GetUserByID(userID)
}

// GetUserByEmail retrieves a user by email
func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	return s.userRepo.FindByEmail(email)
}

// UpdateUserRole updates a user's role
func (s *UserService) UpdateUserRole(userID string, role models.Role, operation string) error {
	// Verify user exists
	_, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return err
	}

	if operation == "add" {
		return s.roleService.AssignRoleToUser(userID, role)
	} else if operation == "remove" {
		return s.roleService.RemoveRoleFromUser(userID, role)
	}

	return nil
}

// SetUserAsAdmin grants admin role to a user
func (s *UserService) SetUserAsAdmin(userID string) error {
	return s.UpdateUserRole(userID, models.RoleAdmin, "add")
}

// RemoveAdminFromUser removes admin role from a user
func (s *UserService) RemoveAdminFromUser(userID string) error {
	return s.UpdateUserRole(userID, models.RoleAdmin, "remove")
}

// GetUserRoles gets all roles for a user
func (s *UserService) GetUserRoles(userID string) ([]models.Role, error) {
	return s.roleService.GetUserRoles(userID)
}

// CreateUserWithRole creates a new user and assigns a role
func (s *UserService) CreateUserWithRole(email string, name string, password string, role models.Role) (*models.User, error) {
	// Check if user already exists
	existingUser, err := s.userRepo.FindByEmail(email)
	if err == nil && existingUser != nil {
		return nil, repository.ErrUserExists
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create the user
	user := &models.User{
		Email:        email,
		Name:         name,
		PasswordHash: string(hashedPassword),
		AuthProvider: "local",
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	// Assign the role
	err = s.roleService.AssignRoleToUser(user.ID, role)
	if err != nil {
		return user, err
	}

	return user, nil
}

// IsUserAdmin checks if a user has admin role
func (s *UserService) IsUserAdmin(userID string) (bool, error) {
	return s.roleService.HasRole(userID, models.RoleAdmin)
}
