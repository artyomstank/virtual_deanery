// internal/app/app.go
package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	grpc_server "github.com/artyomstank/virtual_deanery/internal/transport/grpc"
	http_server "github.com/artyomstank/virtual_deanery/internal/transport/http"
	"github.com/artyomstank/virtual_deanery/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

// App represents application instance with all dependencies.
type App struct {
	httpServer *http_server.Server
	grpcServer *grpc_server.Server
	dbPool     *pgxpool.Pool
	logger     logger.Logger
}

// Config holds application configuration.
type Config struct {
	HTTPPort           int
	GRPCPort           int
	DatabaseURL        string
	LogLevel           string
	JWTSecretKey       string
	JWTAlgo            string
	JWTAccessTTL       time.Duration
	JWTRefreshTTL      time.Duration
	BCryptCost         int
	CORSAllowedOrigins []string
}

// New creates and initializes application with all dependencies.
func New(cfg *Config) (*App, error) {
	// TODO: Initialize logger
	// log, err := logger.NewZapLogger(cfg.LogLevel)

	// TODO: Initialize database connection pool
	// pool, err := postgres.NewPool(context.Background(), cfg.DatabaseURL)

	// TODO: Initialize JWT client
	// jwtClient := jwt.NewJWTClient(
	//     cfg.JWTSecretKey,
	//     cfg.JWTAlgo,
	//     cfg.JWTAccessTTL,
	//     cfg.JWTRefreshTTL,
	// )

	// TODO: Initialize repositories
	// var userRepo repository.UserRepository
	// userRepo = postgres.NewUserRepository(pool, log)

	// TODO: Initialize services
	// userService := service.NewUserService(userRepo, jwtClient, log, cfg.BCryptCost)

	// TODO: Initialize HTTP handlers
	// userHandler := handler.NewUserHandler(userService, log)

	// TODO: Initialize HTTP router
	// httpRouter := router.NewRouter(userHandler, jwtClient, log, cfg.CORSAllowedOrigins)

	// TODO: Initialize gRPC handlers
	// grpcUserHandler := grpc_handler.NewUserServiceServer(userService, log)

	// TODO: Initialize servers
	// httpSrv := http_server.NewServer(httpRouter.Setup(), cfg.HTTPPort, log)
	// grpcSrv := grpc_server.NewServer(cfg.GRPCPort, grpcUserHandler, jwtClient, log)

	return &App{
		// httpServer: httpSrv,
		// grpcServer: grpcSrv,
		// dbPool:     pool,
		// logger:     log,
	}, nil
}

// Run starts both HTTP and gRPC servers with graceful shutdown.
func (a *App) Run(ctx context.Context) error {
	// TODO: Create channel for server errors
	errChan := make(chan error, 2)

	// TODO: Start HTTP server in goroutine
	go func() {
		if err := a.httpServer.Start(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	// TODO: Start gRPC server in goroutine
	go func() {
		if err := a.grpcServer.Start(); err != nil {
			errChan <- err
		}
	}()

	// TODO: Create signal channel for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// TODO: Wait for error or shutdown signal
	select {
	case err := <-errChan:
		a.logger.Error("server error", err)
		return err
	case sig := <-sigChan:
		a.logger.Info("shutdown signal received", sig)
		return a.Shutdown(ctx)
	}
}

// Shutdown gracefully shuts down all servers and closes connections.
func (a *App) Shutdown(ctx context.Context) error {
	// TODO: Create shutdown context with timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// TODO: Shutdown HTTP server
	a.logger.Info("shutting down HTTP server")
	if err := a.httpServer.Stop(shutdownCtx); err != nil {
		a.logger.Error("HTTP server shutdown error", err)
	}

	// TODO: Shutdown gRPC server
	a.logger.Info("shutting down gRPC server")
	a.grpcServer.Stop(shutdownCtx)

	// TODO: Close database pool
	if a.dbPool != nil {
		a.logger.Info("closing database pool")
		a.dbPool.Close()
	}

	// TODO: Sync logger
	if syncLogger, ok := a.logger.(*logger.ZapLogger); ok {
		syncLogger.Sync()
	}

	a.logger.Info("application shut down successfully")
	return nil
}
