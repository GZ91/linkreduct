package handlers

import (
	"context"
	"regexp"

	"github.com/GZ91/linkreduct/internal/models"
)

type HandlerserService interface {
	GetSmallLink(context.Context, string) (string, error)
	GetURL(context.Context, string) (string, bool, error)
	GetURLsUser(context.Context, string) ([]models.ReturnedStructURL, error)
	Ping(ctx context.Context) error
	AddBatchLink(context.Context, []models.IncomingBatchURL) ([]models.ReleasedBatchURL, error)
	DeletedLinks([]string, string)
}

type Handlers struct {
	nodeService HandlerserService
	URLFilter   *regexp.Regexp
}

func New(nodeService HandlerserService) *Handlers {
	return &Handlers{nodeService: nodeService, URLFilter: regexp.MustCompile(`^(?:https?:\/\/)?(?:[^@\/\n]+@)?(?:www\.)?(\w+\.[^:\/\n]+)`)}
}
