package service

import (
	"context"
	"github.com/GZ91/linkreduct/internal/errorsapp"
	"github.com/GZ91/linkreduct/internal/models"
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
}

// Storeger
//
//go:generate mockery --name ConfigerService --with-expecter
type ConfigerService interface {
	GetAddressServerURL() string
}

func New(db Storeger, conf ConfigerService) *NodeService {
	return &NodeService{db: db, conf: conf, URLFormat: regexp.MustCompile(`^(?:https?:\/\/)`),
		URLFilter: regexp.MustCompile(`^(?:https?:\/\/)?(?:[^@\/\n]+@)?(?:www\.)?(\w+\.[^:\/\n]+)`)}
}

type NodeService struct {
	db        Storeger
	conf      ConfigerService
	URLFormat *regexp.Regexp
	URLFilter *regexp.Regexp
}

func (r *NodeService) GetURL(ctx context.Context, id string) (string, bool, error) {
	return r.db.GetURL(ctx, id)
}

func (r *NodeService) addURL(ctx context.Context, link string) (string, error) {
	return r.db.AddURL(ctx, link)
}

func (r *NodeService) GetSmallLink(ctx context.Context, longLink string) (string, error) {

	longLink, err := r.getFormatLongLink(longLink)
	if err != nil {
		return "", err
	}
	id, ok, err := r.db.FindLongURL(ctx, longLink)
	if err != nil {
		return "", err
	}
	if ok {
		return r.conf.GetAddressServerURL() + id, errorsapp.ErrLinkAlreadyExists
	}
	id, err = r.addURL(ctx, longLink)
	if err != nil {
		return "", err
	}
	return r.conf.GetAddressServerURL() + id, nil
}

func (r *NodeService) AddBatchLink(ctx context.Context, batchLink []models.IncomingBatchURL) (releasedBatchURL []models.ReleasedBatchURL, errs error) {

	for _, data := range batchLink {
		link := data.OriginalURL

		if !r.URLFilter.MatchString(link) {
			return nil, errorsapp.ErrInvalidLinkReceived
		}
	}

	releasedBatchURL, errs = r.db.AddBatchLink(ctx, batchLink)
	return
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

func (r *NodeService) GetURLsUser(ctx context.Context, userID string) ([]models.ReturnedStructURL, error) {
	return r.db.GetLinksUser(ctx, userID)
}
