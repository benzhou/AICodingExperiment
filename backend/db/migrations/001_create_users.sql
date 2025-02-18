CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),,
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    password_hash TEXT NOT NULL,
    auth_provider VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Create index on email
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- Create index on id (if not already a primary key)
CREATE INDEX IF NOT EXISTS idx_users_id ON users(id);