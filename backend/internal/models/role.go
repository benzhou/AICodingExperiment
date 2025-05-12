package models

// Role defines access levels in the system
type Role string

const (
	// RolePreparer can upload and match transactions
	RolePreparer Role = "preparer"

	// RoleApprover can review and approve matched transactions
	RoleApprover Role = "approver"

	// RoleAdmin has full system access including user management
	RoleAdmin Role = "admin"
)

// UserRole associates a user with a specific role
type UserRole struct {
	ID     string `json:"id" db:"id"`
	UserID string `json:"user_id" db:"user_id"`
	Role   Role   `json:"role" db:"role"`
}
