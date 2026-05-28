# Stage 1: Build
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/bin/myapp cmd/main.go

# Stage 2: Runtime
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/bin/myapp .

# Copy migrations
COPY migrations/ ./migrations/

# Expose ports
EXPOSE 8080 50051

# Health check
HEALTHCHECK --interval=10s --timeout=5s --retries=5 \
    CMD ["/app/myapp", "--health-check"] || exit 1

CMD ["/app/myapp"]
