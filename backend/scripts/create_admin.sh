#!/bin/bash

# Navigate to the project root
cd "$(dirname "$0")/.."
ROOT_DIR=$(pwd)

# Check if email was provided as command line argument
if [ "$1" != "" ]; then
    ADMIN_EMAIL="$1"
else
    ADMIN_EMAIL="admin@example.com"
fi

# Check if name was provided as command line argument
if [ "$2" != "" ]; then
    ADMIN_NAME="$2"
else
    ADMIN_NAME="Admin User"
fi

# Check if password was provided as command line argument
if [ "$3" != "" ]; then
    ADMIN_PASSWORD="$3"
else
    ADMIN_PASSWORD="admin123"
fi

# Create a small Go program to create an admin user
echo "Creating admin user creation program..."
cat > create_admin.go <<EOL
package main

import (
        "database/sql"
        "log"
        "os"
        "path/filepath"
        "time"

        "github.com/google/uuid"
        "github.com/joho/godotenv"
        _ "github.com/lib/pq"
        "golang.org/x/crypto/bcrypt"
)

func main() {
        // Load .env file from the project root
        if err := godotenv.Load(filepath.Join(".", ".env")); err != nil {
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

        // Create user_roles table if it doesn't exist
        _, err = db.Exec(\`
                CREATE TABLE IF NOT EXISTS user_roles (
                        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                        user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                        role VARCHAR(50) NOT NULL CHECK (role IN ('preparer', 'approver', 'admin')),
                        created_at TIMESTAMP NOT NULL DEFAULT (NOW() AT TIME ZONE 'UTC'),
                        updated_at TIMESTAMP NOT NULL DEFAULT (NOW() AT TIME ZONE 'UTC'),
                        UNIQUE(user_id, role)
                );
        \`)
        if err != nil {
                log.Fatalf("Error creating user_roles table: %v", err)
        }

        // Create index if it doesn't exist
        _, err = db.Exec(\`
                CREATE INDEX IF NOT EXISTS idx_user_roles_user_id ON user_roles(user_id);
        \`)
        if err != nil {
                log.Printf("Warning: Could not create index: %v", err)
        }

        // User details (from command line arguments)
        email := "${ADMIN_EMAIL}"
        name := "${ADMIN_NAME}"
        password := "${ADMIN_PASSWORD}" // Should be changed after first login

        log.Printf("Creating admin user with email: %s, name: %s", email, name)

        // Check if user already exists
        var exists bool
        err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = \$1)", email).Scan(&exists)
        if err != nil {
                log.Fatalf("Error checking if user exists: %v", err)
        }

        var userID string
        if exists {
                log.Printf("User %s already exists, checking if they have admin role", email)
                // Get the user ID
                err = db.QueryRow("SELECT id FROM users WHERE email = \$1", email).Scan(&userID)
                if err != nil {
                        log.Fatalf("Error getting user ID: %v", err)
                }
        } else {
                // Generate password hash
                passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
                if err != nil {
                        log.Fatalf("Error generating password hash: %v", err)
                }

                // Generate UUID for user
                userID = uuid.New().String()

                // Insert user
                now := time.Now().UTC()
                _, err = db.Exec(
                        "INSERT INTO users (id, email, name, password_hash, created_at, updated_at) VALUES (\$1, \$2, \$3, \$4, \$5, \$6)",
                        userID, email, name, string(passwordHash), now, now,
                )
                if err != nil {
                        log.Fatalf("Error creating user: %v", err)
                }
                log.Printf("Created user %s with ID %s", email, userID)
        }

        // Check if user already has admin role
        var hasAdminRole bool
        err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM user_roles WHERE user_id = \$1 AND role = 'admin')", userID).Scan(&hasAdminRole)
        if err != nil {
                log.Fatalf("Error checking if user has admin role: %v", err)
        }

        if hasAdminRole {
                log.Printf("User %s already has admin role", email)
        } else {
                // Add admin role
                now := time.Now().UTC()
                _, err = db.Exec(
                        "INSERT INTO user_roles (user_id, role, created_at, updated_at) VALUES (\$1, 'admin', \$2, \$3)",
                        userID, now, now,
                )
                if err != nil {
                        log.Fatalf("Error adding admin role: %v", err)
                }
                log.Printf("Successfully added admin role to user %s", email)
        }

        // Add role permission if table exists
        var rolePermissionExists bool
        err = db.QueryRow("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'role_permissions')").Scan(&rolePermissionExists)
        if err != nil {
                log.Printf("Warning: Could not check if role_permissions table exists: %v", err)
        }

        if rolePermissionExists {
                // Add basic permissions for admin role
                permissions := []string{"create:datasource", "read:datasource", "update:datasource", "delete:datasource",
                        "create:user", "read:user", "update:user", "delete:user"}

                for _, perm := range permissions {
                        // Check if permission already exists
                        var permExists bool
                        err = db.QueryRow(\`
                                SELECT EXISTS(
                                        SELECT 1 FROM role_permissions 
                                        WHERE role_name = 'admin' AND permission = \$1
                                )
                        \`, perm).Scan(&permExists)

                        if err != nil {
                                log.Printf("Warning: Could not check if permission %s exists: %v", perm, err)
                                continue
                        }

                        if !permExists {
                                _, err = db.Exec(\`
                                        INSERT INTO role_permissions (role_name, permission, created_at) 
                                        VALUES ('admin', \$1, \$2)
                                \`, perm, time.Now().UTC())

                                if err != nil {
                                        log.Printf("Warning: Could not add permission %s: %v", perm, err)
                                } else {
                                        log.Printf("Added permission %s for admin role", perm)
                                }
                        }
                }
        }

        log.Printf("Successfully created admin user %s with password %s", email, password)
}
EOL

# Run the program
echo "Creating admin user..."
go run create_admin.go

# Clean up
rm create_admin.go

echo "Admin user creation complete!" 