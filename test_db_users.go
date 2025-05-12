package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v\n", err)
	}

	// Get database connection string
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable not set")
	}

	// Connect to database
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}

	// Get all users
	rows, err := db.Query("SELECT id, email, name, password_hash FROM users")
	if err != nil {
		log.Fatalf("Error querying users: %v", err)
	}
	defer rows.Close()

	fmt.Println("Users in the database:")
	fmt.Println("======================")

	passwords := []string{
		"securepassword123",
		"admin123",
		"password",
		"123456",
	}

	// Process each user
	for rows.Next() {
		var id, email, name, passwordHash string
		if err := rows.Scan(&id, &email, &name, &passwordHash); err != nil {
			log.Fatalf("Error scanning row: %v", err)
		}

		fmt.Printf("ID: %s\n", id)
		fmt.Printf("Email: %s\n", email)
		fmt.Printf("Name: %s\n", name)
		fmt.Printf("Password Hash: %s\n", passwordHash)

		// Try common passwords
		fmt.Println("Testing passwords:")
		for _, password := range passwords {
			err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
			if err == nil {
				fmt.Printf("  ✓ Password '%s' MATCHES!\n", password)
			} else {
				fmt.Printf("  ✗ Password '%s' does not match\n", password)
			}
		}

		fmt.Println("======================")
	}

	if err := rows.Err(); err != nil {
		log.Fatalf("Error iterating rows: %v", err)
	}
}
