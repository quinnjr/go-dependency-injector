package main

import (
	"fmt"
	"time"

	"github.com/quinnjr/go-dependency-injector/di"
)

// =============================================================================
// Domain Interfaces
// =============================================================================

// Logger defines the logging contract.
type Logger interface {
	Log(message string)
	LogError(message string)
}

// Config holds application configuration.
type Config interface {
	DatabaseURL() string
	CacheEnabled() bool
}

// Database represents a database connection.
type Database interface {
	Query(sql string) ([]map[string]any, error)
	Close() error
}

// Cache represents a caching layer.
type Cache interface {
	Get(key string) (any, bool)
	Set(key string, value any, ttl time.Duration)
}

// UserRepository handles user data access.
type UserRepository interface {
	FindByID(id int) (*User, error)
	FindAll() ([]*User, error)
}

// UserService handles user business logic.
type UserService interface {
	GetUser(id int) (*User, error)
	ListUsers() ([]*User, error)
}

// =============================================================================
// Domain Models
// =============================================================================

// User represents a user entity.
type User struct {
	ID    int
	Name  string
	Email string
}

// =============================================================================
// Implementations
// =============================================================================

// ConsoleLogger logs to stdout.
type ConsoleLogger struct {
	prefix string
}

// NewConsoleLogger creates a new console logger.
func NewConsoleLogger() Logger {
	return &ConsoleLogger{prefix: "[APP]"}
}

// Log outputs an info message to stdout.
func (l *ConsoleLogger) Log(message string) {
	fmt.Printf("%s %s INFO: %s\n", l.prefix, time.Now().Format("15:04:05"), message)
}

// LogError outputs an error message to stdout.
func (l *ConsoleLogger) LogError(message string) {
	fmt.Printf("%s %s ERROR: %s\n", l.prefix, time.Now().Format("15:04:05"), message)
}

// AppConfig holds app configuration.
type AppConfig struct {
	dbURL        string
	cacheEnabled bool
}

// NewAppConfig creates a new application configuration.
func NewAppConfig() Config {
	return &AppConfig{
		dbURL:        "postgres://localhost:5432/myapp",
		cacheEnabled: true,
	}
}

// DatabaseURL returns the database connection URL.
func (c *AppConfig) DatabaseURL() string { return c.dbURL }

// CacheEnabled returns whether caching is enabled.
func (c *AppConfig) CacheEnabled() bool { return c.cacheEnabled }

// PostgresDatabase simulates a postgres connection.
type PostgresDatabase struct {
	logger Logger
	config Config
}

// NewPostgresDatabase creates a new database connection.
func NewPostgresDatabase(logger Logger, config Config) (Database, error) {
	logger.Log(fmt.Sprintf("Connecting to database: %s", config.DatabaseURL()))
	return &PostgresDatabase{logger: logger, config: config}, nil
}

// Query executes a SQL query and returns results.
func (db *PostgresDatabase) Query(sql string) ([]map[string]any, error) {
	db.logger.Log(fmt.Sprintf("Executing query: %s", sql))
	// Simulated query result
	return []map[string]any{
		{"id": 1, "name": "Alice", "email": "alice@example.com"},
		{"id": 2, "name": "Bob", "email": "bob@example.com"},
	}, nil
}

// Close closes the database connection.
func (db *PostgresDatabase) Close() error {
	db.logger.Log("Closing database connection")
	return nil
}

// InMemoryCache is a simple in-memory cache.
type InMemoryCache struct {
	logger Logger
	data   map[string]any
}

// NewInMemoryCache creates a new in-memory cache.
func NewInMemoryCache(logger Logger) Cache {
	logger.Log("Initializing in-memory cache")
	return &InMemoryCache{
		logger: logger,
		data:   make(map[string]any),
	}
}

// Get retrieves a value from the cache.
func (c *InMemoryCache) Get(key string) (any, bool) {
	val, ok := c.data[key]
	return val, ok
}

// Set stores a value in the cache with the given TTL.
func (c *InMemoryCache) Set(key string, value any, ttl time.Duration) {
	c.data[key] = value
}

// DefaultUserRepository implements UserRepository.
type DefaultUserRepository struct {
	db     Database
	cache  Cache
	logger Logger
}

// NewUserRepository creates a new user repository.
func NewUserRepository(db Database, cache Cache, logger Logger) UserRepository {
	logger.Log("Creating user repository")
	return &DefaultUserRepository{db: db, cache: cache, logger: logger}
}

// FindByID finds a user by their ID.
func (r *DefaultUserRepository) FindByID(id int) (*User, error) {
	cacheKey := fmt.Sprintf("user:%d", id)

	// Check cache first
	if cached, ok := r.cache.Get(cacheKey); ok {
		r.logger.Log(fmt.Sprintf("Cache hit for user %d", id))
		return cached.(*User), nil
	}

	r.logger.Log(fmt.Sprintf("Cache miss for user %d, querying database", id))
	results, err := r.db.Query(fmt.Sprintf("SELECT * FROM users WHERE id = %d", id))
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("user %d not found", id)
	}

	user := &User{
		ID:    results[0]["id"].(int),
		Name:  results[0]["name"].(string),
		Email: results[0]["email"].(string),
	}

	r.cache.Set(cacheKey, user, 5*time.Minute)
	return user, nil
}

// FindAll retrieves all users from the database.
func (r *DefaultUserRepository) FindAll() ([]*User, error) {
	results, err := r.db.Query("SELECT * FROM users")
	if err != nil {
		return nil, err
	}

	users := make([]*User, len(results))
	for i, row := range results {
		users[i] = &User{
			ID:    row["id"].(int),
			Name:  row["name"].(string),
			Email: row["email"].(string),
		}
	}
	return users, nil
}

// DefaultUserService implements UserService.
type DefaultUserService struct {
	repo   UserRepository
	logger Logger
}

// NewUserService creates a new user service.
func NewUserService(repo UserRepository, logger Logger) UserService {
	logger.Log("Creating user service")
	return &DefaultUserService{repo: repo, logger: logger}
}

// GetUser retrieves a user by their ID.
func (s *DefaultUserService) GetUser(id int) (*User, error) {
	s.logger.Log(fmt.Sprintf("Getting user %d", id))
	return s.repo.FindByID(id)
}

// ListUsers retrieves all users.
func (s *DefaultUserService) ListUsers() ([]*User, error) {
	s.logger.Log("Listing all users")
	return s.repo.FindAll()
}

// =============================================================================
// Application Bootstrap
// =============================================================================

func main() {
	fmt.Println("╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║     Go Dependency Injection Demo                             ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
	fmt.Println()

	// Create the DI container
	container := di.New()

	// Register dependencies with appropriate lifetimes
	registerDependencies(container)

	fmt.Println("\n─── Resolving UserService (will auto-resolve all dependencies) ───")
	fmt.Println()

	// Resolve the top-level service - all dependencies are resolved automatically!
	userService, err := di.Resolve[UserService](container)
	if err != nil {
		fmt.Printf("Failed to resolve UserService: %v\n", err)
		return
	}

	fmt.Println("\n─── Using the resolved service ───")
	fmt.Println()

	// Use the service
	users, err := userService.ListUsers()
	if err != nil {
		fmt.Printf("Failed to list users: %v\n", err)
		return
	}

	fmt.Println("\n─── Results ───")
	fmt.Println()
	for _, user := range users {
		fmt.Printf("  → User: %s (%s)\n", user.Name, user.Email)
	}

	// Demonstrate singleton behavior
	fmt.Println("\n─── Demonstrating Singleton Behavior ───")
	fmt.Println()

	logger1 := di.MustResolve[Logger](container)
	logger2 := di.MustResolve[Logger](container)

	logger1.Log("This is logger1")
	logger2.Log("This is logger2 (same instance as logger1)")

	// Demonstrate scoped resolution
	fmt.Println("\n─── Demonstrating Scoped Resolution ───")
	fmt.Println()

	demonstrateScopedResolution(container)

	fmt.Println("\n─── Demo Complete ───")
}

func registerDependencies(c *di.Container) {
	fmt.Println("─── Registering Dependencies ───")
	fmt.Println()

	// Config - Singleton (one config for entire app)
	if err := di.Register[Config](c, NewAppConfig, di.AsSingleton()); err != nil {
		panic(err)
	}
	fmt.Println("  ✓ Config registered as Singleton")

	// Logger - Singleton (share across app)
	if err := di.Register[Logger](c, NewConsoleLogger, di.AsSingleton()); err != nil {
		panic(err)
	}
	fmt.Println("  ✓ Logger registered as Singleton")

	// Database - Singleton (one connection pool)
	if err := di.Register[Database](c, NewPostgresDatabase, di.AsSingleton()); err != nil {
		panic(err)
	}
	fmt.Println("  ✓ Database registered as Singleton")

	// Cache - Singleton (shared cache)
	if err := di.Register[Cache](c, NewInMemoryCache, di.AsSingleton()); err != nil {
		panic(err)
	}
	fmt.Println("  ✓ Cache registered as Singleton")

	// UserRepository - Transient (new instance per resolution)
	if err := di.Register[UserRepository](c, NewUserRepository, di.AsTransient()); err != nil {
		panic(err)
	}
	fmt.Println("  ✓ UserRepository registered as Transient")

	// UserService - Transient
	if err := di.Register[UserService](c, NewUserService, di.AsTransient()); err != nil {
		panic(err)
	}
	fmt.Println("  ✓ UserService registered as Transient")
}

// RequestContext simulates a request-scoped dependency
type RequestContext struct {
	RequestID string
	UserID    int
	StartTime time.Time
}

func demonstrateScopedResolution(c *di.Container) {
	// Register a scoped dependency
	err := di.Register[*RequestContext](c, func() *RequestContext {
		return &RequestContext{
			RequestID: fmt.Sprintf("req-%d", time.Now().UnixNano()),
			StartTime: time.Now(),
		}
	}, di.AsScoped())
	if err != nil {
		fmt.Printf("Failed to register RequestContext: %v\n", err)
		return
	}

	// Create a scope for "request 1"
	scope1 := c.CreateScope("request-1")
	ctx1a, err := di.ResolveInScope[*RequestContext](c, scope1)
	if err != nil {
		fmt.Printf("Failed to resolve context: %v\n", err)
		return
	}
	ctx1b, err := di.ResolveInScope[*RequestContext](c, scope1)
	if err != nil {
		fmt.Printf("Failed to resolve context: %v\n", err)
		return
	}

	fmt.Printf("  Scope 'request-1' context A: %s\n", ctx1a.RequestID)
	fmt.Printf("  Scope 'request-1' context B: %s\n", ctx1b.RequestID)
	fmt.Printf("  Same instance? %v\n", ctx1a.RequestID == ctx1b.RequestID)

	// Create a different scope for "request 2"
	scope2 := c.CreateScope("request-2")
	ctx2, err := di.ResolveInScope[*RequestContext](c, scope2)
	if err != nil {
		fmt.Printf("Failed to resolve context: %v\n", err)
		return
	}

	fmt.Printf("\n  Scope 'request-2' context: %s\n", ctx2.RequestID)
	fmt.Printf("  Different from request-1? %v\n", ctx1a.RequestID != ctx2.RequestID)
}
