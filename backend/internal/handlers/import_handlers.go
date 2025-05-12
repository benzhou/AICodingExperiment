package handlers

import (
	"backend/internal/repository"
	"backend/internal/utils"
	"net/http"

	"github.com/gorilla/mux"
)

// ImportHandlers contains handlers for import-related endpoints
type ImportHandlers struct {
	importRepo repository.ImportRepository
}

// NewImportHandlers creates a new instance of ImportHandlers
func NewImportHandlers(importRepo repository.ImportRepository) *ImportHandlers {
	return &ImportHandlers{
		importRepo: importRepo,
	}
}

// RegisterRoutes registers the import handlers with the router
func (h *ImportHandlers) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/api/v1/datasources/{dataSourceId}/imports", h.GetImportsByDataSource).Methods("GET")
	r.HandleFunc("/api/v1/imports/{importId}", h.GetImportByID).Methods("GET")
	r.HandleFunc("/api/v1/imports/{importId}", h.DeleteImport).Methods("DELETE")
	r.HandleFunc("/api/v1/imports/{importId}/raw-transactions", h.GetRawTransactionsByImport).Methods("GET")
	r.HandleFunc("/api/v1/raw-transactions/{rawTransactionId}", h.GetRawTransactionByID).Methods("GET")
}

// GetImportsByDataSource returns all imports for a data source with pagination
func (h *ImportHandlers) GetImportsByDataSource(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	dataSourceID := vars["dataSourceId"]

	// Parse pagination parameters
	limit, offset := utils.GetPaginationParams(r)

	// Get imports from repository
	imports, total, err := h.importRepo.GetImportsByDataSource(dataSourceID, limit, offset)
	if err != nil {
		utils.WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Build response with pagination metadata
	response := map[string]interface{}{
		"data":       imports,
		"pagination": utils.BuildPaginationResponse(limit, offset, total),
	}

	utils.WriteJSON(w, response, http.StatusOK)
}

// GetImportByID returns a single import by ID
func (h *ImportHandlers) GetImportByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	importID := vars["importId"]

	importRecord, err := h.importRepo.GetImportByID(importID)
	if err != nil {
		if err == repository.ErrImportNotFound {
			utils.WriteError(w, "Import not found", http.StatusNotFound)
		} else {
			utils.WriteError(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	utils.WriteJSON(w, importRecord, http.StatusOK)
}

// DeleteImport deletes an import and its associated raw transactions
func (h *ImportHandlers) DeleteImport(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	importID := vars["importId"]

	err := h.importRepo.DeleteImport(importID)
	if err != nil {
		if err == repository.ErrImportNotFound {
			utils.WriteError(w, "Import not found", http.StatusNotFound)
		} else {
			utils.WriteError(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetRawTransactionsByImport returns all raw transactions for an import with pagination
func (h *ImportHandlers) GetRawTransactionsByImport(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	importID := vars["importId"]

	// Parse pagination parameters
	limit, offset := utils.GetPaginationParams(r)

	// Get raw transactions from repository
	transactions, total, err := h.importRepo.GetRawTransactionsByImport(importID, limit, offset)
	if err != nil {
		utils.WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Build response with pagination metadata
	response := map[string]interface{}{
		"data":       transactions,
		"pagination": utils.BuildPaginationResponse(limit, offset, total),
	}

	utils.WriteJSON(w, response, http.StatusOK)
}

// GetRawTransactionByID returns a single raw transaction by ID
func (h *ImportHandlers) GetRawTransactionByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	rawTransactionID := vars["rawTransactionId"]

	transaction, err := h.importRepo.GetRawTransactionByID(rawTransactionID)
	if err != nil {
		utils.WriteError(w, err.Error(), http.StatusNotFound)
		return
	}

	utils.WriteJSON(w, transaction, http.StatusOK)
}
