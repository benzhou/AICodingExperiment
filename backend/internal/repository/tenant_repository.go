package repository

import (
	"backend/internal/db"
	"backend/internal/models"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrTenantNotFound = errors.New("tenant not found")
	ErrTenantExists   = errors.New("tenant with this name already exists")
)

// TenantRepository defines operations for managing tenants
type TenantRepository interface {
	CreateTenant(tenant *models.Tenant) error
	GetTenantByID(id string) (*models.Tenant, error)
	GetTenantByName(name string) (*models.Tenant, error)
	UpdateTenant(tenant *models.Tenant) error
	DeleteTenant(id string) error
	GetAllTenants() ([]models.Tenant, error)
	AssignUserToTenant(userID, tenantID string) error
	RemoveUserFromTenant(userID, tenantID string) error
	GetUserTenants(userID string) ([]models.Tenant, error)
	GetTenantUsers(tenantID string) ([]string, error) // Returns list of user IDs
}

// PostgresTenantRepository implements TenantRepository for PostgreSQL
type PostgresTenantRepository struct {
	db *sql.DB
}

// NewTenantRepository creates a new tenant repository
func NewTenantRepository() TenantRepository {
	if db.DB == nil {
		// Return a mock repository for development
		return &MockTenantRepository{
			tenants:     make(map[string]*models.Tenant),
			tenantUsers: make(map[string][]string), // tenantID -> []userID
			userTenants: make(map[string][]string), // userID -> []tenantID
		}
	}
	return &PostgresTenantRepository{
		db: db.DB,
	}
}

// CreateTenant creates a new tenant
func (r *PostgresTenantRepository) CreateTenant(tenant *models.Tenant) error {
	// Check if tenant with the same name already exists
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM tenants WHERE name = $1)", tenant.Name).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		return ErrTenantExists
	}

	query := `
		INSERT INTO tenants (name, description, active)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`

	return r.db.QueryRow(
		query,
		tenant.Name,
		tenant.Description,
		tenant.Active,
	).Scan(&tenant.ID, &tenant.CreatedAt, &tenant.UpdatedAt)
}

// GetTenantByID retrieves a tenant by ID
func (r *PostgresTenantRepository) GetTenantByID(id string) (*models.Tenant, error) {
	query := `
		SELECT id, name, description, active, created_at, updated_at
		FROM tenants
		WHERE id = $1
	`

	var tenant models.Tenant
	err := r.db.QueryRow(query, id).Scan(
		&tenant.ID,
		&tenant.Name,
		&tenant.Description,
		&tenant.Active,
		&tenant.CreatedAt,
		&tenant.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrTenantNotFound
	}

	if err != nil {
		return nil, err
	}

	return &tenant, nil
}

// GetTenantByName retrieves a tenant by name
func (r *PostgresTenantRepository) GetTenantByName(name string) (*models.Tenant, error) {
	query := `
		SELECT id, name, description, active, created_at, updated_at
		FROM tenants
		WHERE name = $1
	`

	var tenant models.Tenant
	err := r.db.QueryRow(query, name).Scan(
		&tenant.ID,
		&tenant.Name,
		&tenant.Description,
		&tenant.Active,
		&tenant.CreatedAt,
		&tenant.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrTenantNotFound
	}

	if err != nil {
		return nil, err
	}

	return &tenant, nil
}

// UpdateTenant updates a tenant
func (r *PostgresTenantRepository) UpdateTenant(tenant *models.Tenant) error {
	// Check if another tenant with the same name already exists
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM tenants WHERE name = $1 AND id != $2", tenant.Name, tenant.ID).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		return ErrTenantExists
	}

	query := `
		UPDATE tenants
		SET name = $1, description = $2, active = $3, updated_at = NOW()
		WHERE id = $4
		RETURNING updated_at
	`

	var updatedAt time.Time
	err = r.db.QueryRow(query, tenant.Name, tenant.Description, tenant.Active, tenant.ID).Scan(&updatedAt)

	if err == sql.ErrNoRows {
		return ErrTenantNotFound
	}

	if err != nil {
		return err
	}

	tenant.UpdatedAt = updatedAt
	return nil
}

// DeleteTenant deletes a tenant
func (r *PostgresTenantRepository) DeleteTenant(id string) error {
	query := "DELETE FROM tenants WHERE id = $1"
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrTenantNotFound
	}

	return nil
}

// GetAllTenants retrieves all tenants
func (r *PostgresTenantRepository) GetAllTenants() ([]models.Tenant, error) {
	query := `
		SELECT id, name, description, active, created_at, updated_at
		FROM tenants
		ORDER BY name
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tenants []models.Tenant
	for rows.Next() {
		var tenant models.Tenant
		err := rows.Scan(
			&tenant.ID,
			&tenant.Name,
			&tenant.Description,
			&tenant.Active,
			&tenant.CreatedAt,
			&tenant.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		tenants = append(tenants, tenant)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tenants, nil
}

// AssignUserToTenant assigns a user to a tenant
func (r *PostgresTenantRepository) AssignUserToTenant(userID, tenantID string) error {
	// Check if the association already exists
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM tenant_users WHERE tenant_id = $1 AND user_id = $2)", tenantID, userID).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		return nil // Already assigned, no error
	}

	// Add the association
	query := `
		INSERT INTO tenant_users (tenant_id, user_id)
		VALUES ($1, $2)
	`
	_, err = r.db.Exec(query, tenantID, userID)
	return err
}

// RemoveUserFromTenant removes a user from a tenant
func (r *PostgresTenantRepository) RemoveUserFromTenant(userID, tenantID string) error {
	query := `
		DELETE FROM tenant_users
		WHERE tenant_id = $1 AND user_id = $2
	`
	result, err := r.db.Exec(query, tenantID, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("user not associated with this tenant")
	}

	return nil
}

// GetUserTenants gets all tenants for a user
func (r *PostgresTenantRepository) GetUserTenants(userID string) ([]models.Tenant, error) {
	query := `
		SELECT t.id, t.name, t.description, t.active, t.created_at, t.updated_at
		FROM tenants t
		JOIN tenant_users tu ON t.id = tu.tenant_id
		WHERE tu.user_id = $1
		ORDER BY t.name
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tenants []models.Tenant
	for rows.Next() {
		var tenant models.Tenant
		err := rows.Scan(
			&tenant.ID,
			&tenant.Name,
			&tenant.Description,
			&tenant.Active,
			&tenant.CreatedAt,
			&tenant.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		tenants = append(tenants, tenant)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tenants, nil
}

// GetTenantUsers gets all user IDs for a tenant
func (r *PostgresTenantRepository) GetTenantUsers(tenantID string) ([]string, error) {
	query := `
		SELECT user_id
		FROM tenant_users
		WHERE tenant_id = $1
	`

	rows, err := r.db.Query(query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userIDs []string
	for rows.Next() {
		var userID string
		if err := rows.Scan(&userID); err != nil {
			return nil, err
		}
		userIDs = append(userIDs, userID)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return userIDs, nil
}

// MockTenantRepository is a mock implementation for development
type MockTenantRepository struct {
	tenants     map[string]*models.Tenant
	nameIndex   map[string]string   // name -> id mapping
	tenantUsers map[string][]string // tenantID -> []userID
	userTenants map[string][]string // userID -> []tenantID
}

// CreateTenant creates a tenant in the mock repository
func (r *MockTenantRepository) CreateTenant(tenant *models.Tenant) error {
	// Initialize nameIndex if it doesn't exist
	if r.nameIndex == nil {
		r.nameIndex = make(map[string]string)
	}

	// Check if tenant with the same name exists
	if _, exists := r.nameIndex[tenant.Name]; exists {
		return ErrTenantExists
	}

	if tenant.ID == "" {
		// Generate a dummy ID - in real implementation we'd use UUID
		tenant.ID = "mock-" + time.Now().Format("20060102150405")
	}
	tenant.CreatedAt = time.Now()
	tenant.UpdatedAt = time.Now()

	r.tenants[tenant.ID] = tenant
	r.nameIndex[tenant.Name] = tenant.ID

	return nil
}

// GetTenantByID retrieves a tenant by ID from the mock repository
func (r *MockTenantRepository) GetTenantByID(id string) (*models.Tenant, error) {
	if tenant, exists := r.tenants[id]; exists {
		return tenant, nil
	}
	return nil, ErrTenantNotFound
}

// GetTenantByName retrieves a tenant by name from the mock repository
func (r *MockTenantRepository) GetTenantByName(name string) (*models.Tenant, error) {
	// Initialize nameIndex if it doesn't exist
	if r.nameIndex == nil {
		r.nameIndex = make(map[string]string)
	}

	id, exists := r.nameIndex[name]
	if !exists {
		return nil, ErrTenantNotFound
	}

	return r.tenants[id], nil
}

// UpdateTenant updates a tenant in the mock repository
func (r *MockTenantRepository) UpdateTenant(tenant *models.Tenant) error {
	// Initialize nameIndex if it doesn't exist
	if r.nameIndex == nil {
		r.nameIndex = make(map[string]string)
	}

	existing, exists := r.tenants[tenant.ID]
	if !exists {
		return ErrTenantNotFound
	}

	// Check if another tenant with the same name exists
	if id, nameExists := r.nameIndex[tenant.Name]; nameExists && id != tenant.ID {
		return ErrTenantExists
	}

	// Update the name index if the name has changed
	if existing.Name != tenant.Name {
		delete(r.nameIndex, existing.Name)
		r.nameIndex[tenant.Name] = tenant.ID
	}

	tenant.UpdatedAt = time.Now()
	r.tenants[tenant.ID] = tenant

	return nil
}

// DeleteTenant deletes a tenant from the mock repository
func (r *MockTenantRepository) DeleteTenant(id string) error {
	// Initialize nameIndex if it doesn't exist
	if r.nameIndex == nil {
		r.nameIndex = make(map[string]string)
	}

	tenant, exists := r.tenants[id]
	if !exists {
		return ErrTenantNotFound
	}

	delete(r.nameIndex, tenant.Name)
	delete(r.tenants, id)

	// Clean up tenant user associations
	if users, exists := r.tenantUsers[id]; exists {
		for _, userID := range users {
			// Remove tenant from user's tenant list
			if tenants, userExists := r.userTenants[userID]; userExists {
				var newTenants []string
				for _, tenantID := range tenants {
					if tenantID != id {
						newTenants = append(newTenants, tenantID)
					}
				}
				r.userTenants[userID] = newTenants
			}
		}
		delete(r.tenantUsers, id)
	}

	return nil
}

// GetAllTenants retrieves all tenants from the mock repository
func (r *MockTenantRepository) GetAllTenants() ([]models.Tenant, error) {
	var tenants []models.Tenant
	for _, tenant := range r.tenants {
		tenants = append(tenants, *tenant)
	}
	return tenants, nil
}

// AssignUserToTenant assigns a user to a tenant in the mock repository
func (r *MockTenantRepository) AssignUserToTenant(userID, tenantID string) error {
	// Check if tenant exists
	if _, exists := r.tenants[tenantID]; !exists {
		return ErrTenantNotFound
	}

	// Initialize maps if they don't exist
	if r.tenantUsers == nil {
		r.tenantUsers = make(map[string][]string)
	}
	if r.userTenants == nil {
		r.userTenants = make(map[string][]string)
	}

	// Check if user is already assigned to tenant
	users, exists := r.tenantUsers[tenantID]
	if exists {
		for _, uid := range users {
			if uid == userID {
				return nil // Already assigned
			}
		}
	}

	// Add user to tenant
	r.tenantUsers[tenantID] = append(r.tenantUsers[tenantID], userID)

	// Add tenant to user
	r.userTenants[userID] = append(r.userTenants[userID], tenantID)

	return nil
}

// RemoveUserFromTenant removes a user from a tenant in the mock repository
func (r *MockTenantRepository) RemoveUserFromTenant(userID, tenantID string) error {
	// Initialize maps if they don't exist
	if r.tenantUsers == nil {
		r.tenantUsers = make(map[string][]string)
	}
	if r.userTenants == nil {
		r.userTenants = make(map[string][]string)
	}

	// Check if tenant has users
	users, exists := r.tenantUsers[tenantID]
	if !exists {
		return errors.New("user not associated with this tenant")
	}

	// Remove user from tenant
	userFound := false
	var newUsers []string
	for _, uid := range users {
		if uid != userID {
			newUsers = append(newUsers, uid)
		} else {
			userFound = true
		}
	}

	if !userFound {
		return errors.New("user not associated with this tenant")
	}

	r.tenantUsers[tenantID] = newUsers

	// Remove tenant from user
	tenants, exists := r.userTenants[userID]
	if exists {
		var newTenants []string
		for _, tid := range tenants {
			if tid != tenantID {
				newTenants = append(newTenants, tid)
			}
		}
		r.userTenants[userID] = newTenants
	}

	return nil
}

// GetUserTenants gets all tenants for a user from the mock repository
func (r *MockTenantRepository) GetUserTenants(userID string) ([]models.Tenant, error) {
	// Initialize maps if they don't exist
	if r.userTenants == nil {
		r.userTenants = make(map[string][]string)
	}

	var tenants []models.Tenant
	tenantIDs, exists := r.userTenants[userID]
	if !exists {
		return tenants, nil
	}

	for _, id := range tenantIDs {
		if tenant, exists := r.tenants[id]; exists {
			tenants = append(tenants, *tenant)
		}
	}

	return tenants, nil
}

// GetTenantUsers gets all user IDs for a tenant from the mock repository
func (r *MockTenantRepository) GetTenantUsers(tenantID string) ([]string, error) {
	// Initialize maps if they don't exist
	if r.tenantUsers == nil {
		r.tenantUsers = make(map[string][]string)
	}

	users, exists := r.tenantUsers[tenantID]
	if !exists {
		return []string{}, nil
	}

	return users, nil
}
