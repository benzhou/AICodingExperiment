#!/bin/bash

# Function to check if a command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Check if air is installed
if ! command_exists air; then
    echo "Installing air for Go hot-reloading..."
    go install github.com/cosmtrek/air@latest
fi

# Create a new tmux session
if command_exists tmux; then
    # Kill existing session if it exists
    tmux kill-session -t dev 2>/dev/null

    # Create new session
    tmux new-session -d -s dev

    # Split window horizontally
    tmux split-window -h

    # Select first pane and start backend
    tmux select-pane -t 0
    tmux send-keys "cd backend && air" C-m

    # Select second pane and start frontend
    tmux select-pane -t 1
    tmux send-keys "cd frontend && npm start" C-m

    # Attach to the session
    tmux attach-session -t dev
else
    # If tmux is not installed, run in separate terminals
    echo "For better experience, install tmux"
    echo "Starting applications in separate terminals..."
    
    # Start backend
    osascript -e 'tell app "Terminal" to do script "cd '$(pwd)'/backend && air"' &
    
    # Start frontend
    osascript -e 'tell app "Terminal" to do script "cd '$(pwd)'/frontend && npm start"' &
fi 