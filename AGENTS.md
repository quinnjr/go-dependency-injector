<!-- converted from Cursor rules -->

## Cursor rule: `.cursor/rules/auto-commit.mdc`

# Auto-Commit After Work

Always commit changes to git after completing work.

## Rule

When the agent finishes a task that modifies files, it must:

1. Stage all relevant changes with `git add`
2. Commit with a descriptive message following conventional commits format
3. Push to the remote repository (if configured)

## Commit Message Format

```
<type>(<scope>): <short description>
```

**Types:**
- `feat` — New feature
- `fix` — Bug fix
- `docs` — Documentation changes
- `style` — Formatting, whitespace (no code change)
- `refactor` — Code restructuring (no behavior change)
- `test` — Adding or updating tests
- `chore` — Maintenance, dependencies, configs

**Examples:**
```
feat(auth): add JWT token validation
fix(api): resolve null pointer in user lookup
docs(readme): update installation instructions
chore(deps): upgrade to Go 1.22
test(container): add circular dependency tests
```

## Workflow

```bash
# After completing work
git add <files>
git commit -m "<type>(<scope>): <description>"
git push
```

## Exceptions

Do NOT auto-commit when:
- The user explicitly asks not to commit
- Changes are incomplete or broken
- Tests are failing
- The user is reviewing changes before committing


## Cursor rule: `.cursor/rules/git-flow.mdc`

# Git-Flow Branching Pattern

Follow the git-flow branching model for all version control operations.

## Branch Structure

### Primary Branches (Long-lived)

| Branch | Purpose |
|--------|---------|
| `main` | Production-ready code. Only receives merges from `release/*` and `hotfix/*` |
| `develop` | Integration branch. Features merge here first |

### Supporting Branches (Short-lived)

| Branch | Pattern | Created From | Merges Into |
|--------|---------|--------------|-------------|
| Feature | `feature/<name>` | `develop` | `develop` |
| Release | `release/<version>` | `develop` | `main` and `develop` |
| Hotfix | `hotfix/<name>` | `main` | `main` and `develop` |
| Bugfix | `bugfix/<name>` | `develop` | `develop` |

## Branch Naming Conventions

```
feature/add-user-authentication
feature/ISSUE-123-implement-caching
release/v1.2.0
hotfix/fix-critical-login-bug
bugfix/resolve-memory-leak
```

- Use lowercase with hyphens (kebab-case)
- Be descriptive but concise
- Include issue/ticket numbers when applicable

## Workflow Rules

### Starting New Work

1. **Features**: Always branch from `develop`
   ```bash
   git checkout develop
   git pull origin develop
   git checkout -b feature/<name>
   ```

2. **Hotfixes**: Always branch from `main`
   ```bash
   git checkout main
   git pull origin main
   git checkout -b hotfix/<name>
   ```

### Completing Work

1. **Features**: Create PR to `develop`, squash merge preferred
2. **Releases**: Merge to both `main` (with tag) and back to `develop`
3. **Hotfixes**: Merge to both `main` (with tag) and back to `develop`

### Version Tagging

- Tag releases on `main` using semantic versioning: `v1.0.0`, `v1.1.0`, `v2.0.0`
- Hotfixes increment patch version: `v1.0.1`

## Commit Message Format

```
<type>(<scope>): <subject>

<body>

<footer>
```

**Types**: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`

**Examples**:
```
feat(auth): add JWT token refresh mechanism
fix(api): resolve race condition in connection pool
docs(readme): update installation instructions
```

## Commands Reference

```bash
# Start feature
git flow feature start <name>
# or manually:
git checkout -b feature/<name> develop

# Finish feature
git flow feature finish <name>
# or manually:
git checkout develop && git merge --no-ff feature/<name>

# Start release
git flow release start <version>

# Finish release (tags and merges to main + develop)
git flow release finish <version>

# Start hotfix
git flow hotfix start <name>

# Finish hotfix
git flow hotfix finish <name>
```

## Pull Request Guidelines

- Feature PRs target `develop`
- Release PRs target `main`
- Hotfix PRs target `main` (then cherry-pick or merge to `develop`)
- Require at least one approval before merging
- Ensure CI passes before merging
- Use squash merges for features to keep history clean
- Use merge commits for releases/hotfixes to preserve history


## Cursor rule: `.cursor/rules/go-standards.mdc`

# Modern Go Programming Standards

Follow these standards for all Go code in this project.

## Go Version

Use Go 1.22+ features including:
- Generics (`[T any]`, `[T comparable]`)
- Range over integers (`for i := range 10`)
- Enhanced type inference
- `slices` and `maps` packages from stdlib

## Code Organization

### Package Structure

```
project/
├── cmd/              # Executables (main packages)
├── internal/         # Private application code
├── pkg/              # Public library code
├── api/              # API definitions (proto, OpenAPI)
├── web/              # Web assets
└── scripts/          # Build/dev scripts
```

### File Naming

- Use lowercase with underscores: `user_service.go`, `http_handler.go`
- Test files: `*_test.go`
- Platform-specific: `*_linux.go`, `*_windows.go`
- One primary type per file when sensible

## Naming Conventions

### Variables & Functions

```go
// Good: camelCase for local, PascalCase for exported
userID := 42
func GetUserByID(id int) (*User, error)

// Bad: underscores, all caps (except constants)
user_id := 42
func get_user_by_id(id int)
```

### Interfaces

- Name by behavior, not implementation
- Use `-er` suffix for single-method interfaces

```go
// Good
type Reader interface { Read(p []byte) (n int, err error) }
type UserRepository interface { ... }

// Bad
type IUserRepository interface { ... }  // No "I" prefix
type UserRepositoryInterface interface { ... }
```

### Acronyms

Keep acronyms uppercase: `HTTPServer`, `userID`, `parseJSON`

## Error Handling

### Return Errors, Don't Panic

```go
// Good
func Parse(s string) (Result, error) {
    if s == "" {
        return Result{}, errors.New("empty string")
    }
    // ...
}

// Bad: panics for recoverable errors
func Parse(s string) Result {
    if s == "" {
        panic("empty string")
    }
}
```

### Wrap Errors with Context

```go
// Good: adds context
if err != nil {
    return fmt.Errorf("failed to parse config: %w", err)
}

// Bad: loses context
if err != nil {
    return err
}
```

### Use Sentinel Errors or Custom Types

```go
var ErrNotFound = errors.New("not found")

type ValidationError struct {
    Field   string
    Message string
}

func (e ValidationError) Error() string {
    return fmt.Sprintf("%s: %s", e.Field, e.Message)
}
```

### Check Errors with `errors.Is` and `errors.As`

```go
if errors.Is(err, ErrNotFound) { ... }

var valErr *ValidationError
if errors.As(err, &valErr) { ... }
```

## Generics

### Use Generics for Type-Safe Collections and Utilities

```go
func Map[T, U any](items []T, fn func(T) U) []U {
    result := make([]U, len(items))
    for i, item := range items {
        result[i] = fn(item)
    }
    return result
}
```

### Prefer Constraints Over `any`

```go
// Good: constrained
func Max[T cmp.Ordered](a, b T) T

// Less good: unconstrained when constraint is possible
func Max[T any](a, b T) T
```

## Concurrency

### Prefer Channels for Communication

```go
// Good: channel for signaling
done := make(chan struct{})
go func() {
    defer close(done)
    // work...
}()
<-done
```

### Use `sync` Package Appropriately

```go
var mu sync.RWMutex
var cache map[string]string

func Get(key string) string {
    mu.RLock()
    defer mu.RUnlock()
    return cache[key]
}
```

### Use `context.Context` for Cancellation

```go
func FetchData(ctx context.Context, url string) ([]byte, error) {
    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return nil, err
    }
    // ...
}
```

## Testing

### Table-Driven Tests

```go
func TestAdd(t *testing.T) {
    tests := []struct {
        name     string
        a, b     int
        expected int
    }{
        {"positive", 2, 3, 5},
        {"negative", -1, -1, -2},
        {"zero", 0, 0, 0},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := Add(tt.a, tt.b)
            if got != tt.expected {
                t.Errorf("Add(%d, %d) = %d; want %d", tt.a, tt.b, got, tt.expected)
            }
        })
    }
}
```

### Use `t.Parallel()` When Safe

```go
func TestSomething(t *testing.T) {
    t.Parallel()
    // ...
}
```

### Prefer stdlib `testing` Over Assertion Libraries

Use `t.Errorf`, `t.Fatalf`, `t.Helper()` from stdlib.

## Documentation

### Document All Exported Identifiers

```go
// UserService handles user-related business logic.
// It provides methods for creating, retrieving, and managing users.
type UserService struct { ... }

// GetByID retrieves a user by their unique identifier.
// Returns ErrNotFound if the user does not exist.
func (s *UserService) GetByID(id int) (*User, error) { ... }
```

### Use Examples in Tests

```go
func ExampleUserService_GetByID() {
    svc := NewUserService(db)
    user, _ := svc.GetByID(42)
    fmt.Println(user.Name)
    // Output: Alice
}
```

## Performance

### Preallocate Slices When Size is Known

```go
// Good
users := make([]User, 0, len(ids))
for _, id := range ids {
    users = append(users, fetchUser(id))
}

// Bad: repeated reallocations
var users []User
for _, id := range ids {
    users = append(users, fetchUser(id))
}
```

### Use `strings.Builder` for String Concatenation

```go
var b strings.Builder
for _, s := range items {
    b.WriteString(s)
}
result := b.String()
```

### Avoid Premature Optimization

Profile first with `pprof`, optimize second.

## Linting & Formatting

- Run `gofmt` or `goimports` on all code
- Use `golangci-lint` with project config
- Zero tolerance for linter warnings in CI

```bash
gofmt -w .
golangci-lint run
```

## Dependencies

- Minimize external dependencies
- Prefer stdlib when reasonable
- Vet dependencies for maintenance and security
- Use `go mod tidy` regularly

