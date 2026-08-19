# Go Dependency Injector

[![Go Reference](https://pkg.go.dev/badge/github.com/quinnjr/go-dependency-injector.svg)](https://pkg.go.dev/github.com/quinnjr/go-dependency-injector)
[![Go Report Card](https://goreportcard.com/badge/github.com/quinnjr/go-dependency-injector)](https://goreportcard.com/report/github.com/quinnjr/go-dependency-injector)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/github/go-mod/go-version/quinnjr/go-dependency-injector)](https://github.com/quinnjr/go-dependency-injector)

A lightweight, type-safe **dependency injection (DI) container** and **IoC (Inversion of Control)** framework for Go. This service container leverages Go generics for compile-time safety, automatic constructor injection, and a clean, intuitive API.

**Keywords**: dependency injection, DI, IoC, inversion of control, service container, service locator, dependency container, golang DI, go dependency injection, constructor injection

## Features

- **Type-safe generics** — Compile-time type checking with `Register[T]()` and `Resolve[T]()`
- **Multiple lifetimes** — Transient, Singleton, and Scoped dependency management
- **Automatic resolution** — Constructor parameters are automatically resolved from the container
- **Circular dependency detection** — Fails fast with clear error messages
- **Named registrations** — Register multiple implementations of the same interface
- **Thread-safe** — Safe for concurrent access across goroutines
- **Zero external dependencies** — Uses only the Go standard library

## Installation

```bash
go get github.com/quinnjr/go-dependency-injector
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/quinnjr/go-dependency-injector/di"
)

// Define interfaces
type Logger interface {
    Log(message string)
}

type UserService interface {
    GetUser(id int) string
}

// Implement them
type ConsoleLogger struct{}

func (l *ConsoleLogger) Log(message string) {
    fmt.Println("[LOG]", message)
}

type DefaultUserService struct {
    logger Logger
}

func (s *DefaultUserService) GetUser(id int) string {
    s.logger.Log(fmt.Sprintf("Fetching user %d", id))
    return fmt.Sprintf("User-%d", id)
}

func main() {
    // Create container
    container := di.New()

    // Register dependencies
    di.Register[Logger](container, func() Logger {
        return &ConsoleLogger{}
    }, di.AsSingleton())

    di.Register[UserService](container, func(log Logger) UserService {
        return &DefaultUserService{logger: log}
    })

    // Resolve and use
    service := di.MustResolve[UserService](container)
    user := service.GetUser(42)
    fmt.Println("Got:", user)
}
```

## API Reference

### Creating a Container

```go
container := di.New()
```

### Registering Dependencies

#### Basic Registration with Factory

```go
// Factory with no dependencies
di.Register[Logger](container, func() Logger {
    return &ConsoleLogger{}
})

// Factory with dependencies (auto-resolved)
di.Register[UserService](container, func(logger Logger, db Database) UserService {
    return &DefaultUserService{logger: logger, db: db}
})

// Factory that can return an error
di.Register[Database](container, func(config Config) (Database, error) {
    return NewPostgresDB(config.DatabaseURL())
})
```

#### Register an Existing Instance

```go
config := &AppConfig{Port: 8080}
di.RegisterInstance[Config](container, config)
```

#### Register Interface → Implementation Mapping

```go
di.RegisterType[UserRepository, PostgresUserRepository](container, di.AsSingleton())
```

### Lifetime Options

| Lifetime | Behavior |
|----------|----------|
| `Transient` | New instance created on every resolution (default) |
| `Singleton` | Single instance shared across all resolutions |
| `Scoped` | Single instance per scope (e.g., per HTTP request) |

```go
// Transient (default)
di.Register[Service](c, factory)
di.Register[Service](c, factory, di.AsTransient())

// Singleton
di.Register[Service](c, factory, di.AsSingleton())

// Scoped
di.Register[Service](c, factory, di.AsScoped())

// Using WithLifetime
di.Register[Service](c, factory, di.WithLifetime(di.Singleton))
```

### Resolving Dependencies

```go
// Returns (T, error)
service, err := di.Resolve[UserService](container)
if err != nil {
    log.Fatal(err)
}

// Panics on error (use when you know registration exists)
service := di.MustResolve[UserService](container)
```

### Named Registrations

Register multiple implementations of the same interface:

```go
// Register with names
di.Register[Logger](c, func() Logger {
    return &ConsoleLogger{}
}, di.WithName("console"))

di.Register[Logger](c, func() Logger {
    return &FileLogger{path: "/var/log/app.log"}
}, di.WithName("file"))

// Resolve by name
consoleLogger, _ := di.ResolveNamed[Logger](c, "console")
fileLogger, _ := di.ResolveNamed[Logger](c, "file")
```

### Scoped Resolution

Scopes are useful for request-scoped dependencies in web applications:

```go
di.Register[*RequestContext](c, func() *RequestContext {
    return &RequestContext{
        RequestID: uuid.New().String(),
        StartTime: time.Now(),
    }
}, di.AsScoped())

// Per HTTP request
func handler(w http.ResponseWriter, r *http.Request) {
    scope := container.CreateScope("request-" + r.Header.Get("X-Request-ID"))

    // Same instance within this scope
    ctx1, _ := di.ResolveInScope[*RequestContext](container, scope)
    ctx2, _ := di.ResolveInScope[*RequestContext](container, scope)
    // ctx1 == ctx2 ✓
}
```

### Utility Methods

```go
// Check if a type is registered
if di.Has[Logger](container) {
    // ...
}

// Check named registration
if di.HasNamed[Logger](container, "file") {
    // ...
}

// Clear all registrations
container.Clear()
```

## Error Handling

The library provides typed errors for precise error handling:

```go
service, err := di.Resolve[UserService](container)
if err != nil {
    switch e := err.(type) {
    case di.ErrNotRegistered:
        log.Printf("Type %s is not registered", e.Type)
    case di.ErrCircularDependency:
        log.Printf("Circular dependency: %v", e.Chain)
    case di.ErrResolutionFailed:
        log.Printf("Failed to resolve %s: %v", e.Type, e.Cause)
    case di.ErrInvalidFactory:
        log.Printf("Invalid factory for %s: %s", e.Type, e.Message)
    }
}
```

| Error Type | When It Occurs |
|------------|----------------|
| `ErrNotRegistered` | Attempting to resolve an unregistered type |
| `ErrCircularDependency` | A → B → A dependency chain detected |
| `ErrResolutionFailed` | Factory returned an error or dependency failed |
| `ErrInvalidFactory` | Factory signature is invalid |
| `ErrScopeNotFound` | Referenced scope doesn't exist |

## Complete Example

Here's a realistic example with a layered architecture:

```go
package main

import (
    "database/sql"
    "log"
    "net/http"

    "github.com/quinnjr/go-dependency-injector/di"
)

// Interfaces
type Config interface {
    DatabaseURL() string
    Port() string
}

type Logger interface {
    Info(msg string)
    Error(msg string)
}

type UserRepository interface {
    FindByID(id int) (*User, error)
}

type UserService interface {
    GetUser(id int) (*User, error)
}

type UserHandler interface {
    ServeHTTP(w http.ResponseWriter, r *http.Request)
}

// Bootstrap the application
func main() {
    container := di.New()

    // Infrastructure layer
    di.Register[Config](container, NewEnvConfig, di.AsSingleton())
    di.Register[Logger](container, NewZapLogger, di.AsSingleton())
    di.Register[*sql.DB](container, func(cfg Config) (*sql.DB, error) {
        return sql.Open("postgres", cfg.DatabaseURL())
    }, di.AsSingleton())

    // Data layer
    di.Register[UserRepository](container, NewPostgresUserRepo)

    // Business layer
    di.Register[UserService](container, NewUserService)

    // Presentation layer
    di.Register[UserHandler](container, NewUserHandler)

    // Start server
    handler := di.MustResolve[UserHandler](container)
    config := di.MustResolve[Config](container)

    log.Printf("Starting server on %s", config.Port())
    http.ListenAndServe(config.Port(), handler)
}
```

## Best Practices

### 1. Register at Startup

Register all dependencies during application startup, not at runtime:

```go
func main() {
    container := di.New()
    registerDependencies(container)  // All registrations here
    runApplication(container)
}
```

### 2. Depend on Interfaces

Always register and resolve interfaces, not concrete types:

```go
// ✓ Good
di.Register[Logger](c, func() Logger { return &ConsoleLogger{} })

// ✗ Avoid
di.Register[*ConsoleLogger](c, func() *ConsoleLogger { return &ConsoleLogger{} })
```

### 3. Use Appropriate Lifetimes

- **Singleton**: Stateless services, configuration, connection pools
- **Transient**: Stateful objects, request-specific data
- **Scoped**: Per-request context, unit of work patterns

### 4. Keep Factories Simple

Factories should only create objects, not perform business logic:

```go
// ✓ Good
di.Register[UserService](c, func(repo UserRepository, log Logger) UserService {
    return &DefaultUserService{repo: repo, logger: log}
})

// ✗ Avoid
di.Register[UserService](c, func(repo UserRepository) UserService {
    users, _ := repo.FindAll()  // Don't do this!
    return &UserService{cachedUsers: users}
})
```

## Testing

The DI container makes testing easy by allowing you to swap implementations:

```go
func TestUserService(t *testing.T) {
    container := di.New()

    // Register mock dependencies
    di.RegisterInstance[Logger](container, &MockLogger{})
    di.RegisterInstance[UserRepository](container, &MockUserRepo{
        users: map[int]*User{1: {ID: 1, Name: "Test"}},
    })

    // Register the real service
    di.Register[UserService](container, NewUserService)

    // Test
    service := di.MustResolve[UserService](container)
    user, err := service.GetUser(1)

    assert.NoError(t, err)
    assert.Equal(t, "Test", user.Name)
}
```

## Thread Safety

The container is fully thread-safe. You can safely:

- Register dependencies from multiple goroutines
- Resolve dependencies concurrently
- Create scopes in parallel

## Requirements

- Go 1.22 or later (uses generics)

## License

MIT License - see [LICENSE](LICENSE) for details.

## Package Discovery

This package is automatically indexed by [pkg.go.dev](https://pkg.go.dev/github.com/quinnjr/go-dependency-injector) once you create a version tag. Users can:

- Browse documentation at https://pkg.go.dev/github.com/quinnjr/go-dependency-injector
- Import directly: `go get github.com/quinnjr/go-dependency-injector`
- View examples and API documentation

### For Maintainers

To publish a new version:

```bash
# Create and push a new semantic version tag
git tag v1.0.0
git push origin v1.0.0

# The package will be automatically indexed by pkg.go.dev
# You can verify at: https://pkg.go.dev/github.com/quinnjr/go-dependency-injector@v1.0.0
```

## Related Projects

- [uber-go/dig](https://github.com/uber-go/dig) - Reflection-based DI framework
- [google/wire](https://github.com/google/wire) - Compile-time DI code generator
- [samber/do](https://github.com/samber/do) - Generic-based DI container

## Search Terms

If you're looking for any of these, you've found the right package:

- **Dependency Injection in Go** / **Golang DI** - This package provides full DI support
- **IoC Container for Go** / **Inversion of Control** - Implements the IoC pattern
- **Service Container** / **Service Locator** - Acts as a service container for your application
- **Constructor Injection** - Automatically injects dependencies via constructor functions
- **Lifetime Management** - Supports Transient, Singleton, and Scoped lifetimes
- **Interface-based DI** - Designed for programming to interfaces

## FAQ

### Why use this over uber-go/dig or google/wire?

- **Type Safety**: Uses generics for compile-time type checking
- **Simplicity**: Minimal API surface with intuitive methods
- **Zero Dependencies**: No external dependencies
- **Modern Go**: Built for Go 1.22+ with generics

### Can I use this in production?

Yes! The library is fully tested, thread-safe, and follows Go best practices.

### How does performance compare?

See `di/benchmark_test.go` for benchmarks. Performance is comparable to other DI solutions with minimal overhead.

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

Please feel free to submit a Pull Request or open an issue for bugs, features, or questions.

