-- +migrate Up
-- Add schema_definition to data_sources table
ALTER TABLE data_sources 
ADD COLUMN schema_definition JSONB DEFAULT NULL;

-- Create import_records table to track imports
CREATE TABLE IF NOT EXISTS import_records (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    data_source_id UUID NOT NULL REFERENCES data_sources(id),
    file_name VARCHAR(255) NOT NULL,
    file_size BIGINT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'Processing' CHECK (status IN ('Processing', 'Completed', 'Failed')),
    row_count INTEGER DEFAULT 0,
    success_count INTEGER DEFAULT 0,
    error_count INTEGER DEFAULT 0,
    imported_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMP NOT NULL DEFAULT set_utc_timestamp(),
    updated_at TIMESTAMP NOT NULL DEFAULT set_utc_timestamp(),
    metadata JSONB DEFAULT NULL
);

-- Create raw_transactions table to store all incoming data in its original form
CREATE TABLE IF NOT EXISTS raw_transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    import_id UUID NOT NULL REFERENCES import_records(id),
    data_source_id UUID NOT NULL REFERENCES data_sources(id),
    row_number INTEGER NOT NULL,
    data JSONB NOT NULL,
    error_message TEXT DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT set_utc_timestamp()
);

-- Update unmatched_transactions to reference raw_transactions
ALTER TABLE unmatched_transactions
ADD COLUMN raw_transaction_id UUID REFERENCES raw_transactions(id);

-- Add indices for performance
CREATE INDEX IF NOT EXISTS idx_import_records_data_source_id ON import_records(data_source_id);
CREATE INDEX IF NOT EXISTS idx_import_records_status ON import_records(status);
CREATE INDEX IF NOT EXISTS idx_raw_transactions_import_id ON raw_transactions(import_id);
CREATE INDEX IF NOT EXISTS idx_raw_transactions_data_source_id ON raw_transactions(data_source_id);
CREATE INDEX IF NOT EXISTS idx_raw_transactions_data_gin ON raw_transactions USING gin (data);

-- +migrate Down
-- Remove indices
DROP INDEX IF EXISTS idx_raw_transactions_data_gin;
DROP INDEX IF EXISTS idx_raw_transactions_data_source_id;
DROP INDEX IF EXISTS idx_raw_transactions_import_id;
DROP INDEX IF EXISTS idx_import_records_status;
DROP INDEX IF EXISTS idx_import_records_data_source_id;

-- Remove column from unmatched_transactions
ALTER TABLE unmatched_transactions
DROP COLUMN IF EXISTS raw_transaction_id;

-- Drop raw_transactions table
DROP TABLE IF EXISTS raw_transactions;

-- Drop import_records table
DROP TABLE IF EXISTS import_records;

-- Remove schema_definition from data_sources
ALTER TABLE data_sources
DROP COLUMN IF EXISTS schema_definition; 