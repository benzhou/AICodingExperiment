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

	connStr := os.Getenv("DATABASE_URL")
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
