package server

import (
	"errors"
	"github.com/GZ91/linkreduct/internal/api/http/handlers"
	"github.com/GZ91/linkreduct/internal/api/http/middleware"
	"github.com/GZ91/linkreduct/internal/app/config"
	"github.com/GZ91/linkreduct/internal/service"
	"github.com/GZ91/linkreduct/internal/storage/inmemory"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func Start(conf *config.Config) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(r.(string))
		}
	}()

	NodeStorage := inmemory.New(conf)
	NodeService := service.New(NodeStorage, conf)
	handls := handlers.New(NodeService)

	router := chi.NewRouter()
	router.Get("/{id}", middleware.WithLogging(handls.GetShortURL))
	router.Post("/", middleware.WithLogging(handls.AddLongLink))
	router.Post("/api/shorten", middleware.WithLogging(handls.AddLongLinkJSON))
	return http.ListenAndServe(conf.GetAddressServer(), router)

}
