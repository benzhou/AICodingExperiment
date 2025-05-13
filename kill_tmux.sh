#!/bin/bash

# Script to kill tmux sessions
# filepath: /Users/ben.zhou/workspace/AICodingExperiment/kill_tmux.sh

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

# Check if tmux is installed
if ! command -v tmux >/dev/null 2>&1; then
    print_error "Tmux is not installed."
    exit 1
fi

# Check if there are any sessions
if ! tmux ls >/dev/null 2>&1; then
    print_info "No tmux sessions running."
    exit 0
fi

# List all sessions
print_info "Current tmux sessions:"
tmux ls

# Function to kill a specific session
kill_specific_session() {
    local session_name="$1"
    if tmux has-session -t "$session_name" 2>/dev/null; then
        print_info "Killing tmux session '$session_name'..."
        tmux kill-session -t "$session_name"
        print_success "Tmux session '$session_name' killed."
    else
        print_error "Tmux session '$session_name' not found."
    fi
}

# Function to kill all sessions
kill_all_sessions() {
    print_info "Killing all tmux sessions..."
    tmux ls | cut -d ':' -f 1 | xargs -I{} tmux kill-session -t {}
    print_success "All tmux sessions killed."
}

# Main menu
echo
echo "1. Kill the 'dev' session (used by start-dev.sh)"
echo "2. Kill a specific session"
echo "3. Kill all sessions"
echo "4. Exit without killing any sessions"
echo

read -p "Select an option (1-4): " option

case "$option" in
    1)
        kill_specific_session "dev"
        ;;
    2)
        read -p "Enter the name of the session to kill: " session_name
        kill_specific_session "$session_name"
        ;;
    3)
        kill_all_sessions
        ;;
    4)
        print_info "Exiting without killing any sessions."
        ;;
    *)
        print_error "Invalid option. Exiting."
        ;;
esac

# Check if there are any remaining sessions
if tmux ls >/dev/null 2>&1; then
    print_info "Remaining tmux sessions:"
    tmux ls
else
    print_success "No tmux sessions remaining."
fi
