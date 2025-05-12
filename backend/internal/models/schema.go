package models

import "time"

// FieldType represents the data type of a schema field
type FieldType string

const (
	// Field types for schema definition
	FieldTypeString  FieldType = "string"
	FieldTypeNumber  FieldType = "number"
	FieldTypeBoolean FieldType = "boolean"
	FieldTypeDate    FieldType = "date"
	FieldTypeTime    FieldType = "time"
	FieldTypeObject  FieldType = "object"
	FieldTypeArray   FieldType = "array"
)

// SchemaField represents a field in a data source schema
type SchemaField struct {
	ID           string    `json:"id" db:"id"`
	SchemaID     string    `json:"schema_id" db:"schema_id"`
	Name         string    `json:"name" db:"name"`
	DisplayName  string    `json:"display_name" db:"display_name"`
	Type         FieldType `json:"type" db:"type"`
	Required     bool      `json:"required" db:"required"`
	DefaultValue string    `json:"default_value,omitempty" db:"default_value"`
	Order        int       `json:"order" db:"order"`
}

// DataSourceSchema defines the structure of a data source
type DataSourceSchema struct {
	ID          string        `json:"id" db:"id"`
	Name        string        `json:"name" db:"name"`
	Description string        `json:"description" db:"description"`
	TenantID    string        `json:"tenant_id" db:"tenant_id"`
	CreatedBy   string        `json:"created_by" db:"created_by"`
	CreatedAt   time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at" db:"updated_at"`
	Fields      []SchemaField `json:"fields" db:"-"` // Not stored directly in this table
}

// SchemaMapping defines how fields from a source map to the system's internal format
type SchemaMapping struct {
	ID              string    `json:"id" db:"id"`
	SchemaID        string    `json:"schema_id" db:"schema_id"`
	SourceFieldName string    `json:"source_field_name" db:"source_field_name"`
	TargetFieldName string    `json:"target_field_name" db:"target_field_name"`
	Transformation  string    `json:"transformation,omitempty" db:"transformation"` // Optional transformation script/expression
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

// ColumnMapping represents a mapping between CSV columns and schema fields
type ColumnMapping struct {
	ColumnIndex int    `json:"column_index"`
	FieldName   string `json:"field_name"`
}

// FileParsingConfig contains configuration for parsing file uploads
type FileParsingConfig struct {
	ID             string    `json:"id" db:"id"`
	SchemaID       string    `json:"schema_id" db:"schema_id"`
	FileType       string    `json:"file_type" db:"file_type"` // CSV, Excel, JSON, etc.
	HasHeaderRow   bool      `json:"has_header_row" db:"has_header_row"`
	Delimiter      string    `json:"delimiter,omitempty" db:"delimiter"` // For CSV
	DateFormat     string    `json:"date_format,omitempty" db:"date_format"`
	TimeFormat     string    `json:"time_format,omitempty" db:"time_format"`
	NumberFormat   string    `json:"number_format,omitempty" db:"number_format"`
	EncapsulatedBy string    `json:"encapsulated_by,omitempty" db:"encapsulated_by"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}
