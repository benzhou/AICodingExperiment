package services

import (
	"backend/internal/models"
	"backend/internal/repository"
)

// RuleService provides methods for managing match rules
type RuleService struct {
	ruleRepo repository.RuleRepository
}

// NewRuleService creates a new rule service
func NewRuleService(ruleRepo repository.RuleRepository) *RuleService {
	return &RuleService{
		ruleRepo: ruleRepo,
	}
}

// CreateRule creates a new match rule
func (s *RuleService) CreateRule(
	name, description string,
	matchByAmount, matchByDate, matchByReference bool,
	dateTolerance int,
	createdBy string,
) (*models.MatchRule, error) {
	rule := &models.MatchRule{
		Name:             name,
		Description:      description,
		MatchByAmount:    matchByAmount,
		MatchByDate:      matchByDate,
		DateTolerance:    dateTolerance,
		MatchByReference: matchByReference,
		Active:           true,
		CreatedBy:        createdBy,
	}

	if err := s.ruleRepo.CreateRule(rule); err != nil {
		return nil, err
	}

	return rule, nil
}

// GetRuleByID retrieves a match rule by ID
func (s *RuleService) GetRuleByID(id string) (*models.MatchRule, error) {
	return s.ruleRepo.GetRuleByID(id)
}

// GetRuleByName retrieves a match rule by name
func (s *RuleService) GetRuleByName(name string) (*models.MatchRule, error) {
	return s.ruleRepo.GetRuleByName(name)
}

// UpdateRule updates a match rule
func (s *RuleService) UpdateRule(
	id, name, description string,
	matchByAmount, matchByDate, matchByReference, active bool,
	dateTolerance int,
) (*models.MatchRule, error) {
	rule, err := s.ruleRepo.GetRuleByID(id)
	if err != nil {
		return nil, err
	}

	rule.Name = name
	rule.Description = description
	rule.MatchByAmount = matchByAmount
	rule.MatchByDate = matchByDate
	rule.DateTolerance = dateTolerance
	rule.MatchByReference = matchByReference
	rule.Active = active

	if err := s.ruleRepo.UpdateRule(rule); err != nil {
		return nil, err
	}

	return rule, nil
}

// DeleteRule deletes a match rule
func (s *RuleService) DeleteRule(id string) error {
	return s.ruleRepo.DeleteRule(id)
}

// GetAllRules retrieves all match rules
func (s *RuleService) GetAllRules() ([]models.MatchRule, error) {
	return s.ruleRepo.GetAllRules()
}

// GetActiveRules retrieves all active match rules
func (s *RuleService) GetActiveRules() ([]models.MatchRule, error) {
	return s.ruleRepo.GetActiveRules()
}

// ToggleRuleActive toggles the active status of a rule
func (s *RuleService) ToggleRuleActive(id string, active bool) (*models.MatchRule, error) {
	rule, err := s.ruleRepo.GetRuleByID(id)
	if err != nil {
		return nil, err
	}

	rule.Active = active

	if err := s.ruleRepo.UpdateRule(rule); err != nil {
		return nil, err
	}

	return rule, nil
}
