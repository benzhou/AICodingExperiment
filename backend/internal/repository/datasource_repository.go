package repository

import (
	"backend/internal/db"
	"backend/internal/models"
	"backend/internal/utils"
	"database/sql"
	"errors"
	"strings"
	"time"
)

var (
	ErrDataSourceNotFound = errors.New("data source not found")
	ErrDataSourceExists   = errors.New("data source with this name already exists")
)

// DataSourceRepository defines operations for managing data sources
type DataSourceRepository interface {
	CreateDataSource(source *models.DataSource) error
	GetDataSourceByID(id string) (*models.DataSource, error)
	GetDataSourceByName(name string) (*models.DataSource, error)
	UpdateDataSource(source *models.DataSource) error
	DeleteDataSource(id string) error
	GetAllDataSources() ([]models.DataSource, error)
	SearchDataSources(query string, limit, offset int) ([]models.DataSource, int, error)
}

// PostgresDataSourceRepository implements DataSourceRepository for PostgreSQL
type PostgresDataSourceRepository struct {
	db *sql.DB
}

// NewDataSourceRepository creates a new data source repository
func NewDataSourceRepository() DataSourceRepository {
	if db.DB == nil {
		// Return a mock repository for development
		return &MockDataSourceRepository{
			dataSources: make(map[string]*models.DataSource),
		}
	}
	return &PostgresDataSourceRepository{
		db: db.DB,
	}
}

// CreateDataSource creates a new data source
func (r *PostgresDataSourceRepository) CreateDataSource(source *models.DataSource) error {
	// Check if data source with the same name already exists
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM data_sources WHERE name = $1)", source.Name).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		return ErrDataSourceExists
	}

	query := `
		INSERT INTO data_sources (name, description)
		VALUES ($1, $2)
		RETURNING id, created_at, updated_at
	`

	err = r.db.QueryRow(
		query,
		source.Name,
		source.Description,
	).Scan(&source.ID, &source.CreatedAt, &source.UpdatedAt)

	if err != nil {
		return err
	}

	// Set epoch timestamps after creation
	source.CreatedAtEpoch = utils.TimeToMillis(source.CreatedAt)
	source.UpdatedAtEpoch = utils.TimeToMillis(source.UpdatedAt)

	return nil
}

// GetDataSourceByID retrieves a data source by ID
func (r *PostgresDataSourceRepository) GetDataSourceByID(id string) (*models.DataSource, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM data_sources
		WHERE id = $1
	`

	var source models.DataSource
	err := r.db.QueryRow(query, id).Scan(
		&source.ID,
		&source.Name,
		&source.Description,
		&source.CreatedAt,
		&source.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrDataSourceNotFound
	}

	if err != nil {
		return nil, err
	}

	// Set epoch timestamps
	source.CreatedAtEpoch = utils.TimeToMillis(source.CreatedAt)
	source.UpdatedAtEpoch = utils.TimeToMillis(source.UpdatedAt)

	return &source, nil
}

// GetDataSourceByName retrieves a data source by name
func (r *PostgresDataSourceRepository) GetDataSourceByName(name string) (*models.DataSource, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM data_sources
		WHERE name = $1
	`

	var source models.DataSource
	err := r.db.QueryRow(query, name).Scan(
		&source.ID,
		&source.Name,
		&source.Description,
		&source.CreatedAt,
		&source.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrDataSourceNotFound
	}

	if err != nil {
		return nil, err
	}

	// Set epoch timestamps
	source.CreatedAtEpoch = utils.TimeToMillis(source.CreatedAt)
	source.UpdatedAtEpoch = utils.TimeToMillis(source.UpdatedAt)

	return &source, nil
}

// UpdateDataSource updates a data source
func (r *PostgresDataSourceRepository) UpdateDataSource(source *models.DataSource) error {
	// Check if another data source with the same name already exists
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM data_sources WHERE name = $1 AND id != $2", source.Name, source.ID).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		return ErrDataSourceExists
	}

	query := `
		UPDATE data_sources
		SET name = $1, description = $2, updated_at = NOW()
		WHERE id = $3
		RETURNING updated_at
	`

	var updatedAt time.Time
	err = r.db.QueryRow(query, source.Name, source.Description, source.ID).Scan(&updatedAt)

	if err == sql.ErrNoRows {
		return ErrDataSourceNotFound
	}

	if err != nil {
		return err
	}

	source.UpdatedAt = updatedAt
	// Set epoch timestamp
	source.UpdatedAtEpoch = utils.TimeToMillis(updatedAt)

	return nil
}

// DeleteDataSource deletes a data source
func (r *PostgresDataSourceRepository) DeleteDataSource(id string) error {
	// Check if there are any existing transactions with this data source
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM transactions WHERE data_source_id = $1", id).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New("cannot delete data source with existing transactions")
	}

	query := "DELETE FROM data_sources WHERE id = $1"
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrDataSourceNotFound
	}

	return nil
}

// GetAllDataSources retrieves all data sources
func (r *PostgresDataSourceRepository) GetAllDataSources() ([]models.DataSource, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM data_sources
		ORDER BY name
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sources []models.DataSource
	for rows.Next() {
		var source models.DataSource
		err := rows.Scan(
			&source.ID,
			&source.Name,
			&source.Description,
			&source.CreatedAt,
			&source.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		// Set epoch timestamps
		source.CreatedAtEpoch = utils.TimeToMillis(source.CreatedAt)
		source.UpdatedAtEpoch = utils.TimeToMillis(source.UpdatedAt)

		sources = append(sources, source)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return sources, nil
}

// SearchDataSources searches for data sources matching the query
func (r *PostgresDataSourceRepository) SearchDataSources(query string, limit, offset int) ([]models.DataSource, int, error) {
	// First get total count
	countQuery := `
		SELECT COUNT(*) 
		FROM data_sources
		WHERE name ILIKE $1 OR description ILIKE $1
	`

	searchPattern := "%" + query + "%"
	var totalCount int
	err := r.db.QueryRow(countQuery, searchPattern).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	// Then get the actual results with pagination
	searchQuery := `
		SELECT id, name, description, created_at, updated_at
		FROM data_sources
		WHERE name ILIKE $1 OR description ILIKE $1
		ORDER BY name
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(searchQuery, searchPattern, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var sources []models.DataSource
	for rows.Next() {
		var source models.DataSource
		err := rows.Scan(
			&source.ID,
			&source.Name,
			&source.Description,
			&source.CreatedAt,
			&source.UpdatedAt,
		)

		if err != nil {
			return nil, 0, err
		}

		// Set epoch timestamps
		source.CreatedAtEpoch = utils.TimeToMillis(source.CreatedAt)
		source.UpdatedAtEpoch = utils.TimeToMillis(source.UpdatedAt)

		sources = append(sources, source)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return sources, totalCount, nil
}

// MockDataSourceRepository is a mock implementation for development
type MockDataSourceRepository struct {
	dataSources map[string]*models.DataSource
	nameIndex   map[string]string // name -> id mapping
}

// CreateDataSource creates a data source in the mock repository
func (r *MockDataSourceRepository) CreateDataSource(source *models.DataSource) error {
	// Initialize nameIndex if it doesn't exist
	if r.nameIndex == nil {
		r.nameIndex = make(map[string]string)
	}

	// Check if data source with the same name exists
	if _, exists := r.nameIndex[source.Name]; exists {
		return ErrDataSourceExists
	}

	if source.ID == "" {
		// Generate a dummy ID - in real implementation we'd use UUID
		source.ID = "mock-" + time.Now().Format("20060102150405")
	}
	source.CreatedAt = time.Now()
	source.UpdatedAt = time.Now()

	// Set epoch timestamps
	source.CreatedAtEpoch = utils.TimeToMillis(source.CreatedAt)
	source.UpdatedAtEpoch = utils.TimeToMillis(source.UpdatedAt)

	r.dataSources[source.ID] = source
	r.nameIndex[source.Name] = source.ID

	return nil
}

// GetDataSourceByID retrieves a data source by ID from the mock repository
func (r *MockDataSourceRepository) GetDataSourceByID(id string) (*models.DataSource, error) {
	if source, exists := r.dataSources[id]; exists {
		return source, nil
	}
	return nil, ErrDataSourceNotFound
}

// GetDataSourceByName retrieves a data source by name from the mock repository
func (r *MockDataSourceRepository) GetDataSourceByName(name string) (*models.DataSource, error) {
	// Initialize nameIndex if it doesn't exist
	if r.nameIndex == nil {
		r.nameIndex = make(map[string]string)
	}

	id, exists := r.nameIndex[name]
	if !exists {
		return nil, ErrDataSourceNotFound
	}

	return r.dataSources[id], nil
}

// UpdateDataSource updates a data source in the mock repository
func (r *MockDataSourceRepository) UpdateDataSource(source *models.DataSource) error {
	// Initialize nameIndex if it doesn't exist
	if r.nameIndex == nil {
		r.nameIndex = make(map[string]string)
	}

	existing, exists := r.dataSources[source.ID]
	if !exists {
		return ErrDataSourceNotFound
	}

	// Check if another data source with the same name exists
	if id, nameExists := r.nameIndex[source.Name]; nameExists && id != source.ID {
		return ErrDataSourceExists
	}

	// Update the name index if the name has changed
	if existing.Name != source.Name {
		delete(r.nameIndex, existing.Name)
		r.nameIndex[source.Name] = source.ID
	}

	source.UpdatedAt = time.Now()
	r.dataSources[source.ID] = source

	return nil
}

// DeleteDataSource deletes a data source from the mock repository
func (r *MockDataSourceRepository) DeleteDataSource(id string) error {
	// Initialize nameIndex if it doesn't exist
	if r.nameIndex == nil {
		r.nameIndex = make(map[string]string)
	}

	source, exists := r.dataSources[id]
	if !exists {
		return ErrDataSourceNotFound
	}

	delete(r.nameIndex, source.Name)
	delete(r.dataSources, id)

	return nil
}

// GetAllDataSources retrieves all data sources from the mock repository
func (r *MockDataSourceRepository) GetAllDataSources() ([]models.DataSource, error) {
	var sources []models.DataSource
	for _, source := range r.dataSources {
		sources = append(sources, *source)
	}
	return sources, nil
}

// SearchDataSources searches for data sources in the mock repository
func (r *MockDataSourceRepository) SearchDataSources(query string, limit, offset int) ([]models.DataSource, int, error) {
	var matchingSources []models.DataSource
	query = strings.ToLower(query)

	// Find all matching sources
	for _, source := range r.dataSources {
		if strings.Contains(strings.ToLower(source.Name), query) ||
			(source.Description != "" && strings.Contains(strings.ToLower(source.Description), query)) {
			matchingSources = append(matchingSources, *source)
		}
	}

	// Get total count
	totalCount := len(matchingSources)

	// Apply pagination
	start := offset
	end := offset + limit
	if start >= totalCount {
		return []models.DataSource{}, totalCount, nil
	}
	if end > totalCount {
		end = totalCount
	}

	return matchingSources[start:end], totalCount, nil
}
