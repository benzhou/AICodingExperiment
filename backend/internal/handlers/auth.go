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
}

func NewAuthHandler(userRepo repository.UserRepository) *AuthHandler {
	return &AuthHandler{
		authService: services.NewAuthService(userRepo),
		jwtService:  services.NewJWTService(),
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

	json.NewEncoder(w).Encode(map[string]interface{}{
		"token":      tokenInfo.Token,
		"expires_in": tokenInfo.ExpiresIn,
		"user":       user,
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
