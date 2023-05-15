package service

import (
	"regexp"
)

// Storeger
//
//go:generate mockery --name Storeger --with-expecter
type Storeger interface {
	AddURL(string) string
	GetURL(string) (string, bool)
	Ping() error
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

func (r *NodeService) GetSmallLink(longLink string) string {
	if !r.URLFilter.MatchString(longLink) {
		longLink = "http://" + longLink
	}
	id := r.addURL(longLink)
	return r.conf.GetAddressServerURL() + id
}

func (r *NodeService) Ping() error {
	return r.db.Ping()
}
