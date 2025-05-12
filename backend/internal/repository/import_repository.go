package repository

import (
	"backend/internal/db"
	"backend/internal/models"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

var (
	ErrImportNotFound = errors.New("import record not found")
)

// ImportRepository defines operations for managing import records
type ImportRepository interface {
	CreateImport(importRecord *models.ImportRecord) error
	GetImportByID(id string) (*models.ImportRecord, error)
	UpdateImportStatus(id string, status string, rowCount, successCount, errorCount int) error
	GetImportsByDataSource(dataSourceID string, limit, offset int) ([]models.ImportRecord, int, error)
	DeleteImport(id string) error

	// Raw transactions operations
	CreateRawTransaction(rawTx *models.RawTransaction) error
	GetRawTransactionsByImport(importID string, limit, offset int) ([]models.RawTransaction, int, error)
	GetRawTransactionByID(id string) (*models.RawTransaction, error)
}

// PostgresImportRepository implements ImportRepository for PostgreSQL
type PostgresImportRepository struct {
	db *sql.DB
}

// NewImportRepository creates a new import repository
func NewImportRepository() ImportRepository {
	if db.DB == nil {
		// Return a mock repository for development
		return &MockImportRepository{
			imports:         make(map[string]*models.ImportRecord),
			rawTransactions: make(map[string]*models.RawTransaction),
		}
	}
	return &PostgresImportRepository{
		db: db.DB,
	}
}

// CreateImport creates a new import record
func (r *PostgresImportRepository) CreateImport(importRecord *models.ImportRecord) error {
	query := `
		INSERT INTO import_records (data_source_id, file_name, file_size, status, row_count, 
			success_count, error_count, imported_by, metadata)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at, updated_at
	`

	var metadataJSON []byte
	if importRecord.Metadata != nil {
		metadataJSON = importRecord.Metadata
	}

	return r.db.QueryRow(
		query,
		importRecord.DataSourceID,
		importRecord.FileName,
		importRecord.FileSize,
		importRecord.Status,
		importRecord.RowCount,
		importRecord.SuccessCount,
		importRecord.ErrorCount,
		importRecord.ImportedBy,
		metadataJSON,
	).Scan(&importRecord.ID, &importRecord.CreatedAt, &importRecord.UpdatedAt)
}

// GetImportByID retrieves an import record by ID
func (r *PostgresImportRepository) GetImportByID(id string) (*models.ImportRecord, error) {
	query := `
		SELECT id, data_source_id, file_name, file_size, status, row_count, 
			success_count, error_count, imported_by, created_at, updated_at, metadata
		FROM import_records
		WHERE id = $1
	`

	var importRecord models.ImportRecord
	var metadata []byte
	err := r.db.QueryRow(query, id).Scan(
		&importRecord.ID,
		&importRecord.DataSourceID,
		&importRecord.FileName,
		&importRecord.FileSize,
		&importRecord.Status,
		&importRecord.RowCount,
		&importRecord.SuccessCount,
		&importRecord.ErrorCount,
		&importRecord.ImportedBy,
		&importRecord.CreatedAt,
		&importRecord.UpdatedAt,
		&metadata,
	)

	if err == sql.ErrNoRows {
		return nil, ErrImportNotFound
	}

	if err != nil {
		return nil, err
	}

	importRecord.Metadata = metadata
	// Set epoch timestamps
	importRecord.CreatedAtEpoch = importRecord.CreatedAt.UTC().UnixNano() / int64(time.Millisecond)
	importRecord.UpdatedAtEpoch = importRecord.UpdatedAt.UTC().UnixNano() / int64(time.Millisecond)

	return &importRecord, nil
}

// UpdateImportStatus updates the status and counts of an import
func (r *PostgresImportRepository) UpdateImportStatus(id string, status string, rowCount, successCount, errorCount int) error {
	query := `
		UPDATE import_records
		SET status = $2, row_count = $3, success_count = $4, error_count = $5, updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at
	`

	var updatedAt time.Time
	err := r.db.QueryRow(query, id, status, rowCount, successCount, errorCount).Scan(&updatedAt)

	if err == sql.ErrNoRows {
		return ErrImportNotFound
	}

	return err
}

// GetImportsByDataSource retrieves import records for a data source with pagination
func (r *PostgresImportRepository) GetImportsByDataSource(dataSourceID string, limit, offset int) ([]models.ImportRecord, int, error) {
	// Get total count first
	var totalCount int
	countQuery := `SELECT COUNT(*) FROM import_records WHERE data_source_id = $1`
	err := r.db.QueryRow(countQuery, dataSourceID).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	// Then get the paginated results
	query := `
		SELECT id, data_source_id, file_name, file_size, status, row_count, 
			success_count, error_count, imported_by, created_at, updated_at, metadata
		FROM import_records
		WHERE data_source_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, dataSourceID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var imports []models.ImportRecord
	for rows.Next() {
		var importRecord models.ImportRecord
		var metadata []byte
		err := rows.Scan(
			&importRecord.ID,
			&importRecord.DataSourceID,
			&importRecord.FileName,
			&importRecord.FileSize,
			&importRecord.Status,
			&importRecord.RowCount,
			&importRecord.SuccessCount,
			&importRecord.ErrorCount,
			&importRecord.ImportedBy,
			&importRecord.CreatedAt,
			&importRecord.UpdatedAt,
			&metadata,
		)

		if err != nil {
			return nil, 0, err
		}

		importRecord.Metadata = metadata
		// Set epoch timestamps
		importRecord.CreatedAtEpoch = importRecord.CreatedAt.UTC().UnixNano() / int64(time.Millisecond)
		importRecord.UpdatedAtEpoch = importRecord.UpdatedAt.UTC().UnixNano() / int64(time.Millisecond)

		imports = append(imports, importRecord)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return imports, totalCount, nil
}

// DeleteImport deletes an import record and its associated raw transactions
func (r *PostgresImportRepository) DeleteImport(id string) error {
	// Start a transaction
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	// First delete the raw transactions
	_, err = tx.Exec("DELETE FROM raw_transactions WHERE import_id = $1", id)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Then delete the import record
	result, err := tx.Exec("DELETE FROM import_records WHERE id = $1", id)
	if err != nil {
		tx.Rollback()
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		return err
	}

	if rowsAffected == 0 {
		tx.Rollback()
		return ErrImportNotFound
	}

	return tx.Commit()
}

// CreateRawTransaction creates a new raw transaction
func (r *PostgresImportRepository) CreateRawTransaction(rawTx *models.RawTransaction) error {
	query := `
		INSERT INTO raw_transactions (import_id, data_source_id, row_number, data, error_message)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`

	return r.db.QueryRow(
		query,
		rawTx.ImportID,
		rawTx.DataSourceID,
		rawTx.RowNumber,
		rawTx.Data,
		rawTx.ErrorMessage,
	).Scan(&rawTx.ID, &rawTx.CreatedAt)
}

// GetRawTransactionsByImport retrieves raw transactions for an import with pagination
func (r *PostgresImportRepository) GetRawTransactionsByImport(importID string, limit, offset int) ([]models.RawTransaction, int, error) {
	// Get total count first
	var totalCount int
	countQuery := `SELECT COUNT(*) FROM raw_transactions WHERE import_id = $1`
	err := r.db.QueryRow(countQuery, importID).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	// Then get the paginated results
	query := `
		SELECT id, import_id, data_source_id, row_number, data, error_message, created_at
		FROM raw_transactions
		WHERE import_id = $1
		ORDER BY row_number
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, importID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var transactions []models.RawTransaction
	for rows.Next() {
		var tx models.RawTransaction
		err := rows.Scan(
			&tx.ID,
			&tx.ImportID,
			&tx.DataSourceID,
			&tx.RowNumber,
			&tx.Data,
			&tx.ErrorMessage,
			&tx.CreatedAt,
		)

		if err != nil {
			return nil, 0, err
		}

		// Set epoch timestamp
		tx.CreatedAtEpoch = tx.CreatedAt.UTC().UnixNano() / int64(time.Millisecond)

		transactions = append(transactions, tx)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return transactions, totalCount, nil
}

// GetRawTransactionByID retrieves a raw transaction by ID
func (r *PostgresImportRepository) GetRawTransactionByID(id string) (*models.RawTransaction, error) {
	query := `
		SELECT id, import_id, data_source_id, row_number, data, error_message, created_at
		FROM raw_transactions
		WHERE id = $1
	`

	var tx models.RawTransaction
	err := r.db.QueryRow(query, id).Scan(
		&tx.ID,
		&tx.ImportID,
		&tx.DataSourceID,
		&tx.RowNumber,
		&tx.Data,
		&tx.ErrorMessage,
		&tx.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("raw transaction not found: %s", id)
	}

	if err != nil {
		return nil, err
	}

	// Set epoch timestamp
	tx.CreatedAtEpoch = tx.CreatedAt.UTC().UnixNano() / int64(time.Millisecond)

	return &tx, nil
}

// MockImportRepository is a mock implementation for development
type MockImportRepository struct {
	imports         map[string]*models.ImportRecord
	rawTransactions map[string]*models.RawTransaction
	lastID          int
}

// CreateImport creates an import record in the mock repository
func (r *MockImportRepository) CreateImport(importRecord *models.ImportRecord) error {
	r.lastID++
	importRecord.ID = fmt.Sprintf("mock-import-%d", r.lastID)
	importRecord.CreatedAt = time.Now()
	importRecord.UpdatedAt = time.Now()
	r.imports[importRecord.ID] = importRecord
	return nil
}

// GetImportByID retrieves an import record by ID from the mock repository
func (r *MockImportRepository) GetImportByID(id string) (*models.ImportRecord, error) {
	if importRecord, ok := r.imports[id]; ok {
		return importRecord, nil
	}
	return nil, ErrImportNotFound
}

// UpdateImportStatus updates the status of an import in the mock repository
func (r *MockImportRepository) UpdateImportStatus(id string, status string, rowCount, successCount, errorCount int) error {
	if importRecord, ok := r.imports[id]; ok {
		importRecord.Status = status
		importRecord.RowCount = rowCount
		importRecord.SuccessCount = successCount
		importRecord.ErrorCount = errorCount
		importRecord.UpdatedAt = time.Now()
		return nil
	}
	return ErrImportNotFound
}

// GetImportsByDataSource retrieves import records for a data source from the mock repository
func (r *MockImportRepository) GetImportsByDataSource(dataSourceID string, limit, offset int) ([]models.ImportRecord, int, error) {
	var imports []models.ImportRecord
	var filtered []models.ImportRecord

	for _, importRecord := range r.imports {
		if importRecord.DataSourceID == dataSourceID {
			filtered = append(filtered, *importRecord)
		}
	}

	// Sort by created_at desc (simple mock implementation)
	// In a real implementation we would use a proper sorting algorithm
	// This is just a simple mock for development

	total := len(filtered)
	end := offset + limit
	if end > total {
		end = total
	}
	if offset < total {
		imports = filtered[offset:end]
	}

	return imports, total, nil
}

// DeleteImport deletes an import record from the mock repository
func (r *MockImportRepository) DeleteImport(id string) error {
	if _, ok := r.imports[id]; ok {
		delete(r.imports, id)
		// Also delete associated raw transactions
		for txID, tx := range r.rawTransactions {
			if tx.ImportID == id {
				delete(r.rawTransactions, txID)
			}
		}
		return nil
	}
	return ErrImportNotFound
}

// CreateRawTransaction creates a raw transaction in the mock repository
func (r *MockImportRepository) CreateRawTransaction(rawTx *models.RawTransaction) error {
	r.lastID++
	rawTx.ID = fmt.Sprintf("mock-rawtx-%d", r.lastID)
	rawTx.CreatedAt = time.Now()
	r.rawTransactions[rawTx.ID] = rawTx
	return nil
}

// GetRawTransactionsByImport retrieves raw transactions for an import from the mock repository
func (r *MockImportRepository) GetRawTransactionsByImport(importID string, limit, offset int) ([]models.RawTransaction, int, error) {
	var transactions []models.RawTransaction
	var filtered []models.RawTransaction

	for _, tx := range r.rawTransactions {
		if tx.ImportID == importID {
			filtered = append(filtered, *tx)
		}
	}

	// Simple pagination
	total := len(filtered)
	end := offset + limit
	if end > total {
		end = total
	}
	if offset < total {
		transactions = filtered[offset:end]
	}

	return transactions, total, nil
}

// GetRawTransactionByID retrieves a raw transaction by ID from the mock repository
func (r *MockImportRepository) GetRawTransactionByID(id string) (*models.RawTransaction, error) {
	if tx, ok := r.rawTransactions[id]; ok {
		return tx, nil
	}
	return nil, fmt.Errorf("raw transaction not found: %s", id)
}
