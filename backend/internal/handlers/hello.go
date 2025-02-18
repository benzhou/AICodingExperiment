package handlers

import (
	"encoding/json"
	"net/http"
)

type HelloResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Version string `json:"version"`
}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	response := HelloResponse{
		Status:  "success",
		Message: "Hello from Go Backend!",
		Version: "v1",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
