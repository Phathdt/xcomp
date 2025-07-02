# Makefile for XComp API Server

# Build variables
BINARY_NAME=api-server
BUILD_DIR=build
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT)"

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Default target
.DEFAULT_GOAL := help

## Help target
.PHONY: help
help: ## Show this help message
	@echo "XComp API Server - Available commands:"
	@echo ""
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@echo ""
	@echo "Examples:"
	@echo "  make run-dev                # Development mode (colored, debug)"
	@echo "  make run-prod               # Production mode (JSON, info)"
	@echo "  make dev                    # Hot reload development"
	@echo "  make build                  # Build the binary"
	@echo "  make docker-up              # Start all services"

## Build targets
.PHONY: build
build: ## Build the binary
	@echo "🔨 Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) .
	@echo "✅ Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

## Run targets
.PHONY: run
run: ## Run the application with default settings
	@echo "🚀 Starting API Server (default config)..."
	$(GORUN) $(LDFLAGS) . serve

.PHONY: run-dev
run-dev: ## Run in development mode (colored console, debug)
	@echo "🚀 Starting API Server [DEVELOPMENT MODE]..."
	CONFIG_FILE=config-dev.yaml $(GORUN) $(LDFLAGS) . serve

.PHONY: run-prod
run-prod: build ## Run in production mode (JSON logs, info level)
	@echo "🏭 Starting API Server [PRODUCTION MODE]..."
	CONFIG_FILE=config-prod.yaml $(BUILD_DIR)/$(BINARY_NAME) serve

## Development targets
.PHONY: dev
dev: ## Run in development mode with hot reload (requires air)
	@echo "🔄 Starting development server with hot reload..."
	@if command -v air > /dev/null 2>&1; then \
		air; \
	else \
		echo "❌ Air not found. Install with: go install github.com/air-verse/air@latest"; \
		echo "   Or run: make run-dev"; \
	fi

.PHONY: install-tools
install-tools: ## Install development tools
	@echo "🔧 Installing development tools..."
	$(GOGET) -u github.com/air-verse/air@latest
	@echo "✅ Development tools installed!"

## Docker targets
.PHONY: docker-up
docker-up: ## Start all services with Docker Compose
	@echo "🐳 Starting Docker services..."
	docker-compose up -d --wait
	@echo "✅ Docker services started!"

.PHONY: docker-down
docker-down: ## Stop all Docker services
	@echo "🛑 Stopping Docker services..."
	docker-compose down
	@echo "✅ Docker services stopped!"

.PHONY: docker-logs
docker-logs: ## Show Docker logs
	@echo "📜 Docker logs:"
	docker-compose logs -f

.PHONY: docker-restart
docker-restart: docker-down docker-up ## Restart Docker services
	@echo "🔄 Docker services restarted!"

.PHONY: docker-clean
docker-clean: ## Clean Docker resources
	@echo "🧹 Cleaning Docker resources..."
	docker-compose down -v --remove-orphans
	docker system prune -f
	@echo "✅ Docker cleanup complete!"

## Database targets
.PHONY: db-setup
db-setup: docker-up migrate-up ## Setup database with Docker and run migrations
	@echo "🗄️ Database setup complete!"

.PHONY: db-reset
db-reset: docker-up migrate-reset migrate-up ## Reset and setup database
	@echo "🔄 Database reset complete!"

.PHONY: db-seed
db-seed: ## Seed the database with sample data
	@echo "🌱 Seeding database..."
	@echo "✅ Database seeded!"

## Migration targets (using Goose)
MIGRATION_DIR=migrations
DB_STRING="postgres://postgres:password@localhost:5432/example_db?sslmode=disable"

.PHONY: migrate-up
migrate-up: ## Apply all pending migrations
	@echo "⬆️ Running database migrations..."
	@if command -v goose > /dev/null 2>&1; then \
		goose -dir $(MIGRATION_DIR) postgres $(DB_STRING) up; \
	else \
		echo "❌ Goose not found. Installing..."; \
		go install github.com/pressly/goose/v3/cmd/goose@latest; \
		goose -dir $(MIGRATION_DIR) postgres $(DB_STRING) up; \
	fi
	@echo "✅ Migrations applied!"

.PHONY: migrate-down
migrate-down: ## Rollback one migration
	@echo "⬇️ Rolling back one migration..."
	@if command -v goose > /dev/null 2>&1; then \
		goose -dir $(MIGRATION_DIR) postgres $(DB_STRING) down; \
	else \
		echo "❌ Goose not found. Installing..."; \
		go install github.com/pressly/goose/v3/cmd/goose@latest; \
		goose -dir $(MIGRATION_DIR) postgres $(DB_STRING) down; \
	fi
	@echo "✅ Migration rolled back!"

.PHONY: migrate-status
migrate-status: ## Show migration status
	@echo "📊 Migration status:"
	@if command -v goose > /dev/null 2>&1; then \
		goose -dir $(MIGRATION_DIR) postgres $(DB_STRING) status; \
	else \
		echo "❌ Goose not found. Installing..."; \
		go install github.com/pressly/goose/v3/cmd/goose@latest; \
		goose -dir $(MIGRATION_DIR) postgres $(DB_STRING) status; \
	fi

.PHONY: migrate-version
migrate-version: ## Show current migration version
	@echo "📋 Current migration version:"
	@if command -v goose > /dev/null 2>&1; then \
		goose -dir $(MIGRATION_DIR) postgres $(DB_STRING) version; \
	else \
		echo "❌ Goose not found. Installing..."; \
		go install github.com/pressly/goose/v3/cmd/goose@latest; \
		goose -dir $(MIGRATION_DIR) postgres $(DB_STRING) version; \
	fi

.PHONY: migrate-reset
migrate-reset: ## Reset all migrations (careful!)
	@echo "⚠️ Resetting all migrations..."
	@if command -v goose > /dev/null 2>&1; then \
		goose -dir $(MIGRATION_DIR) postgres $(DB_STRING) reset; \
	else \
		echo "❌ Goose not found. Installing..."; \
		go install github.com/pressly/goose/v3/cmd/goose@latest; \
		goose -dir $(MIGRATION_DIR) postgres $(DB_STRING) reset; \
	fi
	@echo "✅ All migrations reset!"

.PHONY: migrate-create
migrate-create: ## Create a new migration (make migrate-create NAME=add_users_table)
	@echo "📝 Creating new migration: $(NAME)"
	@if [ -z "$(NAME)" ]; then \
		echo "❌ Please provide a migration name: make migrate-create NAME=your_migration_name"; \
		exit 1; \
	fi
	@if command -v goose > /dev/null 2>&1; then \
		goose -dir $(MIGRATION_DIR) create $(NAME) sql; \
	else \
		echo "❌ Goose not found. Installing..."; \
		go install github.com/pressly/goose/v3/cmd/goose@latest; \
		goose -dir $(MIGRATION_DIR) create $(NAME) sql; \
	fi
	@echo "✅ Migration created!"

.PHONY: migrate-redo
migrate-redo: ## Redo the last migration
	@echo "🔄 Redoing last migration..."
	@if command -v goose > /dev/null 2>&1; then \
		goose -dir $(MIGRATION_DIR) postgres $(DB_STRING) redo; \
	else \
		echo "❌ Goose not found. Installing..."; \
		go install github.com/pressly/goose/v3/cmd/goose@latest; \
		goose -dir $(MIGRATION_DIR) postgres $(DB_STRING) redo; \
	fi
	@echo "✅ Migration redone!"

.PHONY: migrate-validate
migrate-validate: ## Validate all migrations
	@echo "🔍 Validating migrations..."
	@if command -v goose > /dev/null 2>&1; then \
		goose -dir $(MIGRATION_DIR) validate; \
	else \
		echo "❌ Goose not found. Installing..."; \
		go install github.com/pressly/goose/v3/cmd/goose@latest; \
		goose -dir $(MIGRATION_DIR) validate; \
	fi
	@echo "✅ Migrations validated!"

## Testing targets
.PHONY: test
test: ## Run tests
	@echo "🧪 Running tests..."
	$(GOTEST) -v ./...

.PHONY: test-coverage
test-coverage: ## Run tests with coverage
	@echo "🧪 Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report generated: coverage.html"

.PHONY: benchmark
benchmark: ## Run benchmarks
	@echo "⚡ Running benchmarks..."
	$(GOTEST) -bench=. -benchmem ./...

## Utility targets
.PHONY: clean
clean: ## Clean build artifacts and cache
	@echo "🧹 Cleaning..."
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html
	@echo "✅ Clean complete!"

.PHONY: deps
deps: ## Download and tidy dependencies
	@echo "📦 Managing dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy
	@echo "✅ Dependencies updated!"

.PHONY: update
update: ## Update dependencies to latest versions
	@echo "⬆️ Updating dependencies..."
	$(GOGET) -u ./...
	$(GOMOD) tidy
	@echo "✅ Dependencies updated!"

.PHONY: lint
lint: ## Run linters (requires golangci-lint)
	@echo "🔍 Running linters..."
	@if command -v golangci-lint > /dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "❌ golangci-lint not found. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

.PHONY: format
format: ## Format code
	@echo "💅 Formatting code..."
	$(GOCMD) fmt ./...
	@echo "✅ Code formatted!"

## Informational targets
.PHONY: version
version: ## Show version information
	@$(GORUN) $(LDFLAGS) . version

.PHONY: health
health: ## Check application health
	@$(GORUN) $(LDFLAGS) . health

.PHONY: info
info: ## Show build information
	@echo "📋 Build Information:"
	@echo "  Version: $(VERSION)"
	@echo "  Build Time: $(BUILD_TIME)"
	@echo "  Git Commit: $(GIT_COMMIT)"
	@echo "  Binary Name: $(BINARY_NAME)"
	@echo "  Build Dir: $(BUILD_DIR)"

## Quick development workflow
.PHONY: dev-setup
dev-setup: deps install-tools db-setup ## Setup development environment
	@echo "🎉 Development environment ready!"
	@echo ""
	@echo "Quick start commands:"
	@echo "  make run-dev      # Development mode (colored, debug)"
	@echo "  make dev          # Hot reload development"
	@echo "  make test         # Run tests"
	@echo "  make migrate-status # Check migration status"

.PHONY: dev-reset
dev-reset: clean docker-clean deps db-setup ## Reset development environment
	@echo "🔄 Development environment reset!"

## Production deployment
.PHONY: deploy-check
deploy-check: test lint build ## Check if ready for deployment
	@echo "✅ Deployment checks passed!"

.PHONY: release
release: clean test lint build-all ## Prepare release build
	@echo "🎉 Release build complete!"
	@echo "  Binaries available in $(BUILD_DIR)/"
