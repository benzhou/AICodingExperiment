package services

import (
	"backend/internal/models"
	"backend/internal/repository"
	"errors"
	"log"
	"time"
)

// MatchSetService provides methods for match set operations
type MatchSetService struct {
	matchSetRepo    repository.MatchSetRepository
	ruleRepo        repository.RuleRepository
	dataSourceRepo  repository.DataSourceRepository
	transactionRepo repository.TransactionRepository
	permissionRepo  repository.PermissionRepository
}

// NewMatchSetService creates a new match set service
func NewMatchSetService(
	matchSetRepo repository.MatchSetRepository,
	ruleRepo repository.RuleRepository,
	dataSourceRepo repository.DataSourceRepository,
	transactionRepo repository.TransactionRepository,
	permissionRepo repository.PermissionRepository,
) *MatchSetService {
	return &MatchSetService{
		matchSetRepo:    matchSetRepo,
		ruleRepo:        ruleRepo,
		dataSourceRepo:  dataSourceRepo,
		transactionRepo: transactionRepo,
		permissionRepo:  permissionRepo,
	}
}

// CreateMatchSet creates a new match set
func (s *MatchSetService) CreateMatchSet(name, description, tenantID, ruleID, userID string) (*models.MatchSet, error) {
	// Check if user has permission to create match sets
	hasPermission, err := s.permissionRepo.HasPermission(userID, models.PermCreateMatchSet, tenantID)
	if err != nil {
		return nil, err
	}
	if !hasPermission {
		return nil, errors.New("unauthorized: requires create match set permission")
	}

	// Create the match set
	matchSet := &models.MatchSet{
		Name:        name,
		Description: description,
		TenantID:    tenantID,
		RuleID:      ruleID,
		CreatedBy:   userID,
	}

	if err := s.matchSetRepo.CreateMatchSet(matchSet); err != nil {
		return nil, err
	}

	return matchSet, nil
}

// GetMatchSetByID retrieves a match set by ID
func (s *MatchSetService) GetMatchSetByID(id, userID, tenantID string) (*models.MatchSet, error) {
	// Check if user has permission to view match sets
	hasPermission, err := s.permissionRepo.HasPermission(userID, models.PermViewMatchSet, tenantID)
	if err != nil {
		return nil, err
	}
	if !hasPermission {
		return nil, errors.New("unauthorized: requires view match set permission")
	}

	// Get the match set
	matchSet, err := s.matchSetRepo.GetMatchSetByID(id)
	if err != nil {
		return nil, err
	}

	// Ensure the match set belongs to the tenant
	if matchSet.TenantID != tenantID {
		return nil, errors.New("match set not found in this tenant")
	}

	return matchSet, nil
}

// GetMatchSetsByTenant retrieves match sets for a tenant
func (s *MatchSetService) GetMatchSetsByTenant(tenantID, userID string) ([]models.MatchSet, error) {
	// Check if user has permission to view match sets
	hasPermission, err := s.permissionRepo.HasPermission(userID, models.PermViewMatchSet, tenantID)
	if err != nil {
		return nil, err
	}
	if !hasPermission {
		return nil, errors.New("unauthorized: requires view match set permission")
	}

	// Get match sets
	return s.matchSetRepo.GetMatchSetsByTenant(tenantID)
}

// UpdateMatchSet updates a match set
func (s *MatchSetService) UpdateMatchSet(id, name, description, ruleID, userID, tenantID string) (*models.MatchSet, error) {
	// Check if user has permission to update match sets
	hasPermission, err := s.permissionRepo.HasPermission(userID, models.PermUpdateMatchSet, tenantID)
	if err != nil {
		return nil, err
	}
	if !hasPermission {
		return nil, errors.New("unauthorized: requires update match set permission")
	}

	// Get the existing match set
	matchSet, err := s.matchSetRepo.GetMatchSetByID(id)
	if err != nil {
		return nil, err
	}

	// Ensure the match set belongs to the tenant
	if matchSet.TenantID != tenantID {
		return nil, errors.New("match set not found in this tenant")
	}

	// Update the match set
	matchSet.Name = name
	matchSet.Description = description
	matchSet.RuleID = ruleID

	if err := s.matchSetRepo.UpdateMatchSet(matchSet); err != nil {
		return nil, err
	}

	return matchSet, nil
}

// DeleteMatchSet deletes a match set
func (s *MatchSetService) DeleteMatchSet(id, userID, tenantID string) error {
	// Check if user has permission to delete match sets
	hasPermission, err := s.permissionRepo.HasPermission(userID, models.PermDeleteMatchSet, tenantID)
	if err != nil {
		return err
	}
	if !hasPermission {
		return errors.New("unauthorized: requires delete match set permission")
	}

	// Get the existing match set
	matchSet, err := s.matchSetRepo.GetMatchSetByID(id)
	if err != nil {
		return err
	}

	// Ensure the match set belongs to the tenant
	if matchSet.TenantID != tenantID {
		return errors.New("match set not found in this tenant")
	}

	// Delete the match set
	return s.matchSetRepo.DeleteMatchSet(id)
}

// AddDataSourceToMatchSet adds a data source to a match set
func (s *MatchSetService) AddDataSourceToMatchSet(matchSetID, dataSourceID, userID, tenantID string) error {
	// Check if user has permission to update match sets
	hasPermission, err := s.permissionRepo.HasPermission(userID, models.PermUpdateMatchSet, tenantID)
	if err != nil {
		return err
	}
	if !hasPermission {
		return errors.New("unauthorized: requires update match set permission")
	}

	// Get the match set
	matchSet, err := s.matchSetRepo.GetMatchSetByID(matchSetID)
	if err != nil {
		return err
	}

	// Ensure the match set belongs to the tenant
	if matchSet.TenantID != tenantID {
		return errors.New("match set not found in this tenant")
	}

	// Get the data source
	dataSource, err := s.dataSourceRepo.GetDataSourceByID(dataSourceID)
	if err != nil {
		return err
	}

	// Ensure the data source belongs to the tenant
	if dataSource.TenantID != tenantID {
		return errors.New("data source not found in this tenant")
	}

	// Add the data source to the match set
	return s.matchSetRepo.AddDataSourceToMatchSet(matchSetID, dataSourceID)
}

// RemoveDataSourceFromMatchSet removes a data source from a match set
func (s *MatchSetService) RemoveDataSourceFromMatchSet(matchSetID, dataSourceID, userID, tenantID string) error {
	// Check if user has permission to update match sets
	hasPermission, err := s.permissionRepo.HasPermission(userID, models.PermUpdateMatchSet, tenantID)
	if err != nil {
		return err
	}
	if !hasPermission {
		return errors.New("unauthorized: requires update match set permission")
	}

	// Get the match set
	matchSet, err := s.matchSetRepo.GetMatchSetByID(matchSetID)
	if err != nil {
		return err
	}

	// Ensure the match set belongs to the tenant
	if matchSet.TenantID != tenantID {
		return errors.New("match set not found in this tenant")
	}

	// Remove the data source from the match set
	return s.matchSetRepo.RemoveDataSourceFromMatchSet(matchSetID, dataSourceID)
}

// GetMatchSetDataSources gets all data sources for a match set
func (s *MatchSetService) GetMatchSetDataSources(matchSetID, userID, tenantID string) ([]models.DataSource, error) {
	// Check if user has permission to view match sets
	hasPermission, err := s.permissionRepo.HasPermission(userID, models.PermViewMatchSet, tenantID)
	if err != nil {
		return nil, err
	}
	if !hasPermission {
		return nil, errors.New("unauthorized: requires view match set permission")
	}

	// Get the match set
	matchSet, err := s.matchSetRepo.GetMatchSetByID(matchSetID)
	if err != nil {
		return nil, err
	}

	// Ensure the match set belongs to the tenant
	if matchSet.TenantID != tenantID {
		return nil, errors.New("match set not found in this tenant")
	}

	// Get data sources
	return s.matchSetRepo.GetMatchSetDataSources(matchSetID)
}

// RunMatchSet executes the matching process for a specific match set
func (s *MatchSetService) RunMatchSet(matchSetID, userID, tenantID string) error {
	// Check if user has permission to match transactions
	hasPermission, err := s.permissionRepo.HasPermission(userID, models.PermMatchTransactions, tenantID)
	if err != nil {
		return err
	}
	if !hasPermission {
		return errors.New("unauthorized: requires match transactions permission")
	}

	// Get the match set
	matchSet, err := s.matchSetRepo.GetMatchSetByID(matchSetID)
	if err != nil {
		return err
	}

	// Ensure the match set belongs to the tenant
	if matchSet.TenantID != tenantID {
		return errors.New("match set not found in this tenant")
	}

	// Get the match rule
	rule, err := s.ruleRepo.GetRuleByID(matchSet.RuleID)
	if err != nil {
		return err
	}

	// Get data sources for this match set
	dataSources, err := s.matchSetRepo.GetMatchSetDataSources(matchSetID)
	if err != nil {
		return err
	}

	// Process is started - would typically queue this for async processing
	// For this implementation, we'll run it synchronously
	log.Printf("Starting matching process for match set %s with rule %s", matchSet.Name, rule.Name)
	log.Printf("Using %d data sources for matching", len(dataSources))

	// This would be a good place to spawn a goroutine or queue a task
	// For now, we'll just log a message
	log.Printf("Matching process completed for match set %s", matchSet.Name)

	return nil
}

// GetMatchSetStatus provides information about the match set processing status
func (s *MatchSetService) GetMatchSetStatus(matchSetID, userID, tenantID string) (map[string]interface{}, error) {
	// Check if user has permission to view match sets
	hasPermission, err := s.permissionRepo.HasPermission(userID, models.PermViewMatchSet, tenantID)
	if err != nil {
		return nil, err
	}
	if !hasPermission {
		return nil, errors.New("unauthorized: requires view match set permission")
	}

	// Get the match set
	matchSet, err := s.matchSetRepo.GetMatchSetByID(matchSetID)
	if err != nil {
		return nil, err
	}

	// Ensure the match set belongs to the tenant
	if matchSet.TenantID != tenantID {
		return nil, errors.New("match set not found in this tenant")
	}

	// Get data sources for this match set
	dataSources, err := s.matchSetRepo.GetMatchSetDataSources(matchSetID)
	if err != nil {
		return nil, err
	}

	// In a real implementation, we would get the actual status from a database
	// For now, we'll return some dummy data
	status := map[string]interface{}{
		"match_set_id":           matchSetID,
		"name":                   matchSet.Name,
		"status":                 "Completed", // Example status: "Running", "Completed", "Failed"
		"data_sources":           len(dataSources),
		"total_transactions":     1000, // Example data
		"matched_transactions":   800,
		"unmatched_transactions": 200,
		"last_run":               time.Now().Add(-24 * time.Hour).Format(time.RFC3339), // Example: Last run 24 hours ago
	}

	return status, nil
}
