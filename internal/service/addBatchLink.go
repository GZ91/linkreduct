package service

import (
	"context"

	"github.com/GZ91/linkreduct/internal/errorsapp"
	"github.com/GZ91/linkreduct/internal/models"
)

// AddBatchLink добавляет пакет URL в хранилище и возвращает информацию о результатах обработки.
func (r *NodeService) AddBatchLink(ctx context.Context, batchLink []models.IncomingBatchURL) (releasedBatchURL []models.ReleasedBatchURL, errs error) {
	// Проверка каждого URL в пакете
	for _, data := range batchLink {
		link := data.OriginalURL
		if !r.URLFilter.MatchString(link) {
			return nil, errorsapp.ErrInvalidLinkReceived
		}
	}

	// Добавление пакета URL в хранилище
	releasedBatchURL, errs = r.db.AddBatchLink(ctx, batchLink)

	// Обновление ShortURL в возвращаемых результатах, добавляя к ним адрес сервера
	for index, val := range releasedBatchURL {
		releasedBatchURL[index].ShortURL = r.conf.GetAddressServerURL() + val.ShortURL
	}
	return
}
