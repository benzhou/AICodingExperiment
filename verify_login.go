package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
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

	// Test login credentials
	email := "admin@example.com"
	password := "securepassword123"

	// Find user by email
	fmt.Printf("Looking up user with email: %s\n", email)

	var (
		id           string
		name         string
		passwordHash string
		authProvider string
	)

	// Normalized email lookup
	normalizedEmail := strings.ToLower(strings.TrimSpace(email))

	query := `
		SELECT id, name, password_hash, auth_provider
		FROM users
		WHERE LOWER(email) = LOWER($1)
	`

	err = db.QueryRow(query, normalizedEmail).Scan(&id, &name, &passwordHash, &authProvider)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Fatalf("User not found with email: %s", email)
		}
		log.Fatalf("Database error: %v", err)
	}

	fmt.Printf("✓ Found user: %s (ID: %s)\n", name, id)
	fmt.Printf("  Auth Provider: %s\n", authProvider)

	// Verify password
	fmt.Println("Verifying password...")

	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	if err != nil {
		log.Fatalf("× Password verification failed: %v", err)
	}

	fmt.Println("✓ Password verification successful!")
	fmt.Println("Login credentials are valid.")

	// Additional checks
	fmt.Println("\nChecking user roles...")

	var roleCount int
	err = db.QueryRow(`
		SELECT COUNT(*) 
		FROM user_roles 
		WHERE user_id = $1
	`, id).Scan(&roleCount)

	if err != nil {
		fmt.Printf("Error checking user roles: %v\n", err)
	} else {
		fmt.Printf("User has %d roles assigned\n", roleCount)

		// List the roles
		rows, err := db.Query(`
			SELECT role 
			FROM user_roles 
			WHERE user_id = $1
		`, id)

		if err != nil {
			fmt.Printf("Error retrieving roles: %v\n", err)
		} else {
			defer rows.Close()

			fmt.Println("Roles:")
			for rows.Next() {
				var role string
				if err := rows.Scan(&role); err != nil {
					fmt.Printf("  Error scanning role: %v\n", err)
					continue
				}
				fmt.Printf("  - %s\n", role)
			}
		}
	}

	fmt.Println("\nSummary:")
	fmt.Println("✓ Database connection: OK")
	fmt.Println("✓ User exists: YES")
	fmt.Println("✓ Password matches: YES")
	fmt.Println("✓ Login should succeed with these credentials")

	fmt.Println("\nIf login is still failing through the API, check:")
	fmt.Println("1. API URL format - ensure it's matching what the backend expects")
	fmt.Println("2. CORS configuration - ensure the frontend origin is allowed")
	fmt.Println("3. Request headers - Content-Type should be application/json")
	fmt.Println("4. Request payload format - ensure it matches what the backend expects")
}
