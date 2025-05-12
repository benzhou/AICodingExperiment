package handlers

import (
	"context"
	"errors"
	"net/http"
	"strings"
)

// contextKey is a custom type for context keys
type contextKey string

// Context keys
const (
	UserIDKey   contextKey = "userID"
	TenantIDKey contextKey = "tenantID"
	RolesKey    contextKey = "roles"
)

// GetUserIDFromContext retrieves the user ID from the request context
func GetUserIDFromContext(ctx context.Context) string {
	if userID, ok := ctx.Value(UserIDKey).(string); ok {
		return userID
	}
	return ""
}

// GetTenantIDFromContext retrieves the tenant ID from the request context
func GetTenantIDFromContext(ctx context.Context) string {
	if tenantID, ok := ctx.Value(TenantIDKey).(string); ok {
		return tenantID
	}
	return ""
}

// GetRolesFromContext retrieves the user roles from the request context
func GetRolesFromContext(ctx context.Context) []string {
	if roles, ok := ctx.Value(RolesKey).([]string); ok {
		return roles
	}
	return []string{}
}

// ContextWithUserID adds a user ID to the context
func ContextWithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}

// ContextWithTenantID adds a tenant ID to the context
func ContextWithTenantID(ctx context.Context, tenantID string) context.Context {
	return context.WithValue(ctx, TenantIDKey, tenantID)
}

// ContextWithRoles adds user roles to the context
func ContextWithRoles(ctx context.Context, roles []string) context.Context {
	return context.WithValue(ctx, RolesKey, roles)
}

// handleServiceError handles errors from service calls
func handleServiceError(w http.ResponseWriter, err error) {
	// Check for specific error types
	switch {
	case strings.Contains(err.Error(), "unauthorized"):
		http.Error(w, err.Error(), http.StatusUnauthorized)
	case errors.Is(err, errors.New("not found")) ||
		strings.Contains(err.Error(), "not found"):
		http.Error(w, err.Error(), http.StatusNotFound)
	case strings.Contains(err.Error(), "already exists"):
		http.Error(w, err.Error(), http.StatusConflict)
	case strings.Contains(err.Error(), "invalid"):
		http.Error(w, err.Error(), http.StatusBadRequest)
	default:
		// Internal server error for all other cases
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
