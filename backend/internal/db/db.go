package db

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v\n", err)
	}

	// Check for a direct DATABASE_URL first
	connStr := os.Getenv("DATABASE_URL")
	
	// If no direct DATABASE_URL, construct from components
	if connStr == "" {
		dbUser := os.Getenv("DB_USER")
		dbPassword := os.Getenv("DB_PASSWORD")
		dbHost := os.Getenv("DB_HOST")
		if dbHost == "" {
			dbHost = "localhost"
		}
		dbPort := os.Getenv("DB_PORT")
		if dbPort == "" {
			dbPort = "5432"
		}
		dbName := os.Getenv("DB_NAME")
		if dbName == "" {
			dbName = "myapp"
		}
		
		// Construct connection string
		if dbUser != "" {
			if dbPassword != "" {
				connStr = "postgres://" + dbUser + ":" + dbPassword + "@" + dbHost + ":" + dbPort + "/" + dbName + "?sslmode=disable"
			} else {
				connStr = "postgres://" + dbUser + "@" + dbHost + ":" + dbPort + "/" + dbName + "?sslmode=disable"
			}
		}
	}
	
	if connStr == "" {
		log.Println("No DATABASE_URL provided, skipping database initialization")
		return
	}

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Printf("Error opening database: %v\n", err)
		return
	}

	if err = DB.Ping(); err != nil {
		log.Printf("Error connecting to the database: %v\n", err)
		return
	}

	log.Println("Successfully connected to database")

	// Run migrations
	if err = RunMigrations(); err != nil {
		log.Printf("Error running migrations: %v\n", err)
		return
	}
}
