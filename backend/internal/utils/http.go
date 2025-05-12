package utils

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

// GetPaginationParams extracts limit and offset query parameters from the request
func GetPaginationParams(r *http.Request) (int, int) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 10 // Default limit
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	offset := 0 // Default offset
	if offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	return limit, offset
}

// BuildPaginationResponse creates a pagination response object
func BuildPaginationResponse(limit, offset, total int) map[string]interface{} {
	return map[string]interface{}{
		"limit":   limit,
		"offset":  offset,
		"total":   total,
		"hasMore": offset+limit < total,
	}
}

// WriteJSON writes a JSON response
func WriteJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			// Log error but don't expose it to the client
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}
}

// WriteError writes an error response in JSON format
func WriteError(w http.ResponseWriter, message string, statusCode int) {
	errorResponse := map[string]string{
		"error": message,
	}

	WriteJSON(w, errorResponse, statusCode)
}

// TimeToEpochMS converts a time.Time to epoch milliseconds for API responses
func TimeToEpochMS(t time.Time) int64 {
	// Ensure time is in UTC
	return t.UTC().UnixNano() / int64(time.Millisecond)
}
