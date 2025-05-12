package repository

import (
	"backend/internal/db"
	"backend/internal/models"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrTransactionNotFound = errors.New("transaction not found")
)

// TransactionRepository defines operations for managing transactions
type TransactionRepository interface {
	CreateTransaction(transaction *models.Transaction) error
	CreateTransactions(transactions []models.Transaction) error
	GetTransactionByID(id string) (*models.Transaction, error)
	GetTransactionsByDataSourceID(dataSourceID string) ([]models.Transaction, error)
	GetTransactionsByUserID(userID string) ([]models.Transaction, error)
	GetRecentTransactions(limit int) ([]models.Transaction, error)
	DeleteTransaction(id string) error
	DeleteTransactionsByDataSourceID(dataSourceID string) error
}

// PostgresTransactionRepository implements TransactionRepository for PostgreSQL
type PostgresTransactionRepository struct {
	db *sql.DB
}

// NewTransactionRepository creates a new transaction repository
func NewTransactionRepository() TransactionRepository {
	if db.DB == nil {
		// Return a mock repository for development
		return &MockTransactionRepository{
			transactions: make(map[string]*models.Transaction),
		}
	}
	return &PostgresTransactionRepository{
		db: db.DB,
	}
}

// CreateTransaction creates a new transaction in the database
func (r *PostgresTransactionRepository) CreateTransaction(transaction *models.Transaction) error {
	query := `
		INSERT INTO transactions (
			id, data_source_id, transaction_date, post_date, 
			description, amount, currency, reference,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $9)
	`

	_, err := r.db.Exec(
		query,
		transaction.ID,
		transaction.DataSourceID,
		transaction.TransactionDate,
		transaction.PostDate,
		transaction.Description,
		transaction.Amount,
		transaction.Currency,
		transaction.Reference,
		time.Now(),
	)

	return err
}

// CreateTransactions creates multiple transactions in a batch
func (r *PostgresTransactionRepository) CreateTransactions(transactions []models.Transaction) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	// Prepare the statement for efficient batch insertion
	stmt, err := tx.Prepare(`
		INSERT INTO transactions (
			id, data_source_id, transaction_date, post_date, 
			description, amount, currency, reference,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $9)
	`)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	now := time.Now()
	for _, transaction := range transactions {
		_, err := stmt.Exec(
			transaction.ID,
			transaction.DataSourceID,
			transaction.TransactionDate,
			transaction.PostDate,
			transaction.Description,
			transaction.Amount,
			transaction.Currency,
			transaction.Reference,
			now,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

// GetTransactionByID retrieves a transaction by ID
func (r *PostgresTransactionRepository) GetTransactionByID(id string) (*models.Transaction, error) {
	query := `
		SELECT id, data_source_id, transaction_date, post_date, 
			   description, amount, currency, reference,
			   created_at, updated_at
		FROM transactions
		WHERE id = $1
	`

	var transaction models.Transaction
	err := r.db.QueryRow(query, id).Scan(
		&transaction.ID,
		&transaction.DataSourceID,
		&transaction.TransactionDate,
		&transaction.PostDate,
		&transaction.Description,
		&transaction.Amount,
		&transaction.Currency,
		&transaction.Reference,
		&transaction.CreatedAt,
		&transaction.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrTransactionNotFound
	}

	if err != nil {
		return nil, err
	}

	return &transaction, nil
}

// GetTransactionsByDataSourceID retrieves transactions for a data source
func (r *PostgresTransactionRepository) GetTransactionsByDataSourceID(dataSourceID string) ([]models.Transaction, error) {
	query := `
		SELECT id, data_source_id, transaction_date, post_date, 
			   description, amount, currency, reference,
			   created_at, updated_at
		FROM transactions
		WHERE data_source_id = $1
		ORDER BY transaction_date DESC
	`

	rows, err := r.db.Query(query, dataSourceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var transaction models.Transaction
		err := rows.Scan(
			&transaction.ID,
			&transaction.DataSourceID,
			&transaction.TransactionDate,
			&transaction.PostDate,
			&transaction.Description,
			&transaction.Amount,
			&transaction.Currency,
			&transaction.Reference,
			&transaction.CreatedAt,
			&transaction.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		transactions = append(transactions, transaction)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}

// GetTransactionsByUserID retrieves transactions uploaded by a user
func (r *PostgresTransactionRepository) GetTransactionsByUserID(userID string) ([]models.Transaction, error) {
	query := `
		SELECT t.id, t.data_source_id, t.transaction_date, t.post_date, 
			   t.description, t.amount, t.currency, t.reference,
			   t.created_at, t.updated_at
		FROM transactions t
		JOIN data_sources ds ON t.data_source_id = ds.id
		WHERE ds.created_by = $1
		ORDER BY t.created_at DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var transaction models.Transaction
		err := rows.Scan(
			&transaction.ID,
			&transaction.DataSourceID,
			&transaction.TransactionDate,
			&transaction.PostDate,
			&transaction.Description,
			&transaction.Amount,
			&transaction.Currency,
			&transaction.Reference,
			&transaction.CreatedAt,
			&transaction.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		transactions = append(transactions, transaction)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}

// GetRecentTransactions retrieves recent transactions up to a limit
func (r *PostgresTransactionRepository) GetRecentTransactions(limit int) ([]models.Transaction, error) {
	query := `
		SELECT id, data_source_id, transaction_date, post_date, 
			   description, amount, currency, reference,
			   created_at, updated_at
		FROM transactions
		ORDER BY created_at DESC
		LIMIT $1
	`

	rows, err := r.db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var transaction models.Transaction
		err := rows.Scan(
			&transaction.ID,
			&transaction.DataSourceID,
			&transaction.TransactionDate,
			&transaction.PostDate,
			&transaction.Description,
			&transaction.Amount,
			&transaction.Currency,
			&transaction.Reference,
			&transaction.CreatedAt,
			&transaction.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		transactions = append(transactions, transaction)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}

// DeleteTransaction deletes a transaction
func (r *PostgresTransactionRepository) DeleteTransaction(id string) error {
	query := "DELETE FROM transactions WHERE id = $1"
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrTransactionNotFound
	}

	return nil
}

// DeleteTransactionsByDataSourceID deletes all transactions for a data source
func (r *PostgresTransactionRepository) DeleteTransactionsByDataSourceID(dataSourceID string) error {
	query := "DELETE FROM transactions WHERE data_source_id = $1"
	_, err := r.db.Exec(query, dataSourceID)
	return err
}

// MockTransactionRepository implements TransactionRepository for testing/development
type MockTransactionRepository struct {
	transactions map[string]*models.Transaction
}

// CreateTransaction creates a transaction in the mock repository
func (r *MockTransactionRepository) CreateTransaction(transaction *models.Transaction) error {
	if transaction.ID == "" {
		transaction.ID = uuid.New().String()
	}

	now := time.Now()
	transaction.CreatedAt = now
	transaction.UpdatedAt = now

	r.transactions[transaction.ID] = transaction
	return nil
}

// CreateTransactions creates multiple transactions in the mock repository
func (r *MockTransactionRepository) CreateTransactions(transactions []models.Transaction) error {
	for i := range transactions {
		if transactions[i].ID == "" {
			transactions[i].ID = uuid.New().String()
		}

		now := time.Now()
		transactions[i].CreatedAt = now
		transactions[i].UpdatedAt = now

		r.transactions[transactions[i].ID] = &transactions[i]
	}
	return nil
}

// GetTransactionByID retrieves a transaction by ID from the mock repository
func (r *MockTransactionRepository) GetTransactionByID(id string) (*models.Transaction, error) {
	transaction, exists := r.transactions[id]
	if !exists {
		return nil, ErrTransactionNotFound
	}
	return transaction, nil
}

// GetTransactionsByDataSourceID retrieves transactions for a data source from the mock repository
func (r *MockTransactionRepository) GetTransactionsByDataSourceID(dataSourceID string) ([]models.Transaction, error) {
	var transactions []models.Transaction
	for _, transaction := range r.transactions {
		if transaction.DataSourceID == dataSourceID {
			transactions = append(transactions, *transaction)
		}
	}
	return transactions, nil
}

// GetTransactionsByUserID retrieves transactions uploaded by a user from the mock repository
func (r *MockTransactionRepository) GetTransactionsByUserID(userID string) ([]models.Transaction, error) {
	// In a real implementation, this would join with data_sources to filter by created_by
	// For mock, we'll simply return all transactions
	var transactions []models.Transaction
	for _, transaction := range r.transactions {
		transactions = append(transactions, *transaction)
	}
	return transactions, nil
}

// GetRecentTransactions retrieves recent transactions from the mock repository
func (r *MockTransactionRepository) GetRecentTransactions(limit int) ([]models.Transaction, error) {
	var transactions []models.Transaction
	for _, transaction := range r.transactions {
		transactions = append(transactions, *transaction)
		if len(transactions) >= limit {
			break
		}
	}
	return transactions, nil
}

// DeleteTransaction deletes a transaction from the mock repository
func (r *MockTransactionRepository) DeleteTransaction(id string) error {
	if _, exists := r.transactions[id]; !exists {
		return ErrTransactionNotFound
	}
	delete(r.transactions, id)
	return nil
}

// DeleteTransactionsByDataSourceID deletes all transactions for a data source from the mock repository
func (r *MockTransactionRepository) DeleteTransactionsByDataSourceID(dataSourceID string) error {
	for id, transaction := range r.transactions {
		if transaction.DataSourceID == dataSourceID {
			delete(r.transactions, id)
		}
	}
	return nil
}
