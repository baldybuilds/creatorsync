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

.PHONY: all build run test clean watch docker-run docker-down itest