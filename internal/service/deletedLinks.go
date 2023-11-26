package service

import "github.com/GZ91/linkreduct/internal/models"

// DeletedLinks добавляет URL в канал для последующего удаления.
func (r *NodeService) DeletedLinks(listURLs []string, userID string) {
	var dataForDel []models.StructDelURLs

	// Создание структуры данных для каждого URL и добавление их в срез
	for _, val := range listURLs {
		data := models.StructDelURLs{URL: val, UserID: userID}
		dataForDel = append(dataForDel, data)
	}

	// Отправка данных для удаления в канал
	r.ChsURLForDel <- dataForDel
}
