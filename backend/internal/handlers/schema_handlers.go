package handlers

import (
	"backend/internal/models"
	"backend/internal/services"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// SchemaHandlers provides handlers for schema management
type SchemaHandlers struct {
	schemaService *services.SchemaService
}

// NewSchemaHandlers creates a new schema handlers instance
func NewSchemaHandlers(schemaService *services.SchemaService) *SchemaHandlers {
	return &SchemaHandlers{
		schemaService: schemaService,
	}
}

// RegisterRoutes registers the schema handlers with the router
func (h *SchemaHandlers) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/schemas", h.CreateSchema).Methods("POST")
	router.HandleFunc("/schemas", h.GetSchemas).Methods("GET")
	router.HandleFunc("/schemas/{id}", h.GetSchema).Methods("GET")
	router.HandleFunc("/schemas/{id}", h.UpdateSchema).Methods("PUT")
	router.HandleFunc("/schemas/{id}", h.DeleteSchema).Methods("DELETE")

	router.HandleFunc("/schemas/{schema_id}/fields", h.AddSchemaField).Methods("POST")
	router.HandleFunc("/schemas/{schema_id}/fields", h.GetSchemaFields).Methods("GET")
	router.HandleFunc("/schemas/{schema_id}/fields/{field_id}", h.UpdateSchemaField).Methods("PUT")
	router.HandleFunc("/schemas/{schema_id}/fields/{field_id}", h.DeleteSchemaField).Methods("DELETE")

	router.HandleFunc("/schemas/{schema_id}/mappings", h.CreateSchemaMapping).Methods("POST")
	router.HandleFunc("/schemas/{schema_id}/mappings", h.GetSchemaMappings).Methods("GET")
	router.HandleFunc("/schemas/{schema_id}/mappings/{mapping_id}", h.UpdateSchemaMapping).Methods("PUT")
	router.HandleFunc("/schemas/{schema_id}/mappings/{mapping_id}", h.DeleteSchemaMapping).Methods("DELETE")

	router.HandleFunc("/schemas/{schema_id}/parsing-configs", h.CreateFileParsingConfig).Methods("POST")
	router.HandleFunc("/schemas/{schema_id}/parsing-configs/{file_type}", h.GetFileParsingConfig).Methods("GET")
	router.HandleFunc("/schemas/{schema_id}/parsing-configs/{config_id}", h.UpdateFileParsingConfig).Methods("PUT")
	router.HandleFunc("/schemas/{schema_id}/parsing-configs/{config_id}", h.DeleteFileParsingConfig).Methods("DELETE")
}

// CreateSchema handles the creation of a new schema
func (h *SchemaHandlers) CreateSchema(w http.ResponseWriter, r *http.Request) {
	// Parse tenant ID from context (would be set by authentication middleware)
	tenantID := GetTenantIDFromContext(r.Context())
	if tenantID == "" {
		http.Error(w, "Tenant ID is required", http.StatusBadRequest)
		return
	}

	// Parse user ID from context (would be set by authentication middleware)
	userID := GetUserIDFromContext(r.Context())
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	// Parse request body
	var schema models.DataSourceSchema
	if err := json.NewDecoder(r.Body).Decode(&schema); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Set tenant ID
	schema.TenantID = tenantID

	// Create the schema
	createdSchema, err := h.schemaService.CreateSchema(&schema, userID)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	// Return the created schema
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdSchema)
}

// GetSchemas handles retrieving all schemas for a tenant
func (h *SchemaHandlers) GetSchemas(w http.ResponseWriter, r *http.Request) {
	// Parse tenant ID from context (would be set by authentication middleware)
	tenantID := GetTenantIDFromContext(r.Context())
	if tenantID == "" {
		http.Error(w, "Tenant ID is required", http.StatusBadRequest)
		return
	}

	// Parse user ID from context (would be set by authentication middleware)
	userID := GetUserIDFromContext(r.Context())
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	// Get schemas
	schemas, err := h.schemaService.GetSchemasByTenant(tenantID, userID)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	// Return the schemas
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(schemas)
}

// GetSchema handles retrieving a single schema
func (h *SchemaHandlers) GetSchema(w http.ResponseWriter, r *http.Request) {
	// Parse tenant ID from context (would be set by authentication middleware)
	tenantID := GetTenantIDFromContext(r.Context())
	if tenantID == "" {
		http.Error(w, "Tenant ID is required", http.StatusBadRequest)
		return
	}

	// Parse user ID from context (would be set by authentication middleware)
	userID := GetUserIDFromContext(r.Context())
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	// Get schema ID from URL
	vars := mux.Vars(r)
	schemaID := vars["id"]
	if schemaID == "" {
		http.Error(w, "Schema ID is required", http.StatusBadRequest)
		return
	}

	// Get the schema
	schema, err := h.schemaService.GetSchemaByID(schemaID, userID, tenantID)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	// Return the schema
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(schema)
}

// UpdateSchema handles updating a schema
func (h *SchemaHandlers) UpdateSchema(w http.ResponseWriter, r *http.Request) {
	// Parse tenant ID from context (would be set by authentication middleware)
	tenantID := GetTenantIDFromContext(r.Context())
	if tenantID == "" {
		http.Error(w, "Tenant ID is required", http.StatusBadRequest)
		return
	}

	// Parse user ID from context (would be set by authentication middleware)
	userID := GetUserIDFromContext(r.Context())
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	// Get schema ID from URL
	vars := mux.Vars(r)
	schemaID := vars["id"]
	if schemaID == "" {
		http.Error(w, "Schema ID is required", http.StatusBadRequest)
		return
	}

	// Parse request body
	var schema models.DataSourceSchema
	if err := json.NewDecoder(r.Body).Decode(&schema); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Set ID and tenant ID
	schema.ID = schemaID
	schema.TenantID = tenantID

	// Update the schema
	updatedSchema, err := h.schemaService.UpdateSchema(&schema, userID)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	// Return the updated schema
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedSchema)
}

// DeleteSchema handles deleting a schema
func (h *SchemaHandlers) DeleteSchema(w http.ResponseWriter, r *http.Request) {
	// Parse tenant ID from context (would be set by authentication middleware)
	tenantID := GetTenantIDFromContext(r.Context())
	if tenantID == "" {
		http.Error(w, "Tenant ID is required", http.StatusBadRequest)
		return
	}

	// Parse user ID from context (would be set by authentication middleware)
	userID := GetUserIDFromContext(r.Context())
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	// Get schema ID from URL
	vars := mux.Vars(r)
	schemaID := vars["id"]
	if schemaID == "" {
		http.Error(w, "Schema ID is required", http.StatusBadRequest)
		return
	}

	// Delete the schema
	if err := h.schemaService.DeleteSchema(schemaID, userID, tenantID); err != nil {
		handleServiceError(w, err)
		return
	}

	// Return success
	w.WriteHeader(http.StatusNoContent)
}

// AddSchemaField handles adding a field to a schema
func (h *SchemaHandlers) AddSchemaField(w http.ResponseWriter, r *http.Request) {
	// Parse tenant ID from context (would be set by authentication middleware)
	tenantID := GetTenantIDFromContext(r.Context())
	if tenantID == "" {
		http.Error(w, "Tenant ID is required", http.StatusBadRequest)
		return
	}

	// Parse user ID from context (would be set by authentication middleware)
	userID := GetUserIDFromContext(r.Context())
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	// Get schema ID from URL
	vars := mux.Vars(r)
	schemaID := vars["schema_id"]
	if schemaID == "" {
		http.Error(w, "Schema ID is required", http.StatusBadRequest)
		return
	}

	// Parse request body
	var field models.SchemaField
	if err := json.NewDecoder(r.Body).Decode(&field); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Set schema ID
	field.SchemaID = schemaID

	// Add the field
	addedField, err := h.schemaService.AddFieldToSchema(&field, userID, tenantID)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	// Return the added field
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(addedField)
}

// GetSchemaFields handles retrieving all fields for a schema
func (h *SchemaHandlers) GetSchemaFields(w http.ResponseWriter, r *http.Request) {
	// Parse tenant ID from context (would be set by authentication middleware)
	tenantID := GetTenantIDFromContext(r.Context())
	if tenantID == "" {
		http.Error(w, "Tenant ID is required", http.StatusBadRequest)
		return
	}

	// Parse user ID from context (would be set by authentication middleware)
	userID := GetUserIDFromContext(r.Context())
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	// Get schema ID from URL
	vars := mux.Vars(r)
	schemaID := vars["schema_id"]
	if schemaID == "" {
		http.Error(w, "Schema ID is required", http.StatusBadRequest)
		return
	}

	// Get the schema first to check permissions
	schema, err := h.schemaService.GetSchemaByID(schemaID, userID, tenantID)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	// Return the fields
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(schema.Fields)
}

// UpdateSchemaField handles updating a schema field
func (h *SchemaHandlers) UpdateSchemaField(w http.ResponseWriter, r *http.Request) {
	// Parse tenant ID from context (would be set by authentication middleware)
	tenantID := GetTenantIDFromContext(r.Context())
	if tenantID == "" {
		http.Error(w, "Tenant ID is required", http.StatusBadRequest)
		return
	}

	// Parse user ID from context (would be set by authentication middleware)
	userID := GetUserIDFromContext(r.Context())
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	// Get schema ID and field ID from URL
	vars := mux.Vars(r)
	schemaID := vars["schema_id"]
	if schemaID == "" {
		http.Error(w, "Schema ID is required", http.StatusBadRequest)
		return
	}
	fieldID := vars["field_id"]
	if fieldID == "" {
		http.Error(w, "Field ID is required", http.StatusBadRequest)
		return
	}

	// Parse request body
	var field models.SchemaField
	if err := json.NewDecoder(r.Body).Decode(&field); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Set IDs
	field.ID = fieldID
	field.SchemaID = schemaID

	// Update the field
	updatedField, err := h.schemaService.UpdateSchemaField(&field, userID, tenantID)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	// Return the updated field
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedField)
}

// DeleteSchemaField handles deleting a schema field
func (h *SchemaHandlers) DeleteSchemaField(w http.ResponseWriter, r *http.Request) {
	// Parse tenant ID from context (would be set by authentication middleware)
	tenantID := GetTenantIDFromContext(r.Context())
	if tenantID == "" {
		http.Error(w, "Tenant ID is required", http.StatusBadRequest)
		return
	}

	// Parse user ID from context (would be set by authentication middleware)
	userID := GetUserIDFromContext(r.Context())
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	// Get schema ID and field ID from URL
	vars := mux.Vars(r)
	schemaID := vars["schema_id"]
	if schemaID == "" {
		http.Error(w, "Schema ID is required", http.StatusBadRequest)
		return
	}
	fieldID := vars["field_id"]
	if fieldID == "" {
		http.Error(w, "Field ID is required", http.StatusBadRequest)
		return
	}

	// Delete the field
	if err := h.schemaService.DeleteSchemaField(fieldID, schemaID, userID, tenantID); err != nil {
		handleServiceError(w, err)
		return
	}

	// Return success
	w.WriteHeader(http.StatusNoContent)
}

// CreateSchemaMapping handles creating a schema mapping
func (h *SchemaHandlers) CreateSchemaMapping(w http.ResponseWriter, r *http.Request) {
	// Parse tenant ID from context (would be set by authentication middleware)
	tenantID := GetTenantIDFromContext(r.Context())
	if tenantID == "" {
		http.Error(w, "Tenant ID is required", http.StatusBadRequest)
		return
	}

	// Parse user ID from context (would be set by authentication middleware)
	userID := GetUserIDFromContext(r.Context())
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	// Get schema ID from URL
	vars := mux.Vars(r)
	schemaID := vars["schema_id"]
	if schemaID == "" {
		http.Error(w, "Schema ID is required", http.StatusBadRequest)
		return
	}

	// Parse request body
	var mapping models.SchemaMapping
	if err := json.NewDecoder(r.Body).Decode(&mapping); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Set schema ID
	mapping.SchemaID = schemaID

	// Create the mapping
	createdMapping, err := h.schemaService.CreateSchemaMapping(&mapping, userID, tenantID)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	// Return the created mapping
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdMapping)
}

// GetSchemaMappings handles retrieving all mappings for a schema
func (h *SchemaHandlers) GetSchemaMappings(w http.ResponseWriter, r *http.Request) {
	// Parse tenant ID from context (would be set by authentication middleware)
	tenantID := GetTenantIDFromContext(r.Context())
	if tenantID == "" {
		http.Error(w, "Tenant ID is required", http.StatusBadRequest)
		return
	}

	// Parse user ID from context (would be set by authentication middleware)
	userID := GetUserIDFromContext(r.Context())
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	// Get schema ID from URL
	vars := mux.Vars(r)
	schemaID := vars["schema_id"]
	if schemaID == "" {
		http.Error(w, "Schema ID is required", http.StatusBadRequest)
		return
	}

	// Get mappings
	mappings, err := h.schemaService.GetSchemaMappings(schemaID, userID, tenantID)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	// Return the mappings
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mappings)
}

// UpdateSchemaMapping handles updating a schema mapping
func (h *SchemaHandlers) UpdateSchemaMapping(w http.ResponseWriter, r *http.Request) {
	// Parse tenant ID from context (would be set by authentication middleware)
	tenantID := GetTenantIDFromContext(r.Context())
	if tenantID == "" {
		http.Error(w, "Tenant ID is required", http.StatusBadRequest)
		return
	}

	// Parse user ID from context (would be set by authentication middleware)
	userID := GetUserIDFromContext(r.Context())
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	// Get schema ID and mapping ID from URL
	vars := mux.Vars(r)
	schemaID := vars["schema_id"]
	if schemaID == "" {
		http.Error(w, "Schema ID is required", http.StatusBadRequest)
		return
	}
	mappingID := vars["mapping_id"]
	if mappingID == "" {
		http.Error(w, "Mapping ID is required", http.StatusBadRequest)
		return
	}

	// Parse request body
	var mapping models.SchemaMapping
	if err := json.NewDecoder(r.Body).Decode(&mapping); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Set IDs
	mapping.ID = mappingID
	mapping.SchemaID = schemaID

	// Update the mapping
	updatedMapping, err := h.schemaService.UpdateSchemaMapping(&mapping, userID, tenantID)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	// Return the updated mapping
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedMapping)
}

// DeleteSchemaMapping handles deleting a schema mapping
func (h *SchemaHandlers) DeleteSchemaMapping(w http.ResponseWriter, r *http.Request) {
	// Parse tenant ID from context (would be set by authentication middleware)
	tenantID := GetTenantIDFromContext(r.Context())
	if tenantID == "" {
		http.Error(w, "Tenant ID is required", http.StatusBadRequest)
		return
	}

	// Parse user ID from context (would be set by authentication middleware)
	userID := GetUserIDFromContext(r.Context())
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	// Get schema ID and mapping ID from URL
	vars := mux.Vars(r)
	schemaID := vars["schema_id"]
	if schemaID == "" {
		http.Error(w, "Schema ID is required", http.StatusBadRequest)
		return
	}
	mappingID := vars["mapping_id"]
	if mappingID == "" {
		http.Error(w, "Mapping ID is required", http.StatusBadRequest)
		return
	}

	// Delete the mapping
	if err := h.schemaService.DeleteSchemaMapping(mappingID, schemaID, userID, tenantID); err != nil {
		handleServiceError(w, err)
		return
	}

	// Return success
	w.WriteHeader(http.StatusNoContent)
}

// CreateFileParsingConfig handles creating a file parsing configuration
func (h *SchemaHandlers) CreateFileParsingConfig(w http.ResponseWriter, r *http.Request) {
	// Parse tenant ID from context (would be set by authentication middleware)
	tenantID := GetTenantIDFromContext(r.Context())
	if tenantID == "" {
		http.Error(w, "Tenant ID is required", http.StatusBadRequest)
		return
	}

	// Parse user ID from context (would be set by authentication middleware)
	userID := GetUserIDFromContext(r.Context())
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	// Get schema ID from URL
	vars := mux.Vars(r)
	schemaID := vars["schema_id"]
	if schemaID == "" {
		http.Error(w, "Schema ID is required", http.StatusBadRequest)
		return
	}

	// Parse request body
	var config models.FileParsingConfig
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Set schema ID
	config.SchemaID = schemaID

	// Create the config
	createdConfig, err := h.schemaService.CreateFileParsingConfig(&config, userID, tenantID)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	// Return the created config
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdConfig)
}

// GetFileParsingConfig handles retrieving a file parsing configuration
func (h *SchemaHandlers) GetFileParsingConfig(w http.ResponseWriter, r *http.Request) {
	// Parse tenant ID from context (would be set by authentication middleware)
	tenantID := GetTenantIDFromContext(r.Context())
	if tenantID == "" {
		http.Error(w, "Tenant ID is required", http.StatusBadRequest)
		return
	}

	// Parse user ID from context (would be set by authentication middleware)
	userID := GetUserIDFromContext(r.Context())
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	// Get schema ID and file type from URL
	vars := mux.Vars(r)
	schemaID := vars["schema_id"]
	if schemaID == "" {
		http.Error(w, "Schema ID is required", http.StatusBadRequest)
		return
	}
	fileType := vars["file_type"]
	if fileType == "" {
		http.Error(w, "File type is required", http.StatusBadRequest)
		return
	}

	// Get the config
	config, err := h.schemaService.GetFileParsingConfig(schemaID, fileType, userID, tenantID)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	// Return the config
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(config)
}

// UpdateFileParsingConfig handles updating a file parsing configuration
func (h *SchemaHandlers) UpdateFileParsingConfig(w http.ResponseWriter, r *http.Request) {
	// Parse tenant ID from context (would be set by authentication middleware)
	tenantID := GetTenantIDFromContext(r.Context())
	if tenantID == "" {
		http.Error(w, "Tenant ID is required", http.StatusBadRequest)
		return
	}

	// Parse user ID from context (would be set by authentication middleware)
	userID := GetUserIDFromContext(r.Context())
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	// Get schema ID and config ID from URL
	vars := mux.Vars(r)
	schemaID := vars["schema_id"]
	if schemaID == "" {
		http.Error(w, "Schema ID is required", http.StatusBadRequest)
		return
	}
	configID := vars["config_id"]
	if configID == "" {
		http.Error(w, "Config ID is required", http.StatusBadRequest)
		return
	}

	// Parse request body
	var config models.FileParsingConfig
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Set IDs
	config.ID = configID
	config.SchemaID = schemaID

	// Update the config
	updatedConfig, err := h.schemaService.UpdateFileParsingConfig(&config, userID, tenantID)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	// Return the updated config
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedConfig)
}

// DeleteFileParsingConfig handles deleting a file parsing configuration
func (h *SchemaHandlers) DeleteFileParsingConfig(w http.ResponseWriter, r *http.Request) {
	// Parse tenant ID from context (would be set by authentication middleware)
	tenantID := GetTenantIDFromContext(r.Context())
	if tenantID == "" {
		http.Error(w, "Tenant ID is required", http.StatusBadRequest)
		return
	}

	// Parse user ID from context (would be set by authentication middleware)
	userID := GetUserIDFromContext(r.Context())
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	// Get schema ID and config ID from URL
	vars := mux.Vars(r)
	schemaID := vars["schema_id"]
	if schemaID == "" {
		http.Error(w, "Schema ID is required", http.StatusBadRequest)
		return
	}
	configID := vars["config_id"]
	if configID == "" {
		http.Error(w, "Config ID is required", http.StatusBadRequest)
		return
	}

	// Delete the config
	if err := h.schemaService.DeleteFileParsingConfig(configID, schemaID, userID, tenantID); err != nil {
		handleServiceError(w, err)
		return
	}

	// Return success
	w.WriteHeader(http.StatusNoContent)
}
