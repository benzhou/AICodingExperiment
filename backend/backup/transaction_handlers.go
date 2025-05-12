package handlers

import (
	"backend/internal/models"
	"backend/internal/services"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// TransactionHandler handles transaction-related API endpoints
type TransactionHandler struct {
	transactionService *services.TransactionService
	roleService        *services.RoleService
}

// NewTransactionHandler creates a new transaction handler
func NewTransactionHandler(
	transactionService *services.TransactionService,
	roleService *services.RoleService,
) *TransactionHandler {
	return &TransactionHandler{
		transactionService: transactionService,
		roleService:        roleService,
	}
}

// GetTransactionByID retrieves a transaction by ID
func (h *TransactionHandler) GetTransactionByID(w http.ResponseWriter, r *http.Request) {
	// Extract transaction ID from URL path
	vars := mux.Vars(r)
	id := vars["id"]

	// Get transaction
	transaction, err := h.transactionService.GetTransactionByID(id)
	if err != nil {
		http.Error(w, "Transaction not found", http.StatusNotFound)
		return
	}

	// Return transaction
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transaction)
}

// SearchTransactions searches for transactions using filters
func (h *TransactionHandler) SearchTransactions(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	r.ParseForm()

	// Build filters
	filters := make(map[string]interface{})

	// Date range filter
	if dateFrom := r.Form.Get("dateFrom"); dateFrom != "" {
		if date, err := time.Parse("2006-01-02", dateFrom); err == nil {
			filters["dateFrom"] = date
		}
	}

	if dateTo := r.Form.Get("dateTo"); dateTo != "" {
		if date, err := time.Parse("2006-01-02", dateTo); err == nil {
			filters["dateTo"] = date
		}
	}

	// Amount range filter
	if amountMin := r.Form.Get("amountMin"); amountMin != "" {
		if amount, err := strconv.ParseFloat(amountMin, 64); err == nil {
			filters["amountMin"] = amount
		}
	}

	if amountMax := r.Form.Get("amountMax"); amountMax != "" {
		if amount, err := strconv.ParseFloat(amountMax, 64); err == nil {
			filters["amountMax"] = amount
		}
	}

	// Data source filter
	if dataSourceID := r.Form.Get("dataSourceId"); dataSourceID != "" {
		filters["dataSourceId"] = dataSourceID
	}

	// Status filter
	if status := r.Form.Get("status"); status != "" {
		filters["status"] = status
	}

	// Search term
	if searchTerm := r.Form.Get("searchTerm"); searchTerm != "" {
		filters["searchTerm"] = searchTerm
	}

	// Pagination
	page := 1
	if pageStr := r.Form.Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	pageSize := 20
	if pageSizeStr := r.Form.Get("pageSize"); pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 {
			pageSize = ps
		}
	}

	// Search transactions
	transactions, total, err := h.transactionService.SearchTransactions(filters, page, pageSize)
	if err != nil {
		http.Error(w, "Error searching transactions", http.StatusInternalServerError)
		return
	}

	// Create response
	response := map[string]interface{}{
		"transactions": transactions,
		"total":        total,
		"page":         page,
		"pageSize":     pageSize,
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// CreateManualMatch creates a manual match between transactions
func (h *TransactionHandler) CreateManualMatch(w http.ResponseWriter, r *http.Request) {
	// Get user ID from JWT token
	userID := r.Context().Value("userID").(string)

	// Check if user has the preparer role
	hasRole, err := h.roleService.HasRole(userID, models.RolePreparer)
	if err != nil || !hasRole {
		http.Error(w, "Unauthorized: requires preparer role", http.StatusUnauthorized)
		return
	}

	// Parse request body
	var request struct {
		TransactionIDs []string `json:"transactionIds"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Create manual match
	match, err := h.transactionService.CreateManualMatch(request.TransactionIDs, userID)
	if err != nil {
		http.Error(w, "Error creating match: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Return match
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(match)
}

// FindPotentialMatches finds potential matching transactions for a given transaction
func (h *TransactionHandler) FindPotentialMatches(w http.ResponseWriter, r *http.Request) {
	// Get user ID from JWT token
	userID := r.Context().Value("userID").(string)

	// Check if user has the preparer role
	hasRole, err := h.roleService.HasRole(userID, models.RolePreparer)
	if err != nil || !hasRole {
		http.Error(w, "Unauthorized: requires preparer role", http.StatusUnauthorized)
		return
	}

	// Extract transaction ID from URL path
	vars := mux.Vars(r)
	id := vars["id"]

	// Find potential matches
	matches, err := h.transactionService.FindPotentialMatches(id)
	if err != nil {
		http.Error(w, "Error finding potential matches: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return matches
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(matches)
}

// AutoMatchTransactions attempts to automatically match unmatched transactions
func (h *TransactionHandler) AutoMatchTransactions(w http.ResponseWriter, r *http.Request) {
	// Get user ID from JWT token
	userID := r.Context().Value("userID").(string)

	// Check if user has the preparer role
	hasRole, err := h.roleService.HasRole(userID, models.RolePreparer)
	if err != nil || !hasRole {
		http.Error(w, "Unauthorized: requires preparer role", http.StatusUnauthorized)
		return
	}

	// Auto match transactions
	matchCount, err := h.transactionService.AutoMatchTransactions(userID)
	if err != nil {
		http.Error(w, "Error auto-matching transactions: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return match count
	response := map[string]interface{}{
		"matchCount": matchCount,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
