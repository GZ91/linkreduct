package service

import (
	"context"

	"github.com/GZ91/linkreduct/internal/errorsapp"
)

// addURL добавляет длинный URL в хранилище и возвращает соответствующий короткий URL.
func (r *NodeService) addURL(ctx context.Context, link string) (string, error) {
	return r.db.AddURL(ctx, link)
}

// GetSmallLink возвращает короткий URL для переданного длинного URL, добавляя его в хранилище при необходимости.
func (r *NodeService) GetSmallLink(ctx context.Context, longLink string) (string, error) {
	// Форматирование длинного URL
	longLink, err := r.getFormatLongLink(longLink)
	if err != nil {
		return "", err
	}

	// Поиск длинного URL в хранилище
	id, ok, err := r.db.FindLongURL(ctx, longLink)
	if err != nil {
		return "", err
	}

	// Если длинный URL уже существует, возвращается короткий URL
	if ok {
		return r.conf.GetAddressServerURL() + id, errorsapp.ErrLinkAlreadyExists
	}

	// Если длинного URL нет в хранилище, он добавляется, и возвращается соответствующий короткий URL
	id, err = r.addURL(ctx, longLink)
	if err != nil {
		return "", err
	}
	return r.conf.GetAddressServerURL() + id, nil
}
