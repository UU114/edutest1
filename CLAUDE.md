# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Environment

- **Platform**: Windows
- **Language**: Go 1.21+
- **Framework**: go-zero - a web and RPC framework with built-in engineering practices

## Key Commands

### Building and Testing
```bash
# Install goctl tool
go install github.com/zeromicro/go-zero/tools/goctl@latest

# Build goctl tool
cd tools/goctl
go build -ldflags="-s -w" goctl.go

# Run tests with race detection
go test -race ./...

# Run tests for specific package
go test -race ./core/...

# Format code
gofmt -w .
```

### Code Generation
```bash
# Generate API service from .api file
goctl api go -api greet.api -dir greet

# Generate new API service
goctl api new greet

# Generate RPC service
goctl rpc new greet

# Generate model code from SQL
goctl model mysql ddl -src *.sql -dir ./model
```

## Architecture Overview

go-zero is a microservices framework with three main components:

### 1. Core Framework (`core/`)
Contains essential utilities and modules:
- **`logx/`**: Structured logging with multiple modes (console, file, volume)
- **`breaker/`**: Circuit breaker implementation for resilience
- **`discov/`**: Service discovery
- **`stores/`**: Data stores (SQL, MongoDB, Redis)
- **`syncx/`**: Concurrency utilities
- **`fx/`**: Dependency injection
- **`conf/`**: Configuration management

### 2. Service Frameworks
- **`rest/`**: HTTP/REST server implementation
- **`zrpc/`**: RPC client and server based on gRPC
- **`gateway/`**: API gateway for routing and load balancing

### 3. Code Generation Tool (`tools/goctl/`)
Powerful CLI tool for code generation:
- **API generation**: Creates Go servers from .api files
- **Model generation**: Generates data access layer from SQL schemas
- **RPC generation**: Creates RPC services from proto files
- **Multi-language support**: Can generate client code for various platforms

## Project Structure

Generated service structure:
```
service/
├── etc/
│   └── service.yaml        # Configuration file
├── service.go             # Main entry point
└── internal/
    ├── config/
    │   └── config.go       # Configuration definition
    ├── handler/
    │   ├── routes.go      # Route definitions
    │   └── *.go           # Request handlers
    ├── logic/
    │   └── *.go           # Business logic
    ├── svc/
    │   └── servicecontext.go # Service context with dependencies
    └── types/
        └── types.go       # Request/response types
```

## Key Features

### Resilience Design
- **Circuit breaker**: Automatic fault isolation
- **Rate limiting**: Built-in concurrency control
- **Load shedding**: Adaptive system protection
- **Timeout cascading**: Prevents resource exhaustion

### Development Workflow
1. Define API in `.api` file with simple syntax
2. Generate boilerplate code using `goctl`
3. Implement business logic in `logic/` package
4. Configure services in YAML files
5. Run with automatic service discovery

### Configuration
Services use YAML configuration with support for:
- Database connections (MySQL, PostgreSQL, MongoDB)
- Redis caching
- Service discovery (etcd)
- Logging settings
- Circuit breaker parameters
- Rate limiting rules

## Testing Strategy

- Unit tests: `go test ./...`
- Integration tests: Use provided test utilities
- Race detection: Always include `-race` flag
- Mock support: Use `go.uber.org/mock` for generating mocks

## Code Style Guidelines

- Follow standard Go formatting (`gofmt`)
- Use structured logging via `core/logx`
- Implement proper error handling with `core/errorx`
- Write comprehensive tests with race detection
- Use dependency injection via `core/fx`

## Important Notes

- go-zero emphasizes resilience and high availability
- The framework handles many cross-cutting concerns automatically
- Code generation is central to the development workflow
- Services are designed to be deployed in containerized environments
- Windows development is fully supported

说中文
