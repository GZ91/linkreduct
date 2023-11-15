package service

import (
	"context"
	"regexp"

	"github.com/GZ91/linkreduct/internal/app/logger"
	"github.com/GZ91/linkreduct/internal/models"
	"go.uber.org/zap"
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

func New(ctx context.Context, opts ...func(service *NodeService)) *NodeService {
	Node := &NodeService{
		URLFormat: regexp.MustCompile(`^(?:https?:\/\/)`),
		URLFilter: regexp.MustCompile(`^(?:https?:\/\/)?(?:[^@\/\n]+@)?(?:www\.)?(\w+\.[^:\/\n]+)`),
	}
	for _, opt := range opts {
		opt(Node)
	}
	return Node
}

func AddDB(db Storeger) func(service *NodeService) {
	return func(n *NodeService) {
		n.db = db
	}
}

func AddConf(conf ConfigerService) func(service *NodeService) {
	return func(n *NodeService) {
		n.conf = conf
	}
}

func AddChsURLForDel(ctx context.Context, ChsURLForDel chan []models.StructDelURLs) func(service *NodeService) {
	return func(n *NodeService) {
		n.ChsURLForDel = ChsURLForDel
		err := n.db.InitializingRemovalChannel(ctx, n.ChsURLForDel)
		if err != nil {
			logger.Log.Error("error when initializing the delete link channel", zap.Error(err))
		}
	}

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
