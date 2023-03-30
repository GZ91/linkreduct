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

var Config config.Config

func Start(config config.Config) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(r.(string))
		}
	}()
	Config = config

	mux := &http.ServeMux{}
	mux.HandleFunc("/", defaultHandler)

	return http.ListenAndServe(config.GetAddressServer(), mux)
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if err := methodGet(w, r); err != nil {
			w.Write([]byte(fmt.Sprintf("%v /t/n", err)))
		}
	case http.MethodPost:
		if err := methodPost(w, r); err != nil {
			w.Write([]byte(fmt.Sprintf("%v /t/n", err)))
		}
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}

func methodGet(w http.ResponseWriter, r *http.Request) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(r.(string))
		}
	}()
	id := strings.TrimPrefix(r.URL.Path, "/")
	if link, ok := database.GetUrl(id); ok {
		w.Write([]byte(link))
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
	return nil
}

func methodPost(w http.ResponseWriter, r *http.Request) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	link, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}
	id := database.AddUrl(string(link), Config)
	bodyText := "http://" + Config.GetAddressServer() + "/" + id
	w.Write([]byte(bodyText))
	w.WriteHeader(http.StatusCreated)
	return nil
}
