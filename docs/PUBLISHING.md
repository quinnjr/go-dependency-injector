# Publishing Guide for Go Dependency Injector

This guide explains how to publish and maintain the go-dependency-injector package in the Go ecosystem.

## How Go Package Discovery Works

Unlike npm, PyPI, or Maven, Go doesn't use a centralized package registry. Instead:

1. **Packages are hosted on version control systems** (GitHub, GitLab, etc.)
2. **Users import directly from the repository URL**
3. **pkg.go.dev automatically indexes public Go modules**

## Publishing Checklist

### 1. Ensure Repository is Ready

- [x] Module path matches GitHub URL in `go.mod`
- [x] Repository is public
- [x] LICENSE file exists
- [x] README.md is comprehensive
- [x] Tests pass: `go test ./...`
- [x] No lint errors: `golangci-lint run`

### 2. Create a Version Tag

Go modules use semantic versioning with a `v` prefix:

```bash
# For initial release
git tag v1.0.0
git push origin v1.0.0

# For patches (bug fixes)
git tag v1.0.1
git push origin v1.0.1

# For minor versions (new features, backward compatible)
git tag v1.1.0
git push origin v1.1.0

# For major versions (breaking changes)
git tag v2.0.0
git push origin v2.0.0
```

### 3. Automatic Indexing

Once you push a tag, pkg.go.dev will automatically index your package within minutes. The process:

1. Someone (or CI) requests the module via `go get`
2. Go's module proxy (proxy.golang.org) fetches it
3. pkg.go.dev indexes and publishes the documentation

You can manually trigger indexing by visiting:
```
https://pkg.go.dev/github.com/quinnjr/go-dependency-injector@v1.0.0
```

Or via curl:
```bash
curl "https://proxy.golang.org/github.com/quinnjr/go-dependency-injector/@v/v1.0.0.info"
```

### 4. Verify Publication

Check that your package appears at:
- https://pkg.go.dev/github.com/quinnjr/go-dependency-injector
- https://pkg.go.dev/github.com/quinnjr/go-dependency-injector@v1.0.0

## Version Management

### Semantic Versioning Rules

- **MAJOR** (v2.0.0): Breaking changes, incompatible API changes
- **MINOR** (v1.1.0): New features, backward compatible
- **PATCH** (v1.0.1): Bug fixes, backward compatible

### Pre-release Versions

For beta/alpha releases:

```bash
git tag v1.0.0-beta.1
git tag v1.0.0-rc.1
git push origin --tags
```

### Major Version Upgrades

For v2+ modules, update `go.mod`:

```go
module github.com/quinnjr/go-dependency-injector/v2

go 1.22
```

## Continuous Integration

The repository includes GitHub Actions workflows:

### CI Workflow (`.github/workflows/ci.yml`)
- Runs on every push and PR
- Tests on multiple Go versions
- Generates coverage reports
- Runs linter
- Executes benchmarks

### Release Workflow (`.github/workflows/release.yml`)
- Triggers on tag push
- Runs full test suite
- Creates GitHub release
- Triggers pkg.go.dev indexing

## Making Your Package Discoverable

### 1. Good Documentation

- Clear README with examples
- Package documentation in `doc.go`
- Exported functions have comments
- Examples in `example_test.go`

### 2. Add Badges to README

```markdown
[![Go Reference](https://pkg.go.dev/badge/github.com/quinnjr/go-dependency-injector.svg)](https://pkg.go.dev/github.com/quinnjr/go-dependency-injector)
[![Go Report Card](https://goreportcard.com/badge/github.com/quinnjr/go-dependency-injector)](https://goreportcard.com/report/github.com/quinnjr/go-dependency-injector)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
```

### 3. Submit to Go Package Directories

- [Awesome Go](https://github.com/avelino/awesome-go) - Submit a PR
- [Go Report Card](https://goreportcard.com/) - Automatic after first index
- Social media and blog posts

### 4. SEO and Keywords

Add topics to GitHub repository:
- `go`
- `golang`
- `dependency-injection`
- `di-container`
- `generics`
- `type-safe`

## Release Process

1. **Update CHANGELOG.md** with changes
2. **Update version references** in documentation if needed
3. **Run full test suite**: `go test -race -coverprofile=coverage.txt ./...`
4. **Run linter**: `golangci-lint run`
5. **Commit changes**: `git commit -am "Prepare for v1.0.0"`
6. **Create tag**: `git tag v1.0.0`
7. **Push everything**: `git push && git push --tags`
8. **Verify on pkg.go.dev**: Check documentation appears correctly
9. **Create GitHub release**: Add release notes via GitHub UI

## Troubleshooting

### Package Not Appearing on pkg.go.dev

- Ensure module path in `go.mod` matches repository URL
- Verify tag format is correct (`v1.0.0`, not `1.0.0`)
- Check that repository is public
- Wait up to 15 minutes for indexing
- Try manual trigger via curl command above

### Documentation Not Rendering Correctly

- Ensure all exported types/functions have comments
- Comments should start with the name: `// Container represents...`
- Use standard Go doc format
- Check for valid Go code in examples

### Import Issues

- Module path must exactly match `go.mod`
- Case sensitivity matters
- Use full import path in all files

## Best Practices

1. **Never force-push tags** - Tags are immutable in Go's module system
2. **Test before tagging** - Once published, versions are permanent
3. **Use conventional commits** - Helps with changelog generation
4. **Maintain backward compatibility** - Only break on major versions
5. **Write migration guides** - For major version upgrades
6. **Keep dependencies minimal** - This library has zero dependencies

## Support Channels

- GitHub Issues: Bug reports and feature requests
- GitHub Discussions: Questions and community support
- Pull Requests: Code contributions

## Resources

- [Go Modules Reference](https://go.dev/ref/mod)
- [Publishing Go Modules](https://go.dev/blog/publishing-go-modules)
- [pkg.go.dev About](https://pkg.go.dev/about)
- [Semantic Versioning](https://semver.org/)



