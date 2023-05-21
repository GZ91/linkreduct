package server

import (
	"context"
	errors2 "errors"
	"github.com/GZ91/linkreduct/internal/api/http/handlers"
	"github.com/GZ91/linkreduct/internal/api/http/middleware/compress"
	"github.com/GZ91/linkreduct/internal/api/http/middleware/logger"
	"github.com/GZ91/linkreduct/internal/api/http/middleware/size"
	"github.com/GZ91/linkreduct/internal/app/config"
	"github.com/GZ91/linkreduct/internal/app/logger"
	"github.com/GZ91/linkreduct/internal/app/signalreception"
	"github.com/GZ91/linkreduct/internal/service"
	"github.com/GZ91/linkreduct/internal/service/genrunes"
	"github.com/GZ91/linkreduct/internal/storage/infile"
	"github.com/GZ91/linkreduct/internal/storage/inmemory"
	"github.com/GZ91/linkreduct/internal/storage/postgresql"
	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"net"
	"net/http"
	"sync"
)

type NodeStorager interface {
	service.Storeger
	Close() error
}

func Start(ctx context.Context, conf *config.Config) (er error) {
	defer func() {
		if r := recover(); r != nil {
			er = errors2.Join(er, r.(error))
			logger.Log.Error("Panic", zap.Error(er))
		}
	}()

	var NodeStorage NodeStorager
	GeneratorRunes := genrunes.New()
	if !conf.GetConfDB().Empty() {
		var err error
		NodeStorage, err = postgresql.New(ctx, conf, GeneratorRunes)
		if err != nil {
			panic(err)
		}
	} else if conf.GetNameFileStorage() != "" {
		NodeStorage = infile.New(ctx, conf, GeneratorRunes)
	} else {
		NodeStorage = inmemory.New(ctx, conf, GeneratorRunes)
	}

	NodeService := service.New(NodeStorage, conf)
	handls := handlers.New(NodeService)

	router := chi.NewRouter()
	router.Use(sizemiddleware.CalculateSize)
	router.Use(loggermiddleware.WithLogging)
	router.Use(compressmiddleware.Compress)

	router.Get("/ping", handls.PingDataBase)
	router.Get("/{id}", handls.GetLongURL)
	router.Post("/", handls.AddLongLink)
	router.Post("/api/shorten/batch", handls.AddBatchLinks)
	router.Post("/api/shorten", handls.AddLongLinkJSON)

	Server := http.Server{}
	Server.Addr = conf.GetAddressServer()
	Server.Handler = router
	Server.BaseContext = func(listener net.Listener) context.Context {
		return ctx
	}

	wg := sync.WaitGroup{}

	go signalreception.OnClose([]signalreception.Closer{
		&signalreception.Stopper{CloserInterf: &Server, Name: "server"},
		&signalreception.Stopper{CloserInterf: NodeStorage, Name: "node storage"}},
		&wg)

	if err := Server.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			logger.Log.Error("server startup error", zap.String("error", err.Error()))
		}
	}
	wg.Wait()
	return

}

func GracefulShotdownServer() {

}
