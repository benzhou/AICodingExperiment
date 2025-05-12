package repository

import (
	"backend/internal/db"
	"backend/internal/models"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrRuleNotFound = errors.New("match rule not found")
	ErrRuleExists   = errors.New("match rule with this name already exists")
)

// RuleRepository defines operations for managing matching rules
type RuleRepository interface {
	CreateRule(rule *models.MatchRule) error
	GetRuleByID(id string) (*models.MatchRule, error)
	GetRuleByName(name string) (*models.MatchRule, error)
	UpdateRule(rule *models.MatchRule) error
	DeleteRule(id string) error
	GetAllRules() ([]models.MatchRule, error)
	GetActiveRules() ([]models.MatchRule, error)
}

// PostgresRuleRepository implements RuleRepository for PostgreSQL
type PostgresRuleRepository struct {
	db *sql.DB
}

// NewRuleRepository creates a new rule repository
func NewRuleRepository() RuleRepository {
	if db.DB == nil {
		// Return a mock repository for development
		return &MockRuleRepository{
			rules: make(map[string]*models.MatchRule),
		}
	}
	return &PostgresRuleRepository{
		db: db.DB,
	}
}

// CreateRule creates a new match rule
func (r *PostgresRuleRepository) CreateRule(rule *models.MatchRule) error {
	// Check if rule with the same name already exists
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM match_rules WHERE name = $1)", rule.Name).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		return ErrRuleExists
	}

	query := `
		INSERT INTO match_rules (
			name, description, match_by_amount, match_by_date, 
			date_tolerance, match_by_reference, active, created_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8
		) RETURNING id, created_at, updated_at
	`

	return r.db.QueryRow(
		query,
		rule.Name,
		rule.Description,
		rule.MatchByAmount,
		rule.MatchByDate,
		rule.DateTolerance,
		rule.MatchByReference,
		rule.Active,
		rule.CreatedBy,
	).Scan(&rule.ID, &rule.CreatedAt, &rule.UpdatedAt)
}

// GetRuleByID retrieves a match rule by ID
func (r *PostgresRuleRepository) GetRuleByID(id string) (*models.MatchRule, error) {
	query := `
		SELECT 
			id, name, description, match_by_amount, match_by_date, 
			date_tolerance, match_by_reference, active, created_by, 
			created_at, updated_at
		FROM match_rules
		WHERE id = $1
	`

	var rule models.MatchRule
	err := r.db.QueryRow(query, id).Scan(
		&rule.ID,
		&rule.Name,
		&rule.Description,
		&rule.MatchByAmount,
		&rule.MatchByDate,
		&rule.DateTolerance,
		&rule.MatchByReference,
		&rule.Active,
		&rule.CreatedBy,
		&rule.CreatedAt,
		&rule.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrRuleNotFound
	}

	if err != nil {
		return nil, err
	}

	return &rule, nil
}

// GetRuleByName retrieves a match rule by name
func (r *PostgresRuleRepository) GetRuleByName(name string) (*models.MatchRule, error) {
	query := `
		SELECT 
			id, name, description, match_by_amount, match_by_date, 
			date_tolerance, match_by_reference, active, created_by, 
			created_at, updated_at
		FROM match_rules
		WHERE name = $1
	`

	var rule models.MatchRule
	err := r.db.QueryRow(query, name).Scan(
		&rule.ID,
		&rule.Name,
		&rule.Description,
		&rule.MatchByAmount,
		&rule.MatchByDate,
		&rule.DateTolerance,
		&rule.MatchByReference,
		&rule.Active,
		&rule.CreatedBy,
		&rule.CreatedAt,
		&rule.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrRuleNotFound
	}

	if err != nil {
		return nil, err
	}

	return &rule, nil
}

// UpdateRule updates a match rule
func (r *PostgresRuleRepository) UpdateRule(rule *models.MatchRule) error {
	// Check if another rule with the same name already exists
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM match_rules WHERE name = $1 AND id != $2", rule.Name, rule.ID).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		return ErrRuleExists
	}

	query := `
		UPDATE match_rules
		SET 
			name = $1, 
			description = $2, 
			match_by_amount = $3, 
			match_by_date = $4, 
			date_tolerance = $5, 
			match_by_reference = $6, 
			active = $7,
			updated_at = NOW()
		WHERE id = $8
		RETURNING updated_at
	`

	var updatedAt time.Time
	err = r.db.QueryRow(
		query,
		rule.Name,
		rule.Description,
		rule.MatchByAmount,
		rule.MatchByDate,
		rule.DateTolerance,
		rule.MatchByReference,
		rule.Active,
		rule.ID,
	).Scan(&updatedAt)

	if err == sql.ErrNoRows {
		return ErrRuleNotFound
	}

	if err != nil {
		return err
	}

	rule.UpdatedAt = updatedAt
	return nil
}

// DeleteRule deletes a match rule
func (r *PostgresRuleRepository) DeleteRule(id string) error {
	// Check if there are any existing matches using this rule
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM transaction_matches WHERE match_rule_id = $1", id).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New("cannot delete rule used by existing matches")
	}

	query := "DELETE FROM match_rules WHERE id = $1"
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRuleNotFound
	}

	return nil
}

// GetAllRules retrieves all match rules
func (r *PostgresRuleRepository) GetAllRules() ([]models.MatchRule, error) {
	query := `
		SELECT 
			id, name, description, match_by_amount, match_by_date, 
			date_tolerance, match_by_reference, active, created_by, 
			created_at, updated_at
		FROM match_rules
		ORDER BY name
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rules []models.MatchRule
	for rows.Next() {
		var rule models.MatchRule
		err := rows.Scan(
			&rule.ID,
			&rule.Name,
			&rule.Description,
			&rule.MatchByAmount,
			&rule.MatchByDate,
			&rule.DateTolerance,
			&rule.MatchByReference,
			&rule.Active,
			&rule.CreatedBy,
			&rule.CreatedAt,
			&rule.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		rules = append(rules, rule)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return rules, nil
}

// GetActiveRules retrieves all active match rules
func (r *PostgresRuleRepository) GetActiveRules() ([]models.MatchRule, error) {
	query := `
		SELECT 
			id, name, description, match_by_amount, match_by_date, 
			date_tolerance, match_by_reference, active, created_by, 
			created_at, updated_at
		FROM match_rules
		WHERE active = true
		ORDER BY name
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rules []models.MatchRule
	for rows.Next() {
		var rule models.MatchRule
		err := rows.Scan(
			&rule.ID,
			&rule.Name,
			&rule.Description,
			&rule.MatchByAmount,
			&rule.MatchByDate,
			&rule.DateTolerance,
			&rule.MatchByReference,
			&rule.Active,
			&rule.CreatedBy,
			&rule.CreatedAt,
			&rule.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		rules = append(rules, rule)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return rules, nil
}

// MockRuleRepository is a mock implementation for development
type MockRuleRepository struct {
	rules     map[string]*models.MatchRule
	nameIndex map[string]string // name -> id mapping
}

// CreateRule creates a match rule in the mock repository
func (r *MockRuleRepository) CreateRule(rule *models.MatchRule) error {
	// Initialize nameIndex if it doesn't exist
	if r.nameIndex == nil {
		r.nameIndex = make(map[string]string)
	}

	// Check if rule with the same name exists
	if _, exists := r.nameIndex[rule.Name]; exists {
		return ErrRuleExists
	}

	if rule.ID == "" {
		// Generate a dummy ID - in real implementation we'd use UUID
		rule.ID = "mock-" + time.Now().Format("20060102150405")
	}
	rule.CreatedAt = time.Now()
	rule.UpdatedAt = time.Now()

	r.rules[rule.ID] = rule
	r.nameIndex[rule.Name] = rule.ID

	return nil
}

// GetRuleByID retrieves a match rule by ID from the mock repository
func (r *MockRuleRepository) GetRuleByID(id string) (*models.MatchRule, error) {
	if rule, exists := r.rules[id]; exists {
		return rule, nil
	}
	return nil, ErrRuleNotFound
}

// GetRuleByName retrieves a match rule by name from the mock repository
func (r *MockRuleRepository) GetRuleByName(name string) (*models.MatchRule, error) {
	// Initialize nameIndex if it doesn't exist
	if r.nameIndex == nil {
		r.nameIndex = make(map[string]string)
	}

	id, exists := r.nameIndex[name]
	if !exists {
		return nil, ErrRuleNotFound
	}

	return r.rules[id], nil
}

// UpdateRule updates a match rule in the mock repository
func (r *MockRuleRepository) UpdateRule(rule *models.MatchRule) error {
	// Initialize nameIndex if it doesn't exist
	if r.nameIndex == nil {
		r.nameIndex = make(map[string]string)
	}

	existing, exists := r.rules[rule.ID]
	if !exists {
		return ErrRuleNotFound
	}

	// Check if another rule with the same name exists
	if id, nameExists := r.nameIndex[rule.Name]; nameExists && id != rule.ID {
		return ErrRuleExists
	}

	// Update the name index if the name has changed
	if existing.Name != rule.Name {
		delete(r.nameIndex, existing.Name)
		r.nameIndex[rule.Name] = rule.ID
	}

	rule.UpdatedAt = time.Now()
	r.rules[rule.ID] = rule

	return nil
}

// DeleteRule deletes a match rule from the mock repository
func (r *MockRuleRepository) DeleteRule(id string) error {
	// Initialize nameIndex if it doesn't exist
	if r.nameIndex == nil {
		r.nameIndex = make(map[string]string)
	}

	rule, exists := r.rules[id]
	if !exists {
		return ErrRuleNotFound
	}

	delete(r.nameIndex, rule.Name)
	delete(r.rules, id)

	return nil
}

// GetAllRules retrieves all match rules from the mock repository
func (r *MockRuleRepository) GetAllRules() ([]models.MatchRule, error) {
	var rules []models.MatchRule
	for _, rule := range r.rules {
		rules = append(rules, *rule)
	}
	return rules, nil
}

// GetActiveRules retrieves all active match rules from the mock repository
func (r *MockRuleRepository) GetActiveRules() ([]models.MatchRule, error) {
	var rules []models.MatchRule
	for _, rule := range r.rules {
		if rule.Active {
			rules = append(rules, *rule)
		}
	}
	return rules, nil
}
