package service

import (
	"context"
	"github.com/GZ91/linkreduct/internal/app/logger"
	"github.com/GZ91/linkreduct/internal/models"
	"go.uber.org/zap"
	"regexp"
)

// Storeger
//
//go:generate mockery --name Storeger --with-expecter
type Storeger interface {
	AddURL(context.Context, string) (string, error)
	GetURL(context.Context, string) (string, bool, error)
	Ping(context.Context) error
	AddBatchLink(context.Context, []models.IncomingBatchURL) ([]models.ReleasedBatchURL, error)
	FindLongURL(context.Context, string) (string, bool, error)
	GetLinksUser(context.Context, string) ([]models.ReturnedStructURL, error)
	InitializingRemovalChannel(context.Context, chan []models.StructDelURLs) error
}

// Storeger
//
//go:generate mockery --name ConfigerService --with-expecter
type ConfigerService interface {
	GetAddressServerURL() string
}

type NodeService struct {
	db           Storeger
	conf         ConfigerService
	URLFormat    *regexp.Regexp
	URLFilter    *regexp.Regexp
	ChsURLForDel chan []models.StructDelURLs
}

func New(ctx context.Context, db Storeger, conf ConfigerService, ChsURLForDel chan []models.StructDelURLs) *NodeService {
	Node := &NodeService{
		db:           db,
		conf:         conf,
		URLFormat:    regexp.MustCompile(`^(?:https?:\/\/)`),
		URLFilter:    regexp.MustCompile(`^(?:https?:\/\/)?(?:[^@\/\n]+@)?(?:www\.)?(\w+\.[^:\/\n]+)`),
		ChsURLForDel: ChsURLForDel,
	}

	err := Node.db.InitializingRemovalChannel(ctx, Node.ChsURLForDel)
	if err != nil {
		logger.Log.Error("error when initializing the delete link channel", zap.Error(err))
		return nil
	}
	return Node
}

func (r *NodeService) getFormatLongLink(longLink string) (string, error) {
	if !r.URLFormat.MatchString(longLink) {
		longLink = "http://" + longLink
	}
	return longLink, nil
}

func (r *NodeService) Ping(ctx context.Context) error {
	return r.db.Ping(ctx)
}

func (r *NodeService) Close() error {
	close(r.ChsURLForDel)
	return nil
}
