package models

import "time"

// Permission type for action-based permissions
type Permission string

// System-wide permissions
const (
	// Tenant management permissions
	PermManageTenants Permission = "manage:tenants"
	PermViewTenants   Permission = "view:tenants"

	// User management permissions
	PermManageUsers Permission = "manage:users"
	PermViewUsers   Permission = "view:users"

	// Role management permissions
	PermManageRoles Permission = "manage:roles"
	PermViewRoles   Permission = "view:roles"

	// Schema management permissions
	PermManageSchemas Permission = "manage:schemas"
	PermViewSchemas   Permission = "view:schemas"

	// Match set permissions
	PermCreateMatchSet    Permission = "create:matchset"
	PermViewMatchSet      Permission = "view:matchset"
	PermUpdateMatchSet    Permission = "update:matchset"
	PermDeleteMatchSet    Permission = "delete:matchset"
	PermMatchTransactions Permission = "match:transactions"

	// Data source permissions
	PermCreateDataSource Permission = "create:datasource"
	PermViewDataSource   Permission = "view:datasource"
	PermUpdateDataSource Permission = "update:datasource"
	PermDeleteDataSource Permission = "delete:datasource"
	PermUploadDataSource Permission = "upload:datasource"

	// Rule permissions
	PermCreateRule Permission = "create:rule"
	PermViewRule   Permission = "view:rule"
	PermUpdateRule Permission = "update:rule"
	PermDeleteRule Permission = "delete:rule"

	// Transaction permissions
	PermViewTransactions   Permission = "view:transactions"
	PermManageTransactions Permission = "manage:transactions"
)

// RolePermission represents a permission assigned to a role
type RolePermission struct {
	ID         string     `json:"id" db:"id"`
	RoleName   string     `json:"role_name" db:"role_name"`
	Permission Permission `json:"permission" db:"permission"`
	TenantID   string     `json:"tenant_id" db:"tenant_id"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
}
