#!/bin/bash

# Added debug echo
echo "Starting script..."

# Function to kill tmux sessions
kill_tmux_session() {
    if command -v tmux >/dev/null 2>&1; then
        echo "Checking for existing tmux sessions..."
        
        # Check if there are any sessions
        if ! tmux ls >/dev/null 2>&1; then
            echo "No tmux sessions running."
            return
        fi
        
        # List all sessions
        tmux ls
        
        # Ask if user wants to kill the dev session
        if tmux has-session -t dev 2>/dev/null; then
            read -p "Kill the 'dev' session? (y/n, default: y): " kill_dev
            kill_dev=${kill_dev:-y}
            if [[ "$kill_dev" =~ ^[Yy]$ ]]; then
                echo "Killing tmux session 'dev'..."
                tmux kill-session -t dev
                echo "Tmux session 'dev' killed."
            fi
        fi
        
        # Ask if user wants to kill other sessions
        read -p "Do you want to kill other tmux sessions? (y/n, default: n): " kill_others
        if [[ "$kill_others" =~ ^[Yy]$ ]]; then
            read -p "Enter session name to kill (or 'all' for all remaining sessions): " session_name
            if [ "$session_name" = "all" ]; then
                echo "Killing all tmux sessions..."
                tmux ls | cut -d ':' -f 1 | xargs -I{} tmux kill-session -t {}
                echo "All tmux sessions killed."
            elif [ -n "$session_name" ]; then
                if tmux has-session -t "$session_name" 2>/dev/null; then
                    tmux kill-session -t "$session_name"
                    echo "Tmux session '$session_name' killed."
                else
                    echo "Tmux session '$session_name' not found."
                fi
            fi
        fi
    else
        echo "Tmux is not installed."
    fi
}

# Check if setup_dev_environment.sh has been run
if [ ! -f "./backend/.env" ] || ! command -v air >/dev/null 2>&1 || ! command -v psql >/dev/null 2>&1 || ! command -v node >/dev/null 2>&1; then
    echo "It appears that not all required dependencies are installed or configured."
    echo "Please run ./setup_dev_environment.sh first to set up your development environment."
    
    read -p "Would you like to run the setup script now? (y/n) " choice
    case "$choice" in
        y|Y) 
            ./setup_dev_environment.sh
            ;;
        *) 
            echo "Exiting. Please run ./setup_dev_environment.sh before running this script."
            exit 1
            ;;
    esac
fi

# Check if user wants to manage tmux sessions
read -p "Do you want to check and manage existing tmux sessions? (y/n, default: n): " manage_tmux
if [[ "$manage_tmux" =~ ^[Yy]$ ]]; then
    kill_tmux_session
fi

# Function to check if a command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Check if air is installed
if ! command_exists air; then
    echo "Installing air for Go hot-reloading..."
    go install github.com/cosmtrek/air@v1.40.4
    
    # Add Go bin to PATH for this session if it exists and is not already in PATH
    GO_BIN_PATH=$(go env GOPATH)/bin
    if [ -d "$GO_BIN_PATH" ] && [[ ! ":$PATH:" == *":$GO_BIN_PATH:"* ]]; then
        export PATH="$GO_BIN_PATH:$PATH"
        echo "Added $GO_BIN_PATH to PATH for this session"
        echo "For permanent use, add 'export PATH=\$HOME/go/bin:\$PATH' to your ~/.zshrc file"
    fi
    
    # Check again if air is now in PATH
    if ! command_exists air; then
        echo "Error: air was installed but cannot be found in PATH."
        echo "Please run: export PATH=\$HOME/go/bin:\$PATH"
        exit 1
    fi
    
    echo "Air installed successfully"
else
    echo "Air is already installed"
fi

# Create admin user if needed
read -p "Do you want to create/update the admin user? (y/n, default: n) " create_admin
create_admin=${create_admin:-n} # Set default to 'n' if input is empty
case "$create_admin" in
    y|Y) 
        # Collect admin user details
        read -p "Enter admin email (default: admin@example.com): " admin_email
        admin_email=${admin_email:-admin@example.com}
        
        read -p "Enter admin name (default: Admin User): " admin_name
        admin_name=${admin_name:-Admin User}
        
        read -s -p "Enter admin password (default: admin123): " admin_password
        admin_password=${admin_password:-admin123}
        echo
        
        # Run the create_admin script
        echo "Creating admin user..."
        (cd backend/scripts && chmod +x create_admin_fixed.sh && ./create_admin_fixed.sh "$admin_email" "$admin_name" "$admin_password")
        if [ $? -eq 0 ]; then
            echo "Admin user created/updated successfully."
        else
            echo "Failed to create admin user. Check the errors above."
            read -p "Continue anyway? (y/n) " continue_anyway
            if [[ ! "$continue_anyway" =~ ^[Yy]$ ]]; then
                echo "Exiting."
                exit 1
            fi
        fi
        ;;
    n|N)
        echo "Skipping admin user creation."
        ;;
    *)
        echo "Invalid input '$create_admin'. Skipping admin user creation."
        ;;
esac

# Create a new tmux session
if command_exists tmux; then
    echo "Tmux is installed, using tmux for development"
    
    # Kill existing session if it exists
    echo "Checking for existing tmux session named 'dev'"
    tmux has-session -t dev 2>/dev/null
    if [ $? -eq 0 ]; then
        echo "Killing existing 'dev' session"
        tmux kill-session -t dev
    else
        echo "No existing 'dev' session found"
    fi

    # Create new session
    echo "Creating new tmux session 'dev'"
    tmux new-session -d -s dev
    if [ $? -ne 0 ]; then
        echo "Failed to create tmux session. Error code: $?"
        echo "Trying alternative approach..."
        rm -rf /private/tmp/tmux-*
        echo "Trying to create session again after clearing temp files"
        tmux new-session -d -s dev
        if [ $? -ne 0 ]; then
            echo "Still failing to create tmux session. Starting separate terminals instead."
            # Start backend
            echo "Starting backend in separate terminal"
            osascript -e 'tell app "Terminal" to do script "cd '$(pwd)'/backend && air"' &
            
            # Start frontend
            echo "Starting frontend in separate terminal"
            osascript -e 'tell app "Terminal" to do script "cd '$(pwd)'/frontend && npm start"' &
            exit 0
        fi
    fi

    # Split window horizontally
    echo "Splitting tmux window"
    tmux split-window -h

    # Select first pane and start backend
    echo "Starting backend in first pane"
    tmux select-pane -t 0
    tmux send-keys "cd backend && air" C-m

    # Select second pane and start frontend
    echo "Starting frontend in second pane"
    tmux select-pane -t 1
    tmux send-keys "cd frontend && npm start" C-m

    # Attach to the session
    echo "Attaching to tmux session"
    tmux attach-session -t dev
else
    # If tmux is not installed, run in separate terminals
    echo "Tmux is not installed"
    echo "For better experience, install tmux"
    echo "Starting applications in separate terminals..."
    
    # Start backend
    echo "Starting backend in separate terminal"
    osascript -e 'tell app "Terminal" to do script "cd '$(pwd)'/backend && air"' &
    
    # Start frontend
    echo "Starting frontend in separate terminal"
    osascript -e 'tell app "Terminal" to do script "cd '$(pwd)'/frontend && npm start"' &
fi