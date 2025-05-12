package handlers

import (
	"backend/internal/repository"
	"backend/internal/services"
	"encoding/json"
	"net/http"
	"strings"
)

type AuthHandler struct {
	authService *services.AuthService
	jwtService  *services.JWTService
	roleService *services.RoleService
}

func NewAuthHandler(userRepo repository.UserRepository, roleService *services.RoleService) *AuthHandler {
	return &AuthHandler{
		authService: services.NewAuthService(userRepo),
		jwtService:  services.NewJWTService(),
		roleService: roleService,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.authService.Register(req.Email, req.Name, req.Password)
	if err != nil {
		if err == repository.ErrUserExists {
			http.Error(w, "User already exists", http.StatusConflict)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	tokenInfo, err := h.jwtService.GenerateToken(user)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	// Get user roles
	var roles []string
	if h.roleService != nil {
		userRoles, err := h.roleService.GetUserRoles(user.ID)
		if err == nil {
			// Convert Role type to strings
			for _, role := range userRoles {
				roles = append(roles, string(role))
			}
		}
	}

	// Check for admin role
	isAdmin := false
	for _, role := range roles {
		if role == "admin" {
			isAdmin = true
			break
		}
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"token":      tokenInfo.Token,
		"expires_in": tokenInfo.ExpiresIn,
		"user":       user,
		"roles":      roles,
		"is_admin":   isAdmin,
	})
}

func (h *AuthHandler) GoogleAuth(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement Google OAuth
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (h *AuthHandler) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement Google OAuth callback
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (h *AuthHandler) GetTokenInfo(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Authorization header required", http.StatusUnauthorized)
		return
	}

	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		http.Error(w, "Invalid authorization header", http.StatusUnauthorized)
		return
	}

	tokenInfo, err := h.jwtService.GetTokenInfo(tokenParts[1])
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(tokenInfo)
}

// GetTokenInfoPublic is a public version of GetTokenInfo for debugging/testing
func (h *AuthHandler) GetTokenInfoPublic(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get token from Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "error",
			"message": "No authorization header provided",
		})
		return
	}

	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "error",
			"message": "Invalid authorization format. Use 'Bearer TOKEN'",
		})
		return
	}

	tokenInfo, err := h.jwtService.GetTokenInfo(tokenParts[1])
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "error",
			"message": "Invalid token: " + err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":     "ok",
		"token_info": tokenInfo,
	})
}
