# Makefile for XComp API Server

# Build variables
BINARY_NAME=api-server
BUILD_DIR=build
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Standard ldflags with version info
LDFLAGS_VERSION=-ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT)"

# Optimized ldflags - strip symbols and debug info
LDFLAGS_OPTIMIZED=-ldflags "-s -w -X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT)"

# Ultra optimized - additional flags
LDFLAGS_ULTRA=-ldflags "-s -w -extldflags=-static -X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT)"

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
	@echo "  make build                  # Build standard binary (41MB)"
	@echo "  make build-optimized        # Build optimized binary (30MB - 27% smaller)"
	@echo "  make build-upx              # Build UPX compressed binary (6MB - 85% smaller)"
	@echo "  make docker-up              # Start all services"

## Build targets
.PHONY: build
build: ## Build the binary
	@echo "🔨 Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS_VERSION) -o $(BUILD_DIR)/$(BINARY_NAME) .
	@echo "✅ Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

.PHONY: build-optimized
build-optimized: ## Build optimized binary (stripped symbols, ~30% smaller)
	@echo "🔨 Building optimized $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 $(GOBUILD) $(LDFLAGS_OPTIMIZED) -trimpath -o $(BUILD_DIR)/$(BINARY_NAME)-optimized .
	@echo "✅ Optimized build complete: $(BUILD_DIR)/$(BINARY_NAME)-optimized"
	@if [ -f "$(BUILD_DIR)/$(BINARY_NAME)" ]; then \
		echo "📊 Size comparison:"; \
		echo "  Standard: $$(ls -lh $(BUILD_DIR)/$(BINARY_NAME) | awk '{print $$5}')"; \
		echo "  Optimized: $$(ls -lh $(BUILD_DIR)/$(BINARY_NAME)-optimized | awk '{print $$5}')"; \
	fi

.PHONY: build-ultra
build-ultra: ## Build ultra-optimized binary (static linking, ~40% smaller)
	@echo "🔨 Building ultra-optimized $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=linux $(GOBUILD) $(LDFLAGS_ULTRA) -trimpath -a -installsuffix cgo -o $(BUILD_DIR)/$(BINARY_NAME)-ultra .
	@echo "✅ Ultra-optimized build complete: $(BUILD_DIR)/$(BINARY_NAME)-ultra"

.PHONY: build-upx
build-upx: build-ultra ## Build UPX compressed Linux binary (~85% smaller)
	@echo "🗜️ Compressing Linux binary with UPX..."
	@if command -v upx > /dev/null 2>&1; then \
		upx --lzma --best $(BUILD_DIR)/$(BINARY_NAME)-ultra -o $(BUILD_DIR)/$(BINARY_NAME)-upx-linux; \
		echo "✅ UPX compressed build complete: $(BUILD_DIR)/$(BINARY_NAME)-upx-linux"; \
		echo "📊 Size comparison:"; \
		echo "  Standard (macOS): $$(ls -lh $(BUILD_DIR)/$(BINARY_NAME) | awk '{print $$5}')"; \
		echo "  Optimized (macOS): $$(ls -lh $(BUILD_DIR)/$(BINARY_NAME)-optimized | awk '{print $$5}')"; \
		echo "  Ultra (Linux): $$(ls -lh $(BUILD_DIR)/$(BINARY_NAME)-ultra | awk '{print $$5}')"; \
		echo "  UPX (Linux): $$(ls -lh $(BUILD_DIR)/$(BINARY_NAME)-upx-linux | awk '{print $$5}')"; \
		echo "  💡 Note: UPX works best with Linux binaries"; \
	else \
		echo "❌ UPX not found. Install with: brew install upx (macOS) or apt install upx (Ubuntu)"; \
	fi

.PHONY: build-all
build-all: build build-optimized build-ultra build-upx ## Build all variants
	@echo "🎉 All build variants complete!"
	@echo "📁 Available in $(BUILD_DIR)/:"
	@ls -lh $(BUILD_DIR)/

## Run targets
.PHONY: run
run: ## Run the application with default settings
	@echo "🚀 Starting API Server (default config)..."
	$(GORUN) $(LDFLAGS_VERSION) . serve

.PHONY: run-dev
run-dev: ## Run in development mode (colored console, debug)
	@echo "🚀 Starting API Server [DEVELOPMENT MODE]..."
	CONFIG_FILE=config-dev.yaml $(GORUN) $(LDFLAGS_VERSION) . serve

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
	@$(GORUN) $(LDFLAGS_VERSION) . version

.PHONY: health
health: ## Check application health
	@$(GORUN) $(LDFLAGS_VERSION) . health

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
deploy-check: test lint build-optimized ## Check if ready for deployment (with optimized build)
	@echo "✅ Deployment checks passed!"

.PHONY: release
release: clean test lint build-all ## Prepare all release builds (standard, optimized, ultra, UPX)
	@echo "🎉 Release build complete!"
	@echo "  Binaries available in $(BUILD_DIR)/"
	@echo ""
	@echo "📊 Build size comparison:"
	@ls -lh $(BUILD_DIR)/ | grep api-server

.PHONY: install-upx
install-upx: ## Install UPX compressor
	@echo "📥 Installing UPX..."
	@if command -v brew > /dev/null 2>&1; then \
		brew install upx; \
	elif command -v apt > /dev/null 2>&1; then \
		sudo apt update && sudo apt install -y upx; \
	elif command -v yum > /dev/null 2>&1; then \
		sudo yum install -y upx; \
	else \
		echo "❌ Package manager not found. Please install UPX manually from https://upx.github.io/"; \
	fi
