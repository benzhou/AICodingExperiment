package db

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// RunMigrations executes all SQL migrations in the migrations directory
func RunMigrations() error {
	if DB == nil {
		return fmt.Errorf("database connection not initialized")
	}

	log.Println("Running database migrations...")

	// Get the current working directory
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %v", err)
	}

	// Determine the migrations directory path
	migrationsDir := filepath.Join(wd, "db", "migrations")
	if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
		// Try one level up if not found (in case we're running from internal/db)
		wd = filepath.Dir(filepath.Dir(wd))
		migrationsDir = filepath.Join(wd, "db", "migrations")
		if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
			return fmt.Errorf("migrations directory not found at %s", migrationsDir)
		}
	}

	// Get list of migration files
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %v", err)
	}

	// Filter and sort migration files
	var migrationFiles []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") && file.Name() != "template.sql" {
			migrationFiles = append(migrationFiles, file.Name())
		}
	}
	sort.Strings(migrationFiles)

	// Create migrations table if it doesn't exist
	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMP NOT NULL DEFAULT NOW()
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %v", err)
	}

	// For each migration file
	for _, fileName := range migrationFiles {
		// Check if migration has already been applied
		var exists bool
		err = DB.QueryRow("SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version = $1)", fileName).Scan(&exists)
		if err != nil {
			return fmt.Errorf("failed to check if migration has been applied: %v", err)
		}

		if exists {
			log.Printf("Migration %s already applied, skipping", fileName)
			continue
		}

		// Read migration file
		filePath := filepath.Join(migrationsDir, fileName)
		upSQL, err := extractMigrationSection(filePath, "up")
		if err != nil {
			return fmt.Errorf("failed to extract UP section from migration file %s: %v", fileName, err)
		}

		// Skip empty migrations
		if strings.TrimSpace(upSQL) == "" {
			log.Printf("Migration %s has empty UP section, skipping", fileName)
			continue
		}

		// Execute migration in a transaction
		tx, err := DB.Begin()
		if err != nil {
			return fmt.Errorf("failed to begin transaction: %v", err)
		}

		log.Printf("Applying migration: %s", fileName)
		_, err = tx.Exec(upSQL)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to execute migration %s: %v", fileName, err)
		}

		// Record that migration has been applied
		_, err = tx.Exec("INSERT INTO schema_migrations (version) VALUES ($1)", fileName)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to record migration %s: %v", fileName, err)
		}

		// Commit transaction
		if err = tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit transaction: %v", err)
		}

		log.Printf("Successfully applied migration: %s", fileName)
	}

	log.Println("Database migrations completed successfully")
	return nil
}

// RollbackMigration rolls back the last applied migration
func RollbackMigration() error {
	if DB == nil {
		return fmt.Errorf("database connection not initialized")
	}

	// Get the last applied migration
	var lastMigration string
	err := DB.QueryRow("SELECT version FROM schema_migrations ORDER BY applied_at DESC LIMIT 1").Scan(&lastMigration)
	if err != nil {
		return fmt.Errorf("failed to get last migration: %v", err)
	}

	if lastMigration == "" {
		return fmt.Errorf("no migrations to roll back")
	}

	// Get the current working directory
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %v", err)
	}

	// Determine the migrations directory path
	migrationsDir := filepath.Join(wd, "db", "migrations")
	if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
		// Try one level up if not found (in case we're running from internal/db)
		wd = filepath.Dir(filepath.Dir(wd))
		migrationsDir = filepath.Join(wd, "db", "migrations")
		if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
			return fmt.Errorf("migrations directory not found at %s", migrationsDir)
		}
	}

	// Read migration file
	filePath := filepath.Join(migrationsDir, lastMigration)
	downSQL, err := extractMigrationSection(filePath, "down")
	if err != nil {
		return fmt.Errorf("failed to extract DOWN section from migration file %s: %v", lastMigration, err)
	}

	// Skip empty migrations
	if strings.TrimSpace(downSQL) == "" {
		return fmt.Errorf("migration %s has empty DOWN section", lastMigration)
	}

	// Execute migration in a transaction
	tx, err := DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}

	log.Printf("Rolling back migration: %s", lastMigration)
	_, err = tx.Exec(downSQL)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to execute rollback of migration %s: %v", lastMigration, err)
	}

	// Remove migration from applied list
	_, err = tx.Exec("DELETE FROM schema_migrations WHERE version = $1", lastMigration)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to remove migration %s from applied list: %v", lastMigration, err)
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	log.Printf("Successfully rolled back migration: %s", lastMigration)
	return nil
}

// extractMigrationSection extracts either the "up" or "down" section from a migration file
func extractMigrationSection(filePath string, section string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines []string
	var inSection bool

	upMarker := "-- +migrate Up"
	downMarker := "-- +migrate Down"

	// Convert section to lowercase for case-insensitive comparison
	section = strings.ToLower(section)

	for scanner.Scan() {
		line := scanner.Text()

		// Check for section markers
		if strings.Contains(line, upMarker) {
			inSection = (section == "up")
			continue
		} else if strings.Contains(line, downMarker) {
			inSection = (section == "down")
			continue
		}

		// Collect lines in the requested section
		if inSection {
			lines = append(lines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return strings.Join(lines, "\n"), nil
}
