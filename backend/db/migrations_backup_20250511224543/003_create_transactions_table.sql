-- +migrate Up
-- No-op migration - transactions table already exists in 002_create_transaction_tables.sql
-- This is to prevent conflicts with the existing table structure
SELECT 1;

-- +migrate Down
-- No-op migration - don't drop anything to avoid dependency issues
SELECT 1; 