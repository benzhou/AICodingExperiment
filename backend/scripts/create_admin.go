package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Get path to root directory
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting working directory: %v", err)
	}

	// Check if we're in the scripts directory and adjust path accordingly
	if filepath.Base(wd) == "scripts" {
		wd = filepath.Dir(wd)
	}

	// Load .env file
	envPath := filepath.Join(wd, ".env")
	if err := godotenv.Load(envPath); err != nil {
		log.Printf("Warning: .env file not found at %s: %v\n", envPath, err)
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

	// User details (can be changed or passed as command line arguments)
	email := "admin@example.com"
	name := "Admin User"
	password := "admin123" // Should be changed after first login

	// Check if user already exists
	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", email).Scan(&exists)
	if err != nil {
		log.Fatalf("Error checking if user exists: %v", err)
	}

	if exists {
		log.Printf("User with email %s already exists", email)

		// Check if user already has admin role
		var userId string
		var hasAdminRole bool

		err = db.QueryRow("SELECT id FROM users WHERE email = $1", email).Scan(&userId)
		if err != nil {
			log.Fatalf("Error getting user ID: %v", err)
		}

		err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM user_roles WHERE user_id = $1 AND role = 'admin')", userId).Scan(&hasAdminRole)
		if err != nil {
			log.Fatalf("Error checking admin role: %v", err)
		}

		if hasAdminRole {
			log.Printf("User already has admin role")
			return
		}

		// Assign admin role to existing user
		_, err = db.Exec("INSERT INTO user_roles (user_id, role) VALUES ($1, 'admin')", userId)
		if err != nil {
			log.Fatalf("Error assigning admin role: %v", err)
		}

		log.Printf("Admin role assigned to existing user with email %s", email)
		return
	}

	// Create new user with admin role
	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("Error starting transaction: %v", err)
	}

	// Generate password hash
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		tx.Rollback()
		log.Fatalf("Error hashing password: %v", err)
	}

	// Generate UUID for user
	userId := uuid.New().String()

	// Insert user
	_, err = tx.Exec(
		"INSERT INTO users (id, email, name, password_hash, auth_provider) VALUES ($1, $2, $3, $4, 'local')",
		userId, email, name, string(hashedPassword),
	)
	if err != nil {
		tx.Rollback()
		log.Fatalf("Error creating user: %v", err)
	}

	// Assign admin role
	_, err = tx.Exec(
		"INSERT INTO user_roles (user_id, role) VALUES ($1, 'admin')",
		userId,
	)
	if err != nil {
		tx.Rollback()
		log.Fatalf("Error assigning admin role: %v", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		log.Fatalf("Error committing transaction: %v", err)
	}

	log.Printf("Admin user created successfully with email %s and password %s", email, password)
	fmt.Println("=========================================================")
	fmt.Println("IMPORTANT: Remember to change the admin password after first login!")
	fmt.Println("=========================================================")
}
