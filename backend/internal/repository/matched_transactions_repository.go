package repository

import (
	"backend/internal/db"
	"backend/internal/models"
	"database/sql"
	"errors"
	"time"
)

// MatchedTransactionRepository defines operations for managing matched transactions
type MatchedTransactionRepository interface {
	CreateMatchedTransaction(matchedTx *models.MatchedTransaction) error
	GetMatchedTransactionsByMatchSet(matchSetID string, limit, offset int) ([]models.MatchedTransaction, int, error)
	GetMatchedTransactionsByTenant(tenantID string, limit, offset int) ([]models.MatchedTransaction, int, error)
	GetMatchedTransactionByID(id string) (*models.MatchedTransaction, error)
	GetMatchedTransactionByTransactionID(transactionID string) (*models.MatchedTransaction, error)
	GetMatchedTransactionsByMatchGroup(matchGroupID string) ([]models.MatchedTransaction, error)
}

// UnmatchedTransactionRepository defines operations for managing unmatched transactions
type UnmatchedTransactionRepository interface {
	CreateUnmatchedTransaction(unmatchedTx *models.UnmatchedTransaction) error
	GetUnmatchedTransactionsByMatchSet(matchSetID string, limit, offset int) ([]models.UnmatchedTransaction, int, error)
	GetUnmatchedTransactionsByTenant(tenantID string, limit, offset int) ([]models.UnmatchedTransaction, int, error)
	GetUnmatchedTransactionByID(id string) (*models.UnmatchedTransaction, error)
}

// PostgresMatchedTransactionRepository implements MatchedTransactionRepository for PostgreSQL
type PostgresMatchedTransactionRepository struct {
	db *sql.DB
}

// PostgresUnmatchedTransactionRepository implements UnmatchedTransactionRepository for PostgreSQL
type PostgresUnmatchedTransactionRepository struct {
	db *sql.DB
}

// NewMatchedTransactionRepository creates a new matched transaction repository
func NewMatchedTransactionRepository() MatchedTransactionRepository {
	if db.DB == nil {
		// Return a mock repository for development
		return &MockMatchedTransactionRepository{
			matchedTransactions: make(map[string]*models.MatchedTransaction),
		}
	}
	return &PostgresMatchedTransactionRepository{
		db: db.DB,
	}
}

// NewUnmatchedTransactionRepository creates a new unmatched transaction repository
func NewUnmatchedTransactionRepository() UnmatchedTransactionRepository {
	if db.DB == nil {
		// Return a mock repository for development
		return &MockUnmatchedTransactionRepository{
			unmatchedTransactions: make(map[string]*models.UnmatchedTransaction),
		}
	}
	return &PostgresUnmatchedTransactionRepository{
		db: db.DB,
	}
}

// CreateMatchedTransaction creates a new matched transaction
func (r *PostgresMatchedTransactionRepository) CreateMatchedTransaction(matchedTx *models.MatchedTransaction) error {
	query := `
		INSERT INTO matched_transactions (match_set_id, transaction_id, match_group_id, tenant_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`

	return r.db.QueryRow(
		query,
		matchedTx.MatchSetID,
		matchedTx.TransactionID,
		matchedTx.MatchGroupID,
		matchedTx.TenantID,
	).Scan(&matchedTx.ID, &matchedTx.CreatedAt)
}

// GetMatchedTransactionsByMatchSet retrieves matched transactions for a match set with pagination
func (r *PostgresMatchedTransactionRepository) GetMatchedTransactionsByMatchSet(matchSetID string, limit, offset int) ([]models.MatchedTransaction, int, error) {
	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM matched_transactions
		WHERE match_set_id = $1
	`
	err := r.db.QueryRow(countQuery, matchSetID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get matched transactions
	query := `
		SELECT id, match_set_id, transaction_id, match_group_id, tenant_id, created_at
		FROM matched_transactions
		WHERE match_set_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, matchSetID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var matchedTxs []models.MatchedTransaction
	for rows.Next() {
		var matchedTx models.MatchedTransaction
		err := rows.Scan(
			&matchedTx.ID,
			&matchedTx.MatchSetID,
			&matchedTx.TransactionID,
			&matchedTx.MatchGroupID,
			&matchedTx.TenantID,
			&matchedTx.CreatedAt,
		)

		if err != nil {
			return nil, 0, err
		}

		matchedTxs = append(matchedTxs, matchedTx)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return matchedTxs, total, nil
}

// GetMatchedTransactionsByTenant retrieves matched transactions for a tenant with pagination
func (r *PostgresMatchedTransactionRepository) GetMatchedTransactionsByTenant(tenantID string, limit, offset int) ([]models.MatchedTransaction, int, error) {
	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM matched_transactions
		WHERE tenant_id = $1
	`
	err := r.db.QueryRow(countQuery, tenantID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get matched transactions
	query := `
		SELECT id, match_set_id, transaction_id, match_group_id, tenant_id, created_at
		FROM matched_transactions
		WHERE tenant_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, tenantID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var matchedTxs []models.MatchedTransaction
	for rows.Next() {
		var matchedTx models.MatchedTransaction
		err := rows.Scan(
			&matchedTx.ID,
			&matchedTx.MatchSetID,
			&matchedTx.TransactionID,
			&matchedTx.MatchGroupID,
			&matchedTx.TenantID,
			&matchedTx.CreatedAt,
		)

		if err != nil {
			return nil, 0, err
		}

		matchedTxs = append(matchedTxs, matchedTx)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return matchedTxs, total, nil
}

// GetMatchedTransactionByID retrieves a matched transaction by ID
func (r *PostgresMatchedTransactionRepository) GetMatchedTransactionByID(id string) (*models.MatchedTransaction, error) {
	query := `
		SELECT id, match_set_id, transaction_id, match_group_id, tenant_id, created_at
		FROM matched_transactions
		WHERE id = $1
	`

	var matchedTx models.MatchedTransaction
	err := r.db.QueryRow(query, id).Scan(
		&matchedTx.ID,
		&matchedTx.MatchSetID,
		&matchedTx.TransactionID,
		&matchedTx.MatchGroupID,
		&matchedTx.TenantID,
		&matchedTx.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("matched transaction not found")
	}

	if err != nil {
		return nil, err
	}

	return &matchedTx, nil
}

// GetMatchedTransactionByTransactionID retrieves a matched transaction by transaction ID
func (r *PostgresMatchedTransactionRepository) GetMatchedTransactionByTransactionID(transactionID string) (*models.MatchedTransaction, error) {
	query := `
		SELECT id, match_set_id, transaction_id, match_group_id, tenant_id, created_at
		FROM matched_transactions
		WHERE transaction_id = $1
	`

	var matchedTx models.MatchedTransaction
	err := r.db.QueryRow(query, transactionID).Scan(
		&matchedTx.ID,
		&matchedTx.MatchSetID,
		&matchedTx.TransactionID,
		&matchedTx.MatchGroupID,
		&matchedTx.TenantID,
		&matchedTx.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("matched transaction not found")
	}

	if err != nil {
		return nil, err
	}

	return &matchedTx, nil
}

// GetMatchedTransactionsByMatchGroup retrieves matched transactions by match group ID
func (r *PostgresMatchedTransactionRepository) GetMatchedTransactionsByMatchGroup(matchGroupID string) ([]models.MatchedTransaction, error) {
	query := `
		SELECT id, match_set_id, transaction_id, match_group_id, tenant_id, created_at
		FROM matched_transactions
		WHERE match_group_id = $1
	`

	rows, err := r.db.Query(query, matchGroupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var matchedTxs []models.MatchedTransaction
	for rows.Next() {
		var matchedTx models.MatchedTransaction
		err := rows.Scan(
			&matchedTx.ID,
			&matchedTx.MatchSetID,
			&matchedTx.TransactionID,
			&matchedTx.MatchGroupID,
			&matchedTx.TenantID,
			&matchedTx.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		matchedTxs = append(matchedTxs, matchedTx)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return matchedTxs, nil
}

// CreateUnmatchedTransaction creates a new unmatched transaction
func (r *PostgresUnmatchedTransactionRepository) CreateUnmatchedTransaction(unmatchedTx *models.UnmatchedTransaction) error {
	query := `
		INSERT INTO unmatched_transactions (match_set_id, transaction_id, reason, tenant_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`

	return r.db.QueryRow(
		query,
		unmatchedTx.MatchSetID,
		unmatchedTx.TransactionID,
		unmatchedTx.Reason,
		unmatchedTx.TenantID,
	).Scan(&unmatchedTx.ID, &unmatchedTx.CreatedAt)
}

// GetUnmatchedTransactionsByMatchSet retrieves unmatched transactions for a match set with pagination
func (r *PostgresUnmatchedTransactionRepository) GetUnmatchedTransactionsByMatchSet(matchSetID string, limit, offset int) ([]models.UnmatchedTransaction, int, error) {
	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM unmatched_transactions
		WHERE match_set_id = $1
	`
	err := r.db.QueryRow(countQuery, matchSetID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get unmatched transactions
	query := `
		SELECT id, match_set_id, transaction_id, reason, tenant_id, created_at
		FROM unmatched_transactions
		WHERE match_set_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, matchSetID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var unmatchedTxs []models.UnmatchedTransaction
	for rows.Next() {
		var unmatchedTx models.UnmatchedTransaction
		err := rows.Scan(
			&unmatchedTx.ID,
			&unmatchedTx.MatchSetID,
			&unmatchedTx.TransactionID,
			&unmatchedTx.Reason,
			&unmatchedTx.TenantID,
			&unmatchedTx.CreatedAt,
		)

		if err != nil {
			return nil, 0, err
		}

		unmatchedTxs = append(unmatchedTxs, unmatchedTx)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return unmatchedTxs, total, nil
}

// GetUnmatchedTransactionsByTenant retrieves unmatched transactions for a tenant with pagination
func (r *PostgresUnmatchedTransactionRepository) GetUnmatchedTransactionsByTenant(tenantID string, limit, offset int) ([]models.UnmatchedTransaction, int, error) {
	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM unmatched_transactions
		WHERE tenant_id = $1
	`
	err := r.db.QueryRow(countQuery, tenantID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get unmatched transactions
	query := `
		SELECT id, match_set_id, transaction_id, reason, tenant_id, created_at
		FROM unmatched_transactions
		WHERE tenant_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, tenantID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var unmatchedTxs []models.UnmatchedTransaction
	for rows.Next() {
		var unmatchedTx models.UnmatchedTransaction
		err := rows.Scan(
			&unmatchedTx.ID,
			&unmatchedTx.MatchSetID,
			&unmatchedTx.TransactionID,
			&unmatchedTx.Reason,
			&unmatchedTx.TenantID,
			&unmatchedTx.CreatedAt,
		)

		if err != nil {
			return nil, 0, err
		}

		unmatchedTxs = append(unmatchedTxs, unmatchedTx)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return unmatchedTxs, total, nil
}

// GetUnmatchedTransactionByID retrieves an unmatched transaction by ID
func (r *PostgresUnmatchedTransactionRepository) GetUnmatchedTransactionByID(id string) (*models.UnmatchedTransaction, error) {
	query := `
		SELECT id, match_set_id, transaction_id, reason, tenant_id, created_at
		FROM unmatched_transactions
		WHERE id = $1
	`

	var unmatchedTx models.UnmatchedTransaction
	err := r.db.QueryRow(query, id).Scan(
		&unmatchedTx.ID,
		&unmatchedTx.MatchSetID,
		&unmatchedTx.TransactionID,
		&unmatchedTx.Reason,
		&unmatchedTx.TenantID,
		&unmatchedTx.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("unmatched transaction not found")
	}

	if err != nil {
		return nil, err
	}

	return &unmatchedTx, nil
}

// MockMatchedTransactionRepository is a mock implementation for development
type MockMatchedTransactionRepository struct {
	matchedTransactions map[string]*models.MatchedTransaction // ID -> MatchedTransaction
}

// CreateMatchedTransaction creates a matched transaction in the mock repository
func (r *MockMatchedTransactionRepository) CreateMatchedTransaction(matchedTx *models.MatchedTransaction) error {
	if matchedTx.ID == "" {
		// Generate a dummy ID - in real implementation we'd use UUID
		matchedTx.ID = "mock-" + time.Now().Format("20060102150405")
	}
	matchedTx.CreatedAt = time.Now()
	r.matchedTransactions[matchedTx.ID] = matchedTx
	return nil
}

// GetMatchedTransactionsByMatchSet retrieves matched transactions for a match set from the mock repository
func (r *MockMatchedTransactionRepository) GetMatchedTransactionsByMatchSet(matchSetID string, limit, offset int) ([]models.MatchedTransaction, int, error) {
	var matchedTxs []models.MatchedTransaction
	for _, tx := range r.matchedTransactions {
		if tx.MatchSetID == matchSetID {
			matchedTxs = append(matchedTxs, *tx)
		}
	}

	// Apply pagination
	total := len(matchedTxs)
	start := offset
	if start > total {
		start = total
	}

	end := start + limit
	if end > total {
		end = total
	}

	if start < end {
		return matchedTxs[start:end], total, nil
	}

	return []models.MatchedTransaction{}, total, nil
}

// GetMatchedTransactionsByTenant retrieves matched transactions for a tenant from the mock repository
func (r *MockMatchedTransactionRepository) GetMatchedTransactionsByTenant(tenantID string, limit, offset int) ([]models.MatchedTransaction, int, error) {
	var matchedTxs []models.MatchedTransaction
	for _, tx := range r.matchedTransactions {
		if tx.TenantID == tenantID {
			matchedTxs = append(matchedTxs, *tx)
		}
	}

	// Apply pagination
	total := len(matchedTxs)
	start := offset
	if start > total {
		start = total
	}

	end := start + limit
	if end > total {
		end = total
	}

	if start < end {
		return matchedTxs[start:end], total, nil
	}

	return []models.MatchedTransaction{}, total, nil
}

// GetMatchedTransactionByID retrieves a matched transaction by ID from the mock repository
func (r *MockMatchedTransactionRepository) GetMatchedTransactionByID(id string) (*models.MatchedTransaction, error) {
	if tx, exists := r.matchedTransactions[id]; exists {
		return tx, nil
	}
	return nil, errors.New("matched transaction not found")
}

// GetMatchedTransactionByTransactionID retrieves a matched transaction by transaction ID from the mock repository
func (r *MockMatchedTransactionRepository) GetMatchedTransactionByTransactionID(transactionID string) (*models.MatchedTransaction, error) {
	for _, tx := range r.matchedTransactions {
		if tx.TransactionID == transactionID {
			return tx, nil
		}
	}
	return nil, errors.New("matched transaction not found")
}

// GetMatchedTransactionsByMatchGroup retrieves matched transactions by match group ID from the mock repository
func (r *MockMatchedTransactionRepository) GetMatchedTransactionsByMatchGroup(matchGroupID string) ([]models.MatchedTransaction, error) {
	var matchedTxs []models.MatchedTransaction
	for _, tx := range r.matchedTransactions {
		if tx.MatchGroupID == matchGroupID {
			matchedTxs = append(matchedTxs, *tx)
		}
	}
	return matchedTxs, nil
}

// MockUnmatchedTransactionRepository is a mock implementation for development
type MockUnmatchedTransactionRepository struct {
	unmatchedTransactions map[string]*models.UnmatchedTransaction // ID -> UnmatchedTransaction
}

// CreateUnmatchedTransaction creates an unmatched transaction in the mock repository
func (r *MockUnmatchedTransactionRepository) CreateUnmatchedTransaction(unmatchedTx *models.UnmatchedTransaction) error {
	if unmatchedTx.ID == "" {
		// Generate a dummy ID - in real implementation we'd use UUID
		unmatchedTx.ID = "mock-" + time.Now().Format("20060102150405")
	}
	unmatchedTx.CreatedAt = time.Now()
	r.unmatchedTransactions[unmatchedTx.ID] = unmatchedTx
	return nil
}

// GetUnmatchedTransactionsByMatchSet retrieves unmatched transactions for a match set from the mock repository
func (r *MockUnmatchedTransactionRepository) GetUnmatchedTransactionsByMatchSet(matchSetID string, limit, offset int) ([]models.UnmatchedTransaction, int, error) {
	var unmatchedTxs []models.UnmatchedTransaction
	for _, tx := range r.unmatchedTransactions {
		if tx.MatchSetID == matchSetID {
			unmatchedTxs = append(unmatchedTxs, *tx)
		}
	}

	// Apply pagination
	total := len(unmatchedTxs)
	start := offset
	if start > total {
		start = total
	}

	end := start + limit
	if end > total {
		end = total
	}

	if start < end {
		return unmatchedTxs[start:end], total, nil
	}

	return []models.UnmatchedTransaction{}, total, nil
}

// GetUnmatchedTransactionsByTenant retrieves unmatched transactions for a tenant from the mock repository
func (r *MockUnmatchedTransactionRepository) GetUnmatchedTransactionsByTenant(tenantID string, limit, offset int) ([]models.UnmatchedTransaction, int, error) {
	var unmatchedTxs []models.UnmatchedTransaction
	for _, tx := range r.unmatchedTransactions {
		if tx.TenantID == tenantID {
			unmatchedTxs = append(unmatchedTxs, *tx)
		}
	}

	// Apply pagination
	total := len(unmatchedTxs)
	start := offset
	if start > total {
		start = total
	}

	end := start + limit
	if end > total {
		end = total
	}

	if start < end {
		return unmatchedTxs[start:end], total, nil
	}

	return []models.UnmatchedTransaction{}, total, nil
}

// GetUnmatchedTransactionByID retrieves an unmatched transaction by ID from the mock repository
func (r *MockUnmatchedTransactionRepository) GetUnmatchedTransactionByID(id string) (*models.UnmatchedTransaction, error) {
	if tx, exists := r.unmatchedTransactions[id]; exists {
		return tx, nil
	}
	return nil, errors.New("unmatched transaction not found")
}
