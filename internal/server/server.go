package server

import (
	"errors"
	"github.com/GZ91/linkreduct/internal/config"
	"github.com/GZ91/linkreduct/internal/handlers"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func Start(conf *config.Config) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(r.(string))
		}
	}()
	handlers.InstallConfig(conf)

	router := chi.NewRouter()
	router.Get("/{id}", handlers.MethodGet)
	router.Post("/", handlers.MethodPost)
	return http.ListenAndServe(conf.GetAddressServer(), router)

}
