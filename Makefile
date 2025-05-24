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

# Run backend and frontend together
run:
	@echo "Running backend and frontend..."
	cd $(BACKEND_DIR) && go run ./cmd/api &
	npm install --prefer-offline --no-fund --prefix ./frontend
	npm run dev --prefix ./frontend

# Compose up (supports Docker v2 fallback to v1)
docker-run:
	@if docker compose up --build 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose up --build; \
	fi

# Compose down
docker-down:
	@if docker compose down 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose down; \
	fi

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
.PHONY: docker-run docker-down docker-logs docker-clean
docker-run:
	docker compose up -d

docker-down:
	docker compose down

docker-logs:
	docker compose logs -f

docker-clean:
	docker compose down -v --remove-orphans

# Database Migrations
.PHONY: migrate-up migrate-down migrate-reset
migrate-up:
	cd backend && go run cmd/migrate/main.go up

migrate-down:
	cd backend && go run cmd/migrate/main.go down

migrate-reset: migrate-down migrate-up

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
.PHONY: run run-backend run-frontend dev
run: docker-run
	@echo "Starting backend and frontend in parallel..."
	@make run-backend & make run-frontend & wait

run-backend:
	cd backend && go run cmd/api/main.go

run-frontend:
	cd frontend && npm run dev

dev: run

# Build
.PHONY: build build-backend build-frontend
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
.PHONY: clean clean-deps
clean:
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
	@echo "🐳 Docker & Database:"
	@echo "  make docker-run     - Start PostgreSQL container"
	@echo "  make docker-down    - Stop PostgreSQL container"
	@echo "  make migrate-up     - Run database migrations"
	@echo ""
	@echo "🧪 Testing:"
	@echo "  make test-fast      - Quick tests (pre-commit)"
	@echo "  make test-full      - Comprehensive tests"
	@echo "  make test-integration - Full integration tests (pre-push)"
	@echo ""
	@echo "🚀 Development:"
	@echo "  make dev           - Start both backend and frontend"
	@echo "  make build         - Build both projects"
	@echo ""
	@echo "📊 Staging:"
	@echo "  make staging-check - Validate before staging deployment"
	@echo "  make staging-deploy - Deploy to staging after validation"

.PHONY: all build test clean watch docker-run docker-down itest
