#!/bin/bash

# CreatorSync Analytics Setup Script
# This script sets up the analytics system components

set -e

echo "🚀 Setting up CreatorSync Analytics System..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Helper functions
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if running from project root
if [ ! -f "backend/go.mod" ] || [ ! -f "frontend/package.json" ]; then
    print_error "Please run this script from the project root directory"
    exit 1
fi

# Check prerequisites
print_status "Checking prerequisites..."

# Check Go
if ! command -v go &> /dev/null; then
    print_error "Go is not installed. Please install Go 1.21 or later."
    exit 1
fi

# Check Node.js
if ! command -v node &> /dev/null; then
    print_error "Node.js is not installed. Please install Node.js 18 or later."
    exit 1
fi

# Check PostgreSQL
if ! command -v psql &> /dev/null; then
    print_warning "PostgreSQL client not found. Make sure PostgreSQL is running."
fi

print_success "Prerequisites check completed"

# Install backend dependencies
print_status "Installing backend dependencies..."
cd backend
go mod download
cd ..
print_success "Backend dependencies installed"

# Install frontend dependencies
print_status "Installing frontend dependencies..."
cd frontend
npm install
cd ..
print_success "Frontend dependencies installed"

# Check for environment files
print_status "Checking environment configuration..."

if [ ! -f "backend/.env" ]; then
    print_warning "Backend .env file not found. Creating from example..."
    if [ -f "backend/.env.example" ]; then
        cp backend/.env.example backend/.env
        print_warning "Please edit backend/.env with your configuration"
    else
        print_error "No .env.example file found in backend/"
    fi
fi

if [ ! -f "frontend/.env" ]; then
    print_warning "Frontend .env file not found."
    print_warning "Please create frontend/.env with your Clerk configuration"
fi

# Database setup
print_status "Setting up database..."

# Check if DATABASE_URL is provided (Railway style)
if [ ! -z "$DATABASE_URL" ]; then
    print_success "Using DATABASE_URL for connection"
    DB_CONNECTION_TYPE="DATABASE_URL"
else
    # Fallback to individual environment variables
    DB_NAME=${POSTGRES_DB_DATABASE:-railway}
    DB_HOST=${POSTGRES_DB_HOST:-localhost}
    DB_PORT=${POSTGRES_DB_PORT:-5432}
    DB_USER=${POSTGRES_DB_USERNAME:-postgres}
    DB_CONNECTION_TYPE="INDIVIDUAL_VARS"
    print_status "Using individual database environment variables"
fi

if command -v psql &> /dev/null; then
    # Test database connection
    if [ "$DB_CONNECTION_TYPE" = "DATABASE_URL" ]; then
        # Test with DATABASE_URL
        if psql "$DATABASE_URL" -c '\q' 2>/dev/null; then
            print_success "Database connection successful (using DATABASE_URL)"
            CONNECTION_SUCCESS=true
        else
            print_warning "Could not connect using DATABASE_URL"
            CONNECTION_SUCCESS=false
        fi
    else
        # Test with individual variables
        if PGPASSWORD=$POSTGRES_DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c '\q' 2>/dev/null; then
            print_success "Database connection successful (using individual variables)"
            CONNECTION_SUCCESS=true
        else
            print_warning "Could not connect using individual database variables"
            CONNECTION_SUCCESS=false
        fi
    fi
    
    if [ "$CONNECTION_SUCCESS" = true ]; then
        # Run migrations
        print_status "Running database migrations..."
        cd backend
        if [ -f "cmd/migrate/main.go" ]; then
            go run cmd/migrate/main.go
            print_success "Database migrations completed"
        else
            print_warning "Migration script not found. Please run migrations manually."
        fi
        cd ..
    else
        print_warning "Could not connect to database. Please ensure PostgreSQL is running and configured correctly."
    fi
else
    print_warning "PostgreSQL client not available. Please run migrations manually."
fi

# Build backend
print_status "Building backend..."
cd backend
go build -o ../bin/api cmd/api/main.go
cd ..
print_success "Backend built successfully"

# Build frontend
print_status "Building frontend..."
cd frontend
npm run build
cd ..
print_success "Frontend built successfully"

# Create systemd service files (optional)
create_systemd_services() {
    print_status "Creating systemd service files..."
    
    # Backend service
    cat > /tmp/creatorsync-api.service << EOF
[Unit]
Description=CreatorSync API Server
After=network.target postgresql.service

[Service]
Type=simple
User=www-data
WorkingDirectory=$(pwd)
ExecStart=$(pwd)/bin/api
Restart=always
RestartSec=5
Environment=NODE_ENV=production

[Install]
WantedBy=multi-user.target
EOF

    # Frontend service (if using standalone mode)
    cat > /tmp/creatorsync-frontend.service << EOF
[Unit]
Description=CreatorSync Frontend Server
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=$(pwd)/frontend
ExecStart=/usr/bin/npm start
Restart=always
RestartSec=5
Environment=NODE_ENV=production

[Install]
WantedBy=multi-user.target
EOF

    print_success "Systemd service files created in /tmp/"
    print_warning "To install services, run:"
    print_warning "sudo cp /tmp/creatorsync-*.service /etc/systemd/system/"
    print_warning "sudo systemctl daemon-reload"
    print_warning "sudo systemctl enable creatorsync-api creatorsync-frontend"
}

# Development setup
setup_development() {
    print_status "Setting up development environment..."
    
    # Install Air for Go hot reloading
    if ! command -v air &> /dev/null; then
        print_status "Installing Air for Go hot reloading..."
        go install github.com/cosmtrek/air@latest
    fi
    
    # Create development start script
    cat > start-dev.sh << 'EOF'
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
EOF

    chmod +x start-dev.sh
    print_success "Development start script created: ./start-dev.sh"
}

# Production setup
setup_production() {
    print_status "Setting up production environment..."
    
    # Create production start script
    cat > start-prod.sh << 'EOF'
#!/bin/bash

# Start production servers
echo "Starting CreatorSync production environment..."

# Start backend
./bin/api &
BACKEND_PID=$!

# Start frontend (if using standalone mode)
cd frontend && npm start &
FRONTEND_PID=$!

# Function to cleanup on exit
cleanup() {
    echo "Stopping production servers..."
    kill $BACKEND_PID $FRONTEND_PID 2>/dev/null
    exit
}

# Trap cleanup function on script exit
trap cleanup EXIT

echo "Production servers running:"
echo "- Backend: http://localhost:8080"
echo "- Frontend: http://localhost:3000"
echo "Press Ctrl+C to stop all servers"

wait
EOF

    chmod +x start-prod.sh
    print_success "Production start script created: ./start-prod.sh"
}

# Analytics-specific setup
setup_analytics() {
    print_status "Setting up analytics components..."
    
    # Verify analytics tables exist
    if command -v psql &> /dev/null && PGPASSWORD=$POSTGRES_DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c '\q' 2>/dev/null; then
        # Check if analytics tables exist
        TABLE_COUNT=$(PGPASSWORD=$POSTGRES_DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -t -c "SELECT COUNT(*) FROM information_schema.tables WHERE table_name IN ('channel_analytics', 'stream_sessions', 'video_analytics', 'analytics_jobs');" 2>/dev/null | tr -d ' ')
        
        if [ "$TABLE_COUNT" = "4" ]; then
            print_success "Analytics tables found in database"
        else
            print_warning "Analytics tables not found. Please run migrations."
        fi
    fi
    
    # Create analytics test script
    cat > test-analytics.sh << 'EOF'
#!/bin/bash

# Test analytics endpoints
echo "Testing analytics system..."

API_BASE="http://localhost:8080"

echo "1. Testing health endpoint..."
curl -s "$API_BASE/api/analytics/health" | jq .

echo -e "\n2. Testing overview endpoint (requires auth)..."
echo "Please set TOKEN environment variable with your auth token"
if [ ! -z "$TOKEN" ]; then
    curl -s -H "Authorization: Bearer $TOKEN" "$API_BASE/api/analytics/overview" | jq .
else
    echo "Skipping - no TOKEN set"
fi

echo -e "\nAnalytics test completed"
EOF

    chmod +x test-analytics.sh
    print_success "Analytics test script created: ./test-analytics.sh"
}

# Main setup process
main() {
    print_status "Starting main setup process..."
    
    # Ask user for setup type
    echo ""
    echo "Choose setup type:"
    echo "1) Development setup"
    echo "2) Production setup"
    echo "3) Both"
    read -p "Enter choice (1-3): " setup_choice
    
    case $setup_choice in
        1)
            setup_development
            ;;
        2)
            setup_production
            create_systemd_services
            ;;
        3)
            setup_development
            setup_production
            create_systemd_services
            ;;
        *)
            print_warning "Invalid choice. Setting up development environment."
            setup_development
            ;;
    esac
    
    # Always setup analytics components
    setup_analytics
    
    print_success "Setup completed successfully!"
    
    echo ""
    echo "📋 Next Steps:"
    echo "1. Configure your .env files with API keys and database credentials"
    echo "2. Ensure PostgreSQL is running and accessible"
    echo "3. Set up your Twitch application at https://dev.twitch.tv/console"
    echo "4. Configure Clerk authentication in your frontend .env"
    echo "5. Run database migrations if not done automatically"
    echo ""
    echo "🚀 To start development:"
    echo "   ./start-dev.sh"
    echo ""
    echo "🏭 To start production:"
    echo "   ./start-prod.sh"
    echo ""
    echo "🧪 To test analytics:"
    echo "   export TOKEN=your_auth_token"
    echo "   ./test-analytics.sh"
    echo ""
    echo "📖 For detailed documentation, see: ANALYTICS_ARCHITECTURE.md"
}

# Run main setup
main

print_success "Analytics setup script completed!" 