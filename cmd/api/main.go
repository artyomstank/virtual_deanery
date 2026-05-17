// cmd/api/main.go
package main

import (
	"context"
	"log"

	"github.com/artyomstank/virtual_deanery/internal/app"
)

func main() {
	// TODO: Load configuration from environment variables
	cfg := &app.Config{
		HTTPPort:           8080,
		GRPCPort:           50051,
		DatabaseURL:        "postgres://user:password@localhost:5432/myapp?sslmode=disable",
		LogLevel:           "debug",
		JWTSecretKey:       "your-secret-key",
		JWTAlgo:            "HS256",
		JWTAccessTTL:       15 * 60,       // 15 minutes in seconds
		JWTRefreshTTL:      7 * 24 * 3600, // 7 days in seconds
		BCryptCost:         12,
		CORSAllowedOrigins: []string{"http://localhost:3000", "http://localhost:8080"},
	}

	// TODO: Initialize application
	application, err := app.New(cfg)
	if err != nil {
		log.Fatalf("failed to initialize app: %v", err)
	}

	// TODO: Run application with graceful shutdown
	if err := application.Run(context.Background()); err != nil {
		log.Fatalf("app run error: %v", err)
	}
}
