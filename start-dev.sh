#!/bin/bash

# Start development servers
echo "Starting CreatorSync development environment..."

# Start backend with hot reloading
cd backend && air &
BACKEND_PID=$!

# Start frontend development server
cd frontend && npm run dev &
FRONTEND_PID=$!

# Function to cleanup on exit
cleanup() {
    echo "Stopping development servers..."
    kill $BACKEND_PID $FRONTEND_PID 2>/dev/null
    exit
}

# Trap cleanup function on script exit
trap cleanup EXIT

# Wait for user input to stop
echo "Development servers running:"
echo "- Backend: http://localhost:8080"
echo "- Frontend: http://localhost:3000"
echo "Press Ctrl+C to stop all servers"

wait
