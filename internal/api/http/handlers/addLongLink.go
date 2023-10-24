package handlers

import (
	"errors"
	"github.com/GZ91/linkreduct/internal/app/logger"
	"github.com/GZ91/linkreduct/internal/errorsapp"
	"go.uber.org/zap"
	"io"
	"net/http"
)

func (h *handlers) AddLongLink(w http.ResponseWriter, r *http.Request) {
	StatusReturn := http.StatusCreated

	link, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	if !h.URLFilter.MatchString(string(link)) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	bodyText, err := h.nodeService.GetSmallLink(r.Context(), string(link))
	if err != nil {
		if errors.Is(err, errorsapp.ErrLinkAlreadyExists) {
			StatusReturn = http.StatusConflict
		} else {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
	}

	if bodyText == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Add("Content-Type", "text/plain")
	w.WriteHeader(StatusReturn)
	_, err = w.Write([]byte(bodyText))
	if err != nil {
		logger.Log.Error("response recording error", zap.String("error", err.Error()))
	}
}
