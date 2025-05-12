#!/bin/bash

# Navigate to the project root
cd "$(dirname "$0")/.."
ROOT_DIR=$(pwd)

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

# Migration directory
MIGRATIONS_DIR="db/migrations"

# Ensure schema_migrations table exists
echo "Creating schema_migrations table if it doesn't exist..."
psql "$DATABASE_URL" -c "
    CREATE TABLE IF NOT EXISTS schema_migrations (
        version VARCHAR(255) PRIMARY KEY,
        applied_at TIMESTAMP NOT NULL DEFAULT NOW()
    );
"

# Process each migration file in order
for migration in $(ls -1 $MIGRATIONS_DIR/*.sql | grep -v "template.sql" | sort); do
    filename=$(basename $migration)
    
    # Check if migration has already been applied
    applied=$(psql "$DATABASE_URL" -t -c "SELECT COUNT(*) FROM schema_migrations WHERE version = '$filename'")
    applied=$(echo $applied | tr -d ' ')
    
    if [ "$applied" -eq "0" ]; then
        echo "Applying migration: $filename"
        
        # Execute the migration
        if psql "$DATABASE_URL" -f "$migration"; then
            # Record that migration has been applied
            psql "$DATABASE_URL" -c "INSERT INTO schema_migrations (version) VALUES ('$filename')"
            echo "Successfully applied migration: $filename"
        else
            echo "Error applying migration: $filename"
            exit 1
        fi
    else
        echo "Migration already applied, skipping: $filename"
    fi
done

echo "Migrations completed successfully!" 