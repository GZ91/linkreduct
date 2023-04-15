package server

import (
	"errors"
	"github.com/GZ91/linkreduct/internal/api/http/handlers"
	"github.com/GZ91/linkreduct/internal/config"
	"github.com/GZ91/linkreduct/internal/service"
	"github.com/GZ91/linkreduct/internal/storage"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func Start(conf *config.Config) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(r.(string))
		}
	}()

	NodeStorage := storage.New(conf)
	NodeService := service.New(NodeStorage, conf)
	handls := handlers.New(NodeService)

	router := chi.NewRouter()
	router.Get("/{id}", handls.GetShortURL)
	router.Post("/", handls.AddLongLink)
	return http.ListenAndServe(conf.GetAddressServer(), router)

}
