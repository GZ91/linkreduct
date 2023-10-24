package handlers

import (
	"encoding/json"
	"github.com/GZ91/linkreduct/internal/app/logger"
	"github.com/GZ91/linkreduct/internal/models"
	"go.uber.org/zap"
	"net/http"
)

func (h *handlers) GetURLsUser(w http.ResponseWriter, r *http.Request) {
	var UserID string
	var userIDCTX models.CtxString = "userID"
	UserIDVal := r.Context().Value(userIDCTX)
	if UserIDVal != nil {
		UserID = UserIDVal.(string)
	}
	if UserID == "" {
		logger.Log.Info("trying to execute a method to retrieve a URL by a user by an unauthorized user")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	returnedURLs, err := h.nodeService.GetURLsUser(r.Context(), UserID)
	if err != nil {
		logger.Log.Error("when getting URLs on the user side", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(returnedURLs) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	jsonText, err := json.Marshal(returnedURLs)
	if err != nil {
		logger.Log.Error("when creating a json file in the URL return procedure by user", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonText)
	if err != nil {
		logger.Log.Error("response recording error", zap.Error(err))
	}

}
