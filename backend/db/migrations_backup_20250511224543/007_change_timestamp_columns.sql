-- +migrate Up
-- Change all timestamp columns to use timestamp without time zone

-- Update the default value function to use UTC time
CREATE OR REPLACE FUNCTION set_utc_timestamp()
RETURNS TIMESTAMP AS $$
BEGIN
    RETURN (NOW() AT TIME ZONE 'UTC');
END;
$$ LANGUAGE plpgsql;

-- Helper function to check if a table exists
CREATE OR REPLACE FUNCTION table_exists(tbl_name text) RETURNS boolean AS $$
BEGIN
    RETURN EXISTS (
        SELECT 1 FROM information_schema.tables 
        WHERE table_schema = 'public' 
        AND table_name = tbl_name
    );
END;
$$ LANGUAGE plpgsql;

-- Update timestamp columns for each table if it exists

DO $$
BEGIN
    -- users table
    IF table_exists('users') THEN
        ALTER TABLE users 
            ALTER COLUMN created_at TYPE TIMESTAMP,
            ALTER COLUMN updated_at TYPE TIMESTAMP,
            ALTER COLUMN created_at SET DEFAULT set_utc_timestamp(),
            ALTER COLUMN updated_at SET DEFAULT set_utc_timestamp();
    END IF;

    -- data_sources table
    IF table_exists('data_sources') THEN
        ALTER TABLE data_sources 
            ALTER COLUMN created_at TYPE TIMESTAMP,
            ALTER COLUMN updated_at TYPE TIMESTAMP,
            ALTER COLUMN created_at SET DEFAULT set_utc_timestamp(),
            ALTER COLUMN updated_at SET DEFAULT set_utc_timestamp();
    END IF;

    -- match_rules table
    IF table_exists('match_rules') THEN
        ALTER TABLE match_rules 
            ALTER COLUMN created_at TYPE TIMESTAMP,
            ALTER COLUMN updated_at TYPE TIMESTAMP,
            ALTER COLUMN created_at SET DEFAULT set_utc_timestamp(),
            ALTER COLUMN updated_at SET DEFAULT set_utc_timestamp();
    END IF;

    -- transaction_matches table
    IF table_exists('transaction_matches') THEN
        ALTER TABLE transaction_matches 
            ALTER COLUMN approval_date TYPE TIMESTAMP,
            ALTER COLUMN created_at TYPE TIMESTAMP,
            ALTER COLUMN updated_at TYPE TIMESTAMP,
            ALTER COLUMN created_at SET DEFAULT set_utc_timestamp(),
            ALTER COLUMN updated_at SET DEFAULT set_utc_timestamp();
    END IF;

    -- transactions_audit table
    IF table_exists('transactions_audit') THEN
        ALTER TABLE transactions_audit 
            ALTER COLUMN changed_at TYPE TIMESTAMP,
            ALTER COLUMN changed_at SET DEFAULT set_utc_timestamp();
    END IF;

    -- tenants table
    IF table_exists('tenants') THEN
        ALTER TABLE tenants 
            ALTER COLUMN created_at TYPE TIMESTAMP,
            ALTER COLUMN updated_at TYPE TIMESTAMP,
            ALTER COLUMN created_at SET DEFAULT set_utc_timestamp(),
            ALTER COLUMN updated_at SET DEFAULT set_utc_timestamp();
    END IF;

    -- tenant_users table
    IF table_exists('tenant_users') THEN
        ALTER TABLE tenant_users 
            ALTER COLUMN created_at TYPE TIMESTAMP,
            ALTER COLUMN created_at SET DEFAULT set_utc_timestamp();
    END IF;

    -- role_permissions table
    IF table_exists('role_permissions') THEN
        ALTER TABLE role_permissions 
            ALTER COLUMN created_at TYPE TIMESTAMP,
            ALTER COLUMN created_at SET DEFAULT set_utc_timestamp();
    END IF;

    -- match_sets table
    IF table_exists('match_sets') THEN
        ALTER TABLE match_sets 
            ALTER COLUMN created_at TYPE TIMESTAMP,
            ALTER COLUMN updated_at TYPE TIMESTAMP,
            ALTER COLUMN created_at SET DEFAULT set_utc_timestamp(),
            ALTER COLUMN updated_at SET DEFAULT set_utc_timestamp();
    END IF;

    -- match_set_data_sources table
    IF table_exists('match_set_data_sources') THEN
        ALTER TABLE match_set_data_sources 
            ALTER COLUMN created_at TYPE TIMESTAMP,
            ALTER COLUMN created_at SET DEFAULT set_utc_timestamp();
    END IF;

    -- matched_transactions table
    IF table_exists('matched_transactions') THEN
        ALTER TABLE matched_transactions 
            ALTER COLUMN created_at TYPE TIMESTAMP,
            ALTER COLUMN created_at SET DEFAULT set_utc_timestamp();
    END IF;

    -- unmatched_transactions table
    IF table_exists('unmatched_transactions') THEN
        ALTER TABLE unmatched_transactions 
            ALTER COLUMN created_at TYPE TIMESTAMP,
            ALTER COLUMN created_at SET DEFAULT set_utc_timestamp();
    END IF;

    -- data_source_schemas table
    IF table_exists('data_source_schemas') THEN
        ALTER TABLE data_source_schemas 
            ALTER COLUMN created_at TYPE TIMESTAMP,
            ALTER COLUMN updated_at TYPE TIMESTAMP,
            ALTER COLUMN created_at SET DEFAULT set_utc_timestamp(),
            ALTER COLUMN updated_at SET DEFAULT set_utc_timestamp();
    END IF;

    -- schema_mappings table
    IF table_exists('schema_mappings') THEN
        ALTER TABLE schema_mappings 
            ALTER COLUMN created_at TYPE TIMESTAMP,
            ALTER COLUMN updated_at TYPE TIMESTAMP,
            ALTER COLUMN created_at SET DEFAULT set_utc_timestamp(),
            ALTER COLUMN updated_at SET DEFAULT set_utc_timestamp();
    END IF;

    -- file_parsing_configs table
    IF table_exists('file_parsing_configs') THEN
        ALTER TABLE file_parsing_configs 
            ALTER COLUMN created_at TYPE TIMESTAMP,
            ALTER COLUMN updated_at TYPE TIMESTAMP,
            ALTER COLUMN created_at SET DEFAULT set_utc_timestamp(),
            ALTER COLUMN updated_at SET DEFAULT set_utc_timestamp();
    END IF;

    -- import_records table
    IF table_exists('import_records') THEN
        ALTER TABLE import_records 
            ALTER COLUMN created_at TYPE TIMESTAMP,
            ALTER COLUMN updated_at TYPE TIMESTAMP,
            ALTER COLUMN created_at SET DEFAULT set_utc_timestamp(),
            ALTER COLUMN updated_at SET DEFAULT set_utc_timestamp();
    END IF;

    -- raw_transactions table
    IF table_exists('raw_transactions') THEN
        ALTER TABLE raw_transactions 
            ALTER COLUMN created_at TYPE TIMESTAMP,
            ALTER COLUMN created_at SET DEFAULT set_utc_timestamp();
    END IF;
END $$;

-- +migrate Down
-- Revert back to timestamptz if needed
-- Note: This will lose timezone information when reverting from timestamp to timestamptz

DO $$
BEGIN
    -- users table
    IF table_exists('users') THEN
        ALTER TABLE users 
            ALTER COLUMN created_at TYPE TIMESTAMPTZ,
            ALTER COLUMN updated_at TYPE TIMESTAMPTZ,
            ALTER COLUMN created_at SET DEFAULT NOW(),
            ALTER COLUMN updated_at SET DEFAULT NOW();
    END IF;

    -- data_sources table
    IF table_exists('data_sources') THEN
        ALTER TABLE data_sources 
            ALTER COLUMN created_at TYPE TIMESTAMPTZ,
            ALTER COLUMN updated_at TYPE TIMESTAMPTZ,
            ALTER COLUMN created_at SET DEFAULT NOW(),
            ALTER COLUMN updated_at SET DEFAULT NOW();
    END IF;

    -- match_rules table
    IF table_exists('match_rules') THEN
        ALTER TABLE match_rules 
            ALTER COLUMN created_at TYPE TIMESTAMPTZ,
            ALTER COLUMN updated_at TYPE TIMESTAMPTZ,
            ALTER COLUMN created_at SET DEFAULT NOW(),
            ALTER COLUMN updated_at SET DEFAULT NOW();
    END IF;

    -- transaction_matches table
    IF table_exists('transaction_matches') THEN
        ALTER TABLE transaction_matches 
            ALTER COLUMN approval_date TYPE TIMESTAMPTZ,
            ALTER COLUMN created_at TYPE TIMESTAMPTZ,
            ALTER COLUMN updated_at TYPE TIMESTAMPTZ,
            ALTER COLUMN created_at SET DEFAULT NOW(),
            ALTER COLUMN updated_at SET DEFAULT NOW();
    END IF;

    -- transactions_audit table
    IF table_exists('transactions_audit') THEN
        ALTER TABLE transactions_audit 
            ALTER COLUMN changed_at TYPE TIMESTAMPTZ,
            ALTER COLUMN changed_at SET DEFAULT NOW();
    END IF;

    -- tenants table
    IF table_exists('tenants') THEN
        ALTER TABLE tenants 
            ALTER COLUMN created_at TYPE TIMESTAMPTZ,
            ALTER COLUMN updated_at TYPE TIMESTAMPTZ,
            ALTER COLUMN created_at SET DEFAULT NOW(),
            ALTER COLUMN updated_at SET DEFAULT NOW();
    END IF;

    -- tenant_users table
    IF table_exists('tenant_users') THEN
        ALTER TABLE tenant_users 
            ALTER COLUMN created_at TYPE TIMESTAMPTZ,
            ALTER COLUMN created_at SET DEFAULT NOW();
    END IF;

    -- role_permissions table
    IF table_exists('role_permissions') THEN
        ALTER TABLE role_permissions 
            ALTER COLUMN created_at TYPE TIMESTAMPTZ,
            ALTER COLUMN created_at SET DEFAULT NOW();
    END IF;

    -- match_sets table
    IF table_exists('match_sets') THEN
        ALTER TABLE match_sets 
            ALTER COLUMN created_at TYPE TIMESTAMPTZ,
            ALTER COLUMN updated_at TYPE TIMESTAMPTZ,
            ALTER COLUMN created_at SET DEFAULT NOW(),
            ALTER COLUMN updated_at SET DEFAULT NOW();
    END IF;

    -- match_set_data_sources table
    IF table_exists('match_set_data_sources') THEN
        ALTER TABLE match_set_data_sources 
            ALTER COLUMN created_at TYPE TIMESTAMPTZ,
            ALTER COLUMN created_at SET DEFAULT NOW();
    END IF;

    -- matched_transactions table
    IF table_exists('matched_transactions') THEN
        ALTER TABLE matched_transactions 
            ALTER COLUMN created_at TYPE TIMESTAMPTZ,
            ALTER COLUMN created_at SET DEFAULT NOW();
    END IF;

    -- unmatched_transactions table
    IF table_exists('unmatched_transactions') THEN
        ALTER TABLE unmatched_transactions 
            ALTER COLUMN created_at TYPE TIMESTAMPTZ,
            ALTER COLUMN created_at SET DEFAULT NOW();
    END IF;

    -- data_source_schemas table
    IF table_exists('data_source_schemas') THEN
        ALTER TABLE data_source_schemas 
            ALTER COLUMN created_at TYPE TIMESTAMPTZ,
            ALTER COLUMN updated_at TYPE TIMESTAMPTZ,
            ALTER COLUMN created_at SET DEFAULT NOW(),
            ALTER COLUMN updated_at SET DEFAULT NOW();
    END IF;

    -- schema_mappings table
    IF table_exists('schema_mappings') THEN
        ALTER TABLE schema_mappings 
            ALTER COLUMN created_at TYPE TIMESTAMPTZ,
            ALTER COLUMN updated_at TYPE TIMESTAMPTZ,
            ALTER COLUMN created_at SET DEFAULT NOW(),
            ALTER COLUMN updated_at SET DEFAULT NOW();
    END IF;

    -- file_parsing_configs table
    IF table_exists('file_parsing_configs') THEN
        ALTER TABLE file_parsing_configs 
            ALTER COLUMN created_at TYPE TIMESTAMPTZ,
            ALTER COLUMN updated_at TYPE TIMESTAMPTZ,
            ALTER COLUMN created_at SET DEFAULT NOW(),
            ALTER COLUMN updated_at SET DEFAULT NOW();
    END IF;

    -- import_records table
    IF table_exists('import_records') THEN
        ALTER TABLE import_records 
            ALTER COLUMN created_at TYPE TIMESTAMPTZ,
            ALTER COLUMN updated_at TYPE TIMESTAMPTZ,
            ALTER COLUMN created_at SET DEFAULT NOW(),
            ALTER COLUMN updated_at SET DEFAULT NOW();
    END IF;

    -- raw_transactions table
    IF table_exists('raw_transactions') THEN
        ALTER TABLE raw_transactions 
            ALTER COLUMN created_at TYPE TIMESTAMPTZ,
            ALTER COLUMN created_at SET DEFAULT NOW();
    END IF;
END $$;

-- Drop the functions
DROP FUNCTION IF EXISTS set_utc_timestamp();
DROP FUNCTION IF EXISTS table_exists(text); 