package di_test

import (
	"errors"
	"reflect"
	"sync"
	"testing"

	"github.com/quinnjr/go-dependency-injector/di"
)

// =============================================================================
// Test Interfaces and Implementations
// =============================================================================

type Greeter interface {
	Greet(name string) string
}

type SimpleGreeter struct{}

func (g *SimpleGreeter) Greet(name string) string {
	return "Hello, " + name
}

type Logger interface {
	Log(msg string)
}

type TestLogger struct {
	Messages []string
}

func (l *TestLogger) Log(msg string) {
	l.Messages = append(l.Messages, msg)
}

type Service interface {
	DoWork() string
}

type DefaultService struct {
	logger  Logger
	greeter Greeter
}

func (s *DefaultService) DoWork() string {
	s.logger.Log("doing work")
	return s.greeter.Greet("World")
}

type formalGreeter struct{}

func (g *formalGreeter) Greet(name string) string {
	return "Good day, " + name
}

// =============================================================================
// Container Tests
// =============================================================================

func TestNew(t *testing.T) {
	c := di.New()
	if c == nil {
		t.Error("New() should return a non-nil container")
	}
}

func TestRegisterAndResolve(t *testing.T) {
	c := di.New()

	err := di.Register[Greeter](c, func() Greeter {
		return &SimpleGreeter{}
	})
	if err != nil {
		t.Fatalf("failed to register: %v", err)
	}

	greeter, err := di.Resolve[Greeter](c)
	if err != nil {
		t.Fatalf("failed to resolve: %v", err)
	}

	result := greeter.Greet("Test")
	if result != "Hello, Test" {
		t.Errorf("expected 'Hello, Test', got '%s'", result)
	}
}

func TestSingleton(t *testing.T) {
	c := di.New()

	callCount := 0
	err := di.Register[*TestLogger](c, func() *TestLogger {
		callCount++
		return &TestLogger{}
	}, di.AsSingleton())
	if err != nil {
		t.Fatalf("failed to register: %v", err)
	}

	logger1, _ := di.Resolve[*TestLogger](c)
	logger2, _ := di.Resolve[*TestLogger](c)

	if callCount != 1 {
		t.Errorf("factory should be called once for singleton, called %d times", callCount)
	}

	// Verify they're the same instance
	logger1.Log("test")
	if len(logger2.Messages) != 1 {
		t.Error("expected same instance for singleton")
	}
}

func TestTransient(t *testing.T) {
	c := di.New()

	callCount := 0
	err := di.Register[*TestLogger](c, func() *TestLogger {
		callCount++
		return &TestLogger{}
	}, di.AsTransient())
	if err != nil {
		t.Fatalf("failed to register: %v", err)
	}

	_, _ = di.Resolve[*TestLogger](c)
	_, _ = di.Resolve[*TestLogger](c)

	if callCount != 2 {
		t.Errorf("factory should be called twice for transient, called %d times", callCount)
	}
}

func TestRegisterInstance(t *testing.T) {
	c := di.New()

	instance := &TestLogger{Messages: []string{"pre-existing"}}
	di.RegisterInstance[*TestLogger](c, instance)

	resolved, err := di.Resolve[*TestLogger](c)
	if err != nil {
		t.Fatalf("failed to resolve: %v", err)
	}

	if len(resolved.Messages) != 1 || resolved.Messages[0] != "pre-existing" {
		t.Error("expected pre-existing instance")
	}
}

func TestRegisterInstanceWithName(t *testing.T) {
	c := di.New()

	instance1 := &TestLogger{Messages: []string{"logger1"}}
	instance2 := &TestLogger{Messages: []string{"logger2"}}

	di.RegisterInstance[*TestLogger](c, instance1, di.WithName("primary"))
	di.RegisterInstance[*TestLogger](c, instance2, di.WithName("secondary"))

	resolved1, err := di.ResolveNamed[*TestLogger](c, "primary")
	if err != nil {
		t.Fatalf("failed to resolve primary: %v", err)
	}

	resolved2, err := di.ResolveNamed[*TestLogger](c, "secondary")
	if err != nil {
		t.Fatalf("failed to resolve secondary: %v", err)
	}

	if resolved1.Messages[0] != "logger1" {
		t.Error("expected logger1")
	}
	if resolved2.Messages[0] != "logger2" {
		t.Error("expected logger2")
	}
}

func TestRegisterType(t *testing.T) {
	c := di.New()

	err := di.RegisterType[Greeter, SimpleGreeter](c)
	if err != nil {
		t.Fatalf("failed to register type: %v", err)
	}

	greeter, err := di.Resolve[Greeter](c)
	if err != nil {
		t.Fatalf("failed to resolve: %v", err)
	}

	result := greeter.Greet("Test")
	if result != "Hello, Test" {
		t.Errorf("expected 'Hello, Test', got '%s'", result)
	}
}

func TestRegisterTypeAsSingleton(t *testing.T) {
	c := di.New()

	err := di.RegisterType[Greeter, SimpleGreeter](c, di.AsSingleton())
	if err != nil {
		t.Fatalf("failed to register type: %v", err)
	}

	greeter1, _ := di.Resolve[Greeter](c)
	greeter2, _ := di.Resolve[Greeter](c)

	// Both should be the same instance
	if greeter1 != greeter2 {
		t.Error("singleton should return same instance")
	}
}

func TestRegisterTypeWithName(t *testing.T) {
	c := di.New()

	err := di.RegisterType[Greeter, SimpleGreeter](c, di.WithName("simple"))
	if err != nil {
		t.Fatalf("failed to register type: %v", err)
	}

	greeter, err := di.ResolveNamed[Greeter](c, "simple")
	if err != nil {
		t.Fatalf("failed to resolve: %v", err)
	}

	if greeter == nil {
		t.Error("expected greeter instance")
	}
}

func TestDependencyInjection(t *testing.T) {
	c := di.New()

	// Register Logger
	err := di.Register[Logger](c, func() Logger {
		return &TestLogger{}
	}, di.AsSingleton())
	if err != nil {
		t.Fatalf("failed to register Logger: %v", err)
	}

	// Register Greeter
	err = di.Register[Greeter](c, func() Greeter {
		return &SimpleGreeter{}
	})
	if err != nil {
		t.Fatalf("failed to register Greeter: %v", err)
	}

	// Register Service with dependencies
	err = di.Register[Service](c, func(l Logger, g Greeter) Service {
		return &DefaultService{logger: l, greeter: g}
	})
	if err != nil {
		t.Fatalf("failed to register Service: %v", err)
	}

	// Resolve Service - dependencies should be auto-resolved
	service, err := di.Resolve[Service](c)
	if err != nil {
		t.Fatalf("failed to resolve Service: %v", err)
	}

	result := service.DoWork()
	if result != "Hello, World" {
		t.Errorf("expected 'Hello, World', got '%s'", result)
	}
}

func TestNotRegisteredError(t *testing.T) {
	c := di.New()

	_, err := di.Resolve[Greeter](c)
	if err == nil {
		t.Fatal("expected error for unregistered type")
	}

	var notRegErr di.ErrNotRegistered
	if !errors.As(err, &notRegErr) {
		t.Errorf("expected ErrNotRegistered, got %T", err)
	}

	// Test error message
	errMsg := notRegErr.Error()
	if errMsg == "" {
		t.Error("error message should not be empty")
	}
}

func TestCircularDependencyDetection(t *testing.T) {
	c := di.New()

	// A depends on B
	type ServiceA interface{ A() }
	type ServiceB interface{ B() }

	// Create circular dependency: A -> B -> A
	di.Register[ServiceA](c, func(b ServiceB) ServiceA {
		return nil // Implementation doesn't matter for this test
	})
	di.Register[ServiceB](c, func(a ServiceA) ServiceB {
		return nil
	})

	_, err := di.Resolve[ServiceA](c)
	if err == nil {
		t.Fatal("expected circular dependency error")
	}

	var circErr di.ErrCircularDependency
	if !errors.As(err, &circErr) {
		t.Errorf("expected ErrCircularDependency, got %T: %v", err, err)
	}

	// Test error message
	errMsg := circErr.Error()
	if errMsg == "" {
		t.Error("error message should not be empty")
	}
}

func TestNamedRegistrations(t *testing.T) {
	c := di.New()

	// Register two different Greeter implementations with names
	err := di.Register[Greeter](c, func() Greeter {
		return &SimpleGreeter{}
	}, di.WithName("simple"))
	if err != nil {
		t.Fatalf("failed to register simple: %v", err)
	}

	err = di.Register[Greeter](c, func() Greeter {
		return &formalGreeter{}
	}, di.WithName("formal"))
	if err != nil {
		t.Fatalf("failed to register formal: %v", err)
	}

	simple, err := di.ResolveNamed[Greeter](c, "simple")
	if err != nil {
		t.Fatalf("failed to resolve simple: %v", err)
	}

	formal, err := di.ResolveNamed[Greeter](c, "formal")
	if err != nil {
		t.Fatalf("failed to resolve formal: %v", err)
	}

	if simple.Greet("Test") != "Hello, Test" {
		t.Error("simple greeter incorrect")
	}

	if formal.Greet("Test") != "Good day, Test" {
		t.Error("formal greeter incorrect")
	}
}

func TestHas(t *testing.T) {
	c := di.New()

	if di.Has[Greeter](c) {
		t.Error("expected Has to return false before registration")
	}

	di.Register[Greeter](c, func() Greeter { return &SimpleGreeter{} })

	if !di.Has[Greeter](c) {
		t.Error("expected Has to return true after registration")
	}
}

func TestHasNamed(t *testing.T) {
	c := di.New()

	if di.HasNamed[Greeter](c, "custom") {
		t.Error("expected HasNamed to return false before registration")
	}

	di.Register[Greeter](c, func() Greeter { return &SimpleGreeter{} }, di.WithName("custom"))

	if !di.HasNamed[Greeter](c, "custom") {
		t.Error("expected HasNamed to return true after registration")
	}

	// Default name should still not be registered
	if di.Has[Greeter](c) {
		t.Error("expected Has (no name) to return false")
	}
}

func TestMustResolve(t *testing.T) {
	c := di.New()
	di.Register[Greeter](c, func() Greeter { return &SimpleGreeter{} })

	// Should not panic
	greeter := di.MustResolve[Greeter](c)
	if greeter == nil {
		t.Error("expected greeter instance")
	}
}

func TestMustResolvePanics(t *testing.T) {
	c := di.New()

	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for unregistered type")
		}
	}()

	di.MustResolve[Greeter](c)
}

func TestMustResolveNamed(t *testing.T) {
	c := di.New()
	di.Register[Greeter](c, func() Greeter { return &SimpleGreeter{} }, di.WithName("named"))

	// Should not panic
	greeter := di.MustResolveNamed[Greeter](c, "named")
	if greeter == nil {
		t.Error("expected greeter instance")
	}
}

func TestMustResolveNamedPanics(t *testing.T) {
	c := di.New()

	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for unregistered named type")
		}
	}()

	di.MustResolveNamed[Greeter](c, "nonexistent")
}

func TestFactoryWithErrorReturn(t *testing.T) {
	c := di.New()

	expectedErr := errors.New("factory error")

	err := di.Register[Greeter](c, func() (Greeter, error) {
		return nil, expectedErr
	})
	if err != nil {
		t.Fatalf("failed to register: %v", err)
	}

	_, err = di.Resolve[Greeter](c)
	if err == nil {
		t.Fatal("expected error from factory")
	}

	var resErr di.ErrResolutionFailed
	if !errors.As(err, &resErr) {
		t.Errorf("expected ErrResolutionFailed, got %T", err)
	}

	// Test Unwrap
	if !errors.Is(err, expectedErr) {
		t.Error("expected to unwrap to original error")
	}
}

func TestFactoryWithSuccessfulErrorReturn(t *testing.T) {
	c := di.New()

	err := di.Register[Greeter](c, func() (Greeter, error) {
		return &SimpleGreeter{}, nil
	})
	if err != nil {
		t.Fatalf("failed to register: %v", err)
	}

	greeter, err := di.Resolve[Greeter](c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if greeter.Greet("Test") != "Hello, Test" {
		t.Error("greeter not working correctly")
	}
}

func TestScopedLifetime(t *testing.T) {
	c := di.New()

	callCount := 0
	err := di.Register[*TestLogger](c, func() *TestLogger {
		callCount++
		return &TestLogger{}
	}, di.AsScoped())
	if err != nil {
		t.Fatalf("failed to register: %v", err)
	}

	// Create two scopes
	scope1 := c.CreateScope("scope1")
	scope2 := c.CreateScope("scope2")

	// Resolve in scope1 twice - should return same instance
	logger1a, _ := di.ResolveInScope[*TestLogger](c, scope1)
	logger1b, _ := di.ResolveInScope[*TestLogger](c, scope1)

	// Resolve in scope2 - should return different instance
	logger2, _ := di.ResolveInScope[*TestLogger](c, scope2)

	// Verify scope1 returns same instance
	logger1a.Log("from 1a")
	if len(logger1b.Messages) != 1 {
		t.Error("expected same instance within scope")
	}

	// Verify scope2 has different instance
	if len(logger2.Messages) != 0 {
		t.Error("expected different instance in different scope")
	}

	// Should have created 2 instances (one per scope)
	if callCount != 2 {
		t.Errorf("expected 2 factory calls for 2 scopes, got %d", callCount)
	}
}

func TestScopedWithoutScope(t *testing.T) {
	c := di.New()

	callCount := 0
	err := di.Register[*TestLogger](c, func() *TestLogger {
		callCount++
		return &TestLogger{}
	}, di.AsScoped())
	if err != nil {
		t.Fatalf("failed to register: %v", err)
	}

	// Resolve without scope (nil) - should still work but not cache
	logger1, err := di.ResolveInScope[*TestLogger](c, nil)
	if err != nil {
		t.Fatalf("failed to resolve: %v", err)
	}

	logger2, err := di.ResolveInScope[*TestLogger](c, nil)
	if err != nil {
		t.Fatalf("failed to resolve: %v", err)
	}

	// Without scope, should create new instance each time
	if callCount != 2 {
		t.Errorf("expected 2 factory calls without scope, got %d", callCount)
	}

	// Verify they're different instances
	logger1.Log("test")
	if len(logger2.Messages) != 0 {
		t.Error("expected different instances without scope")
	}
}

func TestScopeName(t *testing.T) {
	c := di.New()
	scope := c.CreateScope("my-scope")

	if scope.Name() != "my-scope" {
		t.Errorf("expected scope name 'my-scope', got '%s'", scope.Name())
	}
}

func TestClear(t *testing.T) {
	c := di.New()

	di.Register[Greeter](c, func() Greeter { return &SimpleGreeter{} })

	if !di.Has[Greeter](c) {
		t.Error("expected registration to exist")
	}

	c.Clear()

	if di.Has[Greeter](c) {
		t.Error("expected registration to be cleared")
	}
}

func TestClearRemovesSingletons(t *testing.T) {
	c := di.New()

	callCount := 0
	di.Register[*TestLogger](c, func() *TestLogger {
		callCount++
		return &TestLogger{}
	}, di.AsSingleton())

	// Resolve to cache singleton
	di.Resolve[*TestLogger](c)
	if callCount != 1 {
		t.Error("expected singleton to be created")
	}

	// Clear and re-register
	c.Clear()
	di.Register[*TestLogger](c, func() *TestLogger {
		callCount++
		return &TestLogger{}
	}, di.AsSingleton())

	// Should create new instance
	di.Resolve[*TestLogger](c)
	if callCount != 2 {
		t.Error("expected new singleton after clear")
	}
}

func TestClearRemovesScopes(t *testing.T) {
	c := di.New()

	di.Register[*TestLogger](c, func() *TestLogger {
		return &TestLogger{}
	}, di.AsScoped())

	scope := c.CreateScope("test")
	di.ResolveInScope[*TestLogger](c, scope)

	c.Clear()

	// After clear, scopes should be gone
	// Re-register and create new scope
	di.Register[*TestLogger](c, func() *TestLogger {
		return &TestLogger{}
	}, di.AsScoped())

	newScope := c.CreateScope("test")
	_, err := di.ResolveInScope[*TestLogger](c, newScope)
	if err != nil {
		t.Fatalf("failed to resolve after clear: %v", err)
	}
}

// =============================================================================
// Factory Validation Tests
// =============================================================================

func TestInvalidFactory(t *testing.T) {
	c := di.New()

	// Not a function
	err := di.Register[Greeter](c, "not a function")
	if err == nil {
		t.Fatal("expected error for non-function factory")
	}

	var invalidErr di.ErrInvalidFactory
	if !errors.As(err, &invalidErr) {
		t.Errorf("expected ErrInvalidFactory, got %T", err)
	}

	// Test error message
	errMsg := invalidErr.Error()
	if errMsg == "" {
		t.Error("error message should not be empty")
	}
}

func TestInvalidFactoryNoReturn(t *testing.T) {
	c := di.New()

	err := di.Register[Greeter](c, func() {})
	if err == nil {
		t.Fatal("expected error for factory with no return")
	}
}

func TestInvalidFactoryWrongReturnType(t *testing.T) {
	c := di.New()

	err := di.Register[Greeter](c, func() string { return "wrong" })
	if err == nil {
		t.Fatal("expected error for factory with wrong return type")
	}
}

func TestInvalidFactoryTooManyReturns(t *testing.T) {
	c := di.New()

	err := di.Register[Greeter](c, func() (Greeter, error, string) {
		return nil, nil, ""
	})
	if err == nil {
		t.Fatal("expected error for factory with more than 2 returns")
	}
}

func TestInvalidFactorySecondReturnNotError(t *testing.T) {
	c := di.New()

	err := di.Register[Greeter](c, func() (Greeter, string) {
		return nil, ""
	})
	if err == nil {
		t.Fatal("expected error for factory with non-error second return")
	}
}

func TestFactoryWithUnregisteredDependency(t *testing.T) {
	c := di.New()

	// Register Service that depends on unregistered Logger
	err := di.Register[Service](c, func(l Logger) Service {
		return &DefaultService{logger: l}
	})
	if err != nil {
		t.Fatalf("failed to register: %v", err)
	}

	_, err = di.Resolve[Service](c)
	if err == nil {
		t.Fatal("expected error for unregistered dependency")
	}

	var resErr di.ErrResolutionFailed
	if !errors.As(err, &resErr) {
		t.Errorf("expected ErrResolutionFailed, got %T", err)
	}
}

// =============================================================================
// Lifetime and Options Tests
// =============================================================================

func TestWithLifetime(t *testing.T) {
	c := di.New()

	callCount := 0
	err := di.Register[*TestLogger](c, func() *TestLogger {
		callCount++
		return &TestLogger{}
	}, di.WithLifetime(di.Singleton))
	if err != nil {
		t.Fatalf("failed to register: %v", err)
	}

	di.Resolve[*TestLogger](c)
	di.Resolve[*TestLogger](c)

	if callCount != 1 {
		t.Errorf("expected singleton behavior with WithLifetime, got %d calls", callCount)
	}
}

func TestLifetimeString(t *testing.T) {
	tests := []struct {
		lifetime di.Lifetime
		expected string
	}{
		{di.Transient, "Transient"},
		{di.Singleton, "Singleton"},
		{di.Scoped, "Scoped"},
		{di.Lifetime(99), "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := tt.lifetime.String(); got != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, got)
			}
		})
	}
}

// =============================================================================
// Error Tests
// =============================================================================

func TestErrNotRegisteredError(t *testing.T) {
	err := di.ErrNotRegistered{Type: reflect.TypeOf("")}
	msg := err.Error()
	if msg == "" {
		t.Error("error message should not be empty")
	}
	if !contains(msg, "not registered") {
		t.Error("error message should mention 'not registered'")
	}
}

func TestErrCircularDependencyError(t *testing.T) {
	err := di.ErrCircularDependency{
		Chain: []reflect.Type{
			reflect.TypeOf(""),
			reflect.TypeOf(0),
		},
	}
	msg := err.Error()
	if msg == "" {
		t.Error("error message should not be empty")
	}
	if !contains(msg, "circular") {
		t.Error("error message should mention 'circular'")
	}
}

func TestErrResolutionFailedError(t *testing.T) {
	cause := errors.New("underlying error")
	err := di.ErrResolutionFailed{
		Type:  reflect.TypeOf(""),
		Cause: cause,
	}
	msg := err.Error()
	if msg == "" {
		t.Error("error message should not be empty")
	}
	if !contains(msg, "failed to resolve") {
		t.Error("error message should mention 'failed to resolve'")
	}

	// Test Unwrap
	if err.Unwrap() != cause {
		t.Error("Unwrap should return cause")
	}
}

func TestErrInvalidFactoryError(t *testing.T) {
	err := di.ErrInvalidFactory{
		Type:    reflect.TypeOf(""),
		Message: "test message",
	}
	msg := err.Error()
	if msg == "" {
		t.Error("error message should not be empty")
	}
	if !contains(msg, "invalid factory") {
		t.Error("error message should mention 'invalid factory'")
	}
}

func TestErrScopeNotFoundError(t *testing.T) {
	err := di.ErrScopeNotFound{Name: "test-scope"}
	msg := err.Error()
	if msg == "" {
		t.Error("error message should not be empty")
	}
	if !contains(msg, "test-scope") {
		t.Error("error message should contain scope name")
	}
}

// =============================================================================
// Concurrency Tests
// =============================================================================

func TestConcurrentRegistration(t *testing.T) {
	c := di.New()

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			di.Register[Greeter](c, func() Greeter {
				return &SimpleGreeter{}
			})
		}(i)
	}
	wg.Wait()

	// Should be able to resolve after concurrent registrations
	greeter, err := di.Resolve[Greeter](c)
	if err != nil {
		t.Fatalf("failed to resolve: %v", err)
	}
	if greeter == nil {
		t.Error("expected greeter instance")
	}
}

func TestConcurrentResolution(t *testing.T) {
	c := di.New()

	callCount := 0
	var mu sync.Mutex

	// Use a shared instance to ensure singleton behavior
	sharedLogger := &TestLogger{}

	di.Register[*TestLogger](c, func() *TestLogger {
		mu.Lock()
		callCount++
		mu.Unlock()
		return sharedLogger
	}, di.AsSingleton())

	var wg sync.WaitGroup
	errs := make([]error, 100)

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			_, err := di.Resolve[*TestLogger](c)
			errs[i] = err
		}(i)
	}
	wg.Wait()

	// Check no errors occurred
	for i, err := range errs {
		if err != nil {
			t.Errorf("goroutine %d failed: %v", i, err)
		}
	}
}

func TestConcurrentScopedResolution(t *testing.T) {
	c := di.New()

	di.Register[*TestLogger](c, func() *TestLogger {
		return &TestLogger{}
	}, di.AsScoped())

	scope := c.CreateScope("concurrent-scope")

	var wg sync.WaitGroup
	errs := make([]error, 50)

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			_, err := di.ResolveInScope[*TestLogger](c, scope)
			errs[i] = err
		}(i)
	}
	wg.Wait()

	// Verify no errors occurred (main goal: thread-safe access)
	for i, err := range errs {
		if err != nil {
			t.Errorf("goroutine %d failed: %v", i, err)
		}
	}
}

// =============================================================================
// Edge Cases
// =============================================================================

func TestOverwriteRegistration(t *testing.T) {
	c := di.New()

	// Register first implementation
	di.Register[Greeter](c, func() Greeter {
		return &SimpleGreeter{}
	})

	// Overwrite with second implementation
	di.Register[Greeter](c, func() Greeter {
		return &formalGreeter{}
	})

	greeter, _ := di.Resolve[Greeter](c)
	result := greeter.Greet("Test")

	// Should use the second registration
	if result != "Good day, Test" {
		t.Errorf("expected 'Good day, Test', got '%s'", result)
	}
}

func TestResolveNamedNotFound(t *testing.T) {
	c := di.New()

	// Register without name
	di.Register[Greeter](c, func() Greeter {
		return &SimpleGreeter{}
	})

	// Try to resolve with name
	_, err := di.ResolveNamed[Greeter](c, "nonexistent")
	if err == nil {
		t.Fatal("expected error for unregistered named type")
	}
}

func TestEmptyContainer(t *testing.T) {
	c := di.New()

	if di.Has[Greeter](c) {
		t.Error("empty container should not have any registrations")
	}

	_, err := di.Resolve[Greeter](c)
	if err == nil {
		t.Error("empty container should error on resolve")
	}
}

func TestNestedDependencies(t *testing.T) {
	c := di.New()

	// Register deep dependency chain: MetaSvc -> Service -> (Logger, Greeter)
	di.Register[Logger](c, func() Logger { return &TestLogger{} })
	di.Register[Greeter](c, func() Greeter { return &SimpleGreeter{} })
	di.Register[Service](c, func(l Logger, g Greeter) Service {
		return &DefaultService{logger: l, greeter: g}
	})
	di.Register[NestedMetaService](c, func(s Service) NestedMetaService {
		return &nestedMetaServiceImpl{svc: s}
	})

	meta, err := di.Resolve[NestedMetaService](c)
	if err != nil {
		t.Fatalf("failed to resolve nested dependencies: %v", err)
	}

	if meta == nil {
		t.Error("expected NestedMetaService instance")
	}
}

type NestedMetaService interface {
	Meta() string
}

type nestedMetaServiceImpl struct {
	svc Service
}

func (m *nestedMetaServiceImpl) Meta() string {
	return m.svc.DoWork()
}

func TestResolveInScopeError(t *testing.T) {
	c := di.New()

	scope := c.CreateScope("test")

	// Try to resolve unregistered type in scope
	_, err := di.ResolveInScope[Greeter](c, scope)
	if err == nil {
		t.Fatal("expected error for unregistered type")
	}
}

// =============================================================================
// Helpers
// =============================================================================

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
