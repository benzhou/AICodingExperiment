#!/bin/bash

# Start the backend server in the background
echo "Starting the backend server..."
cd /Users/benzhou/workspace/AICoding/backend
go run main.go &
SERVER_PID=$!

# Give the server a moment to start up
echo "Waiting for server to start..."
sleep 5

# Now attempt login with curl
echo "Attempting login with admin@example.com..."
curl -v -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com", "password":"securepassword123"}'

# Clean up by killing the server process
echo "Cleaning up..."
kill $SERVER_PID

echo "Test completed." 