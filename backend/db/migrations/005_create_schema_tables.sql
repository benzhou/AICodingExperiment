-- +migrate Up
-- Create data source schema table
CREATE TABLE IF NOT EXISTS data_source_schemas (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMP NOT NULL DEFAULT (NOW() AT TIME ZONE 'UTC'),
    updated_at TIMESTAMP NOT NULL DEFAULT (NOW() AT TIME ZONE 'UTC'),
    UNIQUE(name, tenant_id)
);

-- Create schema fields table
CREATE TABLE IF NOT EXISTS schema_fields (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    schema_id UUID NOT NULL REFERENCES data_source_schemas(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    display_name VARCHAR(100) NOT NULL,
    type VARCHAR(20) NOT NULL CHECK (type IN ('string', 'number', 'boolean', 'date', 'time', 'object', 'array')),
    required BOOLEAN NOT NULL DEFAULT false,
    default_value TEXT,
    "order" INTEGER NOT NULL DEFAULT 0,
    UNIQUE(schema_id, name)
);

-- Create schema mappings table
CREATE TABLE IF NOT EXISTS schema_mappings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    schema_id UUID NOT NULL REFERENCES data_source_schemas(id) ON DELETE CASCADE,
    source_field_name VARCHAR(100) NOT NULL,
    target_field_name VARCHAR(100) NOT NULL,
    transformation TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT (NOW() AT TIME ZONE 'UTC'),
    updated_at TIMESTAMP NOT NULL DEFAULT (NOW() AT TIME ZONE 'UTC'),
    UNIQUE(schema_id, source_field_name)
);

-- Create file parsing configuration table
CREATE TABLE IF NOT EXISTS file_parsing_configs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    schema_id UUID NOT NULL REFERENCES data_source_schemas(id) ON DELETE CASCADE,
    file_type VARCHAR(20) NOT NULL CHECK (file_type IN ('CSV', 'Excel', 'JSON', 'XML')),
    has_header_row BOOLEAN NOT NULL DEFAULT true,
    delimiter VARCHAR(5),
    date_format VARCHAR(50),
    time_format VARCHAR(50),
    number_format VARCHAR(50),
    encapsulated_by VARCHAR(5),
    created_at TIMESTAMP NOT NULL DEFAULT (NOW() AT TIME ZONE 'UTC'),
    updated_at TIMESTAMP NOT NULL DEFAULT (NOW() AT TIME ZONE 'UTC'),
    UNIQUE(schema_id, file_type)
);

-- Add schema_id to data_sources table if doesn't exist
ALTER TABLE data_sources ADD COLUMN IF NOT EXISTS schema_id UUID REFERENCES data_source_schemas(id);

-- Create indices for better performance
CREATE INDEX IF NOT EXISTS idx_data_source_schemas_tenant_id ON data_source_schemas(tenant_id);
CREATE INDEX IF NOT EXISTS idx_schema_fields_schema_id ON schema_fields(schema_id);
CREATE INDEX IF NOT EXISTS idx_schema_mappings_schema_id ON schema_mappings(schema_id);
CREATE INDEX IF NOT EXISTS idx_file_parsing_configs_schema_id ON file_parsing_configs(schema_id);

-- +migrate Down
DROP TABLE IF EXISTS file_parsing_configs CASCADE;
DROP TABLE IF EXISTS schema_mappings CASCADE;
DROP TABLE IF EXISTS schema_fields CASCADE;
DROP TABLE IF EXISTS data_source_schemas CASCADE; 