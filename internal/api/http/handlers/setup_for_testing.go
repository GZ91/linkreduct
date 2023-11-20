package handlers

import (
	"context"
	"testing"

	"github.com/GZ91/linkreduct/internal/app/config"
	"github.com/GZ91/linkreduct/internal/models"
	"github.com/GZ91/linkreduct/internal/service"
	"github.com/GZ91/linkreduct/internal/service/mocks"
	mock_test "github.com/stretchr/testify/mock"
)

var handls *Handlers

// SetupForTesting создает и настраивает компоненты для тестирования.
// Принимает объект тестирования t из пакета testing.
// Создает конфигурацию, контекст, и объект хранилища NodeStorage с использованием мока (mocks.Storeger).
// Инициализирует и настраивает узер-сервис (NodeService) с использованием созданных компонентов.
// Возвращает созданный мок-хранилища NodeStorage для дальнейшего использования в тестах.
func SetupForTesting(t *testing.T) *mocks.Storeger {
	// Создаем объект конфигурации для тестирования
	conf := config.New(true, "localhost:8080", "http://localhost:8080/", 5, 5, "info.txt")

	// Создаем фоновый контекст
	ctx := context.Background()

	// Создаем мок-хранилища NodeStorage
	NodeStorage := mocks.NewStoreger(t)

	// Настраиваем мок-хранилища для ожидания вызова метода InitializingRemovalChannel
	NodeStorage.On("InitializingRemovalChannel", mock_test.Anything, mock_test.Anything).Return(nil).Maybe()

	// Создаем и настраиваем узер-сервис (NodeService) с использованием созданных компонентов
	NodeService := service.New(ctx,
		service.AddDB(NodeStorage),
		service.AddChsURLForDel(ctx, make(chan []models.StructDelURLs)),
		service.AddConf(conf))

	// Создаем и возвращаем объект хендлера (handls) с узер-сервисом NodeService для дальнейшего использования в тестах
	handls = New(NodeService)
	return NodeStorage
}
