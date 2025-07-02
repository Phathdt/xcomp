# Pure ConfigService System (NestJS Style)

## Overview

Xcomp framework uses a **pure ConfigService approach** similar to NestJS:
- **YAML configuration files** for structured configuration
- **Environment variable overrides** with double underscore notation
- **ConfigService injection** into all services and modules
- **Dot notation access** like NestJS ConfigService
- **Thread-safe** config access

## NestJS Pattern Comparison

**NestJS:**
```typescript
@Injectable()
export class UserService {
  constructor(private configService: ConfigService) {}

  getDatabaseUrl(): string {
    return this.configService.get<string>('database.url');
  }

  getAppPort(): number {
    return this.configService.get<number>('app.port', 3000);
  }
}
```

**Xcomp (Same Pattern!):**
```go
type UserService struct {
    ConfigService *xcomp.ConfigService `inject:"ConfigService"`
}

func (s *UserService) getDatabaseURL() string {
    return s.ConfigService.GetString("database.url", "")
}

func (s *UserService) getAppPort() int {
    return s.ConfigService.GetInt("app.port", 3000)
}
```

## Environment Variable Mapping

Environment variables use `__` (double underscore) for nested configuration:

| Code Access | Environment Variable | Example Value |
|-------------|---------------------|---------------|
| `configService.GetString("app.name")` | `APP__NAME` | `"My API Server"` |
| `configService.GetInt("app.port")` | `APP__PORT` | `4000` |
| `configService.GetString("database.url")` | `DATABASE__URL` | `"postgresql://..."` |
| `configService.GetBool("server.cors.enabled")` | `SERVER__CORS__ENABLED` | `true` |
| `configService.GetString("logging.level")` | `LOGGING__LEVEL` | `"debug"` |

## Configuration Structure

```yaml
app:
  name: "XComp API Server"
  version: "1.0.0"
  environment: "development"
  debug: true
  port: 3000

database:
  url: "postgresql://postgres:password@localhost:5432/example_db?sslmode=disable"
  max_connections: 25
  max_idle_connections: 10
  max_lifetime_minutes: 30

logging:
  level: "info"
  format: "json"
  development: true
  time_key: "timestamp"
  level_key: "level"
  message_key: "message"
  enable_caller: true

server:
  port: 3000
  host: "0.0.0.0"
  timeout_seconds: 30
  read_timeout_seconds: 30
  write_timeout_seconds: 30
  prefork: false
  cors:
    enabled: true
    allowed_origins: "*"
    allowed_methods: "GET,POST,PUT,DELETE,OPTIONS,PATCH"
    allowed_headers: "Content-Type,Authorization"

redis:
  host: "localhost"
  port: 6379
  password: ""
  db: 0

async:
  monitor:
    port: 8080
    enabled: true
```

## Usage Examples

### Service Injection
```go
type OrderService struct {
    ConfigService *xcomp.ConfigService `inject:"ConfigService"`
    Logger        xcomp.Logger         `inject:"Logger"`
}

func (s *OrderService) getMaxRetries() int {
    return s.ConfigService.GetInt("order.max_retries", 3)
}

func (s *OrderService) getDatabaseURL() string {
    return s.ConfigService.GetString("database.url", "")
}
```

### Environment Overrides
```bash
# Development
export APP__DEBUG=true
export LOGGING__LEVEL=debug
export DATABASE__URL=postgresql://user:pass@localhost:5432/dev_db

# Production
export APP__ENVIRONMENT=production
export APP__DEBUG=false
export LOGGING__LEVEL=warn
export DATABASE__URL=postgresql://user:pass@prod-db:5432/prod_db?sslmode=require

# CORS Configuration
export SERVER__CORS__ENABLED=true
export SERVER__CORS__ALLOWED_ORIGINS="https://api.com,https://app.com"
```

### Module Factory Pattern
```go
func CreateProductModule() xcomp.Module {
    return xcomp.NewModule().
        AddFactory("ProductService", func(c *xcomp.Container) any {
            service := &services.ProductService{}
            // ConfigService auto-injected via `inject:"ConfigService"` tag
            if err := c.Inject(service); err != nil {
                panic("Failed to inject ProductService: " + err.Error())
            }
            return service
        }).
        Build()
}
```

## Configuration Access Methods

| Method | Return Type | Example |
|--------|-------------|---------|
| `GetString(key, default...)` | string | `configService.GetString("app.name", "API")` |
| `GetInt(key, default...)` | int | `configService.GetInt("app.port", 3000)` |
| `GetBool(key, default...)` | bool | `configService.GetBool("app.debug", false)` |
| `Get(key)` | any | `configService.Get("custom.setting")` |

## Benefits of Pure ConfigService

1. **NestJS Compatibility**: Same patterns and mental model
2. **Type Safety**: With proper default values
3. **Environment Overrides**: Clean `SECTION__KEY` mapping
4. **Injection**: Standard dependency injection pattern
5. **Performance**: No struct binding overhead
6. **Simplicity**: Single source of truth
7. **Hot Reload**: Easy config reloading capability
8. **Testing**: Easy to mock in unit tests

## Testing with ConfigService

```go
func TestUserService(t *testing.T) {
    // Create test config
    configService := xcomp.NewConfigService("test-config.yaml")

    // Or mock it
    mockConfig := &MockConfigService{}
    mockConfig.On("GetString", "database.url").Return("test://...")

    service := &UserService{ConfigService: mockConfig}
    // ... test logic
}
```

## Best Practices

1. **Always provide defaults**: `GetString("key", "default")`
2. **Use consistent naming**: `section.subsection.key`
3. **Environment-specific configs**: `config-dev.yaml`, `config-prod.yaml`
4. **Sensitive data via ENV**: `DATABASE__URL`, `REDIS__PASSWORD`
5. **Validation in service layer**: Check required configs in service constructors
6. **Document config keys**: Comment your config access code

This approach provides the familiar NestJS ConfigService experience while maintaining Go's type safety and performance benefits.
