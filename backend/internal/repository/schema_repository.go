package repository

import (
	"backend/internal/db"
	"backend/internal/models"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrSchemaNotFound        = errors.New("schema not found")
	ErrSchemaExists          = errors.New("schema with this name already exists for this tenant")
	ErrSchemaFieldNotFound   = errors.New("schema field not found")
	ErrSchemaFieldExists     = errors.New("schema field with this name already exists")
	ErrSchemaMappingNotFound = errors.New("schema mapping not found")
	ErrSchemaMappingExists   = errors.New("schema mapping for this source field already exists")
	ErrParsingConfigNotFound = errors.New("file parsing configuration not found")
	ErrParsingConfigExists   = errors.New("file parsing configuration for this file type already exists")
)

// SchemaRepository defines operations for managing data source schemas
type SchemaRepository interface {
	// Schema operations
	CreateSchema(schema *models.DataSourceSchema) error
	GetSchemaByID(id string) (*models.DataSourceSchema, error)
	GetSchemasByTenant(tenantID string) ([]models.DataSourceSchema, error)
	UpdateSchema(schema *models.DataSourceSchema) error
	DeleteSchema(id string) error

	// Schema field operations
	AddFieldToSchema(field *models.SchemaField) error
	GetSchemaFields(schemaID string) ([]models.SchemaField, error)
	UpdateSchemaField(field *models.SchemaField) error
	DeleteSchemaField(id string) error

	// Schema mapping operations
	CreateSchemaMapping(mapping *models.SchemaMapping) error
	GetSchemaMappings(schemaID string) ([]models.SchemaMapping, error)
	UpdateSchemaMapping(mapping *models.SchemaMapping) error
	DeleteSchemaMapping(id string) error

	// File parsing configuration operations
	CreateFileParsingConfig(config *models.FileParsingConfig) error
	GetFileParsingConfig(schemaID, fileType string) (*models.FileParsingConfig, error)
	UpdateFileParsingConfig(config *models.FileParsingConfig) error
	DeleteFileParsingConfig(id string) error
}

// PostgresSchemaRepository implements SchemaRepository for PostgreSQL
type PostgresSchemaRepository struct {
	db *sql.DB
}

// NewSchemaRepository creates a new schema repository
func NewSchemaRepository() SchemaRepository {
	if db.DB == nil {
		// Return a mock repository for development
		return &MockSchemaRepository{
			schemas:        make(map[string]*models.DataSourceSchema),
			schemaFields:   make(map[string][]models.SchemaField),
			schemaMappings: make(map[string][]models.SchemaMapping),
			parsingConfigs: make(map[string]map[string]*models.FileParsingConfig), // schemaID -> fileType -> config
		}
	}
	return &PostgresSchemaRepository{
		db: db.DB,
	}
}

// CreateSchema creates a new data source schema
func (r *PostgresSchemaRepository) CreateSchema(schema *models.DataSourceSchema) error {
	// Check if schema with this name already exists for the tenant
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM data_source_schemas WHERE name = $1 AND tenant_id = $2",
		schema.Name, schema.TenantID).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrSchemaExists
	}

	// Create the schema
	query := `
		INSERT INTO data_source_schemas (name, description, tenant_id, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	now := time.Now()
	schema.CreatedAt = now
	schema.UpdatedAt = now

	err = r.db.QueryRow(
		query,
		schema.Name,
		schema.Description,
		schema.TenantID,
		schema.CreatedBy,
		schema.CreatedAt,
		schema.UpdatedAt,
	).Scan(&schema.ID)

	if err != nil {
		return err
	}

	// Create fields if any are provided
	if len(schema.Fields) > 0 {
		for i := range schema.Fields {
			schema.Fields[i].SchemaID = schema.ID
			if err := r.AddFieldToSchema(&schema.Fields[i]); err != nil {
				return err
			}
		}
	}

	return nil
}

// GetSchemaByID retrieves a schema by ID with its fields
func (r *PostgresSchemaRepository) GetSchemaByID(id string) (*models.DataSourceSchema, error) {
	query := `
		SELECT id, name, description, tenant_id, created_by, created_at, updated_at
		FROM data_source_schemas
		WHERE id = $1
	`

	var schema models.DataSourceSchema
	err := r.db.QueryRow(query, id).Scan(
		&schema.ID,
		&schema.Name,
		&schema.Description,
		&schema.TenantID,
		&schema.CreatedBy,
		&schema.CreatedAt,
		&schema.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrSchemaNotFound
	}

	if err != nil {
		return nil, err
	}

	// Get fields for this schema
	fields, err := r.GetSchemaFields(id)
	if err != nil {
		return nil, err
	}
	schema.Fields = fields

	return &schema, nil
}

// GetSchemasByTenant retrieves all schemas for a tenant
func (r *PostgresSchemaRepository) GetSchemasByTenant(tenantID string) ([]models.DataSourceSchema, error) {
	query := `
		SELECT id, name, description, tenant_id, created_by, created_at, updated_at
		FROM data_source_schemas
		WHERE tenant_id = $1
		ORDER BY name
	`

	rows, err := r.db.Query(query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schemas []models.DataSourceSchema
	for rows.Next() {
		var schema models.DataSourceSchema
		err := rows.Scan(
			&schema.ID,
			&schema.Name,
			&schema.Description,
			&schema.TenantID,
			&schema.CreatedBy,
			&schema.CreatedAt,
			&schema.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		// Get fields for this schema
		fields, err := r.GetSchemaFields(schema.ID)
		if err != nil {
			return nil, err
		}
		schema.Fields = fields

		schemas = append(schemas, schema)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return schemas, nil
}

// UpdateSchema updates a schema
func (r *PostgresSchemaRepository) UpdateSchema(schema *models.DataSourceSchema) error {
	// Check if schema exists
	_, err := r.GetSchemaByID(schema.ID)
	if err != nil {
		return err
	}

	// Check if the new name conflicts with another schema
	var count int
	err = r.db.QueryRow(`
		SELECT COUNT(*) FROM data_source_schemas 
		WHERE name = $1 AND tenant_id = $2 AND id != $3
	`, schema.Name, schema.TenantID, schema.ID).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrSchemaExists
	}

	// Update the schema
	query := `
		UPDATE data_source_schemas
		SET name = $1, description = $2, updated_at = $3
		WHERE id = $4
	`
	schema.UpdatedAt = time.Now()

	_, err = r.db.Exec(
		query,
		schema.Name,
		schema.Description,
		schema.UpdatedAt,
		schema.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

// DeleteSchema deletes a schema
func (r *PostgresSchemaRepository) DeleteSchema(id string) error {
	// Check if schema exists
	_, err := r.GetSchemaByID(id)
	if err != nil {
		return err
	}

	// Delete the schema (fields will be cascaded)
	_, err = r.db.Exec("DELETE FROM data_source_schemas WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}

// AddFieldToSchema adds a field to a schema
func (r *PostgresSchemaRepository) AddFieldToSchema(field *models.SchemaField) error {
	// Check if field with this name already exists in the schema
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM schema_fields WHERE schema_id = $1 AND name = $2",
		field.SchemaID, field.Name).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrSchemaFieldExists
	}

	// Add the field
	query := `
		INSERT INTO schema_fields (schema_id, name, display_name, type, required, default_value, "order")
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	err = r.db.QueryRow(
		query,
		field.SchemaID,
		field.Name,
		field.DisplayName,
		field.Type,
		field.Required,
		field.DefaultValue,
		field.Order,
	).Scan(&field.ID)

	if err != nil {
		return err
	}

	return nil
}

// GetSchemaFields retrieves all fields for a schema
func (r *PostgresSchemaRepository) GetSchemaFields(schemaID string) ([]models.SchemaField, error) {
	query := `
		SELECT id, schema_id, name, display_name, type, required, default_value, "order"
		FROM schema_fields
		WHERE schema_id = $1
		ORDER BY "order"
	`

	rows, err := r.db.Query(query, schemaID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var fields []models.SchemaField
	for rows.Next() {
		var field models.SchemaField
		err := rows.Scan(
			&field.ID,
			&field.SchemaID,
			&field.Name,
			&field.DisplayName,
			&field.Type,
			&field.Required,
			&field.DefaultValue,
			&field.Order,
		)

		if err != nil {
			return nil, err
		}

		fields = append(fields, field)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return fields, nil
}

// UpdateSchemaField updates a schema field
func (r *PostgresSchemaRepository) UpdateSchemaField(field *models.SchemaField) error {
	// Check if field exists
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM schema_fields WHERE id = $1)", field.ID).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return ErrSchemaFieldNotFound
	}

	// Check if the new name conflicts with another field
	var count int
	err = r.db.QueryRow(`
		SELECT COUNT(*) FROM schema_fields 
		WHERE schema_id = $1 AND name = $2 AND id != $3
	`, field.SchemaID, field.Name, field.ID).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrSchemaFieldExists
	}

	// Update the field
	query := `
		UPDATE schema_fields
		SET name = $1, display_name = $2, type = $3, required = $4, default_value = $5, "order" = $6
		WHERE id = $7
	`

	_, err = r.db.Exec(
		query,
		field.Name,
		field.DisplayName,
		field.Type,
		field.Required,
		field.DefaultValue,
		field.Order,
		field.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

// DeleteSchemaField deletes a schema field
func (r *PostgresSchemaRepository) DeleteSchemaField(id string) error {
	// Check if field exists
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM schema_fields WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return ErrSchemaFieldNotFound
	}

	// Delete the field
	_, err = r.db.Exec("DELETE FROM schema_fields WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}

// CreateSchemaMapping creates a new schema mapping
func (r *PostgresSchemaRepository) CreateSchemaMapping(mapping *models.SchemaMapping) error {
	// Check if mapping for this source field already exists
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM schema_mappings WHERE schema_id = $1 AND source_field_name = $2",
		mapping.SchemaID, mapping.SourceFieldName).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrSchemaMappingExists
	}

	// Create the mapping
	query := `
		INSERT INTO schema_mappings (schema_id, source_field_name, target_field_name, transformation, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	now := time.Now()
	mapping.CreatedAt = now
	mapping.UpdatedAt = now

	err = r.db.QueryRow(
		query,
		mapping.SchemaID,
		mapping.SourceFieldName,
		mapping.TargetFieldName,
		mapping.Transformation,
		mapping.CreatedAt,
		mapping.UpdatedAt,
	).Scan(&mapping.ID)

	if err != nil {
		return err
	}

	return nil
}

// GetSchemaMappings retrieves all mappings for a schema
func (r *PostgresSchemaRepository) GetSchemaMappings(schemaID string) ([]models.SchemaMapping, error) {
	query := `
		SELECT id, schema_id, source_field_name, target_field_name, transformation, created_at, updated_at
		FROM schema_mappings
		WHERE schema_id = $1
	`

	rows, err := r.db.Query(query, schemaID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var mappings []models.SchemaMapping
	for rows.Next() {
		var mapping models.SchemaMapping
		err := rows.Scan(
			&mapping.ID,
			&mapping.SchemaID,
			&mapping.SourceFieldName,
			&mapping.TargetFieldName,
			&mapping.Transformation,
			&mapping.CreatedAt,
			&mapping.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		mappings = append(mappings, mapping)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return mappings, nil
}

// UpdateSchemaMapping updates a schema mapping
func (r *PostgresSchemaRepository) UpdateSchemaMapping(mapping *models.SchemaMapping) error {
	// Check if mapping exists
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM schema_mappings WHERE id = $1)", mapping.ID).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return ErrSchemaMappingNotFound
	}

	// Check if the new source field conflicts with another mapping
	var count int
	err = r.db.QueryRow(`
		SELECT COUNT(*) FROM schema_mappings 
		WHERE schema_id = $1 AND source_field_name = $2 AND id != $3
	`, mapping.SchemaID, mapping.SourceFieldName, mapping.ID).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrSchemaMappingExists
	}

	// Update the mapping
	query := `
		UPDATE schema_mappings
		SET source_field_name = $1, target_field_name = $2, transformation = $3, updated_at = $4
		WHERE id = $5
	`
	mapping.UpdatedAt = time.Now()

	_, err = r.db.Exec(
		query,
		mapping.SourceFieldName,
		mapping.TargetFieldName,
		mapping.Transformation,
		mapping.UpdatedAt,
		mapping.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

// DeleteSchemaMapping deletes a schema mapping
func (r *PostgresSchemaRepository) DeleteSchemaMapping(id string) error {
	// Check if mapping exists
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM schema_mappings WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return ErrSchemaMappingNotFound
	}

	// Delete the mapping
	_, err = r.db.Exec("DELETE FROM schema_mappings WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}

// CreateFileParsingConfig creates a new file parsing configuration
func (r *PostgresSchemaRepository) CreateFileParsingConfig(config *models.FileParsingConfig) error {
	// Check if config for this file type already exists
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM file_parsing_configs WHERE schema_id = $1 AND file_type = $2",
		config.SchemaID, config.FileType).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrParsingConfigExists
	}

	// Create the config
	query := `
		INSERT INTO file_parsing_configs (schema_id, file_type, has_header_row, delimiter, date_format, time_format, number_format, encapsulated_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id
	`
	now := time.Now()
	config.CreatedAt = now
	config.UpdatedAt = now

	err = r.db.QueryRow(
		query,
		config.SchemaID,
		config.FileType,
		config.HasHeaderRow,
		config.Delimiter,
		config.DateFormat,
		config.TimeFormat,
		config.NumberFormat,
		config.EncapsulatedBy,
		config.CreatedAt,
		config.UpdatedAt,
	).Scan(&config.ID)

	if err != nil {
		return err
	}

	return nil
}

// GetFileParsingConfig retrieves a file parsing configuration
func (r *PostgresSchemaRepository) GetFileParsingConfig(schemaID, fileType string) (*models.FileParsingConfig, error) {
	query := `
		SELECT id, schema_id, file_type, has_header_row, delimiter, date_format, time_format, number_format, encapsulated_by, created_at, updated_at
		FROM file_parsing_configs
		WHERE schema_id = $1 AND file_type = $2
	`

	var config models.FileParsingConfig
	err := r.db.QueryRow(query, schemaID, fileType).Scan(
		&config.ID,
		&config.SchemaID,
		&config.FileType,
		&config.HasHeaderRow,
		&config.Delimiter,
		&config.DateFormat,
		&config.TimeFormat,
		&config.NumberFormat,
		&config.EncapsulatedBy,
		&config.CreatedAt,
		&config.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrParsingConfigNotFound
	}

	if err != nil {
		return nil, err
	}

	return &config, nil
}

// UpdateFileParsingConfig updates a file parsing configuration
func (r *PostgresSchemaRepository) UpdateFileParsingConfig(config *models.FileParsingConfig) error {
	// Check if config exists
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM file_parsing_configs WHERE id = $1)", config.ID).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return ErrParsingConfigNotFound
	}

	// Check if the new file type conflicts with another config
	var count int
	err = r.db.QueryRow(`
		SELECT COUNT(*) FROM file_parsing_configs 
		WHERE schema_id = $1 AND file_type = $2 AND id != $3
	`, config.SchemaID, config.FileType, config.ID).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrParsingConfigExists
	}

	// Update the config
	query := `
		UPDATE file_parsing_configs
		SET file_type = $1, has_header_row = $2, delimiter = $3, date_format = $4, time_format = $5,
			number_format = $6, encapsulated_by = $7, updated_at = $8
		WHERE id = $9
	`
	config.UpdatedAt = time.Now()

	_, err = r.db.Exec(
		query,
		config.FileType,
		config.HasHeaderRow,
		config.Delimiter,
		config.DateFormat,
		config.TimeFormat,
		config.NumberFormat,
		config.EncapsulatedBy,
		config.UpdatedAt,
		config.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

// DeleteFileParsingConfig deletes a file parsing configuration
func (r *PostgresSchemaRepository) DeleteFileParsingConfig(id string) error {
	// Check if config exists
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM file_parsing_configs WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return ErrParsingConfigNotFound
	}

	// Delete the config
	_, err = r.db.Exec("DELETE FROM file_parsing_configs WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}

// MockSchemaRepository is a mock implementation for development
type MockSchemaRepository struct {
	schemas        map[string]*models.DataSourceSchema             // ID -> Schema
	schemaFields   map[string][]models.SchemaField                 // SchemaID -> []Field
	schemaMappings map[string][]models.SchemaMapping               // SchemaID -> []Mapping
	parsingConfigs map[string]map[string]*models.FileParsingConfig // SchemaID -> FileType -> Config
}
