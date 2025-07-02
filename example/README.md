# XComp API Server Example

A complete example application showcasing the XComp dependency injection framework with:

- **Professional CLI Interface** using `github.com/urfave/cli/v2`
- **Comprehensive Makefile** for easy development and deployment
- **Colored Terminal Logging** with automatic detection
- **Clean Architecture** with domain-driven design
- **Docker Compose** setup with PostgreSQL and Redis
- **Structured Logging** with Zap
- **Hot Reload** development environment

## 🚀 Quick Start

### Using Makefile (Recommended)

```bash
# Show all available commands
make help

# Setup development environment
make dev-setup

# Run with colored console logging (development)
make run-color

# Run with JSON logging (production)
make run-json

# Build the application
make build

# Run tests
make test
```

### Using CLI Directly

```bash
# Show help
go run . --help

# Start server with default config
go run . serve

# Start with custom config and port
go run . serve --config config-color.yaml --port 8080

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

### Development
- `make run` - Run with default settings
- `make run-color` - Run with colored console logging
- `make run-console` - Run with auto-detected colors
- `make run-json` - Run with JSON logging
- `make dev` - Run with hot reload (requires air)
- `make dev-setup` - Setup development environment

### Building
- `make build` - Build the binary
- `make build-linux` - Build for Linux
- `make build-windows` - Build for Windows
- `make build-mac` - Build for macOS
- `make build-all` - Build for all platforms

### Docker
- `make docker-up` - Start all services
- `make docker-down` - Stop all services
- `make docker-restart` - Restart services
- `make docker-clean` - Clean Docker resources

### Testing & Quality
- `make test` - Run tests
- `make test-coverage` - Run tests with coverage
- `make lint` - Run linters
- `make format` - Format code

### Utilities
- `make clean` - Clean build artifacts
- `make deps` - Download dependencies
- `make version` - Show version information
- `make info` - Show build information

## 🌈 Logging Configurations

### config.yaml (Production - JSON)
```yaml
logging:
  level: 'info'
  format: 'json'
  development: false
```

### config-console.yaml (Development - Auto Colors)
```yaml
logging:
  level: 'debug'
  format: 'console'
  development: true
  level_format: 'capital'  # Auto-detects terminal colors
```

### config-color.yaml (Demo - Forced Colors)
```yaml
logging:
  level: 'debug'
  format: 'console'
  development: true
  force_colors: true  # Forces colors even without terminal detection
```

## 🏗️ Build Information

The build system automatically injects version information:

```bash
# Using make
make build    # Injects git info automatically

# Using go directly
go build -ldflags "-X main.Version=v1.0.0 -X main.BuildTime=2024-01-01 -X main.GitCommit=abc123"
```

## 🔄 Hot Reload Development

Install Air for hot reloading:

```bash
make install-tools
make dev
```

Or manually:
```bash
go install github.com/air-verse/air@latest
air
```

## 🐳 Docker Development

```bash
# Start all services (PostgreSQL + Redis)
make docker-up

# Run application
make run-color

# Stop services
make docker-down
```

## 📊 Example Usage

### JSON Logging Output
```json
{"level":"INFO","timestamp":"2025-07-02T08:25:23.244+0700","caller":"xcomp/logger.go:238","message":"Starting API Server","version":"1.0.0","environment":"development"}
```

### Colored Console Output
```
2025-07-02T08:25:24.253+0700    INFO    logger.go:238 Starting API Server    {"version": "1.0.0", "environment": "development", "name": "Colorful API Server 🌈"}
```

## 🔗 API Endpoints

- `GET /health` - Health check
- `GET /api/products` - List products
- `POST /api/products` - Create product
- `GET /api/products/{id}` - Get product by ID
- `PUT /api/products/{id}` - Update product
- `DELETE /api/products/{id}` - Delete product

## 🚀 Deployment

```bash
# Check deployment readiness
make deploy-check

# Build release binaries
make release

# Production build
make build
./build/api-server serve --config config.yaml
```

## 🛠️ Development Workflow

```bash
# 1. Setup environment
make dev-setup

# 2. Start development with hot reload
make dev

# 3. Run tests
make test

# 4. Format code
make format

# 5. Build and test
make build
make run-prod
```

## 📁 Project Structure

```
example/
├── build/              # Build artifacts
├── infrastructure/     # Database, Redis connections
├── modules/           # Business modules (product, order)
│   ├── product/       # Product module
│   └── order/         # Order module
├── migrations/        # Database migrations
├── config*.yaml       # Configuration files
├── docker-compose.yml # Docker services
├── Makefile          # Build automation
├── .air.toml         # Hot reload config
└── main.go           # CLI application entry point
```

This example demonstrates a production-ready Go application with professional development tooling and clean architecture principles.
