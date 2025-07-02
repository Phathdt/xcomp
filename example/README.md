# XComp API Server Example

A complete example application showcasing the XComp dependency injection framework with:

- **Professional CLI Interface** using `github.com/urfave/cli/v2`
- **Comprehensive Makefile** with 40+ targets for development and deployment
- **Goose Database Migrations** with full rollback support
- **Colored Terminal Logging** with automatic detection
- **Clean Architecture** with domain-driven design
- **Docker Compose** setup with PostgreSQL and Redis
- **Structured Logging** with Zap
- **Hot Reload** development environment
- **Production-Ready Configuration** management

## 🚀 Quick Start

### Using Makefile (Recommended)

```bash
# Show all available commands
make help

# Setup complete development environment (Docker + Migrations)
make dev-setup

# Run in development mode (colored console, debug)
make run-dev

# Run in production mode (JSON logs, info level)
make run-prod

# Hot reload development
make dev

# Build the application
make build
```

### Using CLI Directly

```bash
# Show help
go run . --help

# Start server with default config
go run . serve

# Start with custom config and port
go run . serve --config config-dev.yaml --port 8080

# Show version information
go run . version

# Health check
go run . health
```

## 📋 Available Commands

### CLI Commands

| Command | Aliases | Description |
|---------|---------|-------------|
| `serve` | `s`, `start` | Start the HTTP server |
| `version` | `v` | Show version information |
| `health` | | Check application health |

### CLI Flags

| Flag | Short | Environment | Description |
|------|-------|-------------|-------------|
| `--config` | `-c` | `CONFIG_FILE` | Configuration file path |
| `--port` | `-p` | `PORT` | Port to listen on |
| `--verbose` | `-V` | `VERBOSE` | Enable verbose logging |

## 🔨 Makefile Targets

### 🚀 Primary Development Workflow
- `make run-dev` - **Development mode** (colored console, debug level)
- `make run-prod` - **Production mode** (JSON logs, info level)
- `make dev` - **Hot reload** development with Air
- `make dev-setup` - **Complete setup** (dependencies + Docker + migrations)

### 📦 Build Optimization Targets
- `make build` - **Standard build** (41MB) - includes debug symbols
- `make build-optimized` - **Optimized build** (30MB) - stripped symbols, CGO disabled
- `make build-ultra` - **Ultra-optimized** (29MB) - static linking, Linux target
- `make build-upx` - **UPX compressed** (6MB) - Linux binary with compression
- `make build-all` - **All variants** - builds all optimization levels
- `make install-upx` - **Install UPX** compressor tool

### 📊 Size Optimization Results

| Build Type | Size | Reduction | Techniques Used | Best For |
|------------|------|-----------|----------------|----------|
| **Standard** | 41MB | - | Version info only | Development, debugging |
| **Optimized** | 30MB | **27%** | `-s -w`, `CGO_ENABLED=0`, `-trimpath` | Production macOS/general |
| **Ultra** | 29MB | **29%** | Static linking, Linux cross-compile | Docker containers |
| **UPX Linux** | 6MB | **85%** | Binary compression + optimization | Resource-constrained servers |

#### ⚡ **Recommended Usage:**
- **Development**: `make build` (full debug info)
- **Production macOS**: `make build-optimized` (30MB, fast startup)
- **Docker/Cloud**: `make build-upx` (6MB, minimal resources)

#### Optimization Flags Explained:
- `-s` - Strip symbol table and debug info
- `-w` - Strip DWARF debug info
- `CGO_ENABLED=0` - Disable CGO for static binary
- `-trimpath` - Remove absolute paths from binary
- `-extldflags=-static` - Static linking (Linux)
- `upx --lzma --best` - Maximum LZMA compression

### 🗄️ Database & Migrations (Goose)
- `make db-setup` - Setup database with Docker and run migrations
- `make db-reset` - Reset and setup database from scratch
- `make migrate-up` - Apply all pending migrations
- `make migrate-down` - Rollback one migration
- `make migrate-status` - Show migration status
- `make migrate-version` - Show current migration version
- `make migrate-create NAME=migration_name` - Create new migration
- `make migrate-reset` - Reset all migrations (⚠️ destructive)
- `make migrate-redo` - Redo the last migration
- `make migrate-validate` - Validate all migrations

### 🔨 Building & Testing
- `make test` - Run all tests
- `make test-coverage` - Run tests with coverage report
- `make benchmark` - Run performance benchmarks
- `make lint` - Run code linters
- `make format` - Format code with gofmt

### 🐳 Docker Management
- `make docker-up` - Start PostgreSQL + Redis services
- `make docker-down` - Stop all Docker services
- `make docker-restart` - Restart Docker services
- `make docker-clean` - Clean Docker resources and volumes
- `make docker-logs` - Show Docker service logs

### 🧹 Utilities
- `make clean` - Clean build artifacts and cache
- `make deps` - Download and tidy dependencies
- `make update` - Update dependencies to latest versions
- `make version` - Show version information
- `make info` - Show detailed build information
- `make health` - Check application health
- `make deploy-check` - Verify deployment readiness

## 🌈 Configuration Profiles

### config-dev.yaml (Development)
```yaml
app:
  name: 'XComp API Server [DEV] 🚀'
  environment: 'development'

logging:
  level: 'debug'           # Show all messages
  format: 'console'        # Human-readable format
  force_colors: true       # Colored output
  enable_stacktrace: true  # Debug info

database:
  max_connections: 10      # Lower for development
```

### config-prod.yaml (Production)
```yaml
app:
  name: 'XComp API Server'
  environment: 'production'

logging:
  level: 'info'           # Essential messages only
  format: 'json'          # Structured logging
  disable_colors: true    # Clean output

database:
  max_connections: 25     # Higher for production
```

## 🗄️ Database Migrations with Goose

### Migration Commands
```bash
# Check current migration status
make migrate-status

# Create a new migration
make migrate-create NAME=add_users_table

# Apply all pending migrations
make migrate-up

# Rollback one migration
make migrate-down

# Show current version
make migrate-version

# Validate all migration files
make migrate-validate
```

### Migration File Structure
```
migrations/
├── 001_create_products_table.sql
├── 002_create_orders_tables.sql
└── 20250702014452_add_update_trigger_function.sql
```

### Example Migration File
```sql
-- +goose Up
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS users;
```

## 🏗️ Build System & Versioning

The build system automatically injects Git information:

```bash
# Automatic version injection
make build
./build/api-server version
# Output: Version: v1.2.3-abc123, Build Time: 2025-07-02_08:45:00

# Manual version override
VERSION=v2.0.0 make build
```

## 🔄 Hot Reload Development

```bash
# Install development tools
make install-tools

# Start hot reload development
make dev

# Manual Air usage
air
```

The `.air.toml` configuration automatically:
- Restarts on Go file changes
- Uses development config (`config-dev.yaml`)
- Excludes build artifacts and vendor directories

## 🐳 Docker Development Environment

```bash
# Complete development setup
make dev-setup

# Or step by step:
make docker-up        # Start PostgreSQL + Redis
make migrate-up       # Run database migrations
make run-dev          # Start application
```

### Docker Services
- **PostgreSQL 15** on port 5432 with health checks
- **Redis 7** on port 6379 with health checks
- **Automatic service orchestration** with docker-compose

## 📊 Logging Examples

### Development Console Output (Colored)
```
2025-07-02T08:25:24.253+0700    DEBUG   services/product_service.go:28   Getting product    {"product_id": "123e4567-e89b-12d3-a456-426614174000"}
2025-07-02T08:25:24.254+0700    INFO    services/product_service.go:45   Product retrieved successfully    {"product_id": "123e4567-e89b-12d3-a456-426614174000", "cached": false}
```

### Production JSON Output
```json
{"level":"INFO","timestamp":"2025-07-02T08:25:23.244+0700","caller":"services/product_service.go:45","message":"Product retrieved successfully","product_id":"123e4567-e89b-12d3-a456-426614174000","cached":false}
```

## 🔗 API Endpoints

### Health & Info
- `GET /health` - Health check with version info

### Products API
- `GET /api/products` - List products with pagination
- `POST /api/products` - Create new product
- `GET /api/products/{id}` - Get product by ID (with Redis caching)
- `PUT /api/products/{id}` - Update existing product
- `DELETE /api/products/{id}` - Delete product

### Orders API
- `GET /api/orders` - List orders by customer
- `POST /api/orders` - Create new order with items
- `GET /api/orders/{id}` - Get order by ID (with Redis caching)
- `PUT /api/orders/{id}/status` - Update order status

## 🧪 Testing & Quality Assurance

```bash
# Run all tests
make test

# Generate coverage report
make test-coverage

# Run performance benchmarks
make benchmark

# Code quality checks
make lint
make format

# Pre-deployment validation
make deploy-check
```

## 🚀 Production Deployment

### Quick Deploy
```bash
# Build optimized production binary
make build

# Run in production mode
CONFIG_FILE=config-prod.yaml ./build/api-server serve
```

### Advanced Deployment

# Complete release package
make release
```

## 🛠️ Development Workflow

### Daily Development
```bash
# 1. Start development environment
make dev-setup

# 2. Create database migration (if needed)
make migrate-create NAME=add_new_feature

# 3. Start hot reload development
make dev

# 4. Run tests periodically
make test

# 5. Check code quality
make lint format
```

### Feature Development
```bash
# 1. Create feature branch
git checkout -b feature/new-api

# 2. Develop with hot reload
make dev

# 3. Test database changes
make migrate-status
make migrate-up

# 4. Validate before commit
make test lint deploy-check

# 5. Build and test production mode
make build
make run-prod
```

## 📁 Project Architecture

```
example/
├── build/                          # Build artifacts
├── infrastructure/
│   └── database/                   # Database connections
│       ├── interfaces.go           # Database interfaces
│       ├── postgres.go             # PostgreSQL connection
│       └── redis.go                # Redis connection
├── modules/                        # Business domains
│   ├── product/                    # Product module
│   │   ├── domain/                 # Business logic
│   │   │   ├── entities/           # Domain entities
│   │   │   └── interfaces/         # Repository interfaces
│   │   ├── application/            # Application services
│   │   │   ├── dto/                # Data transfer objects
│   │   │   └── services/           # Business services
│   │   ├── infrastructure/         # External concerns
│   │   │   ├── persistence/        # Database repositories
│   │   │   └── http/               # HTTP controllers & routes
│   │   └── product.module.go       # Module definition
│   └── order/                      # Order module (similar structure)
├── migrations/                     # Database migrations (Goose)
│   ├── 001_create_products_table.sql
│   ├── 002_create_orders_tables.sql
│   └── 20250702014452_add_update_trigger_function.sql
├── config-dev.yaml                 # Development configuration
├── config-prod.yaml                # Production configuration
├── docker-compose.yml              # Docker services
├── Makefile                        # Build automation (40+ targets)
├── .air.toml                       # Hot reload configuration
├── sqlc.yaml                       # SQL code generation
└── main.go                         # CLI application entry point
```

## 🎯 XComp Framework Features

### Dependency Injection
- **Automatic injection** via struct tags (`inject:"ServiceName"`)
- **Lazy loading** with singleton pattern
- **Module system** inspired by NestJS
- **Interface-based architecture** for testability

### Configuration Management
- **YAML configuration** with environment variable overrides
- **Multiple environments** (dev/prod configs)
- **Type-safe access** with default values
- **Hot configuration** switching via environment variables

### Structured Logging
- **Zap-based logging** with contextual fields
- **Configurable outputs** (console/JSON)
- **Automatic color detection** for terminals
- **Log level control** per environment

This example demonstrates a **production-ready Go application** with enterprise-grade tooling, clean architecture principles, and comprehensive database migration management using Goose.
