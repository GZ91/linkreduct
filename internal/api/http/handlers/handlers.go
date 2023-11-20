// Модуль handlers содержит в себе ручки вызова сервиса
package handlers

import (
	"context"
	"regexp"

	"github.com/GZ91/linkreduct/internal/models"
)

// HandlerserService - Интерфейс, содержащий методы слоя сервиса, необходимые для работы пакета handlers.
type HandlerserService interface {
	GetSmallLink(context.Context, string) (string, error)
	GetURL(context.Context, string) (string, bool, error)
	GetURLsUser(context.Context, string) ([]models.ReturnedStructURL, error)
	Ping(ctx context.Context) error
	AddBatchLink(context.Context, []models.IncomingBatchURL) ([]models.ReleasedBatchURL, error)
	DeletedLinks([]string, string)
}

// Handlers - Структура, содержащая узер-сервис и фильтр для обработки поступающих ссылок.
type Handlers struct {
	nodeService HandlerserService
	URLFilter   *regexp.Regexp
}

// New - Создает и возвращает ссылку на структуру Handlers с указанным узер-сервисом и фильтром для ссылок.
func New(nodeService HandlerserService) *Handlers {
	return &Handlers{
		nodeService: nodeService,
		URLFilter:   regexp.MustCompile(`^(?:https?:\/\/)?(?:[^@\/\n]+@)?(?:www\.)?(\w+\.[^:\/\n]+)`),
	}
}
