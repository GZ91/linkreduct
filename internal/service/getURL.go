package service

import "context"

// GetURL возвращает длинный URL по переданному короткому идентификатору.
// В случае успеха возвращает длинный URL, true и nil ошибки.
// Если идентификатор не найден, возвращает пустую строку, false и ошибку.
func (r *NodeService) GetURL(ctx context.Context, id string) (string, bool, error) {
	return r.db.GetURL(ctx, id)
}
