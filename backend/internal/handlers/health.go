package handlers

import (
	"encoding/json"
	"net/http"
)

// HealthCheckHandler is a simple handler to check if the server is running
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := map[string]string{
		"status":  "ok",
		"message": "Server is running",
	}

	json.NewEncoder(w).Encode(response)
}
