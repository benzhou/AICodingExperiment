package handlers

import (
	"backend/internal/models"
	"backend/internal/services"
	"backend/internal/utils"
	"encoding/json"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
)

// DataSourceHandler handles data source-related API endpoints
type DataSourceHandler struct {
	dataSourceService *services.DataSourceService
	roleService       *services.RoleService
}

// NewDataSourceHandler creates a new data source handler
func NewDataSourceHandler(
	dataSourceService *services.DataSourceService,
	roleService *services.RoleService,
) *DataSourceHandler {
	return &DataSourceHandler{
		dataSourceService: dataSourceService,
		roleService:       roleService,
	}
}

// CreateDataSource handles data source creation
func (h *DataSourceHandler) CreateDataSource(w http.ResponseWriter, r *http.Request) {
	// Get user claims from JWT token
	userClaims, ok := r.Context().Value("user").(*jwt.MapClaims)
	if !ok || userClaims == nil {
		http.Error(w, "Unauthorized: invalid or missing authentication", http.StatusUnauthorized)
		return
	}

	// Extract user ID from claims
	userIDValue, ok := (*userClaims)["user_id"]
	if !ok || userIDValue == nil {
		http.Error(w, "Unauthorized: user ID not found in token", http.StatusUnauthorized)
		return
	}

	userID := userIDValue.(string)

	// Check if user has admin role
	hasRole, err := h.roleService.HasRole(userID, models.RoleAdmin)
	if err != nil || !hasRole {
		http.Error(w, "Unauthorized: requires admin role", http.StatusUnauthorized)
		return
	}

	// Parse request body
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	// Validate request
	if req.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	// Create data source
	dataSource, err := h.dataSourceService.CreateDataSource(req.Name, req.Description)
	if err != nil {
		if err == services.ErrDataSourceExists {
			http.Error(w, "Data source with this name already exists", http.StatusConflict)
			return
		}
		http.Error(w, "Error creating data source: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return created data source
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dataSource)
}

// GetDataSourceByID retrieves a data source by ID
func (h *DataSourceHandler) GetDataSourceByID(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL path
	vars := mux.Vars(r)
	id := vars["id"]

	// Get data source
	dataSource, err := h.dataSourceService.GetDataSourceByID(id)
	if err != nil {
		http.Error(w, "Data source not found", http.StatusNotFound)
		return
	}

	// Return data source
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dataSource)
}

// UpdateDataSource handles data source updates
func (h *DataSourceHandler) UpdateDataSource(w http.ResponseWriter, r *http.Request) {
	// Get user claims from JWT token
	userClaims, ok := r.Context().Value("user").(*jwt.MapClaims)
	if !ok || userClaims == nil {
		http.Error(w, "Unauthorized: invalid or missing authentication", http.StatusUnauthorized)
		return
	}

	// Extract user ID from claims
	userIDValue, ok := (*userClaims)["user_id"]
	if !ok || userIDValue == nil {
		http.Error(w, "Unauthorized: user ID not found in token", http.StatusUnauthorized)
		return
	}

	userID := userIDValue.(string)

	// Check if user has admin role
	hasRole, err := h.roleService.HasRole(userID, models.RoleAdmin)
	if err != nil || !hasRole {
		http.Error(w, "Unauthorized: requires admin role", http.StatusUnauthorized)
		return
	}

	// Extract ID from URL path
	vars := mux.Vars(r)
	id := vars["id"]

	// Parse request body
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	// Validate request
	if req.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	// Update data source
	dataSource, err := h.dataSourceService.UpdateDataSource(id, req.Name, req.Description)
	if err != nil {
		if err == services.ErrDataSourceNotFound {
			http.Error(w, "Data source not found", http.StatusNotFound)
			return
		}
		if err == services.ErrDataSourceExists {
			http.Error(w, "Data source with this name already exists", http.StatusConflict)
			return
		}
		http.Error(w, "Error updating data source: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return updated data source
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dataSource)
}

// DeleteDataSource handles data source deletion
func (h *DataSourceHandler) DeleteDataSource(w http.ResponseWriter, r *http.Request) {
	// Get user claims from JWT token
	userClaims, ok := r.Context().Value("user").(*jwt.MapClaims)
	if !ok || userClaims == nil {
		http.Error(w, "Unauthorized: invalid or missing authentication", http.StatusUnauthorized)
		return
	}

	// Extract user ID from claims
	userIDValue, ok := (*userClaims)["user_id"]
	if !ok || userIDValue == nil {
		http.Error(w, "Unauthorized: user ID not found in token", http.StatusUnauthorized)
		return
	}

	userID := userIDValue.(string)

	// Check if user has admin role
	hasRole, err := h.roleService.HasRole(userID, models.RoleAdmin)
	if err != nil || !hasRole {
		http.Error(w, "Unauthorized: requires admin role", http.StatusUnauthorized)
		return
	}

	// Extract ID from URL path
	vars := mux.Vars(r)
	id := vars["id"]

	// Delete data source
	err = h.dataSourceService.DeleteDataSource(id)
	if err != nil {
		if err == services.ErrDataSourceNotFound {
			http.Error(w, "Data source not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Error deleting data source: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return success
	w.WriteHeader(http.StatusNoContent)
}

// GetAllDataSources retrieves all data sources
func (h *DataSourceHandler) GetAllDataSources(w http.ResponseWriter, r *http.Request) {
	// Get all data sources
	dataSources, err := h.dataSourceService.GetAllDataSources()
	if err != nil {
		http.Error(w, "Error retrieving data sources: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return data sources
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dataSources)
}

// SearchDataSources handles searching for data sources
func (h *DataSourceHandler) SearchDataSources(w http.ResponseWriter, r *http.Request) {
	// Get query parameters
	query := r.URL.Query().Get("q")
	if query == "" {
		// If no query provided, fall back to listing all data sources
		h.GetAllDataSources(w, r)
		return
	}

	// Get pagination parameters
	limit, offset := utils.GetPaginationParams(r)

	// Search data sources
	dataSources, total, err := h.dataSourceService.SearchDataSources(query, limit, offset)
	if err != nil {
		http.Error(w, "Error searching data sources: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Build response with pagination metadata
	response := map[string]interface{}{
		"data":       dataSources,
		"pagination": utils.BuildPaginationResponse(limit, offset, total),
	}

	// Return search results
	utils.WriteJSON(w, response, http.StatusOK)
}
