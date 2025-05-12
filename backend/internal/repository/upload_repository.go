package repository

import (
	"backend/internal/db"
	"backend/internal/models"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrUploadNotFound = errors.New("transaction upload not found")
)

// UploadRepository defines operations for managing transaction uploads
type UploadRepository interface {
	CreateUpload(upload *models.TransactionUpload) error
	GetUploadByID(id string) (*models.TransactionUpload, error)
	UpdateUploadStatus(id string, status string, recordCount int, errorMessage string) error
	GetUploadsByUser(userID string) ([]models.TransactionUpload, error)
	GetRecentUploads(limit int) ([]models.TransactionUpload, error)
	GetUploadsByDataSource(dataSourceID string) ([]models.TransactionUpload, error)
}

// PostgresUploadRepository implements UploadRepository for PostgreSQL
type PostgresUploadRepository struct {
	db *sql.DB
}

// NewUploadRepository creates a new upload repository
func NewUploadRepository() UploadRepository {
	if db.DB == nil {
		// Return a mock repository for development
		return &MockUploadRepository{
			uploads: make(map[string]*models.TransactionUpload),
		}
	}
	return &PostgresUploadRepository{
		db: db.DB,
	}
}

// CreateUpload creates a new transaction upload
func (r *PostgresUploadRepository) CreateUpload(upload *models.TransactionUpload) error {
	query := `
		INSERT INTO transaction_uploads (
			data_source_id, file_name, file_size, uploaded_by, 
			status, record_count, error_message
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7
		) RETURNING id, upload_date
	`

	return r.db.QueryRow(
		query,
		upload.DataSourceID,
		upload.FileName,
		upload.FileSize,
		upload.UploadedBy,
		upload.Status,
		upload.RecordCount,
		upload.ErrorMessage,
	).Scan(&upload.ID, &upload.UploadDate)
}

// GetUploadByID retrieves a transaction upload by ID
func (r *PostgresUploadRepository) GetUploadByID(id string) (*models.TransactionUpload, error) {
	query := `
		SELECT 
			id, data_source_id, file_name, file_size, uploaded_by, 
			upload_date, status, record_count, error_message
		FROM transaction_uploads
		WHERE id = $1
	`

	var upload models.TransactionUpload
	err := r.db.QueryRow(query, id).Scan(
		&upload.ID,
		&upload.DataSourceID,
		&upload.FileName,
		&upload.FileSize,
		&upload.UploadedBy,
		&upload.UploadDate,
		&upload.Status,
		&upload.RecordCount,
		&upload.ErrorMessage,
	)

	if err == sql.ErrNoRows {
		return nil, ErrUploadNotFound
	}

	if err != nil {
		return nil, err
	}

	return &upload, nil
}

// UpdateUploadStatus updates a transaction upload's status
func (r *PostgresUploadRepository) UpdateUploadStatus(id string, status string, recordCount int, errorMessage string) error {
	query := `
		UPDATE transaction_uploads
		SET status = $1, record_count = $2, error_message = $3
		WHERE id = $4
	`

	result, err := r.db.Exec(query, status, recordCount, errorMessage, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrUploadNotFound
	}

	return nil
}

// GetUploadsByUser retrieves transaction uploads by user
func (r *PostgresUploadRepository) GetUploadsByUser(userID string) ([]models.TransactionUpload, error) {
	query := `
		SELECT 
			id, data_source_id, file_name, file_size, uploaded_by, 
			upload_date, status, record_count, error_message
		FROM transaction_uploads
		WHERE uploaded_by = $1
		ORDER BY upload_date DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var uploads []models.TransactionUpload
	for rows.Next() {
		var upload models.TransactionUpload
		err := rows.Scan(
			&upload.ID,
			&upload.DataSourceID,
			&upload.FileName,
			&upload.FileSize,
			&upload.UploadedBy,
			&upload.UploadDate,
			&upload.Status,
			&upload.RecordCount,
			&upload.ErrorMessage,
		)

		if err != nil {
			return nil, err
		}

		uploads = append(uploads, upload)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return uploads, nil
}

// GetRecentUploads retrieves recent transaction uploads
func (r *PostgresUploadRepository) GetRecentUploads(limit int) ([]models.TransactionUpload, error) {
	query := `
		SELECT 
			id, data_source_id, file_name, file_size, uploaded_by, 
			upload_date, status, record_count, error_message
		FROM transaction_uploads
		ORDER BY upload_date DESC
		LIMIT $1
	`

	rows, err := r.db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var uploads []models.TransactionUpload
	for rows.Next() {
		var upload models.TransactionUpload
		err := rows.Scan(
			&upload.ID,
			&upload.DataSourceID,
			&upload.FileName,
			&upload.FileSize,
			&upload.UploadedBy,
			&upload.UploadDate,
			&upload.Status,
			&upload.RecordCount,
			&upload.ErrorMessage,
		)

		if err != nil {
			return nil, err
		}

		uploads = append(uploads, upload)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return uploads, nil
}

// GetUploadsByDataSource retrieves transaction uploads by data source
func (r *PostgresUploadRepository) GetUploadsByDataSource(dataSourceID string) ([]models.TransactionUpload, error) {
	query := `
		SELECT 
			id, data_source_id, file_name, file_size, uploaded_by, 
			upload_date, status, record_count, error_message
		FROM transaction_uploads
		WHERE data_source_id = $1
		ORDER BY upload_date DESC
	`

	rows, err := r.db.Query(query, dataSourceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var uploads []models.TransactionUpload
	for rows.Next() {
		var upload models.TransactionUpload
		err := rows.Scan(
			&upload.ID,
			&upload.DataSourceID,
			&upload.FileName,
			&upload.FileSize,
			&upload.UploadedBy,
			&upload.UploadDate,
			&upload.Status,
			&upload.RecordCount,
			&upload.ErrorMessage,
		)

		if err != nil {
			return nil, err
		}

		uploads = append(uploads, upload)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return uploads, nil
}

// MockUploadRepository is a mock implementation for development
type MockUploadRepository struct {
	uploads map[string]*models.TransactionUpload
}

// CreateUpload creates a transaction upload in the mock repository
func (r *MockUploadRepository) CreateUpload(upload *models.TransactionUpload) error {
	if upload.ID == "" {
		// Generate a dummy ID - in real implementation we'd use UUID
		upload.ID = "mock-" + time.Now().Format("20060102150405")
	}
	upload.UploadDate = time.Now()
	r.uploads[upload.ID] = upload
	return nil
}

// GetUploadByID retrieves a transaction upload by ID from the mock repository
func (r *MockUploadRepository) GetUploadByID(id string) (*models.TransactionUpload, error) {
	if upload, exists := r.uploads[id]; exists {
		return upload, nil
	}
	return nil, ErrUploadNotFound
}

// UpdateUploadStatus updates a transaction upload's status in the mock repository
func (r *MockUploadRepository) UpdateUploadStatus(id string, status string, recordCount int, errorMessage string) error {
	upload, exists := r.uploads[id]
	if !exists {
		return ErrUploadNotFound
	}

	upload.Status = status
	upload.RecordCount = recordCount
	upload.ErrorMessage = errorMessage

	return nil
}

// GetUploadsByUser retrieves transaction uploads by user from the mock repository
func (r *MockUploadRepository) GetUploadsByUser(userID string) ([]models.TransactionUpload, error) {
	var uploads []models.TransactionUpload
	for _, upload := range r.uploads {
		if upload.UploadedBy == userID {
			uploads = append(uploads, *upload)
		}
	}

	// Sort uploads by upload date (newest first)
	// In a real implementation, we would use a proper sorting function

	return uploads, nil
}

// GetRecentUploads retrieves recent transaction uploads from the mock repository
func (r *MockUploadRepository) GetRecentUploads(limit int) ([]models.TransactionUpload, error) {
	var uploads []models.TransactionUpload
	for _, upload := range r.uploads {
		uploads = append(uploads, *upload)
	}

	// Sort uploads by upload date (newest first)
	// In a real implementation, we would use a proper sorting function

	// Apply limit
	if len(uploads) > limit {
		uploads = uploads[:limit]
	}

	return uploads, nil
}

// GetUploadsByDataSource retrieves transaction uploads by data source from the mock repository
func (r *MockUploadRepository) GetUploadsByDataSource(dataSourceID string) ([]models.TransactionUpload, error) {
	var uploads []models.TransactionUpload
	for _, upload := range r.uploads {
		if upload.DataSourceID == dataSourceID {
			uploads = append(uploads, *upload)
		}
	}

	// Sort uploads by upload date (newest first)
	// In a real implementation, we would use a proper sorting function

	return uploads, nil
}
