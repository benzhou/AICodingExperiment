package repository

import (
	"backend/internal/db"
	"backend/internal/models"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrMatchSetNotFound = errors.New("match set not found")
	ErrMatchSetExists   = errors.New("match set with this name already exists for this tenant")
)

// MatchSetRepository defines operations for managing match sets
type MatchSetRepository interface {
	CreateMatchSet(matchSet *models.MatchSet) error
	GetMatchSetByID(id string) (*models.MatchSet, error)
	GetMatchSetsByTenant(tenantID string) ([]models.MatchSet, error)
	UpdateMatchSet(matchSet *models.MatchSet) error
	DeleteMatchSet(id string) error
	AddDataSourceToMatchSet(matchSetID, dataSourceID string) error
	RemoveDataSourceFromMatchSet(matchSetID, dataSourceID string) error
	GetMatchSetDataSources(matchSetID string) ([]models.DataSource, error)
}

// PostgresMatchSetRepository implements MatchSetRepository for PostgreSQL
type PostgresMatchSetRepository struct {
	db *sql.DB
}

// NewMatchSetRepository creates a new match set repository
func NewMatchSetRepository() MatchSetRepository {
	if db.DB == nil {
		// Return a mock repository for development
		return &MockMatchSetRepository{
			matchSets:           make(map[string]*models.MatchSet),
			matchSetDataSources: make(map[string][]string), // matchSetID -> []dataSourceID
		}
	}
	return &PostgresMatchSetRepository{
		db: db.DB,
	}
}

// CreateMatchSet creates a new match set
func (r *PostgresMatchSetRepository) CreateMatchSet(matchSet *models.MatchSet) error {
	// Check if match set with the same name already exists for this tenant
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM match_sets WHERE name = $1 AND tenant_id = $2)", matchSet.Name, matchSet.TenantID).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		return ErrMatchSetExists
	}

	query := `
		INSERT INTO match_sets (name, description, tenant_id, rule_id, created_by)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`

	return r.db.QueryRow(
		query,
		matchSet.Name,
		matchSet.Description,
		matchSet.TenantID,
		matchSet.RuleID,
		matchSet.CreatedBy,
	).Scan(&matchSet.ID, &matchSet.CreatedAt, &matchSet.UpdatedAt)
}

// GetMatchSetByID retrieves a match set by ID
func (r *PostgresMatchSetRepository) GetMatchSetByID(id string) (*models.MatchSet, error) {
	query := `
		SELECT id, name, description, tenant_id, rule_id, created_by, created_at, updated_at
		FROM match_sets
		WHERE id = $1
	`

	var matchSet models.MatchSet
	err := r.db.QueryRow(query, id).Scan(
		&matchSet.ID,
		&matchSet.Name,
		&matchSet.Description,
		&matchSet.TenantID,
		&matchSet.RuleID,
		&matchSet.CreatedBy,
		&matchSet.CreatedAt,
		&matchSet.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrMatchSetNotFound
	}

	if err != nil {
		return nil, err
	}

	return &matchSet, nil
}

// GetMatchSetsByTenant retrieves match sets for a tenant
func (r *PostgresMatchSetRepository) GetMatchSetsByTenant(tenantID string) ([]models.MatchSet, error) {
	query := `
		SELECT id, name, description, tenant_id, rule_id, created_by, created_at, updated_at
		FROM match_sets
		WHERE tenant_id = $1
		ORDER BY name
	`

	rows, err := r.db.Query(query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var matchSets []models.MatchSet
	for rows.Next() {
		var matchSet models.MatchSet
		err := rows.Scan(
			&matchSet.ID,
			&matchSet.Name,
			&matchSet.Description,
			&matchSet.TenantID,
			&matchSet.RuleID,
			&matchSet.CreatedBy,
			&matchSet.CreatedAt,
			&matchSet.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		matchSets = append(matchSets, matchSet)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return matchSets, nil
}

// UpdateMatchSet updates a match set
func (r *PostgresMatchSetRepository) UpdateMatchSet(matchSet *models.MatchSet) error {
	// Check if another match set with the same name already exists for this tenant
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM match_sets WHERE name = $1 AND tenant_id = $2 AND id != $3", matchSet.Name, matchSet.TenantID, matchSet.ID).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		return ErrMatchSetExists
	}

	query := `
		UPDATE match_sets
		SET name = $1, description = $2, rule_id = $3, updated_at = NOW()
		WHERE id = $4
		RETURNING updated_at
	`

	var updatedAt time.Time
	err = r.db.QueryRow(query, matchSet.Name, matchSet.Description, matchSet.RuleID, matchSet.ID).Scan(&updatedAt)

	if err == sql.ErrNoRows {
		return ErrMatchSetNotFound
	}

	if err != nil {
		return err
	}

	matchSet.UpdatedAt = updatedAt
	return nil
}

// DeleteMatchSet deletes a match set
func (r *PostgresMatchSetRepository) DeleteMatchSet(id string) error {
	query := "DELETE FROM match_sets WHERE id = $1"
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrMatchSetNotFound
	}

	return nil
}

// AddDataSourceToMatchSet adds a data source to a match set
func (r *PostgresMatchSetRepository) AddDataSourceToMatchSet(matchSetID, dataSourceID string) error {
	// Check if the association already exists
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM match_set_data_sources WHERE match_set_id = $1 AND data_source_id = $2)", matchSetID, dataSourceID).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		return nil // Already added, no error
	}

	// Add the association
	query := `
		INSERT INTO match_set_data_sources (match_set_id, data_source_id)
		VALUES ($1, $2)
	`
	_, err = r.db.Exec(query, matchSetID, dataSourceID)
	return err
}

// RemoveDataSourceFromMatchSet removes a data source from a match set
func (r *PostgresMatchSetRepository) RemoveDataSourceFromMatchSet(matchSetID, dataSourceID string) error {
	query := `
		DELETE FROM match_set_data_sources
		WHERE match_set_id = $1 AND data_source_id = $2
	`
	result, err := r.db.Exec(query, matchSetID, dataSourceID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("data source not associated with this match set")
	}

	return nil
}

// GetMatchSetDataSources gets all data sources for a match set
func (r *PostgresMatchSetRepository) GetMatchSetDataSources(matchSetID string) ([]models.DataSource, error) {
	query := `
		SELECT ds.id, ds.name, ds.description, ds.tenant_id, ds.created_at, ds.updated_at
		FROM data_sources ds
		JOIN match_set_data_sources msds ON ds.id = msds.data_source_id
		WHERE msds.match_set_id = $1
		ORDER BY ds.name
	`

	rows, err := r.db.Query(query, matchSetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dataSources []models.DataSource
	for rows.Next() {
		var dataSource models.DataSource
		err := rows.Scan(
			&dataSource.ID,
			&dataSource.Name,
			&dataSource.Description,
			&dataSource.TenantID,
			&dataSource.CreatedAt,
			&dataSource.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		dataSources = append(dataSources, dataSource)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return dataSources, nil
}

// MockMatchSetRepository is a mock implementation for development
type MockMatchSetRepository struct {
	matchSets           map[string]*models.MatchSet
	matchSetDataSources map[string][]string // matchSetID -> []dataSourceID
}

// CreateMatchSet creates a match set in the mock repository
func (r *MockMatchSetRepository) CreateMatchSet(matchSet *models.MatchSet) error {
	// Check if match set with the same name already exists for this tenant
	for _, ms := range r.matchSets {
		if ms.Name == matchSet.Name && ms.TenantID == matchSet.TenantID {
			return ErrMatchSetExists
		}
	}

	if matchSet.ID == "" {
		// Generate a dummy ID - in real implementation we'd use UUID
		matchSet.ID = "mock-" + time.Now().Format("20060102150405")
	}
	matchSet.CreatedAt = time.Now()
	matchSet.UpdatedAt = time.Now()

	r.matchSets[matchSet.ID] = matchSet
	return nil
}

// GetMatchSetByID retrieves a match set by ID from the mock repository
func (r *MockMatchSetRepository) GetMatchSetByID(id string) (*models.MatchSet, error) {
	if matchSet, exists := r.matchSets[id]; exists {
		return matchSet, nil
	}
	return nil, ErrMatchSetNotFound
}

// GetMatchSetsByTenant retrieves match sets for a tenant from the mock repository
func (r *MockMatchSetRepository) GetMatchSetsByTenant(tenantID string) ([]models.MatchSet, error) {
	var matchSets []models.MatchSet
	for _, matchSet := range r.matchSets {
		if matchSet.TenantID == tenantID {
			matchSets = append(matchSets, *matchSet)
		}
	}
	return matchSets, nil
}

// UpdateMatchSet updates a match set in the mock repository
func (r *MockMatchSetRepository) UpdateMatchSet(matchSet *models.MatchSet) error {
	existing, exists := r.matchSets[matchSet.ID]
	if !exists {
		return ErrMatchSetNotFound
	}

	// Check if another match set with the same name already exists for this tenant
	for id, ms := range r.matchSets {
		if id != matchSet.ID && ms.Name == matchSet.Name && ms.TenantID == matchSet.TenantID {
			return ErrMatchSetExists
		}
	}

	matchSet.UpdatedAt = time.Now()
	matchSet.CreatedAt = existing.CreatedAt
	matchSet.CreatedBy = existing.CreatedBy
	matchSet.TenantID = existing.TenantID // Tenant ID cannot be changed
	r.matchSets[matchSet.ID] = matchSet

	return nil
}

// DeleteMatchSet deletes a match set from the mock repository
func (r *MockMatchSetRepository) DeleteMatchSet(id string) error {
	if _, exists := r.matchSets[id]; !exists {
		return ErrMatchSetNotFound
	}

	delete(r.matchSets, id)
	delete(r.matchSetDataSources, id)
	return nil
}

// AddDataSourceToMatchSet adds a data source to a match set in the mock repository
func (r *MockMatchSetRepository) AddDataSourceToMatchSet(matchSetID, dataSourceID string) error {
	if _, exists := r.matchSets[matchSetID]; !exists {
		return ErrMatchSetNotFound
	}

	// Initialize data sources array if it doesn't exist
	if r.matchSetDataSources == nil {
		r.matchSetDataSources = make(map[string][]string)
	}

	// Check if data source is already in match set
	sources, exists := r.matchSetDataSources[matchSetID]
	if exists {
		for _, id := range sources {
			if id == dataSourceID {
				return nil // Already exists
			}
		}
	}

	// Add data source to match set
	r.matchSetDataSources[matchSetID] = append(r.matchSetDataSources[matchSetID], dataSourceID)
	return nil
}

// RemoveDataSourceFromMatchSet removes a data source from a match set in the mock repository
func (r *MockMatchSetRepository) RemoveDataSourceFromMatchSet(matchSetID, dataSourceID string) error {
	if _, exists := r.matchSets[matchSetID]; !exists {
		return ErrMatchSetNotFound
	}

	// Check if match set has data sources
	sources, exists := r.matchSetDataSources[matchSetID]
	if !exists {
		return errors.New("data source not associated with this match set")
	}

	// Find and remove the data source
	sourceFound := false
	var newSources []string
	for _, id := range sources {
		if id != dataSourceID {
			newSources = append(newSources, id)
		} else {
			sourceFound = true
		}
	}

	if !sourceFound {
		return errors.New("data source not associated with this match set")
	}

	r.matchSetDataSources[matchSetID] = newSources
	return nil
}

// GetMatchSetDataSources gets all data sources for a match set from the mock repository
func (r *MockMatchSetRepository) GetMatchSetDataSources(matchSetID string) ([]models.DataSource, error) {
	if _, exists := r.matchSets[matchSetID]; !exists {
		return nil, ErrMatchSetNotFound
	}

	var dataSources []models.DataSource
	sourceIDs, exists := r.matchSetDataSources[matchSetID]
	if !exists {
		return dataSources, nil
	}

	// Note: In a real implementation, we would fetch the actual data sources
	// from the data source repository. For the mock, we'll just create dummy objects.
	for _, id := range sourceIDs {
		dataSources = append(dataSources, models.DataSource{
			ID:          id,
			Name:        "Mock Data Source " + id,
			Description: "Mock data source for testing",
			TenantID:    r.matchSets[matchSetID].TenantID,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		})
	}

	return dataSources, nil
}
