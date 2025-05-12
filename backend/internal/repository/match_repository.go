package repository

import (
	"backend/internal/db"
	"backend/internal/models"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrMatchNotFound = errors.New("match not found")
)

// MatchRepository defines operations for managing transaction matches
type MatchRepository interface {
	CreateMatch(match *models.TransactionMatch) error
	GetMatchByID(id string) (*models.TransactionMatch, error)
	GetMatchesByStatus(status string) ([]models.TransactionMatch, error)
	UpdateMatchStatus(id string, status string, approvedBy string, reason string) error
	GetMatchesByUser(userID string) ([]models.TransactionMatch, error)
	SearchMatches(filters map[string]interface{}, limit, offset int) ([]models.TransactionMatch, int, error)
}

// PostgresMatchRepository implements MatchRepository for PostgreSQL
type PostgresMatchRepository struct {
	db *sql.DB
}

// NewMatchRepository creates a new match repository
func NewMatchRepository() MatchRepository {
	if db.DB == nil {
		// Return a mock repository for development
		return &MockMatchRepository{
			matches: make(map[string]*models.TransactionMatch),
		}
	}
	return &PostgresMatchRepository{
		db: db.DB,
	}
}

// CreateMatch creates a new transaction match
func (r *PostgresMatchRepository) CreateMatch(match *models.TransactionMatch) error {
	query := `
		INSERT INTO transaction_matches (
			match_status, match_type, match_rule_id, matched_by
		) VALUES (
			$1, $2, $3, $4
		) RETURNING id, created_at, updated_at
	`

	var matchRuleIDParam interface{} = nil
	if match.MatchRuleID != "" {
		matchRuleIDParam = match.MatchRuleID
	}

	return r.db.QueryRow(
		query,
		match.MatchStatus,
		match.MatchType,
		matchRuleIDParam,
		match.MatchedBy,
	).Scan(&match.ID, &match.CreatedAt, &match.UpdatedAt)
}

// GetMatchByID retrieves a match by ID
func (r *PostgresMatchRepository) GetMatchByID(id string) (*models.TransactionMatch, error) {
	query := `
		SELECT 
			id, match_status, match_type, match_rule_id, matched_by, 
			approved_by, approval_date, rejection_reason, created_at, updated_at
		FROM transaction_matches
		WHERE id = $1
	`

	var match models.TransactionMatch
	var matchRuleID, approvedBy sql.NullString
	var approvalDate sql.NullTime

	err := r.db.QueryRow(query, id).Scan(
		&match.ID,
		&match.MatchStatus,
		&match.MatchType,
		&matchRuleID,
		&match.MatchedBy,
		&approvedBy,
		&approvalDate,
		&match.RejectionReason,
		&match.CreatedAt,
		&match.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrMatchNotFound
	}

	if err != nil {
		return nil, err
	}

	if matchRuleID.Valid {
		match.MatchRuleID = matchRuleID.String
	}

	if approvedBy.Valid {
		match.ApprovedBy = approvedBy.String
	}

	if approvalDate.Valid {
		match.ApprovalDate = approvalDate.Time
	}

	return &match, nil
}

// GetMatchesByStatus retrieves matches by status
func (r *PostgresMatchRepository) GetMatchesByStatus(status string) ([]models.TransactionMatch, error) {
	query := `
		SELECT 
			id, match_status, match_type, match_rule_id, matched_by, 
			approved_by, approval_date, rejection_reason, created_at, updated_at
		FROM transaction_matches
		WHERE match_status = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var matches []models.TransactionMatch
	for rows.Next() {
		var match models.TransactionMatch
		var matchRuleID, approvedBy sql.NullString
		var approvalDate sql.NullTime

		err := rows.Scan(
			&match.ID,
			&match.MatchStatus,
			&match.MatchType,
			&matchRuleID,
			&match.MatchedBy,
			&approvedBy,
			&approvalDate,
			&match.RejectionReason,
			&match.CreatedAt,
			&match.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		if matchRuleID.Valid {
			match.MatchRuleID = matchRuleID.String
		}

		if approvedBy.Valid {
			match.ApprovedBy = approvedBy.String
		}

		if approvalDate.Valid {
			match.ApprovalDate = approvalDate.Time
		}

		matches = append(matches, match)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return matches, nil
}

// UpdateMatchStatus updates a match's status
func (r *PostgresMatchRepository) UpdateMatchStatus(id string, status string, approvedBy string, reason string) error {
	query := `
		UPDATE transaction_matches
		SET match_status = $1,
			updated_at = NOW(),
	`

	var params []interface{}
	params = append(params, status)

	// Conditionally add approval or rejection fields
	if status == "Approved" {
		query += `
			approved_by = $2,
			approval_date = NOW()
		WHERE id = $3
		RETURNING updated_at
		`
		params = append(params, approvedBy, id)
	} else if status == "Rejected" {
		query += `
			approved_by = $2,
			rejection_reason = $3
		WHERE id = $4
		RETURNING updated_at
		`
		params = append(params, approvedBy, reason, id)
	} else {
		query += `
		WHERE id = $2
		RETURNING updated_at
		`
		params = append(params, id)
	}

	var updatedAt time.Time
	err := r.db.QueryRow(query, params...).Scan(&updatedAt)

	if err == sql.ErrNoRows {
		return ErrMatchNotFound
	}

	return err
}

// GetMatchesByUser retrieves matches created by a specific user
func (r *PostgresMatchRepository) GetMatchesByUser(userID string) ([]models.TransactionMatch, error) {
	query := `
		SELECT 
			id, match_status, match_type, match_rule_id, matched_by, 
			approved_by, approval_date, rejection_reason, created_at, updated_at
		FROM transaction_matches
		WHERE matched_by = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var matches []models.TransactionMatch
	for rows.Next() {
		var match models.TransactionMatch
		var matchRuleID, approvedBy sql.NullString
		var approvalDate sql.NullTime

		err := rows.Scan(
			&match.ID,
			&match.MatchStatus,
			&match.MatchType,
			&matchRuleID,
			&match.MatchedBy,
			&approvedBy,
			&approvalDate,
			&match.RejectionReason,
			&match.CreatedAt,
			&match.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		if matchRuleID.Valid {
			match.MatchRuleID = matchRuleID.String
		}

		if approvedBy.Valid {
			match.ApprovedBy = approvedBy.String
		}

		if approvalDate.Valid {
			match.ApprovalDate = approvalDate.Time
		}

		matches = append(matches, match)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return matches, nil
}

// SearchMatches searches for transaction matches using filters
func (r *PostgresMatchRepository) SearchMatches(filters map[string]interface{}, limit, offset int) ([]models.TransactionMatch, int, error) {
	// Start building the base query
	baseQuery := `
		FROM transaction_matches tm
		WHERE 1=1
	`

	// Build filter conditions and parameters
	var conditions []string
	var params []interface{}
	paramIndex := 1

	// Add filter conditions
	if val, ok := filters["status"]; ok {
		conditions = append(conditions, "tm.match_status = $"+string(paramIndex))
		params = append(params, val)
		paramIndex++
	}

	if val, ok := filters["matchType"]; ok {
		conditions = append(conditions, "tm.match_type = $"+string(paramIndex))
		params = append(params, val)
		paramIndex++
	}

	if val, ok := filters["matchedBy"]; ok {
		conditions = append(conditions, "tm.matched_by = $"+string(paramIndex))
		params = append(params, val)
		paramIndex++
	}

	if val, ok := filters["approvedBy"]; ok {
		conditions = append(conditions, "tm.approved_by = $"+string(paramIndex))
		params = append(params, val)
		paramIndex++
	}

	if val, ok := filters["dateFrom"]; ok {
		conditions = append(conditions, "tm.created_at >= $"+string(paramIndex))
		params = append(params, val)
		paramIndex++
	}

	if val, ok := filters["dateTo"]; ok {
		conditions = append(conditions, "tm.created_at <= $"+string(paramIndex))
		params = append(params, val)
		paramIndex++
	}

	// Apply the conditions
	for _, condition := range conditions {
		baseQuery += " AND " + condition
	}

	// Count query for pagination
	countQuery := "SELECT COUNT(*) " + baseQuery
	var total int
	err := r.db.QueryRow(countQuery, params...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Data query
	dataQuery := `
		SELECT 
			tm.id, tm.match_status, tm.match_type, tm.match_rule_id, tm.matched_by, 
			tm.approved_by, tm.approval_date, tm.rejection_reason, tm.created_at, tm.updated_at
	` + baseQuery + `
		ORDER BY tm.created_at DESC
		LIMIT $` + string(paramIndex) + ` OFFSET $` + string(paramIndex+1)

	params = append(params, limit, offset)

	rows, err := r.db.Query(dataQuery, params...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var matches []models.TransactionMatch
	for rows.Next() {
		var match models.TransactionMatch
		var matchRuleID, approvedBy sql.NullString
		var approvalDate sql.NullTime

		err := rows.Scan(
			&match.ID,
			&match.MatchStatus,
			&match.MatchType,
			&matchRuleID,
			&match.MatchedBy,
			&approvedBy,
			&approvalDate,
			&match.RejectionReason,
			&match.CreatedAt,
			&match.UpdatedAt,
		)

		if err != nil {
			return nil, 0, err
		}

		if matchRuleID.Valid {
			match.MatchRuleID = matchRuleID.String
		}

		if approvedBy.Valid {
			match.ApprovedBy = approvedBy.String
		}

		if approvalDate.Valid {
			match.ApprovalDate = approvalDate.Time
		}

		matches = append(matches, match)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return matches, total, nil
}

// MockMatchRepository is a mock implementation for development
type MockMatchRepository struct {
	matches map[string]*models.TransactionMatch
}

// CreateMatch creates a match in the mock repository
func (r *MockMatchRepository) CreateMatch(match *models.TransactionMatch) error {
	if match.ID == "" {
		// Generate a dummy ID - in real implementation we'd use UUID
		match.ID = "mock-" + time.Now().Format("20060102150405")
	}
	match.CreatedAt = time.Now()
	match.UpdatedAt = time.Now()
	r.matches[match.ID] = match
	return nil
}

// GetMatchByID retrieves a match by ID from the mock repository
func (r *MockMatchRepository) GetMatchByID(id string) (*models.TransactionMatch, error) {
	if match, exists := r.matches[id]; exists {
		return match, nil
	}
	return nil, ErrMatchNotFound
}

// GetMatchesByStatus retrieves matches by status from the mock repository
func (r *MockMatchRepository) GetMatchesByStatus(status string) ([]models.TransactionMatch, error) {
	var matches []models.TransactionMatch
	for _, match := range r.matches {
		if match.MatchStatus == status {
			matches = append(matches, *match)
		}
	}
	return matches, nil
}

// UpdateMatchStatus updates a match's status in the mock repository
func (r *MockMatchRepository) UpdateMatchStatus(id string, status string, approvedBy string, reason string) error {
	match, exists := r.matches[id]
	if !exists {
		return ErrMatchNotFound
	}

	match.MatchStatus = status
	match.UpdatedAt = time.Now()

	if status == "Approved" {
		match.ApprovedBy = approvedBy
		match.ApprovalDate = time.Now()
	} else if status == "Rejected" {
		match.ApprovedBy = approvedBy
		match.RejectionReason = reason
	}

	return nil
}

// GetMatchesByUser retrieves matches created by a specific user from the mock repository
func (r *MockMatchRepository) GetMatchesByUser(userID string) ([]models.TransactionMatch, error) {
	var matches []models.TransactionMatch
	for _, match := range r.matches {
		if match.MatchedBy == userID {
			matches = append(matches, *match)
		}
	}
	return matches, nil
}

// SearchMatches searches for transaction matches using filters in the mock repository
func (r *MockMatchRepository) SearchMatches(filters map[string]interface{}, limit, offset int) ([]models.TransactionMatch, int, error) {
	// Simple implementation for mock - just return all matches
	var matches []models.TransactionMatch
	for _, match := range r.matches {
		matches = append(matches, *match)
	}

	// Apply limit and offset
	start := offset
	if start > len(matches) {
		start = len(matches)
	}

	end := start + limit
	if end > len(matches) {
		end = len(matches)
	}

	if start < end {
		return matches[start:end], len(matches), nil
	}

	return []models.TransactionMatch{}, len(matches), nil
}
