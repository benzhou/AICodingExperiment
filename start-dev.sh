#!/bin/bash

# Added debug echo
echo "Starting script..."

# Function to check if a command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Check if air is installed
if ! command_exists air; then
    echo "Installing air for Go hot-reloading..."
    go install github.com/cosmtrek/air@latest
    echo "Air installed successfully"
else
    echo "Air is already installed"
fi

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