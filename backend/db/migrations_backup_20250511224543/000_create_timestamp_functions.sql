-- +migrate Up
-- Create a function to ensure timestamps are in UTC
CREATE OR REPLACE FUNCTION set_utc_timestamp()
RETURNS TIMESTAMP AS $$
BEGIN
    RETURN (NOW() AT TIME ZONE 'UTC');
END;
$$ LANGUAGE plpgsql;

-- +migrate Down
DROP FUNCTION IF EXISTS set_utc_timestamp(); 