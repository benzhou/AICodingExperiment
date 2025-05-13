#!/bin/bash
# Utility script to check for and kill development processes
# filepath: /Users/ben.zhou/workspace/AICodingExperiment/check_processes.sh

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

# Function to check for running processes
check_processes() {
    print_info "Checking for development processes..."
    
    # Check for Air processes
    air_processes=$(ps aux | grep -E "air.*backend" | grep -v grep)
    if [ -n "$air_processes" ]; then
        print_warning "Found Air processes running:"
        echo "$air_processes"
    else
        print_success "No Air processes found."
    fi
    
    # Check for backend API processes
    backend_processes=$(ps aux | grep -E "(/backend/tmp/main|go run.*backend)" | grep -v grep)
    if [ -n "$backend_processes" ]; then
        print_warning "Found Go backend API processes running:"
        echo "$backend_processes"
        
        # Check port 8080
        port_check=$(lsof -i :8080 2>/dev/null)
        if [ -n "$port_check" ]; then
            print_warning "Port 8080 is in use:"
            echo "$port_check"
        fi
    else
        print_success "No Go backend API processes found."
    fi
    
    # Check for npm processes
    npm_processes=$(ps aux | grep "npm start" | grep -v grep)
    if [ -n "$npm_processes" ]; then
        print_warning "Found npm processes running:"
        echo "$npm_processes"
    else
        print_success "No npm processes found."
    fi
    
    # Check for node processes that might be related to our app
    node_processes=$(ps aux | grep "node.*frontend" | grep -v grep)
    if [ -n "$node_processes" ]; then
        print_warning "Found Node.js processes that might be related to the frontend:"
        echo "$node_processes"
        
        # Check port 3000
        port_check=$(lsof -i :3000 2>/dev/null)
        if [ -n "$port_check" ]; then
            print_warning "Port 3000 is in use:"
            echo "$port_check"
        fi
    else
        print_success "No related Node.js processes found."
    fi
}

# Function to kill processes
kill_processes() {
    print_info "Finding process IDs to kill..."
    
    # Find and kill Air processes
    air_pids=$(ps aux | grep -E "air.*backend" | grep -v grep | awk '{print $2}')
    if [ -n "$air_pids" ]; then
        for pid in $air_pids; do
            print_warning "Killing Air process with PID $pid"
            kill -9 $pid
        done
        print_success "Air processes killed."
    else
        print_info "No Air processes to kill."
    fi
    
    # Find and kill backend API processes
    backend_pids=$(ps aux | grep -E "(/backend/tmp/main|go run.*backend)" | grep -v grep | awk '{print $2}')
    if [ -n "$backend_pids" ]; then
        for pid in $backend_pids; do
            print_warning "Killing Go backend API process with PID $pid"
            kill -9 $pid
        done
        print_success "Go backend API processes killed."
    else
        print_info "No Go backend API processes to kill."
    fi
    
    # Find and kill npm processes
    npm_pids=$(ps aux | grep "npm start" | grep -v grep | awk '{print $2}')
    if [ -n "$npm_pids" ]; then
        for pid in $npm_pids; do
            print_warning "Killing npm process with PID $pid"
            kill -9 $pid
        done
        print_success "npm processes killed."
    else
        print_info "No npm processes to kill."
    fi
    
    # Find and kill related node processes
    node_pids=$(ps aux | grep "node.*frontend" | grep -v grep | awk '{print $2}')
    if [ -n "$node_pids" ]; then
        for pid in $node_pids; do
            print_warning "Killing Node.js process with PID $pid"
            kill -9 $pid
        done
        print_success "Related Node.js processes killed."
    else
        print_info "No related Node.js processes to kill."
    fi
}

# Main menu
echo "======================"
echo "Development Process Manager"
echo "======================"
echo
echo "1. Check for running development processes"
echo "2. Kill running development processes"
echo "3. Exit"
echo

read -p "Select an option (1-3): " option

case "$option" in
    1)
        check_processes
        ;;
    2)
        check_processes
        echo
        read -p "Do you want to kill these processes? (y/n): " confirm
        if [[ "$confirm" =~ ^[Yy]$ ]]; then
            kill_processes
        else
            print_info "No processes were killed."
        fi
        ;;
    3)
        print_info "Exiting."
        ;;
    *)
        print_error "Invalid option. Exiting."
        ;;
esac
