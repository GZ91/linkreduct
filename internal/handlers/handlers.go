package handlers

import (
	"fmt"
	"github.com/GZ91/linkreduct/internal/config"
	"github.com/GZ91/linkreduct/internal/storage"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"strconv"
)

var configHandler *config.Config

func InstallConfig(conf *config.Config) {
	configHandler = conf
}

type Middleware func(http.Handler) http.Handler

func Conveyor(h http.Handler, middlewares ...Middleware) http.Handler {
	for _, middleware := range middlewares {
		h = middleware(h)
	}
	return h
}

func MethodPost(w http.ResponseWriter, r *http.Request) {
	link, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	if string(link) == "" {
		w.WriteHeader(http.StatusBadRequest)
	}
	id := storage.AddURL(string(link), configHandler)
	bodyText := "http://" + configHandler.GetAddressServerURL() + "/" + id
	w.Header().Add("Content-Type", "text/plain")
	w.Header().Add("Content-Length", fmt.Sprint(len(bodyText)))
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(bodyText))
}

func MethodGet(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if link, ok := storage.GetURL(id); ok {
		w.Header().Add("Location", link)
		if configHandler.GetDebug() {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(strconv.Itoa(http.StatusTemporaryRedirect)))
		} else {
			w.WriteHeader(http.StatusTemporaryRedirect)
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}
