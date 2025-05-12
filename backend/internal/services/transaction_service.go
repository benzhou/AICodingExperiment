package services

import (
	"backend/internal/models"
	"backend/internal/repository"
	"errors"
)

// Error definitions
var (
	ErrTransactionNotFound = errors.New("transaction not found")
)

// TransactionService provides methods for managing transactions
type TransactionService struct {
	transactionRepo repository.TransactionRepository
}

// NewTransactionService creates a new transaction service
func NewTransactionService(transactionRepo repository.TransactionRepository) *TransactionService {
	return &TransactionService{
		transactionRepo: transactionRepo,
	}
}

// CreateTransaction creates a new transaction
func (s *TransactionService) CreateTransaction(transaction *models.Transaction) error {
	return s.transactionRepo.CreateTransaction(transaction)
}

// CreateTransactions creates multiple transactions in batch
func (s *TransactionService) CreateTransactions(transactions []models.Transaction) error {
	return s.transactionRepo.CreateTransactions(transactions)
}

// GetTransactionByID retrieves a transaction by ID
func (s *TransactionService) GetTransactionByID(id string) (*models.Transaction, error) {
	transaction, err := s.transactionRepo.GetTransactionByID(id)
	if err != nil {
		if err == repository.ErrTransactionNotFound {
			return nil, ErrTransactionNotFound
		}
		return nil, err
	}
	return transaction, nil
}

// GetTransactionsByDataSourceID retrieves transactions for a data source
func (s *TransactionService) GetTransactionsByDataSourceID(dataSourceID string) ([]models.Transaction, error) {
	return s.transactionRepo.GetTransactionsByDataSourceID(dataSourceID)
}

// GetTransactionsByUserID retrieves transactions uploaded by a user
func (s *TransactionService) GetTransactionsByUserID(userID string) ([]models.Transaction, error) {
	return s.transactionRepo.GetTransactionsByUserID(userID)
}

// GetRecentTransactions retrieves recent transactions up to a limit
func (s *TransactionService) GetRecentTransactions(limit int) ([]models.Transaction, error) {
	return s.transactionRepo.GetRecentTransactions(limit)
}

// DeleteTransaction deletes a transaction
func (s *TransactionService) DeleteTransaction(id string) error {
	err := s.transactionRepo.DeleteTransaction(id)
	if err != nil {
		if err == repository.ErrTransactionNotFound {
			return ErrTransactionNotFound
		}
		return err
	}
	return nil
}

// DeleteTransactionsByDataSourceID deletes all transactions for a data source
func (s *TransactionService) DeleteTransactionsByDataSourceID(dataSourceID string) error {
	return s.transactionRepo.DeleteTransactionsByDataSourceID(dataSourceID)
}
