package services

import (
	"backend/internal/models"
	"backend/internal/repository"
	"errors"
)

// RoleService provides methods for managing user roles
type RoleService struct {
	roleRepo repository.RoleRepository
	userRepo repository.UserRepository
}

// NewRoleService creates a new role service
func NewRoleService(
	roleRepo repository.RoleRepository,
	userRepo repository.UserRepository,
) *RoleService {
	return &RoleService{
		roleRepo: roleRepo,
		userRepo: userRepo,
	}
}

// AssignRoleToUser assigns a role to a user
func (s *RoleService) AssignRoleToUser(userID string, role models.Role) error {
	// Check if user exists
	_, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return err
	}

	// Assign the role
	return s.roleRepo.AssignRoleToUser(userID, role)
}

// RemoveRoleFromUser removes a role from a user
func (s *RoleService) RemoveRoleFromUser(userID string, role models.Role) error {
	// Check if user exists
	_, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return err
	}

	// Remove the role
	return s.roleRepo.RemoveRoleFromUser(userID, role)
}

// GetUserRoles retrieves all roles assigned to a user
func (s *RoleService) GetUserRoles(userID string) ([]models.Role, error) {
	// Check if user exists
	_, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	return s.roleRepo.GetUserRoles(userID)
}

// HasRole checks if a user has a specific role
func (s *RoleService) HasRole(userID string, role models.Role) (bool, error) {
	// Check if user exists
	_, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return false, err
	}

	return s.roleRepo.HasRole(userID, role)
}

// RequireRole ensures a user has a specific role
func (s *RoleService) RequireRole(userID string, role models.Role) error {
	hasRole, err := s.HasRole(userID, role)
	if err != nil {
		return err
	}

	if !hasRole {
		return errors.New("user does not have the required role")
	}

	return nil
}

// UserHasAnyRole checks if a user has any of the specified roles
func (s *RoleService) UserHasAnyRole(userID string, roles []models.Role) (bool, error) {
	// Check if user exists
	_, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return false, err
	}

	// Get user's roles
	userRoles, err := s.roleRepo.GetUserRoles(userID)
	if err != nil {
		return false, err
	}

	// Create a map of the user's roles for efficient lookup
	userRolesMap := make(map[models.Role]bool)
	for _, role := range userRoles {
		userRolesMap[role] = true
	}

	// Check if the user has any of the specified roles
	for _, role := range roles {
		if userRolesMap[role] {
			return true, nil
		}
	}

	return false, nil
}
