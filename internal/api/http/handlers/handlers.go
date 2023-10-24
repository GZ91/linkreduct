package handlers

import (
	"context"
	"github.com/GZ91/linkreduct/internal/models"
	"regexp"
)

type handlerserService interface {
	GetSmallLink(context.Context, string) (string, error)
	GetURL(context.Context, string) (string, bool, error)
	GetURLsUser(context.Context, string) ([]models.ReturnedStructURL, error)
	Ping(ctx context.Context) error
	AddBatchLink(context.Context, []models.IncomingBatchURL) ([]models.ReleasedBatchURL, error)
	DeletedLinks([]string, string)
}

type handlers struct {
	nodeService handlerserService
	URLFilter   *regexp.Regexp
}

func New(nodeService handlerserService) *handlers {
	return &handlers{nodeService: nodeService, URLFilter: regexp.MustCompile(`^(?:https?:\/\/)?(?:[^@\/\n]+@)?(?:www\.)?(\w+\.[^:\/\n]+)`)}
}
