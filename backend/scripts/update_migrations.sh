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

# Remove the timestamp migration
echo "Removing timestamp migration file..."
rm -f "$MIGRATIONS_DIR"/007_change_timestamp_columns.sql

# Rename new migration files to their final names
echo "Renaming new migration files..."
for file in "$MIGRATIONS_DIR"/*.sql.new; do
  if [ -f "$file" ]; then
    newname=$(echo "$file" | sed 's/\.new$//')
    mv "$file" "$newname"
    echo "Renamed: $file to $newname"
  fi
done

echo "Migration files have been updated with UTC timestamps without timezone."
echo "Run './scripts/drop_all_tables.sh' followed by server restart to apply the new schema." 