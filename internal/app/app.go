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
	postgres_repo "github.com/artyomstank/virtual_deanery/internal/repo/postgres"
	"github.com/artyomstank/virtual_deanery/internal/service"
	http_server "github.com/artyomstank/virtual_deanery/internal/transport/http"
	"github.com/artyomstank/virtual_deanery/internal/transport/http/handler"
	"github.com/artyomstank/virtual_deanery/internal/transport/http/router"
	postgres_db "github.com/artyomstank/virtual_deanery/pkg/database/postgres"
	"github.com/artyomstank/virtual_deanery/pkg/jwt"
	"github.com/artyomstank/virtual_deanery/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

// App представляет экземпляр приложения со всеми зависимостями.
type App struct {
	httpServer *http_server.Server
	dbPool     *pgxpool.Pool
	logger     *logger.Logger
}

// New создаёт и инициализирует приложение со всеми зависимостями.
func New(cfg *config.Config) (*App, error) {
	// Инициализируем логгер
	log := logger.New(cfg.LogLevel)

	// Инициализируем пул соединений с БД
	pool, err := postgres_db.New(context.Background(), cfg)
	if err != nil {
		log.Error("failed to connect to database", err, nil)
		return nil, err
	}

	// Инициализируем JWT-клиент
	jwtClient := jwt.NewManager(cfg.JWTSecret, cfg.JWTExpireHours)

	// Инициализируем репозитории
	userRepo := postgres_repo.NewUserRepo(pool)
	roleRepo := postgres_repo.NewRoleRepo(pool)
	aclRepo := postgres_repo.NewACLRepo(pool)

	// Инициализируем сервисы
	userService := service.NewUserService(userRepo, roleRepo, aclRepo, jwtClient, log, cfg.BCryptCost)
	aclService := service.NewACLService(userRepo, roleRepo, aclRepo, log, userService)

	// Инициализируем HTTP-обработчики
	userHandler := handler.NewUserHandler(userService, log)
	aclHandler := handler.NewACLHandler(aclService, log)

	// Инициализируем HTTP-роутер
	corsOrigins := []string{"*"} // TODO: Load from config
	httpRouter := router.NewRouter(userHandler, aclHandler, userService, jwtClient, log, corsOrigins)

	// Инициализируем HTTP-сервер
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

// Run запускает HTTP-сервер с graceful shutdown.
func (a *App) Run(ctx context.Context) error {
	// Создаём канал для ошибок сервера
	errChan := make(chan error, 1)

	// Запускаем HTTP-сервер в отдельной горутине
	go func() {
		if err := a.httpServer.Start(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	a.logger.Info("HTTP server started", nil)

	// Создаём канал для сигналов ОС
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Ждём ошибку или сигнал завершения
	select {
	case err := <-errChan:
		a.logger.Error("server error", err, nil)
		return err
	case sig := <-sigChan:
		a.logger.Info("shutdown signal received", map[string]interface{}{"signal": sig.String()})
		return a.Shutdown(ctx)
	}
}

// Shutdown корректно завершает работу всех сервисов.
func (a *App) Shutdown(ctx context.Context) error {
	// Создаём контекст с timeout для завершения
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Завершаем HTTP-сервер
	a.logger.Info("shutting down HTTP server", nil)
	if err := a.httpServer.Stop(shutdownCtx); err != nil {
		a.logger.Error("HTTP server shutdown error", err, nil)
	}

	// Закрываем пул соединений БД
	if a.dbPool != nil {
		a.logger.Info("closing database pool", nil)
		a.dbPool.Close()
	}

	a.logger.Info("application shut down successfully", nil)
	return nil
}
