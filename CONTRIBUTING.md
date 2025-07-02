# Contributing to XComp

Thank you for your interest in contributing to XComp! We welcome contributions from everyone.

## 🚀 Getting Started

### Prerequisites

- Go 1.24 or higher
- Git
- Docker and Docker Compose (for testing)

### Development Setup

1. **Fork and Clone**
   ```bash
   git clone https://github.com/your-username/xcomp.git
   cd xcomp
   ```

2. **Install Dependencies**
   ```bash
   go mod download
   ```

3. **Run Example Application**
   ```bash
   cd example
   make dev-setup  # Sets up Docker + database + migrations
   make run-dev    # Start development server
   ```

4. **Run Tests**
   ```bash
   make test
   make test-coverage
   ```

## 🛠️ Development Workflow

### Making Changes

1. **Create a Feature Branch**
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make Your Changes**
   - Follow Go conventions and best practices
   - Add tests for new functionality
   - Update documentation as needed

3. **Test Your Changes**
   ```bash
   # Run tests
   make test

   # Check code quality
   make lint
   make format

   # Test example application
   cd example
   make dev-setup
   make run-dev
   ```

4. **Commit Your Changes**
   ```bash
   git add .
   git commit -m "feat: add your feature description"
   ```

### Commit Message Format

We follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

- `feat:` - New feature
- `fix:` - Bug fix
- `docs:` - Documentation changes
- `refactor:` - Code refactoring
- `test:` - Adding or updating tests
- `chore:` - Build process or auxiliary tool changes

Examples:
```
feat: add support for custom service factories
fix: resolve dependency injection circular reference issue
docs: update API documentation for Container methods
refactor: improve logger configuration handling
test: add integration tests for module system
```

## 📋 Code Guidelines

### Go Code Style

- Follow standard Go formatting (`gofmt`)
- Use meaningful variable and function names
- Add comments for exported functions and types
- Keep functions small and focused
- Use interfaces for better testability

### Testing

- Write unit tests for all new functionality
- Use table-driven tests where appropriate
- Mock external dependencies
- Aim for high test coverage

### Documentation

- Update README.md for new features
- Add inline code comments for complex logic
- Include examples in documentation
- Keep documentation up to date

## 🧪 Testing

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run benchmarks
make benchmark

# Run linter
make lint
```

### Test Structure

Place tests in the same directory as the code being tested:

```
container.go
container_test.go
logger.go
logger_test.go
```

### Example Test

```go
func TestContainer_Register(t *testing.T) {
    tests := []struct {
        name        string
        serviceName string
        service     any
        expectError bool
    }{
        {
            name:        "register valid service",
            serviceName: "TestService",
            service:     &TestService{},
            expectError: false,
        },
        // Add more test cases...
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            container := NewContainer()
            container.Register(tt.serviceName, tt.service)

            result := container.Get(tt.serviceName)
            assert.NotNil(t, result)
        })
    }
}
```

## 📚 Areas for Contribution

We welcome contributions in these areas:

### Core Framework
- **Guards System** - Authentication/authorization middleware
- **Pipes System** - Request validation and transformation
- **Interceptors** - Request/response transformation
- **Exception Filters** - Centralized error handling
- **Event System** - Event-driven architecture support

### Infrastructure
- **Database Integrations** - Additional database drivers
- **Cache Providers** - Redis, Memcached, in-memory caching
- **Message Queues** - RabbitMQ, Apache Kafka integration
- **Health Checks** - System health monitoring

### Developer Experience
- **CLI Tools** - Project scaffolding and generators
- **IDE Plugins** - Code completion and snippets
- **Documentation** - Tutorials, guides, and examples
- **Testing Utilities** - Testing helpers and mocks

### Example Applications
- **GraphQL API** - GraphQL server example
- **gRPC Service** - gRPC microservice example
- **WebSocket Gateway** - Real-time communication example
- **Background Jobs** - Task queue processing example

## 🐛 Reporting Bugs

When reporting bugs, please include:

1. **Go version** (`go version`)
2. **XComp version**
3. **Operating system**
4. **Minimal reproduction code**
5. **Expected vs actual behavior**
6. **Error messages or logs**

## 💡 Suggesting Features

For feature requests:

1. **Check existing issues** first
2. **Describe the use case** clearly
3. **Provide examples** of how it would be used
4. **Consider backward compatibility**
5. **Discuss implementation approach**

## 📝 Pull Request Process

1. **Ensure tests pass** and coverage is maintained
2. **Update documentation** as needed
3. **Follow commit message format**
4. **Keep changes focused** - one feature per PR
5. **Respond to review feedback** promptly

### PR Checklist

- [ ] Tests added/updated and passing
- [ ] Documentation updated
- [ ] Code follows style guidelines
- [ ] Commit messages follow convention
- [ ] No breaking changes (or clearly documented)
- [ ] Example application still works

## 🤝 Code of Conduct

- Be respectful and inclusive
- Focus on constructive feedback
- Help others learn and grow
- Maintain a positive community

## 📞 Getting Help

- **Documentation**: Check README.md and example application
- **Issues**: Search existing issues before creating new ones
- **Discussions**: Use for questions and general discussion

## 🙏 Recognition

Contributors will be:
- Added to the contributors list
- Mentioned in release notes for significant contributions
- Invited to join the core team for sustained contributions

Thank you for helping make XComp better! 🚀
