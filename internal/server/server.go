package server

import (
	"errors"
	"fmt"
	"github.com/GZ91/linkreduct/internal/config"
	"github.com/GZ91/linkreduct/internal/database"
	"io"
	"net/http"
	"strings"
)

var Config *config.Config

type Middleware func(http.Handler) http.Handler

func Start(config *config.Config) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(r.(string))
		}
	}()
	Config = config

	mux := &http.ServeMux{}
	mux.Handle("/", Conveyor(http.HandlerFunc(methodGet), methodPost))

	return http.ListenAndServe(config.GetAddressServer(), mux)
}

func Conveyor(h http.Handler, middlewares ...Middleware) http.Handler {
	for _, middleware := range middlewares {
		h = middleware(h)
	}
	return h
}

func methodPost(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			next.ServeHTTP(w, r)
			return
		}

		link, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		id := database.AddURL(string(link), Config)
		bodyText := "http://" + Config.GetAddressServer() + "/" + id
		w.Header().Add("Content-Type", "text/plain")
		w.Header().Add("Content-Length", fmt.Sprint(len(bodyText)))
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(bodyText))
		next.ServeHTTP(w, r)
	})
}

func methodGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id := strings.TrimPrefix(r.URL.Path, "/")
	if link, ok := database.GetURL(id); ok {
		w.Header().Add("Location", link)
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}
