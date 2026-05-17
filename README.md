# MyApp

Production-ready Go backend with HTTP + gRPC using Clean Architecture.

## Features

- **Clean Architecture**: transport → service → repo → domain
- **Dual Protocol Support**: HTTP (Gin) + gRPC
- **Authentication**: JWT with access/refresh token pattern
- **Security**: bcrypt password hashing, CORS, rate limiting
- **Database**: PostgreSQL with pgx/v5
- **Logging**: Structured logging with Zap
- **Validation**: Input validation with validator/v10
- **Error Handling**: Custom error types with HTTP status mapping
- **Graceful Shutdown**: Clean shutdown for all servers
- **CI/CD**: GitHub Actions workflow

## Project Structure

```
myapp/
├── cmd/api/                    # Entry point
├── internal/
│   ├── app/                    # DI container & app lifecycle
│   ├── domain/                 # Domain layer (entities, interfaces)
│   ├── service/                # Business logic
│   ├── repo/                   # Data access layer
│   └── transport/
│       ├── http/               # HTTP handlers, middleware, router
│       └── grpc/               # gRPC handlers, interceptors
├── pkg/                        # Reusable packages
│   ├── jwt/                    # JWT utilities
│   ├── logger/                 # Logging wrapper
│   ├── httputil/               # HTTP response helpers
│   └── database/               # Database utilities
├── apperror/                   # Application error definitions
├── migrations/                 # Database migrations
├── Makefile
├── docker-compose.yml
└── go.mod
```

## Getting Started

### Prerequisites

- Go 1.26+
- PostgreSQL 16+
- Docker & Docker Compose (optional)

### Setup

1. Clone the repository:
```bash
git clone <repo-url>
cd myapp
```

2. Copy environment file:
```bash
cp .env.example .env
```

3. Download dependencies:
```bash
go mod download
```

### Running Locally

#### Without Docker:

1. Start PostgreSQL:
```bash
# Make sure PostgreSQL is running on localhost:5432
```

2. Run migrations:
```bash
make migrate-up
```

3. Start the application:
```bash
make run
```

#### With Docker:

1. Start services:
```bash
make docker-up
```

2. Run migrations:
```bash
docker-compose exec app make migrate-up
```

## API Endpoints

### Authentication

- `POST /api/v1/users/register` - Register new user
- `POST /api/v1/users/login` - Login user
- `POST /api/v1/users/refresh` - Refresh access token

### Users

- `GET /api/v1/users` - List users (protected)
- `GET /api/v1/users/:id` - Get user by ID (protected)
- `PATCH /api/v1/users/:id` - Update user (protected)
- `DELETE /api/v1/users/:id` - Delete user (protected)

## Development

### Format Code

```bash
make fmt
```

### Run Linter

```bash
make lint
```

### Run Tests

```bash
make test
```

### Build Binary

```bash
make build
```

## Configuration

All configuration is read from environment variables. See `.env.example` for available options.

## Deployment

### Build Docker Image

```bash
make docker-build
```

### Docker Compose

```bash
docker-compose up -d
```

## Testing

### Unit Tests

```bash
go test ./...
```

### With Coverage

```bash
make test-coverage
```

## Error Handling

Application errors are mapped to HTTP status codes:

- `NOT_FOUND` → 404
- `CONFLICT` → 409
- `INVALID_CREDENTIALS` → 401
- `FORBIDDEN` → 403
- `BAD_REQUEST` → 400
- `INTERNAL_SERVER` → 500

## License

MIT
