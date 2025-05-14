#!/bin/bash
# Script to kill Air and Go backend server processes

# Print colorful output
print_info() {
    echo -e "\033[1;34m[INFO]\033[0m $1"
}

print_success() {
    echo -e "\033[1;32m[SUCCESS]\033[0m $1"
}

print_warning() {
    echo -e "\033[1;33m[WARNING]\033[0m $1"
}

# Check for tmux sessions
check_tmux_sessions() {
    if command -v tmux >/dev/null 2>&1; then
        print_info "Checking for tmux sessions..."
        tmux list-sessions 2>/dev/null || echo "No tmux sessions found."
    else
        print_warning "tmux is not installed."
    fi
}

# Kill tmux session
kill_tmux_session() {
    if command -v tmux >/dev/null 2>&1; then
        print_info "Checking for dev tmux session..."
        if tmux has-session -t dev 2>/dev/null; then
            print_warning "Killing tmux session 'dev'..."
            tmux kill-session -t dev
            print_success "Tmux session killed."
        else
            print_info "No tmux 'dev' session found."
        fi
    else
        print_warning "tmux is not installed."
    fi
}

# Find and kill backend processes
kill_backend_processes() {
    print_info "Checking for Go and Air processes..."
    
    # Look for the main Go server first (more specific binary name)
    main_pid=$(lsof -i:8080 -t 2>/dev/null)
    if [ -n "$main_pid" ]; then
        print_warning "Found Go server process with PID $main_pid. Killing..."
        kill -9 $main_pid
        print_success "Go server process killed."
    else
        print_info "No Go server process found on port 8080."
    fi
    
    # Look for air processes
    air_pids=$(pgrep -f 'air' 2>/dev/null)
    if [ -n "$air_pids" ]; then
        print_warning "Found air processes. Killing..."
        for pid in $air_pids; do
            print_warning "Killing air process with PID $pid"
            kill -9 $pid
        done
        print_success "Air processes killed."
    else
        print_info "No air processes found."
    fi
    
    # Look for tmp/main processes (compiled by air)
    tmp_main_pids=$(pgrep -f 'tmp/main' 2>/dev/null)
    if [ -n "$tmp_main_pids" ]; then
        print_warning "Found tmp/main processes. Killing..."
        for pid in $tmp_main_pids; do
            print_warning "Killing tmp/main process with PID $pid"
            kill -9 $pid
        done
        print_success "tmp/main processes killed."
    else
        print_info "No tmp/main processes found."
    fi
    
    # Final check for any lingering processes on port 8080
    sleep 1
    if lsof -i:8080 >/dev/null 2>&1; then
        print_warning "Processes still found on port 8080. Attempting final cleanup..."
        lsof -i:8080 -t | xargs kill -9 2>/dev/null
        print_info "Final cleanup completed."
    else
        print_success "Port 8080 is now free."
    fi
}

# Main execution
print_info "Starting cleanup of Go server and Air processes..."

# 1. Check tmux sessions
check_tmux_sessions

# 2. Ask to kill tmux session
read -p "Do you want to kill the tmux session? (y/n, default: y): " kill_tmux
kill_tmux=${kill_tmux:-y}
if [[ "$kill_tmux" =~ ^[Yy]$ ]]; then
    kill_tmux_session
fi

# 3. Kill backend processes
kill_backend_processes

print_success "Cleanup completed!" 