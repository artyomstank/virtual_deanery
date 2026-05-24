// cmd/api/main.go
package main

import (
	"context"
	"log"

	"github.com/artyomstank/virtual_deanery/internal/app"
	"github.com/artyomstank/virtual_deanery/internal/config"
)

func main() {
	// Загружаем конфигурацию из переменных окружения
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Инициализируем приложение
	application, err := app.New(cfg)
	if err != nil {
		log.Fatalf("failed to initialize app: %v", err)
	}

	// Запускаем приложение с graceful shutdown
	if err := application.Run(context.Background()); err != nil {
		log.Fatalf("app run error: %v", err)
	}
}
