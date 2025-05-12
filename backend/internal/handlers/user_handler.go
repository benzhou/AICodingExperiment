package handlers

import (
	"backend/internal/models"
	"backend/internal/services"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type UserHandler struct {
	userService *services.UserService
	roleService *services.RoleService
}

func NewUserHandler(userService *services.UserService, roleService *services.RoleService) *UserHandler {
	return &UserHandler{
		userService: userService,
		roleService: roleService,
	}
}

// GetAllUsers retrieves all users (admin only)
func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	// Check if the requesting user has admin privileges
	requesterID, ok := r.Context().Value("userID").(string)
	if !ok || requesterID == "" {
		http.Error(w, "Unauthorized - could not authenticate user", http.StatusUnauthorized)
		return
	}

	isAdmin, err := h.roleService.HasRole(requesterID, models.RoleAdmin)
	if err != nil || !isAdmin {
		http.Error(w, "Unauthorized - admin privileges required", http.StatusUnauthorized)
		return
	}

	users, err := h.userService.GetAllUsers()
	if err != nil {
		http.Error(w, "Error fetching users: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// GetUserRoles gets all roles for a user
func (h *UserHandler) GetUserRoles(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	// Get the requesting user's ID
	requesterID, ok := r.Context().Value("userID").(string)
	if !ok || requesterID == "" {
		http.Error(w, "Unauthorized - could not authenticate user", http.StatusUnauthorized)
		return
	}

	// Allow users to get their own roles or admins to get any user's roles
	if requesterID != userID {
		isAdmin, err := h.roleService.HasRole(requesterID, models.RoleAdmin)
		if err != nil || !isAdmin {
			http.Error(w, "Unauthorized - admin privileges required to view other users' roles", http.StatusUnauthorized)
			return
		}
	}

	roles, err := h.userService.GetUserRoles(userID)
	if err != nil {
		http.Error(w, "Error fetching user roles", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id": userID,
		"roles":   roles,
	})
}

// UpdateUserRole adds or removes a role from a user
func (h *UserHandler) UpdateUserRole(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Role      string `json:"role"`
		Operation string `json:"operation"` // "add" or "remove"
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	userID := vars["id"]

	// Check if the requesting user has admin privileges
	requesterID, ok := r.Context().Value("userID").(string)
	if !ok || requesterID == "" {
		http.Error(w, "Unauthorized - could not authenticate user", http.StatusUnauthorized)
		return
	}

	isAdmin, err := h.roleService.HasRole(requesterID, models.RoleAdmin)
	if err != nil || !isAdmin {
		http.Error(w, "Unauthorized - admin privileges required", http.StatusUnauthorized)
		return
	}

	// Validate operation
	if req.Operation != "add" && req.Operation != "remove" {
		http.Error(w, "Invalid operation - must be 'add' or 'remove'", http.StatusBadRequest)
		return
	}

	// Validate role
	var role models.Role
	switch req.Role {
	case "admin":
		role = models.RoleAdmin
	case "preparer":
		role = models.RolePreparer
	case "approver":
		role = models.RoleApprover
	default:
		http.Error(w, "Invalid role - must be 'admin', 'preparer', or 'approver'", http.StatusBadRequest)
		return
	}

	// Update the user role
	err = h.userService.UpdateUserRole(userID, role, req.Operation)
	if err != nil {
		http.Error(w, "Error updating user role: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return updated roles
	roles, err := h.userService.GetUserRoles(userID)
	if err != nil {
		http.Error(w, "Error fetching updated user roles", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User role updated successfully",
		"user_id": userID,
		"roles":   roles,
	})
}

// CreateUserWithRole creates a new user with a specific role
func (h *UserHandler) CreateUserWithRole(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
		Role     string `json:"role"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Check if the requesting user has admin privileges
	requesterID, ok := r.Context().Value("userID").(string)
	if !ok || requesterID == "" {
		http.Error(w, "Unauthorized - could not authenticate user", http.StatusUnauthorized)
		return
	}

	isAdmin, err := h.roleService.HasRole(requesterID, models.RoleAdmin)
	if err != nil || !isAdmin {
		http.Error(w, "Unauthorized - admin privileges required", http.StatusUnauthorized)
		return
	}

	// Validate role
	var role models.Role
	switch req.Role {
	case "admin":
		role = models.RoleAdmin
	case "preparer":
		role = models.RolePreparer
	case "approver":
		role = models.RoleApprover
	default:
		// Default to preparer if no valid role provided
		role = models.RolePreparer
	}

	// Create user with role
	user, err := h.userService.CreateUserWithRole(req.Email, req.Name, req.Password, role)
	if err != nil {
		http.Error(w, "Error creating user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the assigned roles
	roles, _ := h.userService.GetUserRoles(user.ID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User created successfully",
		"user":    user,
		"roles":   roles,
	})
}

// GetUserById retrieves a user by ID
func (h *UserHandler) GetUserById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	// Check if the requesting user has admin privileges or is requesting their own info
	requesterID, ok := r.Context().Value("userID").(string)
	if !ok || requesterID == "" {
		http.Error(w, "Unauthorized - could not authenticate user", http.StatusUnauthorized)
		return
	}

	if requesterID != userID {
		isAdmin, err := h.roleService.HasRole(requesterID, models.RoleAdmin)
		if err != nil || !isAdmin {
			http.Error(w, "Unauthorized - admin privileges required", http.StatusUnauthorized)
			return
		}
	}

	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		http.Error(w, "Error fetching user", http.StatusInternalServerError)
		return
	}

	// Get user roles
	roles, _ := h.userService.GetUserRoles(userID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user":  user,
		"roles": roles,
	})
}

// SetAdminRole grants admin role to a user
func (h *UserHandler) SetAdminRole(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	// Check if the requesting user has admin privileges
	requesterID, ok := r.Context().Value("userID").(string)
	if !ok || requesterID == "" {
		http.Error(w, "Unauthorized - could not authenticate user", http.StatusUnauthorized)
		return
	}

	isAdmin, err := h.roleService.HasRole(requesterID, models.RoleAdmin)
	if err != nil || !isAdmin {
		http.Error(w, "Unauthorized - admin privileges required", http.StatusUnauthorized)
		return
	}

	// Set user as admin
	err = h.userService.SetUserAsAdmin(userID)
	if err != nil {
		http.Error(w, "Error granting admin role: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Admin role granted successfully",
		"user_id": userID,
	})
}
