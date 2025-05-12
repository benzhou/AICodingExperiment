package services

import (
	"backend/internal/models"
	"backend/internal/repository"
	"errors"
)

// SchemaService provides methods for schema operations
type SchemaService struct {
	schemaRepo     repository.SchemaRepository
	permissionRepo repository.PermissionRepository
}

// NewSchemaService creates a new schema service
func NewSchemaService(
	schemaRepo repository.SchemaRepository,
	permissionRepo repository.PermissionRepository,
) *SchemaService {
	return &SchemaService{
		schemaRepo:     schemaRepo,
		permissionRepo: permissionRepo,
	}
}

// CreateSchema creates a new data source schema
func (s *SchemaService) CreateSchema(schema *models.DataSourceSchema, userID string) (*models.DataSourceSchema, error) {
	// Check if user has permission to create schemas
	hasPermission, err := s.permissionRepo.HasPermission(userID, models.PermManageSchemas, schema.TenantID)
	if err != nil {
		return nil, err
	}
	if !hasPermission {
		return nil, errors.New("unauthorized: requires manage schemas permission")
	}

	// Set created by
	schema.CreatedBy = userID

	// Create the schema
	if err := s.schemaRepo.CreateSchema(schema); err != nil {
		return nil, err
	}

	return schema, nil
}

// GetSchemaByID retrieves a schema by ID
func (s *SchemaService) GetSchemaByID(id, userID, tenantID string) (*models.DataSourceSchema, error) {
	// Check if user has permission to view schemas
	hasPermission, err := s.permissionRepo.HasPermission(userID, models.PermViewSchemas, tenantID)
	if err != nil {
		return nil, err
	}
	if !hasPermission {
		return nil, errors.New("unauthorized: requires view schemas permission")
	}

	// Get the schema
	schema, err := s.schemaRepo.GetSchemaByID(id)
	if err != nil {
		return nil, err
	}

	// Ensure the schema belongs to the tenant
	if schema.TenantID != tenantID {
		return nil, errors.New("schema not found in this tenant")
	}

	return schema, nil
}

// GetSchemasByTenant retrieves schemas for a tenant
func (s *SchemaService) GetSchemasByTenant(tenantID, userID string) ([]models.DataSourceSchema, error) {
	// Check if user has permission to view schemas
	hasPermission, err := s.permissionRepo.HasPermission(userID, models.PermViewSchemas, tenantID)
	if err != nil {
		return nil, err
	}
	if !hasPermission {
		return nil, errors.New("unauthorized: requires view schemas permission")
	}

	// Get schemas
	return s.schemaRepo.GetSchemasByTenant(tenantID)
}

// UpdateSchema updates a schema
func (s *SchemaService) UpdateSchema(schema *models.DataSourceSchema, userID string) (*models.DataSourceSchema, error) {
	// Check if user has permission to manage schemas
	hasPermission, err := s.permissionRepo.HasPermission(userID, models.PermManageSchemas, schema.TenantID)
	if err != nil {
		return nil, err
	}
	if !hasPermission {
		return nil, errors.New("unauthorized: requires manage schemas permission")
	}

	// Get the existing schema
	existingSchema, err := s.schemaRepo.GetSchemaByID(schema.ID)
	if err != nil {
		return nil, err
	}

	// Ensure the schema belongs to the tenant
	if existingSchema.TenantID != schema.TenantID {
		return nil, errors.New("schema not found in this tenant")
	}

	// Update the schema
	if err := s.schemaRepo.UpdateSchema(schema); err != nil {
		return nil, err
	}

	// Re-fetch the updated schema with fields
	return s.schemaRepo.GetSchemaByID(schema.ID)
}

// DeleteSchema deletes a schema
func (s *SchemaService) DeleteSchema(id, userID, tenantID string) error {
	// Check if user has permission to manage schemas
	hasPermission, err := s.permissionRepo.HasPermission(userID, models.PermManageSchemas, tenantID)
	if err != nil {
		return err
	}
	if !hasPermission {
		return errors.New("unauthorized: requires manage schemas permission")
	}

	// Get the schema
	schema, err := s.schemaRepo.GetSchemaByID(id)
	if err != nil {
		return err
	}

	// Ensure the schema belongs to the tenant
	if schema.TenantID != tenantID {
		return errors.New("schema not found in this tenant")
	}

	// Delete the schema
	return s.schemaRepo.DeleteSchema(id)
}

// AddFieldToSchema adds a field to a schema
func (s *SchemaService) AddFieldToSchema(field *models.SchemaField, userID, tenantID string) (*models.SchemaField, error) {
	// Check if user has permission to manage schemas
	hasPermission, err := s.permissionRepo.HasPermission(userID, models.PermManageSchemas, tenantID)
	if err != nil {
		return nil, err
	}
	if !hasPermission {
		return nil, errors.New("unauthorized: requires manage schemas permission")
	}

	// Get the schema
	schema, err := s.schemaRepo.GetSchemaByID(field.SchemaID)
	if err != nil {
		return nil, err
	}

	// Ensure the schema belongs to the tenant
	if schema.TenantID != tenantID {
		return nil, errors.New("schema not found in this tenant")
	}

	// Add the field
	if err := s.schemaRepo.AddFieldToSchema(field); err != nil {
		return nil, err
	}

	return field, nil
}

// UpdateSchemaField updates a schema field
func (s *SchemaService) UpdateSchemaField(field *models.SchemaField, userID, tenantID string) (*models.SchemaField, error) {
	// Check if user has permission to manage schemas
	hasPermission, err := s.permissionRepo.HasPermission(userID, models.PermManageSchemas, tenantID)
	if err != nil {
		return nil, err
	}
	if !hasPermission {
		return nil, errors.New("unauthorized: requires manage schemas permission")
	}

	// Get the schema
	schema, err := s.schemaRepo.GetSchemaByID(field.SchemaID)
	if err != nil {
		return nil, err
	}

	// Ensure the schema belongs to the tenant
	if schema.TenantID != tenantID {
		return nil, errors.New("schema not found in this tenant")
	}

	// Update the field
	if err := s.schemaRepo.UpdateSchemaField(field); err != nil {
		return nil, err
	}

	return field, nil
}

// DeleteSchemaField deletes a schema field
func (s *SchemaService) DeleteSchemaField(fieldID, schemaID, userID, tenantID string) error {
	// Check if user has permission to manage schemas
	hasPermission, err := s.permissionRepo.HasPermission(userID, models.PermManageSchemas, tenantID)
	if err != nil {
		return err
	}
	if !hasPermission {
		return errors.New("unauthorized: requires manage schemas permission")
	}

	// Get the schema
	schema, err := s.schemaRepo.GetSchemaByID(schemaID)
	if err != nil {
		return err
	}

	// Ensure the schema belongs to the tenant
	if schema.TenantID != tenantID {
		return errors.New("schema not found in this tenant")
	}

	// Delete the field
	return s.schemaRepo.DeleteSchemaField(fieldID)
}

// CreateSchemaMapping creates a schema mapping
func (s *SchemaService) CreateSchemaMapping(mapping *models.SchemaMapping, userID, tenantID string) (*models.SchemaMapping, error) {
	// Check if user has permission to manage schemas
	hasPermission, err := s.permissionRepo.HasPermission(userID, models.PermManageSchemas, tenantID)
	if err != nil {
		return nil, err
	}
	if !hasPermission {
		return nil, errors.New("unauthorized: requires manage schemas permission")
	}

	// Get the schema
	schema, err := s.schemaRepo.GetSchemaByID(mapping.SchemaID)
	if err != nil {
		return nil, err
	}

	// Ensure the schema belongs to the tenant
	if schema.TenantID != tenantID {
		return nil, errors.New("schema not found in this tenant")
	}

	// Create the mapping
	if err := s.schemaRepo.CreateSchemaMapping(mapping); err != nil {
		return nil, err
	}

	return mapping, nil
}

// GetSchemaMappings retrieves mappings for a schema
func (s *SchemaService) GetSchemaMappings(schemaID, userID, tenantID string) ([]models.SchemaMapping, error) {
	// Check if user has permission to view schemas
	hasPermission, err := s.permissionRepo.HasPermission(userID, models.PermViewSchemas, tenantID)
	if err != nil {
		return nil, err
	}
	if !hasPermission {
		return nil, errors.New("unauthorized: requires view schemas permission")
	}

	// Get the schema
	schema, err := s.schemaRepo.GetSchemaByID(schemaID)
	if err != nil {
		return nil, err
	}

	// Ensure the schema belongs to the tenant
	if schema.TenantID != tenantID {
		return nil, errors.New("schema not found in this tenant")
	}

	// Get mappings
	return s.schemaRepo.GetSchemaMappings(schemaID)
}

// UpdateSchemaMapping updates a schema mapping
func (s *SchemaService) UpdateSchemaMapping(mapping *models.SchemaMapping, userID, tenantID string) (*models.SchemaMapping, error) {
	// Check if user has permission to manage schemas
	hasPermission, err := s.permissionRepo.HasPermission(userID, models.PermManageSchemas, tenantID)
	if err != nil {
		return nil, err
	}
	if !hasPermission {
		return nil, errors.New("unauthorized: requires manage schemas permission")
	}

	// Get the schema
	schema, err := s.schemaRepo.GetSchemaByID(mapping.SchemaID)
	if err != nil {
		return nil, err
	}

	// Ensure the schema belongs to the tenant
	if schema.TenantID != tenantID {
		return nil, errors.New("schema not found in this tenant")
	}

	// Update the mapping
	if err := s.schemaRepo.UpdateSchemaMapping(mapping); err != nil {
		return nil, err
	}

	return mapping, nil
}

// DeleteSchemaMapping deletes a schema mapping
func (s *SchemaService) DeleteSchemaMapping(mappingID, schemaID, userID, tenantID string) error {
	// Check if user has permission to manage schemas
	hasPermission, err := s.permissionRepo.HasPermission(userID, models.PermManageSchemas, tenantID)
	if err != nil {
		return err
	}
	if !hasPermission {
		return errors.New("unauthorized: requires manage schemas permission")
	}

	// Get the schema
	schema, err := s.schemaRepo.GetSchemaByID(schemaID)
	if err != nil {
		return err
	}

	// Ensure the schema belongs to the tenant
	if schema.TenantID != tenantID {
		return errors.New("schema not found in this tenant")
	}

	// Delete the mapping
	return s.schemaRepo.DeleteSchemaMapping(mappingID)
}

// CreateFileParsingConfig creates a file parsing configuration
func (s *SchemaService) CreateFileParsingConfig(config *models.FileParsingConfig, userID, tenantID string) (*models.FileParsingConfig, error) {
	// Check if user has permission to manage schemas
	hasPermission, err := s.permissionRepo.HasPermission(userID, models.PermManageSchemas, tenantID)
	if err != nil {
		return nil, err
	}
	if !hasPermission {
		return nil, errors.New("unauthorized: requires manage schemas permission")
	}

	// Get the schema
	schema, err := s.schemaRepo.GetSchemaByID(config.SchemaID)
	if err != nil {
		return nil, err
	}

	// Ensure the schema belongs to the tenant
	if schema.TenantID != tenantID {
		return nil, errors.New("schema not found in this tenant")
	}

	// Create the config
	if err := s.schemaRepo.CreateFileParsingConfig(config); err != nil {
		return nil, err
	}

	return config, nil
}

// GetFileParsingConfig retrieves a file parsing configuration
func (s *SchemaService) GetFileParsingConfig(schemaID, fileType, userID, tenantID string) (*models.FileParsingConfig, error) {
	// Check if user has permission to view schemas
	hasPermission, err := s.permissionRepo.HasPermission(userID, models.PermViewSchemas, tenantID)
	if err != nil {
		return nil, err
	}
	if !hasPermission {
		return nil, errors.New("unauthorized: requires view schemas permission")
	}

	// Get the schema
	schema, err := s.schemaRepo.GetSchemaByID(schemaID)
	if err != nil {
		return nil, err
	}

	// Ensure the schema belongs to the tenant
	if schema.TenantID != tenantID {
		return nil, errors.New("schema not found in this tenant")
	}

	// Get the config
	return s.schemaRepo.GetFileParsingConfig(schemaID, fileType)
}

// UpdateFileParsingConfig updates a file parsing configuration
func (s *SchemaService) UpdateFileParsingConfig(config *models.FileParsingConfig, userID, tenantID string) (*models.FileParsingConfig, error) {
	// Check if user has permission to manage schemas
	hasPermission, err := s.permissionRepo.HasPermission(userID, models.PermManageSchemas, tenantID)
	if err != nil {
		return nil, err
	}
	if !hasPermission {
		return nil, errors.New("unauthorized: requires manage schemas permission")
	}

	// Get the schema
	schema, err := s.schemaRepo.GetSchemaByID(config.SchemaID)
	if err != nil {
		return nil, err
	}

	// Ensure the schema belongs to the tenant
	if schema.TenantID != tenantID {
		return nil, errors.New("schema not found in this tenant")
	}

	// Update the config
	if err := s.schemaRepo.UpdateFileParsingConfig(config); err != nil {
		return nil, err
	}

	return config, nil
}

// DeleteFileParsingConfig deletes a file parsing configuration
func (s *SchemaService) DeleteFileParsingConfig(configID, schemaID, userID, tenantID string) error {
	// Check if user has permission to manage schemas
	hasPermission, err := s.permissionRepo.HasPermission(userID, models.PermManageSchemas, tenantID)
	if err != nil {
		return err
	}
	if !hasPermission {
		return errors.New("unauthorized: requires manage schemas permission")
	}

	// Get the schema
	schema, err := s.schemaRepo.GetSchemaByID(schemaID)
	if err != nil {
		return err
	}

	// Ensure the schema belongs to the tenant
	if schema.TenantID != tenantID {
		return errors.New("schema not found in this tenant")
	}

	// Delete the config
	return s.schemaRepo.DeleteFileParsingConfig(configID)
}
