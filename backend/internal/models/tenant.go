package models

import "time"

// Tenant represents a customer of the system
type Tenant struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Active      bool      `json:"active" db:"active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// TenantUser associates a user with a tenant
type TenantUser struct {
	ID        string    `json:"id" db:"id"`
	TenantID  string    `json:"tenant_id" db:"tenant_id"`
	UserID    string    `json:"user_id" db:"user_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// MatchSet represents a collection of data sources to be matched together
type MatchSet struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	TenantID    string    `json:"tenant_id" db:"tenant_id"`
	RuleID      string    `json:"rule_id" db:"rule_id"`
	CreatedBy   string    `json:"created_by" db:"created_by"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// MatchSetDataSource associates data sources with a match set
type MatchSetDataSource struct {
	ID           string    `json:"id" db:"id"`
	MatchSetID   string    `json:"match_set_id" db:"match_set_id"`
	DataSourceID string    `json:"data_source_id" db:"data_source_id"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}
