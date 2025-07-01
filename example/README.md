# Product API - Clean Architecture with XComp DI

This is a complete example of a Product CRUD API built with clean architecture principles and using the XComp dependency injection system inspired by NestJS.

## 🏗️ Architecture

The project follows clean architecture with NestJS/Nx-like modular structure:

```
example/
├── migrations/                          # Database migrations
│   └── 001_create_products_table.sql
├── modules/                             # Feature modules
│   └── product/                         # Product module
│       ├── domain/                      # Core business logic
│       │   ├── entities/               # Business entities
│       │   └── repositories/           # Repository interfaces
│       ├── application/                 # Application services
│       │   ├── dto/                    # Data Transfer Objects
│       │   └── services/               # Application services
│       ├── infrastructure/             # Infrastructure layer
│       │   ├── persistence/            # Data persistence
│       │   │   ├── queries/            # SQLC queries & generated code
│       │   │   │   ├── products.sql    # SQL queries
│       │   │   │   ├── *.go            # Generated SQLC code
│       │   │   └── product_repository_impl.go
│       │   └── http/                   # HTTP layer
│       │       ├── controllers/        # HTTP controllers
│       │       └── routes/             # Route definitions
│       └── product.module.go           # Module definition
├── infrastructure/                     # Shared infrastructure
│   └── database/
│       └── connection.go               # Database connection
├── config.yaml                        # Configuration file
├── docker-compose.yml                 # PostgreSQL setup
├── sqlc.yaml                          # SQLC configuration
└── main.go                            # Application entry point
```

## 🚀 Features

- **Clean Architecture**: Separation of concerns with clear layer boundaries
- **Modular Structure**: NestJS/Nx-like modules for scalability
- **XComp Dependency Injection**: NestJS-inspired DI system with struct tags
- **PostgreSQL Integration**: Database setup with migrations and SQLC
- **Type-Safe Queries**: SQLC generates type-safe Go code from SQL
- **Fiber Framework**: Fast HTTP framework for Go
- **Docker Support**: PostgreSQL database in Docker
- **Configuration Management**: YAML config with environment variable override
- **Graceful Shutdown**: Proper resource cleanup
- **Error Handling**: Comprehensive error handling and validation
- **CORS Support**: Cross-Origin Resource Sharing
- **Health Check**: Health check endpoint

## 📋 Prerequisites

- Go 1.24+
- Docker & Docker Compose
- SQLC (for generating queries)

## 🛠️ Installation & Setup

1. **Navigate to the example directory:**
   ```bash
   cd example
   ```

2. **Start PostgreSQL database:**
   ```bash
   docker-compose up -d postgres
   ```

3. **Install dependencies:**
   ```bash
   go mod tidy
   ```

4. **Install SQLC (if not installed):**
   ```bash
   go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
   ```

5. **Generate SQLC code:**
   ```bash
   sqlc generate
   ```

6. **Run database migrations:**
   ```bash
   cat migrations/001_create_products_table.sql | docker-compose exec -T postgres psql -U postgres -d productdb
   ```

7. **Run the application:**
   ```bash
   go run main.go
   ```

The API will be available at `http://localhost:3000`

## 🧩 XComp Dependency Injection System

XComp provides a powerful dependency injection system similar to NestJS with the following features:

### Core Concepts

1. **Container**: Central registry for services
2. **Modules**: Organize related services and dependencies
3. **Factories**: Functions that create service instances
4. **Injection Tags**: Struct tags for automatic dependency injection

### Basic Usage

#### 1. Service Definition with Injection Tags

```go
type ProductService struct {
    ProductRepo repositories.ProductRepository `inject:"ProductRepository"`
}

type ProductController struct {
    ProductService *services.ProductService `inject:"ProductService"`
}

type ProductRepositoryImpl struct {
    DbConnection *database.DatabaseConnection `inject:"DatabaseConnection"`
    queries      *queries.Queries
}
```

#### 2. Module Creation

```go
// modules/product/product.module.go
func CreateProductModule() xcomp.Module {
    return xcomp.NewModule().
        AddFactory("ProductRepository", func(container *xcomp.Container) any {
            repo := &persistence.ProductRepositoryImpl{}
            container.Inject(repo)  // Auto-inject dependencies using struct tags
            repo.Initialize()       // Custom initialization
            return repo
        }).
        AddFactory("ProductService", func(container *xcomp.Container) any {
            service := &services.ProductService{}
            container.Inject(service)
            return service
        }).
        AddFactory("ProductController", func(container *xcomp.Container) any {
            controller := &controllers.ProductController{}
            container.Inject(controller)
            return controller
        }).
        AddFactory("ProductRoutes", func(container *xcomp.Container) any {
            routes := &routes.ProductRoutes{}
            container.Inject(routes)
            return routes
        }).
        Build()
}
```

#### 3. Module Registration and Usage

```go
// main.go
func main() {
    container := xcomp.NewContainer()

    // Create and register modules
    infrastructureModule := createInfrastructureModule(container)
    productModule := product.CreateProductModule()

    appModule := xcomp.NewModule().
        Import(infrastructureModule).
        Import(productModule).
        Build()

    if err := container.RegisterModule(appModule); err != nil {
        log.Fatalf("Failed to register app module: %v", err)
    }

    // List registered services
    log.Println("Registered Services:")
    for _, serviceName := range container.ListServices() {
        log.Printf("- %s", serviceName)
    }

    // Get services from container
    configService, ok := container.Get("ConfigService").(*xcomp.ConfigService)
    if !ok {
        log.Fatal("Failed to get ConfigService")
    }

    productRoutes, ok := container.Get("ProductRoutes").(*routes.ProductRoutes)
    if !ok {
        log.Fatal("Failed to get ProductRoutes")
    }

    // Use services...
}
```

### Advanced XComp Features

#### 1. Service Names (Optional)

```go
func (ps *ProductService) GetServiceName() string {
    return "ProductService"  // Optional: helps with debugging
}
```

#### 2. Configuration Service

```go
// Built-in configuration service
configService := xcomp.NewConfigService("config.yaml")

// Usage in services
type DatabaseConnection struct {
    Config xcomp.ConfigService `inject:"ConfigService"`
}

func (dc *DatabaseConnection) Initialize() error {
    host := dc.Config.GetString("database.host", "localhost")
    port := dc.Config.GetInt("database.port", 5432)
    // ...
}
```

#### 3. Module Composition

```go
func createAppModule(container *xcomp.Container) xcomp.Module {
    // Create shared infrastructure
    infrastructureModule := createInfrastructureModule(container)

    // Import feature modules
    productModule := product.CreateProductModule()
    userModule := user.CreateUserModule()  // Future module

    return xcomp.NewModule().
        Import(infrastructureModule).
        Import(productModule).
        Import(userModule).
        Build()
}
```

#### 4. Lazy Loading and Initialization

```go
type ProductRepositoryImpl struct {
    DbConnection *database.DatabaseConnection `inject:"DatabaseConnection"`
    queries      *queries.Queries
}

func (pr *ProductRepositoryImpl) Initialize() {
    // Called after dependency injection
    if pr.DbConnection != nil && pr.DbConnection.GetDB() != nil {
        pr.queries = queries.New(pr.DbConnection.GetDB())
    }
}

func (pr *ProductRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*entities.Product, error) {
    if pr.queries == nil {
        pr.Initialize()  // Lazy initialization
    }
    // ...
}
```

### XComp vs Other DI Systems

| Feature | XComp | Wire | Fx |
|---------|-------|------|-----|
| Runtime Injection | ✅ | ❌ | ✅ |
| Struct Tags | ✅ | ❌ | ❌ |
| Module System | ✅ | ❌ | ✅ |
| NestJS-like API | ✅ | ❌ | ❌ |
| Zero Dependencies | ✅ | ✅ | ❌ |

## 🔧 Configuration

### Environment Variables

Create a `.env` file to override configuration:

```env
APP_PORT=3000
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_USERNAME=postgres
DATABASE_PASSWORD=password
DATABASE_DATABASE=productdb
```

### YAML Configuration

Edit `config.yaml` for application settings:

```yaml
app:
  name: "Product API"
  port: 3000
  debug: true

database:
  host: "localhost"
  port: 5432
  username: "postgres"
  password: "password"
  database: "productdb"
  max_connections: 25

server:
  cors:
    enabled: true
```

## 📚 API Endpoints

### Products

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/products` | List all products with pagination |
| GET | `/api/v1/products?category=Electronics` | List products by category |
| GET | `/api/v1/products/search?q=laptop` | Search products |
| GET | `/api/v1/products/:id` | Get product by ID |
| POST | `/api/v1/products` | Create new product |
| PUT | `/api/v1/products/:id` | Update product |
| PATCH | `/api/v1/products/:id/stock` | Update product stock |
| DELETE | `/api/v1/products/:id` | Delete product (soft delete) |

### Health Check

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Application health check |

## 📝 API Examples

### Create Product
```bash
curl -X POST http://localhost:3000/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "MacBook Pro",
    "description": "Apple MacBook Pro with M3 chip",
    "price": 1999.99,
    "stock_quantity": 10,
    "category": "Electronics"
  }'
```

### Get Products
```bash
curl http://localhost:3000/api/v1/products?page=1&page_size=10
```

### Search Products
```bash
curl "http://localhost:3000/api/v1/products/search?q=laptop&page=1&page_size=5"
```

### Update Product Stock
```bash
curl -X PATCH http://localhost:3000/api/v1/products/{id}/stock \
  -H "Content-Type: application/json" \
  -d '{
    "stock_quantity": 25
  }'
```

## 🗄️ Database Schema

The application uses PostgreSQL with SQLC for type-safe queries:

```sql
-- migrations/001_create_products_table.sql
CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(10,2) NOT NULL CHECK (price >= 0),
    stock_quantity INTEGER NOT NULL DEFAULT 0 CHECK (stock_quantity >= 0),
    category VARCHAR(100),
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

### SQLC Integration

SQLC generates type-safe Go code from SQL queries:

```sql
-- modules/product/infrastructure/persistence/queries/products.sql
-- name: GetProduct :one
SELECT id, name, description, price, stock_quantity, category, is_active, created_at, updated_at
FROM products
WHERE id = $1 AND is_active = true;

-- name: ListProducts :many
SELECT id, name, description, price, stock_quantity, category, is_active, created_at, updated_at
FROM products
WHERE is_active = true
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;
```

Generated Go code:
```go
// Auto-generated by SQLC
func (q *Queries) GetProduct(ctx context.Context, id pgtype.UUID) (Product, error) {
    // Type-safe implementation
}
```

## 🧪 Testing

Access the database via Adminer:
- URL: http://localhost:8080
- System: PostgreSQL
- Server: postgres
- Username: postgres
- Password: password
- Database: productdb

## 🔍 Monitoring

- Health check: `GET /health`
- Application logs show all registered services
- Database connection health monitoring
- Graceful shutdown with proper cleanup

## 📈 Adding New Modules

To add a new module (e.g., User module):

1. **Create module structure:**
   ```
   modules/user/
   ├── domain/
   ├── application/
   ├── infrastructure/
   └── user.module.go
   ```

2. **Define the module:**
   ```go
   // modules/user/user.module.go
   func CreateUserModule() xcomp.Module {
       return xcomp.NewModule().
           AddFactory("UserRepository", func(container *xcomp.Container) any {
               repo := &persistence.UserRepositoryImpl{}
               container.Inject(repo)
               return repo
           }).
           // ... other services
           Build()
   }
   ```

3. **Register in main:**
   ```go
   userModule := user.CreateUserModule()
   appModule := xcomp.NewModule().
       Import(infrastructureModule).
       Import(productModule).
       Import(userModule).
       Build()
   ```

## 🎯 Key Features Demonstrated

1. **Modular Architecture**: NestJS/Nx-like module organization
2. **XComp Dependency Injection**: Struct tag-based DI system
3. **Clean Architecture**: Clear separation of concerns
4. **Type-Safe Database**: SQLC integration with PostgreSQL
5. **Configuration Management**: YAML + Environment variables
6. **HTTP API**: RESTful API with Fiber framework
7. **Module System**: Organized service registration and imports
8. **Graceful Shutdown**: Proper resource cleanup
9. **Error Handling**: Comprehensive error handling throughout

This example showcases how to build a production-ready Go application with modular architecture and dependency injection patterns inspired by NestJS, using the XComp framework.

## 🔗 XComp Framework

XComp is a lightweight dependency injection framework for Go inspired by NestJS:
- **Struct Tag Injection**: `inject:"ServiceName"`
- **Module System**: Organize and compose services
- **Configuration Service**: Built-in YAML + env configuration
- **Zero Dependencies**: No external dependencies
- **Type Safety**: Compile-time type checking
- **Runtime Injection**: Dynamic service resolution
