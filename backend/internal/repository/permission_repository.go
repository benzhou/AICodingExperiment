package repository

import (
	"backend/internal/db"
	"backend/internal/models"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrPermissionNotFound = errors.New("permission not found")
	ErrPermissionExists   = errors.New("permission already exists for this role")
)

// PermissionRepository defines operations for managing role permissions
type PermissionRepository interface {
	AssignPermissionToRole(roleName string, permission models.Permission, tenantID string) error
	RemovePermissionFromRole(roleName string, permission models.Permission, tenantID string) error
	GetRolePermissions(roleName string, tenantID string) ([]models.Permission, error)
	GetAllRolePermissions() ([]models.RolePermission, error)
	HasPermission(userID string, permission models.Permission, tenantID string) (bool, error)
}

// PostgresPermissionRepository implements PermissionRepository for PostgreSQL
type PostgresPermissionRepository struct {
	db       *sql.DB
	roleRepo RoleRepository
}

// NewPermissionRepository creates a new permission repository
func NewPermissionRepository(roleRepo RoleRepository) PermissionRepository {
	if db.DB == nil {
		// Return a mock repository for development
		return &MockPermissionRepository{
			roleRepo:        roleRepo,
			rolePermissions: make(map[string]map[models.Permission]bool), // roleName:tenantID -> permissions
		}
	}
	return &PostgresPermissionRepository{
		db:       db.DB,
		roleRepo: roleRepo,
	}
}

// AssignPermissionToRole assigns a permission to a role
func (r *PostgresPermissionRepository) AssignPermissionToRole(roleName string, permission models.Permission, tenantID string) error {
	// Check if the permission already exists for this role and tenant
	var exists bool
	var query string
	var args []interface{}

	if tenantID == "" {
		query = "SELECT EXISTS(SELECT 1 FROM role_permissions WHERE role_name = $1 AND permission = $2 AND tenant_id IS NULL)"
		args = []interface{}{roleName, permission}
	} else {
		query = "SELECT EXISTS(SELECT 1 FROM role_permissions WHERE role_name = $1 AND permission = $2 AND tenant_id = $3)"
		args = []interface{}{roleName, permission, tenantID}
	}

	err := r.db.QueryRow(query, args...).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		return ErrPermissionExists
	}

	// Add the permission
	if tenantID == "" {
		query = `
			INSERT INTO role_permissions (role_name, permission)
			VALUES ($1, $2)
		`
		_, err = r.db.Exec(query, roleName, permission)
	} else {
		query = `
			INSERT INTO role_permissions (role_name, permission, tenant_id)
			VALUES ($1, $2, $3)
		`
		_, err = r.db.Exec(query, roleName, permission, tenantID)
	}

	return err
}

// RemovePermissionFromRole removes a permission from a role
func (r *PostgresPermissionRepository) RemovePermissionFromRole(roleName string, permission models.Permission, tenantID string) error {
	var query string
	var args []interface{}

	if tenantID == "" {
		query = `
			DELETE FROM role_permissions
			WHERE role_name = $1 AND permission = $2 AND tenant_id IS NULL
		`
		args = []interface{}{roleName, permission}
	} else {
		query = `
			DELETE FROM role_permissions
			WHERE role_name = $1 AND permission = $2 AND tenant_id = $3
		`
		args = []interface{}{roleName, permission, tenantID}
	}

	result, err := r.db.Exec(query, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrPermissionNotFound
	}

	return nil
}

// GetRolePermissions returns all permissions assigned to a role for a specific tenant
func (r *PostgresPermissionRepository) GetRolePermissions(roleName string, tenantID string) ([]models.Permission, error) {
	var query string
	var args []interface{}

	if tenantID == "" {
		query = `
			SELECT permission
			FROM role_permissions
			WHERE role_name = $1 AND tenant_id IS NULL
		`
		args = []interface{}{roleName}
	} else {
		query = `
			SELECT permission
			FROM role_permissions
			WHERE role_name = $1 AND (tenant_id = $2 OR tenant_id IS NULL)
		`
		args = []interface{}{roleName, tenantID}
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []models.Permission
	for rows.Next() {
		var permission models.Permission
		if err := rows.Scan(&permission); err != nil {
			return nil, err
		}
		permissions = append(permissions, permission)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return permissions, nil
}

// GetAllRolePermissions retrieves all role permission mappings
func (r *PostgresPermissionRepository) GetAllRolePermissions() ([]models.RolePermission, error) {
	query := `
		SELECT id, role_name, permission, tenant_id
		FROM role_permissions
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []models.RolePermission
	for rows.Next() {
		var permission models.RolePermission
		var tenantID sql.NullString

		if err := rows.Scan(&permission.ID, &permission.RoleName, &permission.Permission, &tenantID); err != nil {
			return nil, err
		}

		if tenantID.Valid {
			permission.TenantID = tenantID.String
		}

		permissions = append(permissions, permission)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return permissions, nil
}

// HasPermission checks if a user has a specific permission for a tenant
func (r *PostgresPermissionRepository) HasPermission(userID string, permission models.Permission, tenantID string) (bool, error) {
	// Get user roles
	roles, err := r.roleRepo.GetUserRoles(userID)
	if err != nil {
		return false, err
	}

	// For each role, check if it has the required permission
	for _, role := range roles {
		query := `
			SELECT EXISTS(
				SELECT 1
				FROM role_permissions
				WHERE role_name = $1 AND permission = $2 AND (tenant_id = $3 OR tenant_id IS NULL)
			)
		`
		var hasPermission bool
		err := r.db.QueryRow(query, role, permission, tenantID).Scan(&hasPermission)
		if err != nil {
			return false, err
		}

		if hasPermission {
			return true, nil
		}
	}

	return false, nil
}

// MockPermissionRepository is a mock implementation for development
type MockPermissionRepository struct {
	roleRepo        RoleRepository
	rolePermissions map[string]map[models.Permission]bool // roleName:tenantID -> permissions
}

// key generates a unique key for role and tenant
func (r *MockPermissionRepository) key(roleName string, tenantID string) string {
	if tenantID == "" {
		return roleName + ":global"
	}
	return roleName + ":" + tenantID
}

// AssignPermissionToRole assigns a permission to a role in the mock repository
func (r *MockPermissionRepository) AssignPermissionToRole(roleName string, permission models.Permission, tenantID string) error {
	key := r.key(roleName, tenantID)

	// Initialize permissions map for role if it doesn't exist
	if _, exists := r.rolePermissions[key]; !exists {
		r.rolePermissions[key] = make(map[models.Permission]bool)
	}

	// Check if permission already exists
	if r.rolePermissions[key][permission] {
		return ErrPermissionExists
	}

	// Add the permission
	r.rolePermissions[key][permission] = true
	return nil
}

// RemovePermissionFromRole removes a permission from a role in the mock repository
func (r *MockPermissionRepository) RemovePermissionFromRole(roleName string, permission models.Permission, tenantID string) error {
	key := r.key(roleName, tenantID)

	// Check if role has any permissions
	permissions, exists := r.rolePermissions[key]
	if !exists {
		return ErrPermissionNotFound
	}

	// Check if permission exists
	if !permissions[permission] {
		return ErrPermissionNotFound
	}

	// Remove the permission
	delete(r.rolePermissions[key], permission)
	return nil
}

// GetRolePermissions returns all permissions assigned to a role for a specific tenant in the mock repository
func (r *MockPermissionRepository) GetRolePermissions(roleName string, tenantID string) ([]models.Permission, error) {
	var permissions []models.Permission

	// Get global permissions for this role
	globalKey := r.key(roleName, "")
	if globalPerms, exists := r.rolePermissions[globalKey]; exists {
		for perm := range globalPerms {
			permissions = append(permissions, perm)
		}
	}

	// If tenant is specified, also get tenant-specific permissions
	if tenantID != "" {
		tenantKey := r.key(roleName, tenantID)
		if tenantPerms, exists := r.rolePermissions[tenantKey]; exists {
			for perm := range tenantPerms {
				// Check if permission is already in the list (from global)
				found := false
				for _, p := range permissions {
					if p == perm {
						found = true
						break
					}
				}
				if !found {
					permissions = append(permissions, perm)
				}
			}
		}
	}

	return permissions, nil
}

// GetAllRolePermissions retrieves all role permission mappings from the mock repository
func (r *MockPermissionRepository) GetAllRolePermissions() ([]models.RolePermission, error) {
	var permissions []models.RolePermission

	for key, perms := range r.rolePermissions {
		// Parse key to get role name and tenant ID
		var roleName, tenantID string
		for i, c := range key {
			if c == ':' {
				roleName = key[:i]
				tenantID = key[i+1:]
				break
			}
		}

		// Skip invalid keys
		if roleName == "" {
			continue
		}

		// Convert "global" tenant ID to empty string
		if tenantID == "global" {
			tenantID = ""
		}

		// Add each permission
		for perm := range perms {
			permissions = append(permissions, models.RolePermission{
				ID:         "mock-" + time.Now().Format("20060102150405"),
				RoleName:   roleName,
				Permission: perm,
				TenantID:   tenantID,
			})
		}
	}

	return permissions, nil
}

// HasPermission checks if a user has a specific permission for a tenant in the mock repository
func (r *MockPermissionRepository) HasPermission(userID string, permission models.Permission, tenantID string) (bool, error) {
	// Get user roles
	roles, err := r.roleRepo.GetUserRoles(userID)
	if err != nil {
		return false, err
	}

	// For each role, check if it has the required permission
	for _, role := range roles {
		// Check global permissions
		globalKey := r.key(string(role), "")
		if globalPerms, exists := r.rolePermissions[globalKey]; exists && globalPerms[permission] {
			return true, nil
		}

		// Check tenant-specific permissions
		if tenantID != "" {
			tenantKey := r.key(string(role), tenantID)
			if tenantPerms, exists := r.rolePermissions[tenantKey]; exists && tenantPerms[permission] {
				return true, nil
			}
		}
	}

	return false, nil
}
