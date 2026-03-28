# AGENTS.md - Docker GC Project Guidelines

This file provides guidelines for AI agents working on the Docker GC project, a Golang implementation of Docker garbage collection with cron scheduling.

## Project Overview

- **Language**: Go 1.21+
- **Purpose**: Docker garbage collection with cron scheduling
- **Architecture**: Single binary application with Docker CLI integration
- **Deployment**: Docker container with multi-stage build

## Build Commands

### Development Build
```bash
# Build the Go binary
go build -o docker-gc-cron ./cmd/docker-gc-cron

# Run locally
./docker-gc-cron
```

### Docker Build
```bash
# Build Docker image
docker build -t ghcr.io/morya/docker-gc:latest .

# Run with Docker socket
docker run -d \
  -v /var/run/docker.sock:/var/run/docker.sock \
  ghcr.io/morya/docker-gc:latest
```

### Dependency Management
```bash
# Download dependencies
go mod download

# Tidy dependencies
go mod tidy

# Verify dependencies
go mod verify
```

## Testing Commands

**Note**: This project currently has no test suite. When adding tests:

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test
go test -run TestFunctionName ./path/to/package

# Run tests with verbose output
go test -v ./...
```

## Linting and Code Quality

### Go Formatting
```bash
# Format all Go files
go fmt ./...

# Check formatting without applying
gofmt -d .
```

### Go Vet (Static Analysis)
```bash
# Run go vet for static analysis
go vet ./...
```

### Suggested Linting Tools (Not Currently Configured)
- `golangci-lint run` - Comprehensive linter
- `staticcheck ./...` - Advanced static analysis
- `revive ./...` - Fast, configurable linter

## Code Style Guidelines

### Imports Organization
```go
import (
    // Standard library
    "fmt"
    "log"
    "os"
    "os/exec"
    "strings"

    // Third-party packages
    "github.com/robfig/cron/v3"
)
```

### Naming Conventions
- **Packages**: Lowercase, single word (e.g., `main`)
- **Variables**: camelCase (e.g., `containerIDs`, `gracePeriodSeconds`)
- **Constants**: camelCase or UPPER_SNAKE_CASE for environment defaults
- **Types**: PascalCase (e.g., `Config`, `ContainerInfo`)
- **Interfaces**: PascalCase ending with "er" if appropriate (e.g., `Cleaner`)

### Error Handling
- Use `log.Printf` for non-fatal errors with context
- Use `log.Fatalf` only for unrecoverable startup errors
- Return early from functions on error
- Include meaningful error messages with context

Example:
```go
output, err := cmd.Output()
if err != nil {
    log.Printf("Failed to list exited containers: %v", err)
    return
}
```

### Function Structure
- Keep functions focused and small (< 50 lines)
- Use descriptive function names (e.g., `cleanContainers`, `loadConfig`)
- Document public functions with comments
- Group related functions together

### Configuration Pattern
- Use `Config` struct for environment variables
- Provide sensible defaults
- Validate configuration in `loadConfig()` function
- Use boolean flags for feature toggles (e.g., `FORCE_CONTAINER_REMOVAL=1`)

### Logging Guidelines
- Use `log.Println` for informational messages
- Prefix log messages with `[Docker GC]` for important operations
- Include relevant identifiers (container IDs, image IDs)
- Log before and after significant operations

### Docker CLI Integration
- Use `exec.Command` for Docker operations
- Always check command output and errors
- Handle empty results gracefully
- Support dry-run mode for testing

## Project Structure

```
docker-gc/
├── cmd/
│   └── docker-gc-cron/
│       └── main.go          # Main application entry point
├── Dockerfile               # Multi-stage Docker build
├── entrypoint.sh           # Container entrypoint script
├── go.mod                  # Go module definition
├── go.sum                  # Dependency checksums
├── .github/
│   └── workflows/
│       └── docker-build.yml # CI/CD pipeline
└── tasks/
    └── task.md             # Project requirements
```

## Docker Build Guidelines

### Multi-stage Build
- Use `golang:alpine` for build stage
- Use `alpine:latest` for runtime stage
- Set `CGO_ENABLED=0` for static binary
- Copy only necessary files between stages

### Runtime Dependencies
- Include `docker-cli` in runtime image
- Install `bash` and `tzdata` for cron support
- Create necessary directories (`/var/log`)
- Set appropriate file permissions

### Environment Variables
- Document all environment variables in README
- Provide sensible defaults
- Use boolean flags (`0`/`1`) for toggles
- Support configuration via Docker run command

## GitHub Actions Workflow

The project uses GitHub Actions for CI/CD:
- Builds on push to main/master branches
- Builds on tag pushes (`v*`)
- Builds on pull requests (no push)
- Pushes to GitHub Container Registry

**Important**: The workflow uses Node.js 22.x as required by GitHub Actions specifications.

## Development Workflow

1. **Make Changes**: Edit `cmd/docker-gc-cron/main.go` or other files
2. **Test Locally**: `go build -o docker-gc-cron ./cmd/docker-gc-cron && ./docker-gc-cron`
3. **Build Docker**: `docker build -t ghcr.io/morya/docker-gc:latest .`
4. **Test Container**: `docker run -v /var/run/docker.sock:/var/run/docker.sock ghcr.io/morya/docker-gc:latest`
5. **Commit Changes**: Follow conventional commit messages
6. **Create PR**: Changes will be built and tested automatically

## Adding New Features

When adding new features:
1. Add configuration to `Config` struct in `main.go`
2. Add environment variable handling in `loadConfig()`
3. Update `README.md` with new environment variable
4. Update `Dockerfile` with default value if needed
5. Consider backward compatibility
6. Test with various environment configurations

## Common Tasks for Agents

### Adding New Cleanup Functionality
1. Create new function (e.g., `cleanNetworks()`)
2. Add to `runGarbageCollection()` call sequence
3. Add configuration options if needed
4. Update documentation

### Improving Error Handling
1. Add more specific error messages
2. Consider retry logic for transient failures
3. Add metrics/logging for error rates
4. Test error scenarios

### Performance Optimization
1. Consider parallel execution of cleanup tasks
2. Batch Docker operations where possible
3. Add resource usage monitoring
4. Profile memory and CPU usage

## Notes for AI Agents

- This is a relatively simple Go application with minimal dependencies
- The code follows standard Go idioms and patterns
- Focus on clarity and maintainability over premature optimization
- When in doubt, follow the existing code patterns
- Always test Docker builds after making changes
- Consider the containerized deployment environment