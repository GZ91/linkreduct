package service

import (
	"context"

	"github.com/GZ91/linkreduct/internal/errorsapp"
)

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
