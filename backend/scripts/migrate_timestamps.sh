#!/bin/bash

# Navigate to the project root
cd "$(dirname "$0")/.."
ROOT_DIR=$(pwd)

# Create a small Go program to run the migration
echo "Creating temporary migration runner..."
cat > tmp_migration_runner.go <<EOL
package main

import (
	"backend/internal/db"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file from the project root
	if err := godotenv.Load(filepath.Join(".", ".env")); err != nil {
		log.Printf("Warning: .env file not found: %v\n", err)
	}

	// Print current directory to check path
	dir, _ := os.Getwd()
	fmt.Println("Current directory:", dir)
	fmt.Println("Looking for migrations in:", filepath.Join(dir, "db", "migrations"))

	// Initialize DB which runs migrations
	db.InitDB()
}
EOL

# Run the migration
echo "Running timestamp migration..."
go run tmp_migration_runner.go

# Clean up
rm tmp_migration_runner.go

echo "Migration complete!"
echo "All timestamps are now stored as UTC TIMESTAMP without timezone." 