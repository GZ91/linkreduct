package service

import (
	"context"
	"github.com/GZ91/linkreduct/internal/models"
)

func (r *NodeService) GetURLsUser(ctx context.Context, userID string) ([]models.ReturnedStructURL, error) {
	addressServer := r.conf.GetAddressServerURL()
	returnedStructURL, err := r.db.GetLinksUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	for index, val := range returnedStructURL {
		returnedStructURL[index].ShortURL = addressServer + val.ShortURL
	}
	return returnedStructURL, nil
}
