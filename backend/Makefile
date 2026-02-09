# Makefile for ecommerce backend

.PHONY: help build run test clean migrate dev

# Default target
help:
	@echo "Available commands:"
	@echo "  make build    - Build the application"
	@echo "  make run      - Run the application"
	@echo "  make test     - Run tests"
	@echo "  make clean    - Clean build artifacts"
	@echo "  make migrate  - Run database migrations"
	@echo "  make dev      - Run in development mode with hot reload"
	@echo "  make docker   - Build Docker image"
	@echo "  make api-test - Run API tests"

# Build the application
build:
	@echo "Building application..."
	go build -o bin/ecommerce-backend cmd/main.go
	@echo "Build complete: bin/ecommerce-backend"

# Run the application
run: build
	@echo "Starting application..."
	./bin/ecommerce-backend

# Run tests
test:
	@echo "Running tests..."
	go test ./... -v

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	go clean

# Run database migrations
migrate:
	@echo "Running database migrations..."
	@if [ -f "migrations/001_initial_schema.sql" ]; then \
		echo "Found migration file"; \
		read -p "Enter database name [ecommerce]: " dbname; \
		dbname=$${dbname:-ecommerce}; \
		read -p "Enter PostgreSQL username [postgres]: " username; \
		username=$${username:-postgres}; \
		psql -U $$username -d $$dbname -f migrations/001_initial_schema.sql; \
	else \
		echo "No migration files found"; \
	fi

# Development mode with hot reload
dev:
	@echo "Starting development server with hot reload..."
	@if command -v air > /dev/null; then \
		air; \
	elif command -v nodemon > /dev/null; then \
		nodemon --exec "go run cmd/main.go" --ext go; \
	else \
		echo "Installing air for hot reload..."; \
		go install github.com/cosmtrek/air@latest; \
		air; \
	fi

# Build Docker image
docker:
	@echo "Building Docker image..."
	docker build -t ecommerce-backend:latest .

# Run API tests
api-test:
	@echo "Running API tests..."
	@if [ -f "test_api.sh" ]; then \
		chmod +x test_api.sh; \
		./test_api.sh; \
	else \
		echo "API test script not found"; \
	fi

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod download
	go install github.com/cosmtrek/air@latest

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Lint code
lint:
	@echo "Linting code..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "Installing golangci-lint..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
		golangci-lint run; \
	fi

# Generate swagger docs (if using swagger)
swagger:
	@echo "Generating swagger documentation..."
	@if command -v swag > /dev/null; then \
		swag init -g cmd/main.go; \
	else \
		echo "Installing swag..."; \
		go install github.com/swaggo/swag/cmd/swag@latest; \
		swag init -g cmd/main.go; \
	fi

# Create database (for development)
create-db:
	@echo "Creating database..."
	@read -p "Enter database name [ecommerce]: " dbname; \
	dbname=$${dbname:-ecommerce}; \
	read -p "Enter PostgreSQL username [postgres]: " username; \
	username=$${username:-postgres}; \
	createdb -U $$username $$dbname; \
	echo "Database '$$dbname' created successfully"

# Drop database (for development)
drop-db:
	@echo "Dropping database..."
	@read -p "Enter database name [ecommerce]: " dbname; \
	dbname=$${dbname:-ecommerce}; \
	read -p "Enter PostgreSQL username [postgres]: " username; \
	username=$${username:-postgres}; \
	dropdb -U $$username $$dbname; \
	echo "Database '$$dbname' dropped successfully"

# Reset database (drop and recreate)
reset-db: drop-db create-db migrate
	@echo "Database reset complete"

# Check code coverage
coverage:
	@echo "Checking code coverage..."
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Vendor dependencies
vendor:
	@echo "Vendoring dependencies..."
	go mod vendor

# Update dependencies
update-deps:
	@echo "Updating dependencies..."
	go get -u ./...
	go mod tidy

# Show dependency graph
deps-graph:
	@echo "Generating dependency graph..."
	go mod graph | dot -Tpng -o deps.png
	@echo "Dependency graph generated: deps.png"

# Default target
.DEFAULT_GOAL := help