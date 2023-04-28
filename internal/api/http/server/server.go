package server

import (
	"errors"
	"github.com/GZ91/linkreduct/internal/api/http/handlers"
	"github.com/GZ91/linkreduct/internal/api/http/middleware/compress"
	"github.com/GZ91/linkreduct/internal/api/http/middleware/logger"
	"github.com/GZ91/linkreduct/internal/api/http/middleware/size"
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
	router.Use(sizemiddleware.CalculateSize)
	router.Use(loggermiddleware.WithLogging)
	router.Use(compressMiddleware.Compress)

	router.Get("/{id}", handls.GetShortURL)
	router.Post("/", handls.AddLongLink)
	router.Post("/api/shorten", handls.AddLongLinkJSON)
	return http.ListenAndServe(conf.GetAddressServer(), router)

}
