package service

import "github.com/GZ91/linkreduct/internal/models"

func (r *NodeService) DeletedLinks(listURLs []string, userID string) {

	var dataForDel []models.StructDelURLs
	for _, val := range listURLs {
		data := models.StructDelURLs{URL: val, UserID: userID}
		dataForDel = append(dataForDel, data)
	}

	r.ChsURLForDel <- dataForDel
}
