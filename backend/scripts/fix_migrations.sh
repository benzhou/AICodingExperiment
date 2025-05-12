#!/bin/bash

# Navigate to the project root
cd "$(dirname "$0")/.."
ROOT_DIR=$(pwd)

# Directory for migrations
MIGRATIONS_DIR="db/migrations"

# Create a backup of the original migrations
echo "Creating backup of original migrations..."
BACKUP_DIR="${MIGRATIONS_DIR}_backup_$(date +%Y%m%d%H%M%S)"
mkdir -p "$BACKUP_DIR"
cp $MIGRATIONS_DIR/*.sql "$BACKUP_DIR/"
echo "Backup created at: $BACKUP_DIR"

# Fix all migrations to use direct UTC formatting
echo "Updating migrations with direct UTC timestamp..."
find $MIGRATIONS_DIR -name "*.sql" -type f | xargs sed -i "" 's/DEFAULT set_utc_timestamp()/DEFAULT (NOW() AT TIME ZONE '\''UTC'\'')/g'

echo "Migrations have been updated to use direct UTC timestamp conversion."
echo "Run './scripts/drop_all_tables.sh' followed by './scripts/run_migrations.sh' to apply the schema." 