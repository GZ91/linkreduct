package handlers

import (
	"fmt"
	"github.com/GZ91/linkreduct/internal/config"
	"github.com/GZ91/linkreduct/internal/storage"
	"io"
	"net/http"
	"strings"
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

func MethodPost(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			next.ServeHTTP(w, r)
			return
		}

		link, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		if string(link) == "" {
			w.WriteHeader(http.StatusBadRequest)
		}
		id := storage.AddURL(string(link), configHandler)
		bodyText := "http://" + configHandler.GetAddressServer() + "/" + id
		w.Header().Add("Content-Type", "text/plain")
		w.Header().Add("Content-Length", fmt.Sprint(len(bodyText)))
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(bodyText))
		next.ServeHTTP(w, r)
	})
}

func MethodGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id := strings.TrimPrefix(r.URL.Path, "/")
	if link, ok := storage.GetURL(id); ok {
		w.Header().Add("Location", link)
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}
