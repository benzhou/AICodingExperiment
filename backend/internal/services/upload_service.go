package services

import (
	"backend/internal/models"
	"backend/internal/repository"
	"encoding/csv"
	"errors"
	"io"
	"strconv"
	"strings"
	"time"
)

// UploadService provides methods for transaction upload operations
type UploadService struct {
	uploadRepo      repository.UploadRepository
	transactionRepo repository.TransactionRepository
	dataSourceRepo  repository.DataSourceRepository
}

// NewUploadService creates a new upload service
func NewUploadService(
	uploadRepo repository.UploadRepository,
	transactionRepo repository.TransactionRepository,
	dataSourceRepo repository.DataSourceRepository,
) *UploadService {
	return &UploadService{
		uploadRepo:      uploadRepo,
		transactionRepo: transactionRepo,
		dataSourceRepo:  dataSourceRepo,
	}
}

// UploadCSV processes a CSV file and creates transactions
func (s *UploadService) UploadCSV(
	dataSourceID string,
	fileName string,
	fileSize int64,
	userID string,
	reader io.Reader,
	columnMapping map[string]int,
	dateFormat string,
) (*models.TransactionUpload, error) {
	// Create upload record
	upload := &models.TransactionUpload{
		DataSourceID: dataSourceID,
		FileName:     fileName,
		FileSize:     fileSize,
		UploadedBy:   userID,
		Status:       "Processing",
		RecordCount:  0,
	}

	err := s.uploadRepo.CreateUpload(upload)
	if err != nil {
		return nil, err
	}

	// Process the CSV file
	csvReader := csv.NewReader(reader)

	// Skip the header row
	_, err = csvReader.Read()
	if err != nil {
		s.updateUploadStatus(upload.ID, "Failed", 0, "Failed to read CSV header: "+err.Error())
		return upload, err
	}

	// Read records and create transactions
	var transactions []models.Transaction
	recordCount := 0

	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			s.updateUploadStatus(upload.ID, "Failed", recordCount, "Error reading CSV: "+err.Error())
			return upload, err
		}

		transaction, err := s.parseTransaction(record, columnMapping, dateFormat, dataSourceID, userID)
		if err != nil {
			// Log error but continue processing
			continue
		}

		transactions = append(transactions, *transaction)
		recordCount++
	}

	// Save transactions to the database
	for _, transaction := range transactions {
		err := s.transactionRepo.CreateTransaction(&transaction)
		if err != nil {
			// Log error but continue processing
			continue
		}
	}

	// Update upload status
	s.updateUploadStatus(upload.ID, "Completed", recordCount, "")

	// Update the upload with the final status
	upload.Status = "Completed"
	upload.RecordCount = recordCount

	return upload, nil
}

// updateUploadStatus updates the status of an upload
func (s *UploadService) updateUploadStatus(id string, status string, recordCount int, errorMessage string) error {
	return s.uploadRepo.UpdateUploadStatus(id, status, recordCount, errorMessage)
}

// parseTransaction parses a CSV record into a Transaction
func (s *UploadService) parseTransaction(
	record []string,
	columnMapping map[string]int,
	dateFormat string,
	dataSourceID string,
	userID string,
) (*models.Transaction, error) {
	transaction := &models.Transaction{
		DataSourceID: dataSourceID,
		Status:       "Unmatched",
		CreatedBy:    userID,
	}

	// Parse date
	if dateCol, ok := columnMapping["date"]; ok && dateCol < len(record) {
		date, err := time.Parse(dateFormat, strings.TrimSpace(record[dateCol]))
		if err != nil {
			return nil, errors.New("invalid date format: " + err.Error())
		}
		transaction.TransactionDate = date
	} else {
		return nil, errors.New("missing date column")
	}

	// Parse post date (if provided, otherwise use transaction date)
	if postDateCol, ok := columnMapping["postDate"]; ok && postDateCol < len(record) {
		date, err := time.Parse(dateFormat, strings.TrimSpace(record[postDateCol]))
		if err == nil {
			transaction.PostDate = date
		} else {
			transaction.PostDate = transaction.TransactionDate
		}
	} else {
		transaction.PostDate = transaction.TransactionDate
	}

	// Parse description
	if descCol, ok := columnMapping["description"]; ok && descCol < len(record) {
		transaction.Description = strings.TrimSpace(record[descCol])
	}

	// Parse reference
	if refCol, ok := columnMapping["reference"]; ok && refCol < len(record) {
		transaction.Reference = strings.TrimSpace(record[refCol])
	}

	// Parse amount
	if amountCol, ok := columnMapping["amount"]; ok && amountCol < len(record) {
		amountStr := strings.TrimSpace(record[amountCol])
		// Remove currency symbols and commas
		amountStr = strings.ReplaceAll(amountStr, "$", "")
		amountStr = strings.ReplaceAll(amountStr, ",", "")
		amount, err := strconv.ParseFloat(amountStr, 64)
		if err != nil {
			return nil, errors.New("invalid amount format: " + err.Error())
		}
		transaction.Amount = amount
	} else {
		return nil, errors.New("missing amount column")
	}

	// Parse currency
	if currencyCol, ok := columnMapping["currency"]; ok && currencyCol < len(record) {
		transaction.Currency = strings.TrimSpace(record[currencyCol])
	} else {
		// Default currency
		transaction.Currency = "USD"
	}

	return transaction, nil
}

// GetUploadByID retrieves an upload by ID
func (s *UploadService) GetUploadByID(id string) (*models.TransactionUpload, error) {
	return s.uploadRepo.GetUploadByID(id)
}

// GetUploadsByUser retrieves uploads by user
func (s *UploadService) GetUploadsByUser(userID string) ([]models.TransactionUpload, error) {
	return s.uploadRepo.GetUploadsByUser(userID)
}

// GetRecentUploads retrieves recent uploads
func (s *UploadService) GetRecentUploads(limit int) ([]models.TransactionUpload, error) {
	return s.uploadRepo.GetRecentUploads(limit)
}
