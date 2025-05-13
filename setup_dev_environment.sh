#!/bin/bash

# setup_dev_environment.sh
# Script to set up the development environment for AICodingExperiment

# Print colorful text
print_info() {
    echo -e "\033[1;34m[INFO]\033[0m $1"
}

print_error() {
    echo -e "\033[1;31m[ERROR]\033[0m $1"
}

print_success() {
    echo -e "\033[1;32m[SUCCESS]\033[0m $1"
}

print_warning() {
    echo -e "\033[1;33m[WARNING]\033[0m $1"
}

# Check if a command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Check and install Homebrew if needed
install_homebrew() {
    if ! command_exists brew; then
        print_info "Installing Homebrew..."
        /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
        if [[ "$(uname -m)" == "arm64" ]]; then
            print_info "Setting up Homebrew for Apple Silicon..."
            echo 'eval "$(/opt/homebrew/bin/brew shellenv)"' >> ~/.zprofile
            eval "$(/opt/homebrew/bin/brew shellenv)"
        else
            print_info "Setting up Homebrew for Intel Mac..."
            echo 'eval "$(/usr/local/bin/brew shellenv)"' >> ~/.zprofile
            eval "$(/usr/local/bin/brew shellenv)"
        fi
    else
        print_info "Homebrew is already installed."
    fi
}

# Check and install Go if needed
install_go() {
    if ! command_exists go; then
        print_info "Installing Go..."
        brew install go
    else
        print_info "Go is already installed."
    fi
    
    # Check Go version
    GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    REQUIRED_VERSION="1.20"
    
    if [ "$(printf '%s\n' "$REQUIRED_VERSION" "$GO_VERSION" | sort -V | head -n1)" != "$REQUIRED_VERSION" ]; then
        print_warning "Go version $GO_VERSION detected. This project requires Go $REQUIRED_VERSION or higher."
        print_info "Updating Go..."
        brew upgrade go
    else
        print_success "Go version $GO_VERSION is compatible."
    fi
}

# Check and install Node.js and npm if needed
install_node() {
    if ! command_exists node; then
        print_info "Installing Node.js..."
        brew install node
    else
        print_info "Node.js is already installed."
    fi
    
    if ! command_exists npm; then
        print_info "Installing npm..."
        brew install npm
    else
        print_info "npm is already installed."
    fi
}

# Check and install PostgreSQL if needed
install_postgres() {
    if ! command_exists psql; then
        print_info "Installing PostgreSQL..."
        brew install postgresql@15
        brew services start postgresql@15
        
        # Add PostgreSQL to PATH if not already in there
        if [[ ! ":$PATH:" == *":/opt/homebrew/opt/postgresql@15/bin:"* ]]; then
            print_info "Adding PostgreSQL to PATH..."
            echo 'export PATH="/opt/homebrew/opt/postgresql@15/bin:$PATH"' >> ~/.zshrc
            export PATH="/opt/homebrew/opt/postgresql@15/bin:$PATH"
        fi
        
        # Wait for PostgreSQL to start
        sleep 5
    else
        print_info "PostgreSQL is already installed."
        
        # Ensure PostgreSQL is running
        if ! brew services list | grep postgresql | grep started > /dev/null; then
            print_info "Starting PostgreSQL service..."
            brew services start postgresql@15 || brew services start postgresql
        else
            print_success "PostgreSQL service is already running."
        fi
    fi
}

# Check and install air for hot reloading if needed
install_air() {
    if ! command_exists air; then
        print_info "Installing air for Go hot-reloading..."
        go install github.com/cosmtrek/air@v1.40.4
        
        # Get Go binary path
        GO_BIN_PATH=$(go env GOPATH)/bin
        
        # Add Go bin to PATH if not already there
        if [[ ! ":$PATH:" == *":$GO_BIN_PATH:"* ]]; then
            print_info "Adding Go bin directory to PATH..."
            echo 'export PATH="'$GO_BIN_PATH':$PATH"' >> ~/.zshrc
            export PATH="$GO_BIN_PATH:$PATH"
            
            print_info "Please run 'source ~/.zshrc' after this script completes or start a new terminal session."
        fi
    else
        print_info "air is already installed."
    fi
}

# Check and install tmux if needed
install_tmux() {
    if ! command_exists tmux; then
        print_info "Installing tmux..."
        brew install tmux
    else
        print_info "tmux is already installed."
    fi
}

# Set up database for the application
setup_database() {
    print_info "Setting up the database..."
    
    # Get database configuration
    read -p "Enter PostgreSQL username (default: postgres): " DB_USER
    DB_USER=${DB_USER:-postgres}
    
    read -s -p "Enter PostgreSQL password (leave empty for no password): " DB_PASSWORD
    echo
    
    read -p "Enter database name (default: myapp): " DB_NAME
    DB_NAME=${DB_NAME:-myapp}
    
    # Check if the user exists
    if ! psql postgres -c "SELECT 1 FROM pg_roles WHERE rolname='$DB_USER'" | grep -q 1; then
        print_info "Creating database user $DB_USER..."
        if [ -z "$DB_PASSWORD" ]; then
            createuser -s $DB_USER
        else
            psql postgres -c "CREATE ROLE $DB_USER WITH LOGIN CREATEDB SUPERUSER PASSWORD '$DB_PASSWORD';"
        fi
    else
        print_success "Database user $DB_USER already exists."
    fi
    
    # Check if the database exists
    if ! psql -l | grep -q "$DB_NAME"; then
        print_info "Creating database $DB_NAME..."
        createdb -O $DB_USER $DB_NAME
    else
        print_success "Database $DB_NAME already exists."
    fi
    
    # Create or update .env file
    ENV_FILE="./backend/.env"
    if [ -z "$DB_PASSWORD" ]; then
        echo "DATABASE_URL=postgres://$DB_USER@localhost:5432/$DB_NAME?sslmode=disable" > $ENV_FILE
    else
        echo "DATABASE_URL=postgres://$DB_USER:$DB_PASSWORD@localhost:5432/$DB_NAME?sslmode=disable" > $ENV_FILE
    fi
    echo "FRONTEND_URL=http://localhost:3000" >> $ENV_FILE
    echo "JWT_SECRET=development-secret-key-replace-in-production" >> $ENV_FILE
    echo "JWT_EXPIRY_MINUTES=60" >> $ENV_FILE
    
    print_success "Database configuration saved to $ENV_FILE"
    
    # Run database migrations
    print_info "Running database migrations..."
    cd backend && chmod +x scripts/run_migrations.sh && ./scripts/run_migrations.sh
    cd ..
}

# Install frontend dependencies
setup_frontend() {
    print_info "Setting up frontend dependencies..."
    cd frontend && npm install
    cd ..
    print_success "Frontend dependencies installed."
}

# Create admin user
create_admin_user() {
    print_info "Setting up admin user..."
    
    # Ask if admin user should be created
    read -p "Do you want to create an admin user? (y/n, default: y): " CREATE_ADMIN
    CREATE_ADMIN=${CREATE_ADMIN:-y}
    
    if [[ "$CREATE_ADMIN" =~ ^[Yy]$ ]]; then
        # Get admin user details
        read -p "Enter admin email (default: admin@example.com): " ADMIN_EMAIL
        ADMIN_EMAIL=${ADMIN_EMAIL:-admin@example.com}
        
        read -p "Enter admin name (default: Admin User): " ADMIN_NAME
        ADMIN_NAME=${ADMIN_NAME:-Admin User}
        
        read -s -p "Enter admin password (default: admin123): " ADMIN_PASSWORD
        ADMIN_PASSWORD=${ADMIN_PASSWORD:-admin123}
        echo
        
        # Run the create_admin script
        print_info "Creating admin user..."
        (cd backend/scripts && chmod +x create_admin_fixed.sh && ./create_admin_fixed.sh "$ADMIN_EMAIL" "$ADMIN_NAME" "$ADMIN_PASSWORD")
        
        if [ $? -eq 0 ]; then
            print_success "Admin user created successfully."
        else
            print_error "Failed to create admin user. See errors above."
        fi
    else
        print_info "Skipping admin user creation."
    fi
}

# Utility function to kill tmux sessions
kill_tmux_session() {
    if command_exists tmux; then
        print_info "Checking for existing tmux sessions..."
        
        # List all sessions
        if tmux ls >/dev/null 2>&1; then
            tmux ls
            
            # Ask which session to kill
            read -p "Enter session name to kill (or 'all' for all sessions, leave empty to skip): " SESSION_NAME
            
            if [ -n "$SESSION_NAME" ]; then
                if [ "$SESSION_NAME" = "all" ]; then
                    print_info "Killing all tmux sessions..."
                    tmux ls | cut -d ':' -f 1 | xargs -I{} tmux kill-session -t {}
                    print_success "All tmux sessions killed."
                else
                    print_info "Killing tmux session '$SESSION_NAME'..."
                    if tmux has-session -t "$SESSION_NAME" 2>/dev/null; then
                        tmux kill-session -t "$SESSION_NAME"
                        print_success "Tmux session '$SESSION_NAME' killed."
                    else
                        print_error "Tmux session '$SESSION_NAME' not found."
                    fi
                fi
            else
                print_info "Skipping tmux session termination."
            fi
        else
            print_info "No tmux sessions running."
        fi
    else
        print_error "Tmux is not installed."
    fi
}

# Main function to run all setup steps
main() {
    print_info "Starting development environment setup..."
    
    # Check OS
    if [[ "$OSTYPE" != "darwin"* ]]; then
        print_error "This script is designed for macOS. Your OS is $OSTYPE."
        exit 1
    fi
    
    # Check if user wants to kill existing tmux sessions
    read -p "Do you want to check for and kill existing tmux sessions? (y/n, default: n): " KILL_TMUX
    if [[ "$KILL_TMUX" =~ ^[Yy]$ ]]; then
        kill_tmux_session
    fi
    
    # Install dependencies
    install_homebrew
    install_go
    install_node
    install_postgres
    install_air
    install_tmux
    
    # Setup application
    setup_database
    setup_frontend
    create_admin_user
    
    print_success "Setup complete! You can now run ./start-dev.sh to start the application."
}

# Run the main function
main
