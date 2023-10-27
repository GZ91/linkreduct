package service

import "context"

func (r *NodeService) GetURL(ctx context.Context, id string) (string, bool, error) {
	return r.db.GetURL(ctx, id)
}
