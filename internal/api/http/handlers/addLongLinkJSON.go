package handlers

import (
	"encoding/json"
	"errors"
	"github.com/GZ91/linkreduct/internal/app/logger"
	"github.com/GZ91/linkreduct/internal/errorsapp"
	"github.com/GZ91/linkreduct/internal/models"
	"go.uber.org/zap"
	"io"
	"net/http"
)

func (h *Handlers) AddLongLinkJSON(w http.ResponseWriter, r *http.Request) {
	StatusReturn := http.StatusCreated
	textBody, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var data models.RequestData

	err = json.Unmarshal(textBody, &data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	link := data.URL

	if !h.URLFilter.MatchString(link) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	bodyText, err := h.nodeService.GetSmallLink(r.Context(), link)
	if err != nil {
		if errors.Is(err, errorsapp.ErrLinkAlreadyExists) {
			StatusReturn = http.StatusConflict
		} else {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	if bodyText == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	Result := models.ResultReturn{Result: bodyText}

	res, err := json.Marshal(Result)
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
