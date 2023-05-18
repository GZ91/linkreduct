package service

import (
	"github.com/GZ91/linkreduct/internal/errorsapp"
	"regexp"
)

// Storeger
//
//go:generate mockery --name Storeger --with-expecter
type Storeger interface {
	AddURL(string) string
	GetURL(string) (string, bool)
	Ping() error
	FindLongURL(string) (string, bool)
}

// Storeger
//
//go:generate mockery --name ConfigerService --with-expecter
type ConfigerService interface {
	GetAddressServerURL() string
}

func New(db Storeger, conf ConfigerService) *NodeService {
	return &NodeService{db: db, conf: conf, URLFilter: regexp.MustCompile(`^(?:https?:\/\/)`)}
}

type NodeService struct {
	db        Storeger
	conf      ConfigerService
	URLFilter *regexp.Regexp
}

func (r *NodeService) GetURL(id string) (string, bool) {
	return r.db.GetURL(id)
}

func (r *NodeService) addURL(link string) string {
	return r.db.AddURL(link)
}

func (r *NodeService) GetSmallLink(longLink string) (string, error) {
	if !r.URLFilter.MatchString(longLink) {
		longLink = "http://" + longLink
	}
	if id, ok := r.db.FindLongURL(longLink); ok {
		return r.conf.GetAddressServerURL() + id, errorsapp.ErrLinkAlreadyExists
	}
	id := r.addURL(longLink)
	return r.conf.GetAddressServerURL() + id, nil
}

func (r *NodeService) Ping() error {
	return r.db.Ping()
}
