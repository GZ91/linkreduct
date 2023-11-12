package handlers

import (
	"encoding/json"
	"github.com/GZ91/linkreduct/internal/app/logger"
	"github.com/GZ91/linkreduct/internal/models"
	"go.uber.org/zap"
	"io"
	"net/http"
)

func (h *Handlers) DeleteURLs(w http.ResponseWriter, r *http.Request) {
	var listURLs []string

	bodyByte, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Log.Error("error when reading the request body", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(bodyByte, &listURLs)
	if err != nil {
		logger.Log.Error("error when reading the json conversion", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if len(listURLs) == 0 {
		logger.Log.Error("sent an empty list of links to be deleted")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var UserID string
	var userIDCTX models.CtxString = "userID"
	UserIDVal := r.Context().Value(userIDCTX)
	if UserIDVal != nil {
		UserID = UserIDVal.(string)
	}
	if UserID == "" {
		logger.Log.Error("userID is not filled in")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	go h.nodeService.DeletedLinks(listURLs, UserID)
	w.WriteHeader(http.StatusAccepted)
}
