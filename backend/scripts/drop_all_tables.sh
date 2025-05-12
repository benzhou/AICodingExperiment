#!/bin/bash

# Navigate to the project root
cd "$(dirname "$0")/.."
ROOT_DIR=$(pwd)

# Create a PostgreSQL script to drop all tables
echo "Creating SQL script to drop all tables..."
cat > drop_all_tables.sql <<EOL
-- Drop tables in proper order to handle dependencies
DO \$\$
DECLARE
    table_rec RECORD;
BEGIN
    -- Disable triggers temporarily to avoid constraint errors
    SET session_replication_role = 'replica';
    
    -- Drop all public tables (except PostgreSQL system tables) in reverse order of dependencies
    FOR table_rec IN 
        SELECT tablename FROM pg_tables 
        WHERE schemaname = 'public'
        AND tablename != 'schema_migrations'
        ORDER BY tablename
    LOOP
        EXECUTE 'DROP TABLE IF EXISTS ' || quote_ident(table_rec.tablename) || ' CASCADE';
    END LOOP;
    
    -- Also drop the migrations table to start completely fresh
    DROP TABLE IF EXISTS schema_migrations;
    
    -- Re-enable triggers
    SET session_replication_role = 'origin';
    
    RAISE NOTICE 'All tables have been dropped';
END
\$\$;
EOL

# Check if DATABASE_URL environment variable is set
if [ -z "$DATABASE_URL" ]; then
    # Try to read it from .env file
    if [ -f ".env" ]; then
        export DATABASE_URL=$(grep DATABASE_URL .env | cut -d '=' -f2-)
    fi
fi

# Check if we have a DATABASE_URL to use
if [ -z "$DATABASE_URL" ]; then
    echo "Error: DATABASE_URL environment variable not set."
    echo "Please set it in your environment or in a .env file."
    exit 1
fi

# Run the SQL script
echo "Dropping all tables from database..."
psql "$DATABASE_URL" -f drop_all_tables.sql

# Clean up
rm drop_all_tables.sql

echo "All tables have been dropped. You can now start with a fresh database." 