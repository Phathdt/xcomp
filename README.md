# XComp - NestJS-Inspired Dependency Injection Framework for Go

[![Go Version](https://img.shields.io/badge/Go-1.24%2B-00ADD8?style=flat&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
![Build Status](https://img.shields.io/badge/Build-Passing-brightgreen?style=flat)

XComp is a powerful, NestJS-inspired dependency injection framework for Go that brings enterprise-grade architecture patterns to Go applications. Build scalable, maintainable applications with clean separation of concerns, automatic dependency injection, and modular design.

## ✨ Key Features

- 🏗️ **Dependency Injection Container** - Automatic service resolution with lazy loading
- 🧩 **Modular Architecture** - NestJS-style modules with providers and imports
- ⚙️ **Configuration Management** - YAML-based config with environment variable overrides
- 📝 **Structured Logging** - Zap-based logging with contextual fields and colored output
- 🏷️ **Tag-Based Injection** - Simple `inject:"ServiceName"` struct tags
- 🔄 **Lazy Loading** - Services instantiated only when needed
- 🎯 **Interface-Driven** - Clean architecture with proper abstraction layers
- 🛡️ **Type Safety** - Compile-time dependency verification

## 🚀 Quick Start

### Installation

```bash
go get xcomp
# or if using as a dependency in your go.mod:
# require xcomp v1.0.0
```

### Basic Usage

```go
package main

import (
    "fmt"
    "xcomp"
)

// Define a service
type UserService struct {
    Config *xcomp.ConfigService `inject:"ConfigService"`
    Logger xcomp.Logger         `inject:"Logger"`
}

func (us *UserService) GetServiceName() string {
    return "UserService"
}

func (us *UserService) GetUser(id string) *User {
    us.Logger.Info("Getting user", xcomp.Field("user_id", id))
    // Business logic here...
    return &User{ID: id, Name: "John Doe"}
}

// Create and configure container
func main() {
    container := xcomp.NewContainer()

    // Register services
    container.Register("ConfigService", xcomp.NewConfigService("config.yaml"))
    container.Register("Logger", xcomp.NewDevelopmentLogger())
    container.Register("UserService", &UserService{})

    // Auto-inject dependencies
    userService := container.Get("UserService").(*UserService)
    container.Inject(userService)

    // Use the service
    user := userService.GetUser("123")
    fmt.Printf("User: %+v\n", user)
}
```

## 🧩 Module System

Create modular, reusable components inspired by NestJS:

```go
// Define a module
func CreateUserModule() xcomp.Module {
    return xcomp.NewModule().
        AddFactory("UserRepository", func(c *xcomp.Container) any {
            repo := &UserRepositoryImpl{}
            c.Inject(repo)
            return repo
        }).
        AddFactory("UserService", func(c *xcomp.Container) any {
            service := &UserService{}
            c.Inject(service)
            return service
        }).
        AddFactory("UserController", func(c *xcomp.Container) any {
            controller := &UserController{}
            c.Inject(controller)
            return controller
        }).
        Build()
}

// Register module
container := xcomp.NewContainer()
userModule := CreateUserModule()
container.RegisterModule(userModule)
```

## ⚙️ Configuration Management

XComp provides a powerful configuration system with YAML files and environment variable overrides:

```go
// config.yaml
app:
  name: "My Application"
  port: 3000
  debug: true

database:
  host: "localhost"
  port: 5432
  name: "mydb"

logging:
  level: "debug"
  format: "console"
```

```go
// Usage in Go
type DatabaseService struct {
    Config *xcomp.ConfigService `inject:"ConfigService"`
}

func (ds *DatabaseService) Connect() error {
    host := ds.Config.GetString("database.host", "localhost")
    port := ds.Config.GetInt("database.port", 5432)
    dbName := ds.Config.GetString("database.name", "defaultdb")

    // Connection logic...
    return nil
}
```

### Environment Variable Overrides

Environment variables automatically override config file values:

```bash
export DATABASE_HOST=production-db.com
export APP_PORT=8080
export LOGGING_LEVEL=info
```

## 📝 Structured Logging

XComp includes a powerful Zap-based logging system with contextual fields:

```go
type OrderService struct {
    Logger xcomp.Logger `inject:"Logger"`
}

func (os *OrderService) CreateOrder(order *Order) error {
    os.Logger.Info("Creating order",
        xcomp.Field("order_id", order.ID),
        xcomp.Field("customer_id", order.CustomerID),
        xcomp.Field("total", order.Total))

    if err := os.validateOrder(order); err != nil {
        os.Logger.Error("Order validation failed",
            xcomp.Field("order_id", order.ID),
            xcomp.Field("error", err))
        return err
    }

    os.Logger.Info("Order created successfully",
        xcomp.Field("order_id", order.ID))
    return nil
}
```

### Configurable Logging Output

```yaml
# Development - Colored console output
logging:
  level: "debug"
  format: "console"
  development: true
  force_colors: true

# Production - Structured JSON output
logging:
  level: "info"
  format: "json"
  development: false
  disable_colors: true
```

## 🏗️ Clean Architecture Example

XComp promotes clean architecture patterns with proper separation of concerns:

```go
// Domain Layer - Business entities and interfaces
type User struct {
    ID    string
    Email string
    Name  string
}

type UserRepository interface {
    GetByID(id string) (*User, error)
    Create(user *User) error
}

// Application Layer - Business logic
type UserService struct {
    UserRepo UserRepository `inject:"UserRepository"`
    Logger   xcomp.Logger   `inject:"Logger"`
}

func (us *UserService) CreateUser(email, name string) (*User, error) {
    us.Logger.Info("Creating user", xcomp.Field("email", email))

    user := &User{
        ID:    generateID(),
        Email: email,
        Name:  name,
    }

    return user, us.UserRepo.Create(user)
}

// Infrastructure Layer - External concerns
type PostgresUserRepository struct {
    DB     *sql.DB         `inject:"DatabaseConnection"`
    Logger xcomp.Logger    `inject:"Logger"`
}

func (pur *PostgresUserRepository) GetByID(id string) (*User, error) {
    pur.Logger.Debug("Fetching user from database", xcomp.Field("user_id", id))
    // Database query logic...
    return user, nil
}

// HTTP Layer - Controllers
type UserController struct {
    UserService *UserService `inject:"UserService"`
    Logger      xcomp.Logger `inject:"Logger"`
}

func (uc *UserController) CreateUser(c *fiber.Ctx) error {
    var req CreateUserRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
    }

    user, err := uc.UserService.CreateUser(req.Email, req.Name)
    if err != nil {
        uc.Logger.Error("Failed to create user", xcomp.Field("error", err))
        return c.Status(500).JSON(fiber.Map{"error": "Failed to create user"})
    }

    return c.JSON(user)
}
```

## 🔧 Advanced Features

### Lazy Loading & Singletons

Services are created only when first accessed, improving startup performance:

```go
container.RegisterSingleton("ExpensiveService", func(c *xcomp.Container) any {
    // This factory function runs only when service is first requested
    return &ExpensiveService{
        Config: c.Get("ConfigService").(*xcomp.ConfigService),
    }
})
```

### Module Imports

Modules can import other modules, creating a dependency graph:

```go
func CreateAppModule() xcomp.Module {
    return xcomp.NewModule().
        Import(CreateDatabaseModule()).
        Import(CreateUserModule()).
        Import(CreateOrderModule()).
        Build()
}
```

### Service Discovery

List all registered services for debugging:

```go
services := container.ListServices()
fmt.Printf("Registered services: %v\n", services)
```

## 📚 Complete Example Application

See the [`example/`](./example/) directory for a complete application showcasing:

- **REST API** with Fiber framework
- **PostgreSQL** integration with migrations
- **Redis** caching layer
- **Docker Compose** development environment
- **Hot reload** development setup
- **Production-ready** configuration
- **Comprehensive Makefile** with 40+ targets

```bash
cd example
make dev-setup  # Complete setup with Docker + migrations
make run-dev    # Start development server
```

## 🎯 Use Cases

XComp is perfect for:

- 🌐 **Web APIs** - REST, GraphQL, gRPC services
- 🏢 **Enterprise Applications** - Complex business logic with many dependencies
- 🔌 **Microservices** - Modular, testable service architecture
- 🧪 **Testable Code** - Easy mocking and dependency injection for tests
- 📊 **Data Processing** - ETL pipelines with configurable components
- 🛠️ **CLI Tools** - Command-line applications with modular design

## 🧪 Testing

XComp makes testing easy with dependency injection:

```go
func TestUserService(t *testing.T) {
    // Create test container
    container := xcomp.NewContainer()

    // Mock dependencies
    mockRepo := &MockUserRepository{}
    mockLogger := &MockLogger{}

    container.Register("UserRepository", mockRepo)
    container.Register("Logger", mockLogger)

    // Create service with mocked dependencies
    userService := &UserService{}
    container.Inject(userService)

    // Test the service
    user, err := userService.CreateUser("test@example.com", "Test User")
    assert.NoError(t, err)
    assert.Equal(t, "test@example.com", user.Email)
}
```

## 📊 Performance

XComp is designed for performance:

- **Lazy Loading** - Services created only when needed
- **Singleton Pattern** - Services reused across requests
- **Reflection Optimization** - Minimal reflection usage
- **Memory Efficient** - Low overhead dependency injection
- **Concurrent Safe** - Thread-safe service resolution

## 🛠️ API Reference

### Container

```go
// Create new container
container := xcomp.NewContainer()

// Register service instance
container.Register(name string, service any)

// Register factory function (lazy loading)
container.RegisterSingleton(name string, factory func(*Container) any)

// Get service by name
service := container.Get(name string) any

// Get service with type assertion
var userService UserService
if container.GetTyped("UserService", &userService) {
    // userService is now populated
}

// Inject dependencies into struct
container.Inject(target any) error

// Register module
container.RegisterModule(module Module) error

// List all services
services := container.ListServices() []string
```

### Configuration

```go
// Create config service
config := xcomp.NewConfigService("config.yaml")

// Get values with defaults
config.GetString("key", "default")
config.GetInt("key", 0)
config.GetBool("key", false)
config.Get("key") any
```

### Logging

```go
// Create logger
logger := xcomp.NewLogger(configService)
logger := xcomp.NewDevelopmentLogger() // Quick development logger

// Log with contextual fields
logger.Info("Message", xcomp.Field("key", "value"))
logger.Error("Error occurred", xcomp.Field("error", err))

// Create contextual logger
contextLogger := logger.With(xcomp.Field("request_id", "123"))
```

### Modules

```go
// Create module
module := xcomp.NewModule().
    AddFactory("ServiceName", factoryFunc).
    AddService("ServiceName", serviceInstance).
    Import(otherModule).
    Build()
```

## 🤝 Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Inspired by [NestJS](https://nestjs.com/) dependency injection system
- Built with [Zap](https://github.com/uber-go/zap) for high-performance logging
- Example application uses [Fiber](https://github.com/gofiber/fiber) web framework

## 🔗 Links

- [Example Application](./example/) - Complete working example with REST API
- [API Documentation](#-api-reference) - Full API reference guide
- [License](LICENSE) - MIT License details
- [Contributing](#-contributing) - How to contribute to the project
