// internal/app/app.go
package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/artyomstank/virtual_deanery/internal/config"
	postgres_db "github.com/artyomstank/virtual_deanery/internal/db"
	postgres_repo "github.com/artyomstank/virtual_deanery/internal/repo/postgres"
	"github.com/artyomstank/virtual_deanery/internal/service"
	http_server "github.com/artyomstank/virtual_deanery/internal/transport/http"
	"github.com/artyomstank/virtual_deanery/internal/transport/http/handler"
	"github.com/artyomstank/virtual_deanery/internal/transport/http/router"
	"github.com/artyomstank/virtual_deanery/pkg/jwt"
	"github.com/artyomstank/virtual_deanery/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	httpServer *http_server.Server
	dbPool     *pgxpool.Pool
	logger     *logger.Logger
}

func New(cfg *config.Config) (*App, error) {
	log := logger.New(cfg.LogLevel)

	pool, err := postgres_db.New(context.Background(), cfg)
	if err != nil {
		log.Error("failed to connect to database", err, nil)
		return nil, err
	}

	jwtClient := jwt.NewManager(cfg.JWTSecret, cfg.JWTExpireHours)

	userRepo := postgres_repo.NewUserRepo(pool)
	roleRepo := postgres_repo.NewRoleRepo(pool)
	aclRepo := postgres_repo.NewACLRepo(pool)

	// Создаём три независимых сервиса
	userSvc := service.NewUserService(userRepo, roleRepo, aclRepo, jwtClient, log, 12)
	adminSvc := service.NewAdminService(userRepo, roleRepo, log, 12)
	aclSvc := service.NewACLService(aclRepo, roleRepo, log)

	// Обработчики получают каждый свой интерфейс
	userHandler := handler.NewUserHandler(userSvc, jwtClient, log)
	adminHandler := handler.NewAdminHandler(adminSvc, log)
	aclHandler := handler.NewACLHandler(aclSvc, log)

	// Роутер
	corsOrigins := []string{"*"} // TODO: from config
	httpRouter := router.NewRouter(
		userHandler,
		adminHandler,
		aclHandler,
		userSvc, // для middleware аутентификации/ACL
		jwtClient,
		log,
		corsOrigins,
	)

	httpPort, err := strconv.Atoi(cfg.HTTPPort)
	if err != nil {
		httpPort = 8080
	}
	httpSrv := http_server.NewServer(httpRouter.Setup(), httpPort, log)

	log.Info("application initialized successfully", nil)

	return &App{
		httpServer: httpSrv,
		dbPool:     pool,
		logger:     log,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	errChan := make(chan error, 1)

	go func() {
		if err := a.httpServer.Start(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	a.logger.Info("HTTP server started", nil)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errChan:
		a.logger.Error("server error", err, nil)
		return err
	case sig := <-sigChan:
		a.logger.Info("shutdown signal received", map[string]interface{}{"signal": sig.String()})
		return a.Shutdown(ctx)
	}
}

func (a *App) Shutdown(ctx context.Context) error {
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	a.logger.Info("shutting down HTTP server", nil)
	if err := a.httpServer.Stop(shutdownCtx); err != nil {
		a.logger.Error("HTTP server shutdown error", err, nil)
	}

	if a.dbPool != nil {
		a.logger.Info("closing database pool", nil)
		a.dbPool.Close()
	}

	a.logger.Info("application shut down successfully", nil)
	return nil
}
