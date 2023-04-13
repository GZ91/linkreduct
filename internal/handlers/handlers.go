package handlers

import (
	"fmt"
	"github.com/GZ91/linkreduct/internal/config"
	"github.com/GZ91/linkreduct/internal/storage"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
)

var configHandler *config.Config

func InstallConfig(conf *config.Config) {
	configHandler = conf
}

type Middleware func(http.Handler) http.Handler

func AddLongLink(w http.ResponseWriter, r *http.Request) {
	link, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if string(link) == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id := storage.AddURL(string(link), configHandler)
	bodyText := configHandler.GetAddressServerURL() + id
	if bodyText == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Add("Content-Type", "text/plain")
	w.Header().Add("Content-Length", fmt.Sprint(len(bodyText)))
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(bodyText))
}

func GetShortURL(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if link, ok := storage.DB.GetURL(id); ok {
		w.Header().Add("Location", link)
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}
