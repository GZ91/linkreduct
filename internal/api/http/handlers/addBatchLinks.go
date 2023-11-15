package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/GZ91/linkreduct/internal/app/logger"
	"github.com/GZ91/linkreduct/internal/errorsapp"
	"github.com/GZ91/linkreduct/internal/models"
	"go.uber.org/zap"
)

func (h *Handlers) AddBatchLinks(w http.ResponseWriter, r *http.Request) {
	StatusReturn := http.StatusCreated
	textBody, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var incomingBatchURL []models.IncomingBatchURL

	err = json.Unmarshal(textBody, &incomingBatchURL)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	releasedBatchURL, err := h.nodeService.AddBatchLink(r.Context(), incomingBatchURL)
	if err != nil {
		if errors.Is(err, errorsapp.ErrLinkAlreadyExists) {
			StatusReturn = http.StatusConflict
		} else {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
	}

	res, err := json.Marshal(releasedBatchURL)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	if len(res) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(StatusReturn)
	_, err = w.Write(res)
	if err != nil {
		logger.Log.Error("response recording error", zap.String("error", err.Error()))
	}
}
