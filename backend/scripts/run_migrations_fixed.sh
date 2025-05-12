#!/bin/bash

# Navigation to the project root
cd "$(dirname "$0")/.."
ROOT_DIR=$(pwd)

# Default direction is "up"
DIRECTION="up"
if [ "$1" == "down" ]; then
    DIRECTION="down"
    echo "CAUTION: Running migrations in DOWN direction (rollback mode)"
fi

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

if [ "$DIRECTION" == "up" ]; then
    # Process each migration file in order (ascending for up migrations)
    for migration in $(ls -1 $MIGRATIONS_DIR/*.sql | grep -v "template.sql" | sort); do
        filename=$(basename $migration)
        
        # Check if migration has already been applied
        applied=$(psql "$DATABASE_URL" -t -c "SELECT COUNT(*) FROM schema_migrations WHERE version = '$filename'")
        applied=$(echo $applied | tr -d ' ')
        
        if [ "$applied" -eq "0" ]; then
            echo "Applying migration: $filename (UP)"
            
            # Extract and execute only the UP section of the migration
            awk '/-- \+migrate Up/{flag=1;next}/-- \+migrate Down/{flag=0}flag' "$migration" > /tmp/up_migration.sql
            
            # Execute the UP migration
            if psql "$DATABASE_URL" -f "/tmp/up_migration.sql"; then
                # Record that migration has been applied
                psql "$DATABASE_URL" -c "INSERT INTO schema_migrations (version) VALUES ('$filename')"
                echo "Successfully applied migration: $filename (UP)"
            else
                echo "Error applying migration: $filename (UP)"
                exit 1
            fi
        else
            echo "Migration already applied, skipping: $filename"
        fi
    done
else
    # Process each migration file in reverse order (descending for down migrations)
    for migration in $(ls -1 $MIGRATIONS_DIR/*.sql | grep -v "template.sql" | sort -r); do
        filename=$(basename $migration)
        
        # Check if migration has been applied
        applied=$(psql "$DATABASE_URL" -t -c "SELECT COUNT(*) FROM schema_migrations WHERE version = '$filename'")
        applied=$(echo $applied | tr -d ' ')
        
        if [ "$applied" -eq "1" ]; then
            echo "Rolling back migration: $filename (DOWN)"
            
            # Extract and execute only the DOWN section of the migration
            awk '/-- \+migrate Down/{flag=1;next}flag' "$migration" > /tmp/down_migration.sql
            
            # Execute the DOWN migration
            if psql "$DATABASE_URL" -f "/tmp/down_migration.sql"; then
                # Remove migration from the applied list
                psql "$DATABASE_URL" -c "DELETE FROM schema_migrations WHERE version = '$filename'"
                echo "Successfully rolled back migration: $filename (DOWN)"
            else
                echo "Error rolling back migration: $filename (DOWN)"
                exit 1
            fi
        else
            echo "Migration not applied, skipping rollback: $filename"
        fi
    done
fi

echo "Migrations completed successfully!" 