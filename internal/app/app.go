package app

import (
	"context"
	"fmt"

	"github.com/GZ91/linkreduct/internal/api/http/server"
	"github.com/GZ91/linkreduct/internal/app/config"
)

// App представляет структуру приложения.
type App struct {
	config *config.Config
}

// appLink содержит глобальный экземпляр приложения.
var appLink *App

// New создает новый экземпляр приложения или возвращает существующий, если он уже создан.
func New(config *config.Config) *App {
	// Если экземпляр приложения не создан, создаем его и сохраняем в appLink
	if appLink == nil {
		appLink = &App{
			config: config,
		}
		return appLink
	}
	// Если экземпляр приложения уже создан, возвращаем его
	return appLink
}

// Run запускает приложение с использованием переданного контекста.
func (app App) Run(ctx context.Context) {
	// Запуск сервера с переданным конфигом
	if err := server.Start(ctx, app.config); err != nil {
		fmt.Printf("%v \n", err)
	}
}
