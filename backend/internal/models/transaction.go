package models

import (
	"backend/internal/utils"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// TransactionSchemaField represents a field in a data source schema
type TransactionSchemaField struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	Type        string `json:"type"` // string, number, date, boolean
	Required    bool   `json:"required"`
	Format      string `json:"format"` // For date format or number format
	Description string `json:"description"`
}

// SchemaDefinition is a JSON struct for defining the expected columns in a data source
type SchemaDefinition struct {
	Fields          []TransactionSchemaField `json:"fields"`
	DateFormat      string                   `json:"dateFormat"`
	DefaultMappings map[string]string        `json:"defaultMappings"` // Maps standard fields to custom fields
	RequiredFields  []string                 `json:"requiredFields"`
}

// Value implements the driver.Valuer interface for SchemaDefinition
func (sd SchemaDefinition) Value() (driver.Value, error) {
	return json.Marshal(sd)
}

// Scan implements the sql.Scanner interface for SchemaDefinition
func (sd *SchemaDefinition) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &sd)
}

// DataSource represents a source of transaction data
type DataSource struct {
	ID               string           `json:"id" db:"id"`
	Name             string           `json:"name" db:"name"`
	Description      string           `json:"description" db:"description"`
	TenantID         string           `json:"tenant_id" db:"tenant_id"`
	SchemaDefinition SchemaDefinition `json:"schema_definition,omitempty" db:"schema_definition"`
	CreatedAt        time.Time        `json:"-" db:"created_at"`
	UpdatedAt        time.Time        `json:"-" db:"updated_at"`

	CreatedAtEpoch int64 `json:"created_at" db:"-"`
	UpdatedAtEpoch int64 `json:"updated_at" db:"-"`
}

// PrepareMarshal prepares the DataSource for JSON marshaling
func (ds *DataSource) PrepareMarshal() {
	// Convert time.Time to Unix timestamps in milliseconds
	ds.CreatedAtEpoch = utils.TimeToMillis(ds.CreatedAt)
	ds.UpdatedAtEpoch = utils.TimeToMillis(ds.UpdatedAt)
}

// MarshalJSON customizes JSON serialization for DataSource
func (ds *DataSource) MarshalJSON() ([]byte, error) {
	// Ensure epoch timestamps are set
	if ds.CreatedAtEpoch == 0 && !ds.CreatedAt.IsZero() {
		ds.CreatedAtEpoch = utils.TimeToMillis(ds.CreatedAt)
	}

	if ds.UpdatedAtEpoch == 0 && !ds.UpdatedAt.IsZero() {
		ds.UpdatedAtEpoch = utils.TimeToMillis(ds.UpdatedAt)
	}

	type Alias DataSource
	return json.Marshal(&struct {
		*Alias
		CreatedAt int64 `json:"created_at"`
		UpdatedAt int64 `json:"updated_at"`
	}{
		Alias:     (*Alias)(ds),
		CreatedAt: ds.CreatedAtEpoch,
		UpdatedAt: ds.UpdatedAtEpoch,
	})
}

// ImportRecord represents a batch import of transactions
type ImportRecord struct {
	ID           string          `json:"id" db:"id"`
	DataSourceID string          `json:"data_source_id" db:"data_source_id"`
	FileName     string          `json:"file_name" db:"file_name"`
	FileSize     int64           `json:"file_size" db:"file_size"`
	Status       string          `json:"status" db:"status"` // Processing, Completed, Failed
	RowCount     int             `json:"row_count" db:"row_count"`
	SuccessCount int             `json:"success_count" db:"success_count"`
	ErrorCount   int             `json:"error_count" db:"error_count"`
	ImportedBy   string          `json:"imported_by" db:"imported_by"`
	CreatedAt    time.Time       `json:"-" db:"created_at"`
	UpdatedAt    time.Time       `json:"-" db:"updated_at"`
	Metadata     json.RawMessage `json:"metadata,omitempty" db:"metadata"`

	CreatedAtEpoch int64 `json:"created_at" db:"-"`
	UpdatedAtEpoch int64 `json:"updated_at" db:"-"`
}

// MarshalJSON customizes JSON serialization for ImportRecord
func (ir *ImportRecord) MarshalJSON() ([]byte, error) {
	// Ensure epoch timestamps are set
	if ir.CreatedAtEpoch == 0 && !ir.CreatedAt.IsZero() {
		ir.CreatedAtEpoch = utils.TimeToMillis(ir.CreatedAt)
	}

	if ir.UpdatedAtEpoch == 0 && !ir.UpdatedAt.IsZero() {
		ir.UpdatedAtEpoch = utils.TimeToMillis(ir.UpdatedAt)
	}

	type Alias ImportRecord
	return json.Marshal(&struct {
		*Alias
		CreatedAt int64 `json:"created_at"`
		UpdatedAt int64 `json:"updated_at"`
	}{
		Alias:     (*Alias)(ir),
		CreatedAt: ir.CreatedAtEpoch,
		UpdatedAt: ir.UpdatedAtEpoch,
	})
}

// RawTransaction represents a transaction in its original form
type RawTransaction struct {
	ID           string          `json:"id" db:"id"`
	ImportID     string          `json:"import_id" db:"import_id"`
	DataSourceID string          `json:"data_source_id" db:"data_source_id"`
	RowNumber    int             `json:"row_number" db:"row_number"`
	Data         json.RawMessage `json:"data" db:"data"`
	ErrorMessage string          `json:"error_message,omitempty" db:"error_message"`
	CreatedAt    time.Time       `json:"-" db:"created_at"`

	CreatedAtEpoch int64 `json:"created_at" db:"-"`
}

// MarshalJSON customizes JSON serialization for RawTransaction
func (rt *RawTransaction) MarshalJSON() ([]byte, error) {
	// Ensure epoch timestamp is set
	if rt.CreatedAtEpoch == 0 && !rt.CreatedAt.IsZero() {
		rt.CreatedAtEpoch = utils.TimeToMillis(rt.CreatedAt)
	}

	type Alias RawTransaction
	return json.Marshal(&struct {
		*Alias
		CreatedAt int64 `json:"created_at"`
	}{
		Alias:     (*Alias)(rt),
		CreatedAt: rt.CreatedAtEpoch,
	})
}

// Transaction represents a financial transaction
type Transaction struct {
	ID              string         `json:"id" db:"id"`
	DataSourceID    string         `json:"dataSourceId" db:"data_source_id"`
	TransactionDate time.Time      `json:"-" db:"transaction_date"`
	PostDate        time.Time      `json:"-" db:"post_date"`
	Description     string         `json:"description" db:"description"`
	Reference       string         `json:"reference" db:"reference"`
	Amount          float64        `json:"amount" db:"amount"`
	Currency        string         `json:"currency" db:"currency"`
	Status          string         `json:"status" db:"status"`
	MatchID         sql.NullString `json:"matchId,omitempty" db:"match_id"`
	CreatedBy       string         `json:"createdBy" db:"created_by"`
	CreatedAt       time.Time      `json:"-" db:"created_at"`
	UpdatedAt       time.Time      `json:"-" db:"updated_at"`

	TransactionDateEpoch int64 `json:"transactionDate" db:"-"`
	PostDateEpoch        int64 `json:"postDate" db:"-"`
	CreatedAtEpoch       int64 `json:"createdAt" db:"-"`
	UpdatedAtEpoch       int64 `json:"updatedAt" db:"-"`
}

// MarshalJSON customizes JSON serialization for Transaction
func (t *Transaction) MarshalJSON() ([]byte, error) {
	// Ensure epoch timestamps are set
	if t.TransactionDateEpoch == 0 && !t.TransactionDate.IsZero() {
		t.TransactionDateEpoch = utils.TimeToMillis(t.TransactionDate)
	}

	if t.PostDateEpoch == 0 && !t.PostDate.IsZero() {
		t.PostDateEpoch = utils.TimeToMillis(t.PostDate)
	}

	if t.CreatedAtEpoch == 0 && !t.CreatedAt.IsZero() {
		t.CreatedAtEpoch = utils.TimeToMillis(t.CreatedAt)
	}

	if t.UpdatedAtEpoch == 0 && !t.UpdatedAt.IsZero() {
		t.UpdatedAtEpoch = utils.TimeToMillis(t.UpdatedAt)
	}

	type Alias Transaction
	return json.Marshal(&struct {
		*Alias
		TransactionDate int64 `json:"transactionDate"`
		PostDate        int64 `json:"postDate"`
		CreatedAt       int64 `json:"createdAt"`
		UpdatedAt       int64 `json:"updatedAt"`
	}{
		Alias:           (*Alias)(t),
		TransactionDate: t.TransactionDateEpoch,
		PostDate:        t.PostDateEpoch,
		CreatedAt:       t.CreatedAtEpoch,
		UpdatedAt:       t.UpdatedAtEpoch,
	})
}

// TransactionMatch represents a match between two or more transactions
type TransactionMatch struct {
	ID              string    `json:"id" db:"id"`
	MatchStatus     string    `json:"match_status" db:"match_status"` // Pending, Approved, Rejected
	MatchType       string    `json:"match_type" db:"match_type"`     // Automatic, Manual
	MatchRuleID     string    `json:"match_rule_id,omitempty" db:"match_rule_id"`
	TenantID        string    `json:"tenant_id" db:"tenant_id"`
	MatchedBy       string    `json:"matched_by" db:"matched_by"` // User ID
	ApprovedBy      string    `json:"approved_by,omitempty" db:"approved_by"`
	ApprovalDate    time.Time `json:"approval_date,omitempty" db:"approval_date"`
	RejectionReason string    `json:"rejection_reason,omitempty" db:"rejection_reason"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

// MatchRule defines criteria for automatic transaction matching
type MatchRule struct {
	ID               string    `json:"id" db:"id"`
	Name             string    `json:"name" db:"name"`
	Description      string    `json:"description" db:"description"`
	TenantID         string    `json:"tenant_id" db:"tenant_id"`
	MatchByAmount    bool      `json:"match_by_amount" db:"match_by_amount"`
	MatchByDate      bool      `json:"match_by_date" db:"match_by_date"`
	DateTolerance    int       `json:"date_tolerance" db:"date_tolerance"` // Days
	MatchByReference bool      `json:"match_by_reference" db:"match_by_reference"`
	Active           bool      `json:"active" db:"active"`
	CreatedBy        string    `json:"created_by" db:"created_by"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

// TransactionUpload tracks file uploads
type TransactionUpload struct {
	ID           string    `json:"id" db:"id"`
	DataSourceID string    `json:"data_source_id" db:"data_source_id"`
	TenantID     string    `json:"tenant_id" db:"tenant_id"`
	FileName     string    `json:"file_name" db:"file_name"`
	FileSize     int64     `json:"file_size" db:"file_size"`
	UploadedBy   string    `json:"uploaded_by" db:"uploaded_by"`
	UploadDate   time.Time `json:"-" db:"upload_date"`
	Status       string    `json:"status" db:"status"` // Processing, Completed, Failed
	RecordCount  int       `json:"record_count" db:"record_count"`
	ErrorMessage string    `json:"error_message,omitempty" db:"error_message"`

	UploadDateEpoch int64 `json:"upload_date" db:"-"`
}

// MarshalJSON customizes JSON serialization for TransactionUpload
func (tu *TransactionUpload) MarshalJSON() ([]byte, error) {
	// Ensure epoch timestamp is set
	if tu.UploadDateEpoch == 0 && !tu.UploadDate.IsZero() {
		tu.UploadDateEpoch = utils.TimeToMillis(tu.UploadDate)
	}

	type Alias TransactionUpload
	return json.Marshal(&struct {
		*Alias
		UploadDate int64 `json:"upload_date"`
	}{
		Alias:      (*Alias)(tu),
		UploadDate: tu.UploadDateEpoch,
	})
}

// MatchedTransaction represents a successfully matched transaction
type MatchedTransaction struct {
	ID            string    `json:"id" db:"id"`
	MatchSetID    string    `json:"match_set_id" db:"match_set_id"`
	TransactionID string    `json:"transaction_id" db:"transaction_id"`
	MatchGroupID  string    `json:"match_group_id" db:"match_group_id"`
	TenantID      string    `json:"tenant_id" db:"tenant_id"`
	CreatedAt     time.Time `json:"-" db:"created_at"`

	CreatedAtEpoch int64 `json:"created_at" db:"-"`
}

// MarshalJSON customizes JSON serialization for MatchedTransaction
func (mt *MatchedTransaction) MarshalJSON() ([]byte, error) {
	// Ensure epoch timestamp is set
	if mt.CreatedAtEpoch == 0 && !mt.CreatedAt.IsZero() {
		mt.CreatedAtEpoch = utils.TimeToMillis(mt.CreatedAt)
	}

	type Alias MatchedTransaction
	return json.Marshal(&struct {
		*Alias
		CreatedAt int64 `json:"created_at"`
	}{
		Alias:     (*Alias)(mt),
		CreatedAt: mt.CreatedAtEpoch,
	})
}

// UnmatchedTransaction represents a transaction that couldn't be matched
type UnmatchedTransaction struct {
	ID               string    `json:"id" db:"id"`
	MatchSetID       string    `json:"match_set_id" db:"match_set_id"`
	TransactionID    string    `json:"transaction_id" db:"transaction_id"`
	RawTransactionID string    `json:"raw_transaction_id,omitempty" db:"raw_transaction_id"`
	Reason           string    `json:"reason" db:"reason"`
	TenantID         string    `json:"tenant_id" db:"tenant_id"`
	CreatedAt        time.Time `json:"-" db:"created_at"`

	CreatedAtEpoch int64 `json:"created_at" db:"-"`
}

// MarshalJSON customizes JSON serialization for UnmatchedTransaction
func (ut *UnmatchedTransaction) MarshalJSON() ([]byte, error) {
	// Ensure epoch timestamp is set
	if ut.CreatedAtEpoch == 0 && !ut.CreatedAt.IsZero() {
		ut.CreatedAtEpoch = utils.TimeToMillis(ut.CreatedAt)
	}

	type Alias UnmatchedTransaction
	return json.Marshal(&struct {
		*Alias
		CreatedAt int64 `json:"created_at"`
	}{
		Alias:     (*Alias)(ut),
		CreatedAt: ut.CreatedAtEpoch,
	})
}
