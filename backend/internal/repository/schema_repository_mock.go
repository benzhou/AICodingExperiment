package repository

import (
	"backend/internal/models"
	"time"
)

// CreateSchema creates a schema in the mock repository
func (r *MockSchemaRepository) CreateSchema(schema *models.DataSourceSchema) error {
	// Check for existing schema with same name in tenant
	for _, s := range r.schemas {
		if s.Name == schema.Name && s.TenantID == schema.TenantID {
			return ErrSchemaExists
		}
	}

	// Generate an ID if one isn't provided
	if schema.ID == "" {
		schema.ID = "schema-" + time.Now().Format("20060102150405")
	}

	// Set timestamps
	if schema.CreatedAt.IsZero() {
		schema.CreatedAt = time.Now()
	}
	if schema.UpdatedAt.IsZero() {
		schema.UpdatedAt = time.Now()
	}

	// Store the schema
	r.schemas[schema.ID] = schema

	// Store the fields
	if len(schema.Fields) > 0 {
		r.schemaFields[schema.ID] = make([]models.SchemaField, 0)
		for i := range schema.Fields {
			schema.Fields[i].SchemaID = schema.ID
			if schema.Fields[i].ID == "" {
				schema.Fields[i].ID = "field-" + time.Now().Format("20060102150405") + "-" + schema.Fields[i].Name
			}
			r.schemaFields[schema.ID] = append(r.schemaFields[schema.ID], schema.Fields[i])
		}
	}

	return nil
}

// GetSchemaByID retrieves a schema by ID from the mock repository
func (r *MockSchemaRepository) GetSchemaByID(id string) (*models.DataSourceSchema, error) {
	schema, exists := r.schemas[id]
	if !exists {
		return nil, ErrSchemaNotFound
	}

	// Get fields for this schema
	fields, err := r.GetSchemaFields(id)
	if err != nil {
		schema.Fields = []models.SchemaField{}
	} else {
		schema.Fields = fields
	}

	return schema, nil
}

// GetSchemasByTenant retrieves schemas for a tenant from the mock repository
func (r *MockSchemaRepository) GetSchemasByTenant(tenantID string) ([]models.DataSourceSchema, error) {
	var schemas []models.DataSourceSchema

	for _, schema := range r.schemas {
		if schema.TenantID == tenantID {
			// Get fields for this schema
			fields, err := r.GetSchemaFields(schema.ID)
			if err != nil {
				schema.Fields = []models.SchemaField{}
			} else {
				schema.Fields = fields
			}

			schemas = append(schemas, *schema)
		}
	}

	return schemas, nil
}

// UpdateSchema updates a schema in the mock repository
func (r *MockSchemaRepository) UpdateSchema(schema *models.DataSourceSchema) error {
	// Check if schema exists
	_, exists := r.schemas[schema.ID]
	if !exists {
		return ErrSchemaNotFound
	}

	// Check for name conflict
	for id, s := range r.schemas {
		if s.Name == schema.Name && s.TenantID == schema.TenantID && id != schema.ID {
			return ErrSchemaExists
		}
	}

	// Preserve created info
	createdAt := r.schemas[schema.ID].CreatedAt
	createdBy := r.schemas[schema.ID].CreatedBy

	// Update schema
	schema.CreatedAt = createdAt
	schema.CreatedBy = createdBy
	schema.UpdatedAt = time.Now()

	r.schemas[schema.ID] = schema

	return nil
}

// DeleteSchema deletes a schema from the mock repository
func (r *MockSchemaRepository) DeleteSchema(id string) error {
	// Check if schema exists
	_, exists := r.schemas[id]
	if !exists {
		return ErrSchemaNotFound
	}

	// Delete schema
	delete(r.schemas, id)

	// Delete related data
	delete(r.schemaFields, id)
	delete(r.schemaMappings, id)
	delete(r.parsingConfigs, id)

	return nil
}

// AddFieldToSchema adds a field to a schema in the mock repository
func (r *MockSchemaRepository) AddFieldToSchema(field *models.SchemaField) error {
	// Check if schema exists
	_, exists := r.schemas[field.SchemaID]
	if !exists {
		return ErrSchemaNotFound
	}

	// Initialize fields array if needed
	if _, exists := r.schemaFields[field.SchemaID]; !exists {
		r.schemaFields[field.SchemaID] = make([]models.SchemaField, 0)
	}

	// Check for field name conflict
	for _, f := range r.schemaFields[field.SchemaID] {
		if f.Name == field.Name {
			return ErrSchemaFieldExists
		}
	}

	// Generate an ID if one isn't provided
	if field.ID == "" {
		field.ID = "field-" + time.Now().Format("20060102150405") + "-" + field.Name
	}

	// Add the field
	r.schemaFields[field.SchemaID] = append(r.schemaFields[field.SchemaID], *field)

	return nil
}

// GetSchemaFields retrieves fields for a schema from the mock repository
func (r *MockSchemaRepository) GetSchemaFields(schemaID string) ([]models.SchemaField, error) {
	// Check if schema exists
	_, exists := r.schemas[schemaID]
	if !exists {
		return nil, ErrSchemaNotFound
	}

	// Return fields (or empty array if none)
	fields, exists := r.schemaFields[schemaID]
	if !exists {
		return []models.SchemaField{}, nil
	}

	return fields, nil
}

// UpdateSchemaField updates a schema field in the mock repository
func (r *MockSchemaRepository) UpdateSchemaField(field *models.SchemaField) error {
	// Check if schema exists
	_, exists := r.schemas[field.SchemaID]
	if !exists {
		return ErrSchemaNotFound
	}

	// Check if fields exist for this schema
	fields, exists := r.schemaFields[field.SchemaID]
	if !exists {
		return ErrSchemaFieldNotFound
	}

	// Find the field to update
	fieldIndex := -1
	for i, f := range fields {
		if f.ID == field.ID {
			fieldIndex = i
			break
		}
	}

	if fieldIndex == -1 {
		return ErrSchemaFieldNotFound
	}

	// Check for name conflict
	for i, f := range fields {
		if f.Name == field.Name && i != fieldIndex {
			return ErrSchemaFieldExists
		}
	}

	// Update the field
	fields[fieldIndex] = *field
	r.schemaFields[field.SchemaID] = fields

	return nil
}

// DeleteSchemaField deletes a schema field from the mock repository
func (r *MockSchemaRepository) DeleteSchemaField(id string) error {
	// Find the field
	for schemaID, fields := range r.schemaFields {
		for i, field := range fields {
			if field.ID == id {
				// Remove the field
				r.schemaFields[schemaID] = append(fields[:i], fields[i+1:]...)
				return nil
			}
		}
	}

	return ErrSchemaFieldNotFound
}

// CreateSchemaMapping creates a schema mapping in the mock repository
func (r *MockSchemaRepository) CreateSchemaMapping(mapping *models.SchemaMapping) error {
	// Check if schema exists
	_, exists := r.schemas[mapping.SchemaID]
	if !exists {
		return ErrSchemaNotFound
	}

	// Initialize mappings array if needed
	if _, exists := r.schemaMappings[mapping.SchemaID]; !exists {
		r.schemaMappings[mapping.SchemaID] = make([]models.SchemaMapping, 0)
	}

	// Check for source field conflict
	for _, m := range r.schemaMappings[mapping.SchemaID] {
		if m.SourceFieldName == mapping.SourceFieldName {
			return ErrSchemaMappingExists
		}
	}

	// Generate an ID if one isn't provided
	if mapping.ID == "" {
		mapping.ID = "mapping-" + time.Now().Format("20060102150405")
	}

	// Set timestamps
	if mapping.CreatedAt.IsZero() {
		mapping.CreatedAt = time.Now()
	}
	if mapping.UpdatedAt.IsZero() {
		mapping.UpdatedAt = time.Now()
	}

	// Add the mapping
	r.schemaMappings[mapping.SchemaID] = append(r.schemaMappings[mapping.SchemaID], *mapping)

	return nil
}

// GetSchemaMappings retrieves mappings for a schema from the mock repository
func (r *MockSchemaRepository) GetSchemaMappings(schemaID string) ([]models.SchemaMapping, error) {
	// Check if schema exists
	_, exists := r.schemas[schemaID]
	if !exists {
		return nil, ErrSchemaNotFound
	}

	// Return mappings (or empty array if none)
	mappings, exists := r.schemaMappings[schemaID]
	if !exists {
		return []models.SchemaMapping{}, nil
	}

	return mappings, nil
}

// UpdateSchemaMapping updates a schema mapping in the mock repository
func (r *MockSchemaRepository) UpdateSchemaMapping(mapping *models.SchemaMapping) error {
	// Check if schema exists
	_, exists := r.schemas[mapping.SchemaID]
	if !exists {
		return ErrSchemaNotFound
	}

	// Check if mappings exist for this schema
	mappings, exists := r.schemaMappings[mapping.SchemaID]
	if !exists {
		return ErrSchemaMappingNotFound
	}

	// Find the mapping to update
	mappingIndex := -1
	for i, m := range mappings {
		if m.ID == mapping.ID {
			mappingIndex = i
			break
		}
	}

	if mappingIndex == -1 {
		return ErrSchemaMappingNotFound
	}

	// Check for source field conflict
	for i, m := range mappings {
		if m.SourceFieldName == mapping.SourceFieldName && i != mappingIndex {
			return ErrSchemaMappingExists
		}
	}

	// Update timestamps
	mapping.UpdatedAt = time.Now()
	if mapping.CreatedAt.IsZero() {
		mapping.CreatedAt = mappings[mappingIndex].CreatedAt
	}

	// Update the mapping
	mappings[mappingIndex] = *mapping
	r.schemaMappings[mapping.SchemaID] = mappings

	return nil
}

// DeleteSchemaMapping deletes a schema mapping from the mock repository
func (r *MockSchemaRepository) DeleteSchemaMapping(id string) error {
	// Find the mapping
	for schemaID, mappings := range r.schemaMappings {
		for i, mapping := range mappings {
			if mapping.ID == id {
				// Remove the mapping
				r.schemaMappings[schemaID] = append(mappings[:i], mappings[i+1:]...)
				return nil
			}
		}
	}

	return ErrSchemaMappingNotFound
}

// CreateFileParsingConfig creates a file parsing configuration in the mock repository
func (r *MockSchemaRepository) CreateFileParsingConfig(config *models.FileParsingConfig) error {
	// Check if schema exists
	_, exists := r.schemas[config.SchemaID]
	if !exists {
		return ErrSchemaNotFound
	}

	// Initialize configs map if needed
	if _, exists := r.parsingConfigs[config.SchemaID]; !exists {
		r.parsingConfigs[config.SchemaID] = make(map[string]*models.FileParsingConfig)
	}

	// Check for file type conflict
	if _, exists := r.parsingConfigs[config.SchemaID][config.FileType]; exists {
		return ErrParsingConfigExists
	}

	// Generate an ID if one isn't provided
	if config.ID == "" {
		config.ID = "config-" + time.Now().Format("20060102150405")
	}

	// Set timestamps
	if config.CreatedAt.IsZero() {
		config.CreatedAt = time.Now()
	}
	if config.UpdatedAt.IsZero() {
		config.UpdatedAt = time.Now()
	}

	// Store the config
	r.parsingConfigs[config.SchemaID][config.FileType] = config

	return nil
}

// GetFileParsingConfig retrieves a file parsing configuration from the mock repository
func (r *MockSchemaRepository) GetFileParsingConfig(schemaID, fileType string) (*models.FileParsingConfig, error) {
	// Check if schema exists
	_, exists := r.schemas[schemaID]
	if !exists {
		return nil, ErrSchemaNotFound
	}

	// Check if configs exist for this schema
	configs, exists := r.parsingConfigs[schemaID]
	if !exists {
		return nil, ErrParsingConfigNotFound
	}

	// Check if config exists for this file type
	config, exists := configs[fileType]
	if !exists {
		return nil, ErrParsingConfigNotFound
	}

	return config, nil
}

// UpdateFileParsingConfig updates a file parsing configuration in the mock repository
func (r *MockSchemaRepository) UpdateFileParsingConfig(config *models.FileParsingConfig) error {
	// Check if schema exists
	_, exists := r.schemas[config.SchemaID]
	if !exists {
		return ErrSchemaNotFound
	}

	// Check if configs exist for this schema
	configs, exists := r.parsingConfigs[config.SchemaID]
	if !exists {
		return ErrParsingConfigNotFound
	}

	// Find the existing config
	var oldConfig *models.FileParsingConfig
	var oldFileType string
	for ft, c := range configs {
		if c.ID == config.ID {
			oldConfig = c
			oldFileType = ft
			break
		}
	}

	if oldConfig == nil {
		return ErrParsingConfigNotFound
	}

	// Check for file type conflict
	if oldFileType != config.FileType {
		if _, exists := configs[config.FileType]; exists {
			return ErrParsingConfigExists
		}
	}

	// Update timestamps
	config.UpdatedAt = time.Now()
	if config.CreatedAt.IsZero() {
		config.CreatedAt = oldConfig.CreatedAt
	}

	// Remove old entry if file type changed
	if oldFileType != config.FileType {
		delete(configs, oldFileType)
	}

	// Store updated config
	configs[config.FileType] = config
	r.parsingConfigs[config.SchemaID] = configs

	return nil
}

// DeleteFileParsingConfig deletes a file parsing configuration from the mock repository
func (r *MockSchemaRepository) DeleteFileParsingConfig(id string) error {
	// Find the config
	for schemaID, configs := range r.parsingConfigs {
		for fileType, config := range configs {
			if config.ID == id {
				// Remove the config
				delete(r.parsingConfigs[schemaID], fileType)
				return nil
			}
		}
	}

	return ErrParsingConfigNotFound
}
