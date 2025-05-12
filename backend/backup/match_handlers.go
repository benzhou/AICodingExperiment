package handlers

import (
	"backend/internal/models"
	"backend/internal/repository"
	"backend/internal/services"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// MatchHandler handles match-related API endpoints
type MatchHandler struct {
	matchRepo          repository.MatchRepository
	transactionRepo    repository.TransactionRepository
	transactionService *services.TransactionService
	roleService        *services.RoleService
}

// NewMatchHandler creates a new match handler
func NewMatchHandler(
	matchRepo repository.MatchRepository,
	transactionRepo repository.TransactionRepository,
	transactionService *services.TransactionService,
	roleService *services.RoleService,
) *MatchHandler {
	return &MatchHandler{
		matchRepo:          matchRepo,
		transactionRepo:    transactionRepo,
		transactionService: transactionService,
		roleService:        roleService,
	}
}

// GetMatchByID retrieves a match by ID
func (h *MatchHandler) GetMatchByID(w http.ResponseWriter, r *http.Request) {
	// Extract match ID from URL path
	vars := mux.Vars(r)
	id := vars["id"]

	// Get match
	match, err := h.matchRepo.GetMatchByID(id)
	if err != nil {
		http.Error(w, "Match not found", http.StatusNotFound)
		return
	}

	// Get transactions associated with the match
	transactions, err := h.transactionRepo.GetMatchedTransactions(id)
	if err != nil {
		http.Error(w, "Error retrieving matched transactions", http.StatusInternalServerError)
		return
	}

	// Create response
	response := map[string]interface{}{
		"match":        match,
		"transactions": transactions,
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// SearchMatches searches for matches using filters
func (h *MatchHandler) SearchMatches(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	r.ParseForm()

	// Build filters
	filters := make(map[string]interface{})

	// Status filter
	if status := r.Form.Get("status"); status != "" {
		filters["status"] = status
	}

	// Match type filter
	if matchType := r.Form.Get("matchType"); matchType != "" {
		filters["matchType"] = matchType
	}

	// User filters
	if matchedBy := r.Form.Get("matchedBy"); matchedBy != "" {
		filters["matchedBy"] = matchedBy
	}

	if approvedBy := r.Form.Get("approvedBy"); approvedBy != "" {
		filters["approvedBy"] = approvedBy
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

	// Calculate offset
	offset := (page - 1) * pageSize

	// Search matches
	matches, total, err := h.matchRepo.SearchMatches(filters, pageSize, offset)
	if err != nil {
		http.Error(w, "Error searching matches", http.StatusInternalServerError)
		return
	}

	// Create response
	response := map[string]interface{}{
		"matches":  matches,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ApproveMatch approves a match
func (h *MatchHandler) ApproveMatch(w http.ResponseWriter, r *http.Request) {
	// Get user ID from JWT token
	userID := r.Context().Value("userID").(string)

	// Check if user has the approver role
	hasRole, err := h.roleService.HasRole(userID, models.RoleApprover)
	if err != nil || !hasRole {
		http.Error(w, "Unauthorized: requires approver role", http.StatusUnauthorized)
		return
	}

	// Extract match ID from URL path
	vars := mux.Vars(r)
	id := vars["id"]

	// Approve match
	err = h.transactionService.ApproveMatch(id, userID)
	if err != nil {
		http.Error(w, "Error approving match: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Return success
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "approved"})
}

// RejectMatch rejects a match
func (h *MatchHandler) RejectMatch(w http.ResponseWriter, r *http.Request) {
	// Get user ID from JWT token
	userID := r.Context().Value("userID").(string)

	// Check if user has the approver role
	hasRole, err := h.roleService.HasRole(userID, models.RoleApprover)
	if err != nil || !hasRole {
		http.Error(w, "Unauthorized: requires approver role", http.StatusUnauthorized)
		return
	}

	// Extract match ID from URL path
	vars := mux.Vars(r)
	id := vars["id"]

	// Parse request body
	var request struct {
		Reason string `json:"reason"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Reject match
	err = h.transactionService.RejectMatch(id, userID, request.Reason)
	if err != nil {
		http.Error(w, "Error rejecting match: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Return success
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "rejected"})
}

// GetPendingMatches retrieves pending matches for approval
func (h *MatchHandler) GetPendingMatches(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	r.ParseForm()

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

	// Calculate offset
	offset := (page - 1) * pageSize

	// Search pending matches
	filters := map[string]interface{}{
		"status": "Pending",
	}

	matches, total, err := h.matchRepo.SearchMatches(filters, pageSize, offset)
	if err != nil {
		http.Error(w, "Error retrieving pending matches", http.StatusInternalServerError)
		return
	}

	// Create response
	response := map[string]interface{}{
		"matches":  matches,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetMatchDetails retrieves detailed information about a match
func (h *MatchHandler) GetMatchDetails(w http.ResponseWriter, r *http.Request) {
	// Extract match ID from URL path
	vars := mux.Vars(r)
	id := vars["id"]

	// Get match
	match, err := h.matchRepo.GetMatchByID(id)
	if err != nil {
		http.Error(w, "Match not found", http.StatusNotFound)
		return
	}

	// Get transactions associated with the match
	transactions, err := h.transactionRepo.GetMatchedTransactions(id)
	if err != nil {
		http.Error(w, "Error retrieving matched transactions", http.StatusInternalServerError)
		return
	}

	// Create response
	response := map[string]interface{}{
		"match":        match,
		"transactions": transactions,
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
