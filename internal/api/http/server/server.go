package server

import (
	"github.com/GZ91/linkreduct/internal/api/http/handlers"
	"github.com/GZ91/linkreduct/internal/api/http/middleware/compress"
	"github.com/GZ91/linkreduct/internal/api/http/middleware/logger"
	"github.com/GZ91/linkreduct/internal/api/http/middleware/size"
	"github.com/GZ91/linkreduct/internal/app/config"
	"github.com/GZ91/linkreduct/internal/app/logger"
	"github.com/GZ91/linkreduct/internal/app/signalreception"
	"github.com/GZ91/linkreduct/internal/service"
	"github.com/GZ91/linkreduct/internal/storage/inmemory"
	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"net/http"
	"sync"
)

func Start(conf *config.Config) (er error) {
	defer func() {
		if r := recover(); r != nil {
			er = errors.Wrap(er, r.(error).Error())
		}
	}()

	NodeStorage := inmemory.New(conf)
	NodeService := service.New(NodeStorage, conf)
	handls := handlers.New(NodeService)

	router := chi.NewRouter()
	router.Use(sizemiddleware.CalculateSize)
	router.Use(loggermiddleware.WithLogging)
	router.Use(compressmiddleware.Compress)

	router.Get("/{id}", handls.GetShortURL)
	router.Post("/", handls.AddLongLink)
	router.Post("/api/shorten", handls.AddLongLinkJSON)

	Server := http.Server{}
	Server.Addr = conf.GetAddressServer()
	Server.Handler = router

	wg := sync.WaitGroup{}

	go signalreception.OnClose(&Server, &wg, "server")

	if err := Server.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			logger.Log.Error("server startup error", zap.String("error", err.Error()))
		}
	}
	wg.Wait()
	return

}
