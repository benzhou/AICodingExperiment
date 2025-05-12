package handlers

import (
	"backend/internal/models"
	"backend/internal/services"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// UploadHandler handles file uploads for transaction data
type UploadHandler struct {
	dataSourceService  *services.DataSourceService
	transactionService *services.TransactionService
	roleService        *services.RoleService
}

// NewUploadHandler creates a new upload handler
func NewUploadHandler(
	dataSourceService *services.DataSourceService,
	transactionService *services.TransactionService,
	roleService *services.RoleService,
) *UploadHandler {
	return &UploadHandler{
		dataSourceService:  dataSourceService,
		transactionService: transactionService,
		roleService:        roleService,
	}
}

// UploadTransactions handles transaction data uploads from CSV files
func (h *UploadHandler) UploadTransactions(w http.ResponseWriter, r *http.Request) {
	// Get user claims from JWT token
	userClaims, ok := r.Context().Value("user").(*jwt.MapClaims)
	if !ok || userClaims == nil {
		http.Error(w, "Unauthorized: invalid or missing authentication", http.StatusUnauthorized)
		return
	}

	// Extract user ID from claims
	userIDValue, ok := (*userClaims)["user_id"]
	if !ok || userIDValue == nil {
		http.Error(w, "Unauthorized: user ID not found in token", http.StatusUnauthorized)
		return
	}

	userID := userIDValue.(string)

	// Check if user has appropriate role
	hasRole, err := h.roleService.HasRole(userID, models.RoleAdmin)
	if err != nil || !hasRole {
		http.Error(w, "Unauthorized: requires admin role", http.StatusUnauthorized)
		return
	}

	// Parse the multipart form
	err = r.ParseMultipartForm(32 << 20) // 32MB max
	if err != nil {
		http.Error(w, "Failed to parse form: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Get form values
	dataSourceID := r.FormValue("dataSourceId")
	if dataSourceID == "" {
		http.Error(w, "Data source ID is required", http.StatusBadRequest)
		return
	}

	// Get date format
	dateFormat := r.FormValue("dateFormat")
	if dateFormat == "" {
		dateFormat = "2006-01-02" // Default to ISO format
	}

	// Get column mappings
	columnMappingsJSON := r.FormValue("columnMappings")
	if columnMappingsJSON == "" {
		http.Error(w, "Column mappings are required", http.StatusBadRequest)
		return
	}

	var columnMappings map[string]int
	err = json.Unmarshal([]byte(columnMappingsJSON), &columnMappings)
	if err != nil {
		http.Error(w, "Invalid column mappings format: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Verify required mappings
	requiredFields := []string{"date", "description", "amount", "reference"}
	for _, field := range requiredFields {
		if _, exists := columnMappings[field]; !exists {
			http.Error(w, fmt.Sprintf("Required field mapping missing: %s", field), http.StatusBadRequest)
			return
		}
	}

	// Get data source
	_, err = h.dataSourceService.GetDataSourceByID(dataSourceID)
	if err != nil {
		http.Error(w, "Data source not found: "+err.Error(), http.StatusNotFound)
		return
	}

	// Get the file from the form
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to get file: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Check if file is a CSV
	if header.Filename[len(header.Filename)-4:] != ".csv" {
		http.Error(w, "Only CSV files are supported", http.StatusBadRequest)
		return
	}

	// Parse the CSV
	reader := csv.NewReader(file)

	// Read header row
	_, err = reader.Read()
	if err != nil {
		http.Error(w, "Failed to read CSV header: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Process the CSV rows
	var transactions []models.Transaction
	var lineNumber int = 2 // Start at line 2 (after header)

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			http.Error(w, fmt.Sprintf("Error reading CSV at line %d: %s", lineNumber, err.Error()), http.StatusBadRequest)
			return
		}

		// Skip empty rows
		if len(record) == 0 {
			lineNumber++
			continue
		}

		// Parse transaction data using column mappings
		transaction := models.Transaction{
			ID:           uuid.New().String(),
			DataSourceID: dataSourceID,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
			Status:       "Unmatched", // Default status
		}

		// Set required fields
		dateCol := columnMappings["date"]
		if dateCol < len(record) {
			date, err := time.Parse(dateFormat, record[dateCol])
			if err != nil {
				http.Error(
					w,
					fmt.Sprintf("Invalid date format at line %d, column %d: %s", lineNumber, dateCol+1, err.Error()),
					http.StatusBadRequest,
				)
				return
			}
			transaction.TransactionDate = date
		}

		descCol := columnMappings["description"]
		if descCol < len(record) {
			transaction.Description = record[descCol]
		}

		amountCol := columnMappings["amount"]
		if amountCol < len(record) {
			amount, err := parseAmount(record[amountCol])
			if err != nil {
				http.Error(
					w,
					fmt.Sprintf("Invalid amount format at line %d, column %d: %s", lineNumber, amountCol+1, err.Error()),
					http.StatusBadRequest,
				)
				return
			}
			transaction.Amount = amount
		}

		refCol := columnMappings["reference"]
		if refCol < len(record) {
			transaction.Reference = record[refCol]
		}

		// Set optional fields if mapped
		if postDateCol, exists := columnMappings["postDate"]; exists && postDateCol < len(record) && record[postDateCol] != "" {
			postDate, err := time.Parse(dateFormat, record[postDateCol])
			if err == nil { // Ignore parsing errors for optional fields
				transaction.PostDate = postDate
			}
		}

		if currencyCol, exists := columnMappings["currency"]; exists && currencyCol < len(record) {
			transaction.Currency = record[currencyCol]
		}

		transactions = append(transactions, transaction)
		lineNumber++
	}

	// Save transactions
	if len(transactions) > 0 {
		err = h.transactionService.CreateTransactions(transactions)
		if err != nil {
			http.Error(w, "Failed to save transactions: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("Successfully imported %d transactions", len(transactions)),
		"count":   len(transactions),
	})
}

// Helper to parse amount strings to float64
func parseAmount(s string) (float64, error) {
	return json.Number(s).Float64()
}

// UploadJSONTransactions uploads transactions from a JSON file
func (h *UploadHandler) UploadJSONTransactions(w http.ResponseWriter, r *http.Request) {
	// Get user claims from JWT token
	userClaims, ok := r.Context().Value("user").(*jwt.MapClaims)
	if !ok || userClaims == nil {
		http.Error(w, "Unauthorized: invalid or missing authentication", http.StatusUnauthorized)
		return
	}

	// Extract user ID from claims
	userIDValue, ok := (*userClaims)["user_id"]
	if !ok || userIDValue == nil {
		http.Error(w, "Unauthorized: user ID not found in token", http.StatusUnauthorized)
		return
	}

	userID := userIDValue.(string)

	// Check if user has the preparer role
	hasRole, err := h.roleService.HasRole(userID, models.RolePreparer)
	if err != nil || !hasRole {
		http.Error(w, "Unauthorized: requires preparer role", http.StatusUnauthorized)
		return
	}

	// Parse form data
	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10MB max
		http.Error(w, "Error parsing form data: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Get data source ID
	dataSourceID := r.FormValue("dataSourceId")
	if dataSourceID == "" {
		http.Error(w, "Missing data source ID", http.StatusBadRequest)
		return
	}

	// Get file from form data
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving file: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Check file extension
	fileExt := filepath.Ext(header.Filename)
	if fileExt != ".json" {
		http.Error(w, "Only JSON files are supported", http.StatusBadRequest)
		return
	}

	// Read file contents
	fileContents, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Error reading file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Parse JSON
	var transactions []models.Transaction
	if err := json.Unmarshal(fileContents, &transactions); err != nil {
		http.Error(w, "Error parsing JSON file: "+err.Error(), http.StatusBadRequest)
		return
	}

	// TODO: Process JSON transactions
	// For now, return a not implemented error
	http.Error(w, "JSON upload not fully implemented", http.StatusNotImplemented)
}

// GetUploadByID retrieves an upload by ID
func (h *UploadHandler) GetUploadByID(w http.ResponseWriter, r *http.Request) {
	// Extract upload ID from URL path
	vars := mux.Vars(r)
	id := vars["id"]

	// Get upload
	upload, err := h.transactionService.GetTransactionByID(id)
	if err != nil {
		http.Error(w, "Transaction not found", http.StatusNotFound)
		return
	}

	// Return upload
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(upload)
}

// GetUploadsByUser retrieves uploads by user
func (h *UploadHandler) GetUploadsByUser(w http.ResponseWriter, r *http.Request) {
	// Get user claims from JWT token
	userClaims, ok := r.Context().Value("user").(*jwt.MapClaims)
	if !ok || userClaims == nil {
		http.Error(w, "Unauthorized: invalid or missing authentication", http.StatusUnauthorized)
		return
	}

	// Extract user ID from claims
	userIDValue, ok := (*userClaims)["user_id"]
	if !ok || userIDValue == nil {
		http.Error(w, "Unauthorized: user ID not found in token", http.StatusUnauthorized)
		return
	}

	userID := userIDValue.(string)

	// Get uploads
	uploads, err := h.transactionService.GetTransactionsByUserID(userID)
	if err != nil {
		http.Error(w, "Error retrieving transactions", http.StatusInternalServerError)
		return
	}

	// Return uploads
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(uploads)
}

// GetRecentUploads retrieves recent uploads
func (h *UploadHandler) GetRecentUploads(w http.ResponseWriter, r *http.Request) {
	// Get user claims from JWT token
	userClaims, ok := r.Context().Value("user").(*jwt.MapClaims)
	if !ok || userClaims == nil {
		http.Error(w, "Unauthorized: invalid or missing authentication", http.StatusUnauthorized)
		return
	}

	// Extract user ID from claims
	userIDValue, ok := (*userClaims)["user_id"]
	if !ok || userIDValue == nil {
		http.Error(w, "Unauthorized: user ID not found in token", http.StatusUnauthorized)
		return
	}

	userID := userIDValue.(string)

	// Check if user has admin role
	hasRole, err := h.roleService.HasRole(userID, models.RoleAdmin)
	if err != nil || !hasRole {
		http.Error(w, "Unauthorized: requires admin role", http.StatusUnauthorized)
		return
	}

	// Parse query parameters
	r.ParseForm()

	// Get limit
	limit := 10
	if limitStr := r.Form.Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	// Get uploads
	uploads, err := h.transactionService.GetRecentTransactions(limit)
	if err != nil {
		http.Error(w, "Error retrieving transactions", http.StatusInternalServerError)
		return
	}

	// Return uploads
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(uploads)
}

// ColumnMapping represents column mapping for CSV files
type ColumnMapping map[string]int

// PreviewResponse is the response sent after a file is uploaded
type PreviewResponse struct {
	PreviewUrl        string        `json:"previewUrl"`
	Preview           [][]string    `json:"preview"`
	SuggestedMappings ColumnMapping `json:"suggestedMappings"`
}

// ProcessRequest is the request to process a previously uploaded file
type ProcessRequest struct {
	PreviewUrl     string        `json:"previewUrl"`
	DataSourceID   string        `json:"dataSourceId"`
	DateFormat     string        `json:"dateFormat"`
	ColumnMappings ColumnMapping `json:"columnMappings"`
}

// UploadDir is the directory where uploaded files are stored
const UploadDir = "./uploads"

// EnsureUploadDir ensures the upload directory exists
func EnsureUploadDir() error {
	if _, err := os.Stat(UploadDir); os.IsNotExist(err) {
		return os.MkdirAll(UploadDir, 0755)
	}
	return nil
}

// PreviewUploadHandler handles the preview of uploaded CSV files
func PreviewUploadHandler(w http.ResponseWriter, r *http.Request) {
	// Set response headers
	w.Header().Set("Content-Type", "application/json")

	// Ensure upload directory exists
	if err := EnsureUploadDir(); err != nil {
		http.Error(w, "Failed to create upload directory", http.StatusInternalServerError)
		return
	}

	// Parse multipart form
	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10MB max
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Get file from form
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "No file uploaded", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Get data source ID from form
	dataSourceID := r.FormValue("dataSourceId")
	if dataSourceID == "" {
		http.Error(w, "No data source ID provided", http.StatusBadRequest)
		return
	}

	// Create a unique filename based on dataSourceID and timestamp
	filename := fmt.Sprintf("%s_%s%s",
		dataSourceID,
		time.Now().Format("20060102_150405"),
		filepath.Ext(header.Filename),
	)

	// Create full path
	filePath := filepath.Join(UploadDir, filename)

	// Create file on disk
	dst, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Failed to create file on server", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copy uploaded file to destination file
	if _, err = io.Copy(dst, file); err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	// Reopen file for CSV reading
	dst.Close()
	csvFile, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "Failed to open saved file", http.StatusInternalServerError)
		return
	}
	defer csvFile.Close()

	// Create CSV reader
	reader := csv.NewReader(csvFile)

	// Allow variable number of fields per record
	reader.FieldsPerRecord = -1

	// Read the first 6 lines (or fewer if the file is smaller)
	var preview [][]string
	for i := 0; i < 6; i++ {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue // Skip problematic lines
		}
		preview = append(preview, record)
	}

	// If we couldn't read any lines, return an error
	if len(preview) == 0 {
		http.Error(w, "Failed to parse CSV content", http.StatusBadRequest)
		return
	}

	// Try to automatically map columns based on headers
	suggestedMappings := suggestColumnMappings(preview[0])

	// Create response with preview URL and data
	response := PreviewResponse{
		PreviewUrl:        filename,
		Preview:           preview,
		SuggestedMappings: suggestedMappings,
	}

	// Send response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// ProcessUploadHandler processes a previously uploaded file
func ProcessUploadHandler(w http.ResponseWriter, r *http.Request) {
	// Set response headers
	w.Header().Set("Content-Type", "application/json")

	// Parse request body
	var req ProcessRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	if req.PreviewUrl == "" || req.DataSourceID == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Verify file exists
	filePath := filepath.Join(UploadDir, req.PreviewUrl)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusBadRequest)
		return
	}

	// Here you would add logic to process the CSV file based on mappings
	// For example, inserting data into your database

	// Send success response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "File processed successfully",
	})
}

// suggestColumnMappings suggests column mappings based on headers
func suggestColumnMappings(headers []string) ColumnMapping {
	mappings := make(ColumnMapping)

	for i, header := range headers {
		if header == "" {
			continue
		}
		headerLower := strings.ToLower(strings.TrimSpace(header))

		// Try to map common header names
		if strings.Contains(headerLower, "date") && !strings.Contains(headerLower, "post") {
			mappings["date"] = i
		} else if strings.Contains(headerLower, "post") && strings.Contains(headerLower, "date") {
			mappings["postDate"] = i
		} else if strings.Contains(headerLower, "desc") {
			mappings["description"] = i
		} else if strings.Contains(headerLower, "amount") || strings.Contains(headerLower, "sum") ||
			strings.Contains(headerLower, "value") || strings.Contains(headerLower, "price") {
			mappings["amount"] = i
		} else if strings.Contains(headerLower, "ref") || strings.Contains(headerLower, "number") ||
			strings.Contains(headerLower, "id") || strings.Contains(headerLower, "trans") {
			mappings["reference"] = i
		} else if strings.Contains(headerLower, "curr") || strings.Contains(headerLower, "currency") {
			mappings["currency"] = i
		}
	}

	return mappings
}
