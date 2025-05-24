# Project-level variables
APP_NAME=main
BACKEND_DIR=backend
CMD_ENTRY=$(BACKEND_DIR)/cmd/api/main.go
BIN_OUTPUT=$(APP_NAME)

# Default: build + test
all: build test

# Build backend from correct directory
build:
	@echo "Building..."
	cd $(BACKEND_DIR) && go build -o ../$(BIN_OUTPUT) ./cmd/api

# Run backend and frontend together (local development)
run:
	@echo "Running backend and frontend..."
	cd $(BACKEND_DIR) && go run ./cmd/api &
	npm install --prefer-offline --no-fund --prefix ./frontend
	npm run dev --prefix ./frontend

# Unit tests
test:
	@echo "Running tests..."
	cd $(BACKEND_DIR) && go test ./... -v

# Integration tests
itest:
	@echo "Running integration tests..."
	cd $(BACKEND_DIR) && go test ./internal/database -v

# Clean built binary
clean:
	@echo "Cleaning..."
	@rm -f $(BIN_OUTPUT)

# Live reload using Air
watch:
	@if command -v air > /dev/null; then \
		cd $(BACKEND_DIR) && air; \
	else \
		read -p "Go's 'air' is not installed. Install it now? [Y/n] " choice; \
		if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
			go install github.com/air-verse/air@latest; \
			cd $(BACKEND_DIR) && air; \
		else \
			echo "Air not installed. Exiting..."; \
			exit 1; \
		fi; \
	fi

# CreatorSync Development Tools

# Docker & Database
.PHONY: docker-run docker-down docker-logs docker-clean docker-local docker-local-down docker-local-logs docker-local-clean

# Production/Staging Docker commands (using existing docker-compose.yml)
docker-run:
	cd backend && docker compose up -d

docker-down:
	cd backend && docker compose down

docker-logs:
	cd backend && docker compose logs -f

docker-clean:
	cd backend && docker compose down -v --remove-orphans

# LOCAL DEVELOPMENT Docker commands (using docker-compose.local.yml)
docker-local:
	cd backend && docker compose -f docker-compose.local.yml up -d

docker-local-down:
	cd backend && docker compose -f docker-compose.local.yml down

docker-local-logs:
	cd backend && docker compose -f docker-compose.local.yml logs -f

docker-local-clean:
	cd backend && docker compose -f docker-compose.local.yml down -v --remove-orphans

# Database Migrations
.PHONY: migrate-up migrate-down migrate-reset migrate-local-up migrate-local-down migrate-local-reset

# Production/Staging migrations (use current DATABASE_URL from .env)
migrate-up:
	cd backend && go run cmd/migrate/main.go up

migrate-down:
	cd backend && go run cmd/migrate/main.go down

migrate-reset: migrate-down migrate-up

# LOCAL DEVELOPMENT migrations (uses local database)
migrate-local-up:
	cd backend && DATABASE_URL=postgresql://postgres:localdev123@localhost:5432/creatorsync_local go run cmd/migrate/main.go up

migrate-local-down:
	cd backend && DATABASE_URL=postgresql://postgres:localdev123@localhost:5432/creatorsync_local go run cmd/migrate/main.go down

migrate-local-reset: migrate-local-down migrate-local-up

# LOCAL DEVELOPMENT test data seeding
.PHONY: seed-local inspect-local
seed-local:
	cd backend && DATABASE_URL=postgresql://postgres:localdev123@localhost:5432/creatorsync_local go run cmd/seed/main.go

# LOCAL DEVELOPMENT database inspection
inspect-local:
	@echo "🔍 Inspecting local database..."
	@echo "Tables:"
	@PGPASSWORD=localdev123 psql -h localhost -U postgres -d creatorsync_local -c "\dt" 2>/dev/null || echo "❌ Database not accessible. Run: make docker-local"
	@echo ""
	@echo "Users:"
	@PGPASSWORD=localdev123 psql -h localhost -U postgres -d creatorsync_local -c "SELECT clerk_id, username, email, twitch_username FROM users LIMIT 5;" 2>/dev/null || echo "No users table or no data"
	@echo ""
	@echo "Analytics count:"
	@PGPASSWORD=localdev123 psql -h localhost -U postgres -d creatorsync_local -c "SELECT COUNT(*) as analytics_records FROM analytics;" 2>/dev/null || echo "No analytics table"

# Add your actual user to local database
.PHONY: add-my-user
add-my-user:
	@echo "👤 Adding your Clerk user to local database..."
	@PGPASSWORD=localdev123 psql -h localhost -U postgres -d creatorsync_local -f add-my-user.sql
	@echo "✅ User added! You should now be able to log in locally."

# Testing (New Tiered Approach)
.PHONY: test-fast test-full test-integration test-backend test-frontend
test-fast:
	npm run test:fast

test-full:
	npm run test:full

test-integration:
	npm run test:integration

test-backend:
	npm run test:backend:integration

test-frontend:
	npm run test:frontend:full

# Development
.PHONY: run-backend run-frontend dev dev-local dev-with-docker
run-backend:
	cd backend && go run cmd/api/main.go

run-frontend:
	cd frontend && npm run dev

dev: run

# LOCAL DEVELOPMENT with local database (recommended for development)
dev-local: docker-local
	@echo "🚀 Starting LOCAL DEVELOPMENT environment..."
	@echo "📝 Make sure you have .env.local configured with your Clerk development keys!"
	@echo "🗄️  Using LOCAL PostgreSQL database (test data only)"
	@make run-backend & make run-frontend & wait

# Development with staging/production Docker (use carefully)
dev-with-docker: docker-run
	@echo "⚠️  Starting with staging/production database..."
	@make run-backend & make run-frontend & wait

# Build
.PHONY: build-backend build-frontend
build: build-backend build-frontend

build-backend:
	cd backend && go build -o bin/api cmd/api/main.go

build-frontend:
	cd frontend && npm run build

# Production
.PHONY: prod start-prod
prod: build docker-run
	cd backend && ./bin/api

start-prod:
	cd frontend && npm start

# Staging Validation
.PHONY: staging-check staging-deploy
staging-check:
	@echo "🔍 Running comprehensive validation before staging..."
	npm run test:integration
	@echo "✅ All tests passed! Safe to deploy to staging."

staging-deploy: staging-check
	@echo "🚀 Deploying to staging..."
	git push origin staging

# Clean up
.PHONY: clean-build clean-deps
clean-build:
	rm -rf backend/bin/
	rm -rf frontend/.next/
	rm -rf frontend/out/

clean-deps:
	rm -rf node_modules/
	rm -rf frontend/node_modules/
	cd backend && go clean -modcache

# Help
.PHONY: help
help:
	@echo "CreatorSync Development Commands:"
	@echo ""
	@echo "🚀 LOCAL DEVELOPMENT (Recommended):"
	@echo "  make dev-local        - Start LOCAL development (isolated test data)"
	@echo "  make docker-local     - Start LOCAL PostgreSQL only"
	@echo "  make migrate-local-up - Run migrations on LOCAL database"
	@echo "  make seed-local       - Add test data to LOCAL database"
	@echo "  make inspect-local    - Check what's in LOCAL database"
	@echo ""
	@echo "🔧 General Development:"
	@echo "  make run              - Start backend and frontend (no database)"
	@echo "  make run-backend      - Start only backend"
	@echo "  make run-frontend     - Start only frontend"
	@echo ""
	@echo "⚠️  Staging/Production (Use Carefully):"
	@echo "  make dev-with-docker  - Start with staging/production database"
	@echo "  make docker-run       - Start staging/production PostgreSQL"
	@echo "  make migrate-up       - Run migrations on staging/production"
	@echo ""
	@echo "🧪 Testing:"
	@echo "  make test-fast     - Quick tests (pre-commit)"
	@echo "  make test-full     - Comprehensive tests"
	@echo "  make test-integration - Full integration tests (pre-push)"
	@echo ""
	@echo "🔧 Build & Deploy:"
	@echo "  make build         - Build both projects"
	@echo "  make staging-check - Validate before staging deployment"
	@echo "  make staging-deploy - Deploy to staging after validation"

.PHONY: all build test clean watch docker-run docker-down itest
