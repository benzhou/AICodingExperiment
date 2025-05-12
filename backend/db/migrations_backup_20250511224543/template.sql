-- Migration Template
-- Replace <migration_name> with a descriptive name

-- +migrate Up
-- Add your migration statements here

-- IMPORTANT: Always use TIMESTAMP (without timezone) for timestamp columns
-- Example:
-- CREATE TABLE IF NOT EXISTS example_table (
--     id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
--     name VARCHAR(100) NOT NULL,
--     created_at TIMESTAMP NOT NULL DEFAULT set_utc_timestamp(),
--     updated_at TIMESTAMP NOT NULL DEFAULT set_utc_timestamp()
-- );

-- +migrate Down
-- Add statements to revert the migration here 