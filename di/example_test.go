package di_test

import (
	"fmt"

	"github.com/quinnjr/go-dependency-injector/di"
)

// ExampleLogger is an example interface for logging.
type ExampleLogger interface {
	Log(message string)
}

// ExampleConsoleLogger is a simple logger that prints to console.
type ExampleConsoleLogger struct{}

func (l *ExampleConsoleLogger) Log(message string) {
	fmt.Println("[LOG]", message)
}

// ExampleFileLogger logs to a file (simulated for example).
type ExampleFileLogger struct {
	Path string
}

func (l *ExampleFileLogger) Log(message string) {
	fmt.Printf("[FILE:%s] %s\n", l.Path, message)
}

// ExampleUserService is an example service interface.
type ExampleUserService interface {
	GetUser(id int) string
}

// ExampleDefaultUserService is the default implementation of ExampleUserService.
type ExampleDefaultUserService struct {
	logger ExampleLogger
}

func (s *ExampleDefaultUserService) GetUser(id int) string {
	s.logger.Log(fmt.Sprintf("Fetching user %d", id))
	return fmt.Sprintf("User-%d", id)
}

// Example demonstrates basic dependency injection usage.
func Example() {
	// Create a new container
	container := di.New()

	// Register a Logger as a singleton
	di.Register[ExampleLogger](container, func() ExampleLogger {
		return &ExampleConsoleLogger{}
	}, di.AsSingleton())

	// Register UserService with Logger as a dependency (auto-resolved)
	di.Register[ExampleUserService](container, func(log ExampleLogger) ExampleUserService {
		return &ExampleDefaultUserService{logger: log}
	})

	// Resolve UserService - Logger is automatically injected
	service := di.MustResolve[ExampleUserService](container)
	user := service.GetUser(42)
	fmt.Println("Got:", user)

	// Output:
	// [LOG] Fetching user 42
	// Got: User-42
}

// ExampleNew demonstrates creating a new container.
func ExampleNew() {
	container := di.New()
	fmt.Printf("Container created: %T\n", container)

	// Output:
	// Container created: *di.Container
}

// ExampleRegister demonstrates registering a dependency with a factory.
func ExampleRegister() {
	container := di.New()

	// Register with a simple factory
	err := di.Register[ExampleLogger](container, func() ExampleLogger {
		return &ExampleConsoleLogger{}
	})
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Verify registration
	if di.Has[ExampleLogger](container) {
		fmt.Println("Logger registered successfully")
	}

	// Output:
	// Logger registered successfully
}

// ExampleRegister_withDependencies demonstrates registering a service
// that has dependencies which are automatically resolved.
func ExampleRegister_withDependencies() {
	container := di.New()

	// Register Logger first
	di.Register[ExampleLogger](container, func() ExampleLogger {
		return &ExampleConsoleLogger{}
	})

	// Register UserService with Logger dependency - it will be auto-resolved
	di.Register[ExampleUserService](container, func(log ExampleLogger) ExampleUserService {
		return &ExampleDefaultUserService{logger: log}
	})

	// When resolving UserService, Logger is automatically injected
	service := di.MustResolve[ExampleUserService](container)
	service.GetUser(1)

	// Output:
	// [LOG] Fetching user 1
}

// ExampleRegister_withError demonstrates registering a factory that can return an error.
func ExampleRegister_withError() {
	container := di.New()

	// Factory that returns (T, error)
	di.Register[ExampleLogger](container, func() (ExampleLogger, error) {
		// Simulate potential initialization error
		return &ExampleConsoleLogger{}, nil
	})

	logger, err := di.Resolve[ExampleLogger](container)
	if err != nil {
		fmt.Println("Resolution failed:", err)
		return
	}
	logger.Log("Hello from error-aware factory")

	// Output:
	// [LOG] Hello from error-aware factory
}

// ExampleRegisterInstance demonstrates registering a pre-created instance.
func ExampleRegisterInstance() {
	container := di.New()

	// Create an instance manually
	logger := &ExampleConsoleLogger{}

	// Register the existing instance
	di.RegisterInstance[ExampleLogger](container, logger)

	// Resolve returns the same instance
	resolved := di.MustResolve[ExampleLogger](container)
	resolved.Log("Hello from pre-created instance")

	// Output:
	// [LOG] Hello from pre-created instance
}

// ExampleRegisterType demonstrates registering an interface to implementation mapping.
func ExampleRegisterType() {
	container := di.New()

	// Register interface → implementation mapping
	di.RegisterType[ExampleLogger, ExampleConsoleLogger](container)

	// Resolve creates the implementation automatically
	logger := di.MustResolve[ExampleLogger](container)
	logger.Log("Hello from auto-created implementation")

	// Output:
	// [LOG] Hello from auto-created implementation
}

// ExampleAsSingleton demonstrates singleton lifetime.
func ExampleAsSingleton() {
	container := di.New()

	callCount := 0
	di.Register[ExampleLogger](container, func() ExampleLogger {
		callCount++
		return &ExampleConsoleLogger{}
	}, di.AsSingleton())

	// Multiple resolutions return the same instance
	_ = di.MustResolve[ExampleLogger](container)
	_ = di.MustResolve[ExampleLogger](container)
	_ = di.MustResolve[ExampleLogger](container)

	fmt.Printf("Factory called %d time(s)\n", callCount)

	// Output:
	// Factory called 1 time(s)
}

// ExampleAsTransient demonstrates transient lifetime.
func ExampleAsTransient() {
	container := di.New()

	callCount := 0
	di.Register[ExampleLogger](container, func() ExampleLogger {
		callCount++
		return &ExampleConsoleLogger{}
	}, di.AsTransient()) // This is the default

	// Each resolution creates a new instance
	_ = di.MustResolve[ExampleLogger](container)
	_ = di.MustResolve[ExampleLogger](container)
	_ = di.MustResolve[ExampleLogger](container)

	fmt.Printf("Factory called %d time(s)\n", callCount)

	// Output:
	// Factory called 3 time(s)
}

// ExampleAsScoped demonstrates scoped lifetime.
func ExampleAsScoped() {
	container := di.New()

	callCount := 0
	di.Register[ExampleLogger](container, func() ExampleLogger {
		callCount++
		return &ExampleConsoleLogger{}
	}, di.AsScoped())

	// Create a scope (e.g., for an HTTP request)
	scope := container.CreateScope("request-1")

	// Multiple resolutions in the same scope return the same instance
	_, _ = di.ResolveInScope[ExampleLogger](container, scope)
	_, _ = di.ResolveInScope[ExampleLogger](container, scope)

	fmt.Printf("Factory called %d time(s) within scope\n", callCount)

	// Output:
	// Factory called 1 time(s) within scope
}

// ExampleResolve demonstrates resolving a dependency.
func ExampleResolve() {
	container := di.New()

	di.Register[ExampleLogger](container, func() ExampleLogger {
		return &ExampleConsoleLogger{}
	})

	// Resolve returns (T, error)
	logger, err := di.Resolve[ExampleLogger](container)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	logger.Log("Resolved successfully")

	// Output:
	// [LOG] Resolved successfully
}

// ExampleMustResolve demonstrates resolving a dependency that panics on error.
func ExampleMustResolve() {
	container := di.New()

	di.Register[ExampleLogger](container, func() ExampleLogger {
		return &ExampleConsoleLogger{}
	})

	// MustResolve panics if resolution fails
	logger := di.MustResolve[ExampleLogger](container)
	logger.Log("Must resolved successfully")

	// Output:
	// [LOG] Must resolved successfully
}

// ExampleWithName demonstrates named registrations.
func ExampleWithName() {
	container := di.New()

	// Register multiple loggers with different names
	di.Register[ExampleLogger](container, func() ExampleLogger {
		return &ExampleConsoleLogger{}
	}, di.WithName("console"))

	di.Register[ExampleLogger](container, func() ExampleLogger {
		return &ExampleFileLogger{Path: "/var/log/app.log"}
	}, di.WithName("file"))

	// Resolve by name
	consoleLogger, _ := di.ResolveNamed[ExampleLogger](container, "console")
	fileLogger, _ := di.ResolveNamed[ExampleLogger](container, "file")

	consoleLogger.Log("Hello console")
	fileLogger.Log("Hello file")

	// Output:
	// [LOG] Hello console
	// [FILE:/var/log/app.log] Hello file
}

// ExampleResolveNamed demonstrates resolving named dependencies.
func ExampleResolveNamed() {
	container := di.New()

	di.Register[ExampleLogger](container, func() ExampleLogger {
		return &ExampleConsoleLogger{}
	}, di.WithName("console"))

	// ResolveNamed gets a specific named registration
	logger, err := di.ResolveNamed[ExampleLogger](container, "console")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	logger.Log("Named resolution works")

	// Output:
	// [LOG] Named resolution works
}

// ExampleHas demonstrates checking if a type is registered.
func ExampleHas() {
	container := di.New()

	fmt.Println("Before registration:", di.Has[ExampleLogger](container))

	di.Register[ExampleLogger](container, func() ExampleLogger {
		return &ExampleConsoleLogger{}
	})

	fmt.Println("After registration:", di.Has[ExampleLogger](container))

	// Output:
	// Before registration: false
	// After registration: true
}

// ExampleHasNamed demonstrates checking if a named type is registered.
func ExampleHasNamed() {
	container := di.New()

	di.Register[ExampleLogger](container, func() ExampleLogger {
		return &ExampleConsoleLogger{}
	}, di.WithName("console"))

	fmt.Println("Has 'console':", di.HasNamed[ExampleLogger](container, "console"))
	fmt.Println("Has 'file':", di.HasNamed[ExampleLogger](container, "file"))

	// Output:
	// Has 'console': true
	// Has 'file': false
}

// ExampleContainer_CreateScope demonstrates creating a resolution scope.
func ExampleContainer_CreateScope() {
	container := di.New()

	di.Register[ExampleLogger](container, func() ExampleLogger {
		return &ExampleConsoleLogger{}
	}, di.AsScoped())

	// Create a scope (useful for request-scoped dependencies)
	scope := container.CreateScope("request-123")
	fmt.Printf("Scope created: %s\n", scope.Name())

	// Output:
	// Scope created: request-123
}

// ExampleResolveInScope demonstrates resolving within a scope.
func ExampleResolveInScope() {
	container := di.New()

	di.Register[ExampleLogger](container, func() ExampleLogger {
		return &ExampleConsoleLogger{}
	}, di.AsScoped())

	scope := container.CreateScope("request-1")

	logger, err := di.ResolveInScope[ExampleLogger](container, scope)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	logger.Log("Scoped resolution works")

	// Output:
	// [LOG] Scoped resolution works
}

// ExampleContainer_Clear demonstrates clearing all registrations.
func ExampleContainer_Clear() {
	container := di.New()

	di.Register[ExampleLogger](container, func() ExampleLogger {
		return &ExampleConsoleLogger{}
	})

	fmt.Println("Before clear:", di.Has[ExampleLogger](container))

	container.Clear()

	fmt.Println("After clear:", di.Has[ExampleLogger](container))

	// Output:
	// Before clear: true
	// After clear: false
}

// ExampleErrNotRegistered demonstrates handling unregistered type errors.
func ExampleErrNotRegistered() {
	container := di.New()

	// Try to resolve an unregistered type
	_, err := di.Resolve[ExampleLogger](container)
	if err != nil {
		fmt.Println("Error:", err)
	}

	// Output:
	// Error: di: type di_test.ExampleLogger is not registered
}

// Example_layeredArchitecture demonstrates a realistic layered architecture setup.
func Example_layeredArchitecture() {
	container := di.New()

	// Infrastructure layer - singletons
	di.Register[ExampleLogger](container, func() ExampleLogger {
		return &ExampleConsoleLogger{}
	}, di.AsSingleton())

	// Service layer - depends on infrastructure
	di.Register[ExampleUserService](container, func(log ExampleLogger) ExampleUserService {
		return &ExampleDefaultUserService{logger: log}
	})

	// Application startup
	service := di.MustResolve[ExampleUserService](container)
	result := service.GetUser(1)
	fmt.Println("Result:", result)

	// Output:
	// [LOG] Fetching user 1
	// Result: User-1
}

// Example_testing demonstrates how DI makes testing easier.
func Example_testing() {
	// In tests, you can swap implementations easily

	// Production container
	prodContainer := di.New()
	di.Register[ExampleLogger](prodContainer, func() ExampleLogger {
		return &ExampleConsoleLogger{}
	})

	// Test container with mock
	testContainer := di.New()
	di.RegisterInstance[ExampleLogger](testContainer, &ExampleMockLogger{})

	// Same resolution code works with different implementations
	prodLogger := di.MustResolve[ExampleLogger](prodContainer)
	testLogger := di.MustResolve[ExampleLogger](testContainer)

	prodLogger.Log("Production message")
	testLogger.Log("Test message")

	// Output:
	// [LOG] Production message
	// [MOCK] Test message
}

// ExampleMockLogger is a mock implementation for testing.
type ExampleMockLogger struct{}

func (l *ExampleMockLogger) Log(message string) {
	fmt.Println("[MOCK]", message)
}
