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

	// Test credentials
	email := "admin@example.com"
	password := "securepassword123"

	// Debug and trace the login process
	fmt.Printf("===== LOGIN TRACE FOR %s =====\n", email)

	// Step 1: Find the user by email
	fmt.Println("Step 1: Finding user by email")
	user, err := findByEmail(db, email)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("  ✗ User not found")
			return
		}
		fmt.Printf("  ✗ Database error: %v\n", err)
		return
	}
	fmt.Println("  ✓ User found")
	fmt.Printf("  * ID: %s\n", user.ID)
	fmt.Printf("  * Email: %s\n", user.Email)
	fmt.Printf("  * Name: %s\n", user.Name)
	fmt.Printf("  * Hash: %s\n", user.PasswordHash)

	// Step 2: Check password
	fmt.Println("\nStep 2: Verifying password")
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		fmt.Printf("  ✗ Password verification FAILED: %v\n", err)

		// Additional debugging for password
		fmt.Println("\nPassword debugging:")
		fmt.Printf("  * Raw password: '%s'\n", password)
		fmt.Printf("  * Password length: %d\n", len(password))
		fmt.Printf("  * ASCII codes: %v\n", getASCIICodes(password))
		return
	}
	fmt.Println("  ✓ Password verification SUCCEEDED")

	// Step 3: Check if user has roles
	fmt.Println("\nStep 3: Checking user roles")
	roles, err := getUserRoles(db, user.ID)
	if err != nil {
		fmt.Printf("  ✗ Error getting roles: %v\n", err)
	} else {
		if len(roles) == 0 {
			fmt.Println("  ! No roles found for user")
		} else {
			fmt.Printf("  ✓ Found %d roles\n", len(roles))
			for i, role := range roles {
				fmt.Printf("    %d. %s\n", i+1, role)
			}
		}
	}

	fmt.Println("\n===== LOGIN TRACE COMPLETED =====")
	fmt.Println("Authentication should SUCCEED based on this test")
}

// User model
type User struct {
	ID           string
	Email        string
	Name         string
	PasswordHash string
}

// FindByEmail implementation based on repository
func findByEmail(db *sql.DB, email string) (*User, error) {
	user := &User{}
	query := `
		SELECT id, email, name, password_hash
		FROM users
		WHERE email = $1`

	err := db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.PasswordHash,
	)

	return user, err
}

// Get user roles
func getUserRoles(db *sql.DB, userID string) ([]string, error) {
	query := `
		SELECT role
		FROM user_roles
		WHERE user_id = $1`

	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []string
	for rows.Next() {
		var role string
		if err := rows.Scan(&role); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}

	return roles, rows.Err()
}

// Helper to get ASCII codes for debugging
func getASCIICodes(s string) []int {
	var codes []int
	for _, char := range s {
		codes = append(codes, int(char))
	}
	return codes
}

// Check if string has any invisible characters
func hasInvisibleChars(s string) bool {
	for _, char := range s {
		// Check for common invisible characters
		if char < 32 || (char >= 127 && char <= 160) {
			return true
		}
	}
	return false
}

// Helper to trim invisible characters
func cleanString(s string) string {
	return strings.TrimSpace(s)
}
