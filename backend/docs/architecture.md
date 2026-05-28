# Architecture

## Overview

This project follows **Clean Architecture** principles with a clear separation of concerns:

```
Transport Layer (HTTP + gRPC)
    ↓
Service Layer (Business Logic)
    ↓
Repository Layer (Data Access)
    ↓
Domain Layer (Entities & Interfaces)
```

## Layer Descriptions

### Domain Layer (`internal/domain/`)

- **Entities**: Core business objects (User, etc.)
- **Repository Interfaces**: Abstract data access contracts
- **Service Interfaces**: Abstract business logic contracts
- **Zero Dependencies**: No imports from other layers

### Service Layer (`internal/service/`)

- **Business Logic**: Implementation of use cases
- **Orchestration**: Coordinates between repositories
- **Validation**: Input and state validation
- **Dependencies**: Depends on domain layer only

### Repository Layer (`internal/repo/`)

- **Data Access**: Database operations via pgx
- **Concrete Implementations**: PostgreSQL, Memory (for tests)
- **Transaction Management**: Handles database transactions
- **Error Mapping**: Maps DB errors to domain errors

### Transport Layer (`internal/transport/`)

#### HTTP (`internal/transport/http/`)

- **Handlers**: Request processing and response formatting
- **Middleware**: Cross-cutting concerns (auth, logging, recovery, CORS)
- **Router**: Route registration and grouping
- **DTOs**: Request/response data structures

#### gRPC (`internal/transport/grpc/`)

- **Handlers**: RPC method implementations
- **Interceptors**: Cross-cutting concerns (auth, logging, recovery)
- **Server**: gRPC server setup and lifecycle

### Package Layer (`pkg/`)

- **JWT**: Token generation and validation
- **Logger**: Structured logging wrapper
- **Database**: Database connection and utilities
- **HTTP Utilities**: Response helpers

## Dependency Flow

- Transport layers depend on Service layer
- Service layer depends on Repository layer and Domain interfaces
- Repository layer depends on Domain interfaces
- Domain layer has zero dependencies (except stdlib)

This ensures:
- ✅ Easy testing (mock repositories, services)
- ✅ Framework independence
- ✅ Reusability
- ✅ Clear contract definitions

## Error Handling

Application errors are defined in `apperror/` package:

- Custom error types with codes
- Automatic HTTP status mapping (404, 409, 400, 500, etc.)
- Cause chain for debugging

## Authentication

- **Access Token**: Short-lived JWT in Authorization header
- **Refresh Token**: Long-lived token in httpOnly cookie
- **Middleware/Interceptors**: Validate tokens, extract claims, check authorization

## Database

- **Migrations**: Version control for schema changes
- **Connection Pool**: pgx with configurable max connections
- **Transactions**: Interface-based for testability
- **Repositories**: Implement domain repository interfaces

## Deployment

- **Docker**: Multi-stage build for minimal image size
- **Docker Compose**: Local development with PostgreSQL
- **Graceful Shutdown**: Clean shutdown of all components
- **Signals**: SIGINT, SIGTERM handling

## Configuration

- **Environment Variables**: 12-factor app principles
- **Config Validation**: Checked at startup
- **Hot-loadable**: Some configs via environment

## Testing Strategy

- **Unit Tests**: Mock repositories (memory implementation)
- **Integration Tests**: Real database with testcontainers
- **Test Helpers**: Factory functions for test data
