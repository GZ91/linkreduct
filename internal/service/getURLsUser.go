package service

import (
	"context"

	"github.com/GZ91/linkreduct/internal/models"
)

// GetURLsUser возвращает список коротких URL для пользователя с указанным идентификатором.
func (r *NodeService) GetURLsUser(ctx context.Context, userID string) ([]models.ReturnedStructURL, error) {
	// Получение базового URL сервера
	addressServer := r.conf.GetAddressServerURL()

	// Получение списка коротких URL для пользователя из хранилища
	returnedStructURL, err := r.db.GetLinksUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Добавление базового URL к каждому короткому URL в списке
	for index, val := range returnedStructURL {
		returnedStructURL[index].ShortURL = addressServer + val.ShortURL
	}

	return returnedStructURL, nil
}
