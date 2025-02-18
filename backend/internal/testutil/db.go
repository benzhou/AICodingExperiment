package testutil

import (
	"database/sql"
	"os"
	"testing"
)

func SetupTestDB(t *testing.T) *sql.DB {
	connStr := os.Getenv("TEST_DATABASE_URL")
	if connStr == "" {
		connStr = "postgres://postgres:postgres@localhost:5432/myapp_test?sslmode=disable"
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Clear test database
	if _, err := db.Exec(`TRUNCATE users CASCADE`); err != nil {
		t.Fatalf("Failed to clear test database: %v", err)
	}

	return db
}
