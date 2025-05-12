package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	// Get database connection string
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatal("DATABASE_URL environment variable not set")
	}

	// Connect to the database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	// Test database connection
	if err = db.Ping(); err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	fmt.Println("✓ Successfully connected to database")

	// Count users with NULL auth_provider
	var nullCount int
	err = db.QueryRow(`
		SELECT COUNT(*) 
		FROM users 
		WHERE auth_provider IS NULL
	`).Scan(&nullCount)

	if err != nil {
		log.Fatalf("Error counting NULL auth_provider users: %v", err)
	}

	fmt.Printf("Found %d users with NULL auth_provider\n", nullCount)

	// Update users with NULL auth_provider
	if nullCount > 0 {
		result, err := db.Exec(`
			UPDATE users 
			SET auth_provider = 'local' 
			WHERE auth_provider IS NULL
		`)

		if err != nil {
			log.Fatalf("Error updating users: %v", err)
		}

		rowsAffected, _ := result.RowsAffected()
		fmt.Printf("Updated %d users, setting auth_provider to 'local'\n", rowsAffected)
	} else {
		fmt.Println("No users need updating")
	}

	// Verify the fix
	var stillNullCount int
	err = db.QueryRow(`
		SELECT COUNT(*) 
		FROM users 
		WHERE auth_provider IS NULL
	`).Scan(&stillNullCount)

	if err != nil {
		log.Fatalf("Error verifying fix: %v", err)
	}

	if stillNullCount == 0 {
		fmt.Println("✓ Fix successful! No users with NULL auth_provider remain.")
	} else {
		fmt.Printf("! Fix incomplete. %d users still have NULL auth_provider.\n", stillNullCount)
	}
}
