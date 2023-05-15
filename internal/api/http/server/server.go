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
	"github.com/GZ91/linkreduct/internal/service/genrunes"
	"github.com/GZ91/linkreduct/internal/storage/infile"
	"github.com/GZ91/linkreduct/internal/storage/inmemory"
	"github.com/GZ91/linkreduct/internal/storage/postgresql"
	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"net/http"
	"sync"
)

type NodeStorager interface {
	service.Storeger
	Close() error
}

func Start(conf *config.Config) (er error) {
	defer func() {
		if r := recover(); r != nil {
			er = errors.Wrap(er, r.(error).Error())
		}
	}()

	var NodeStorage NodeStorager
	GeneratorRunes := genrunes.New()
	if !conf.GetConfDB().Empty() {
		NodeStorage = postgresql.New(conf, GeneratorRunes)
	} else if conf.GetNameFileStorage() != "" {
		NodeStorage = infile.New(conf, GeneratorRunes)
	} else {
		NodeStorage = inmemory.New(conf, GeneratorRunes)
	}

	NodeService := service.New(NodeStorage, conf)
	handls := handlers.New(NodeService)
	//handls.PstgrSQL = conf.GetConfDB() //Временное решение для выполнения задачи с Ping

	router := chi.NewRouter()
	router.Use(sizemiddleware.CalculateSize)
	router.Use(loggermiddleware.WithLogging)
	router.Use(compressmiddleware.Compress)

	router.Get("/ping", handls.PingDataBase)
	router.Get("/{id}", handls.GetShortURL)
	router.Post("/", handls.AddLongLink)
	router.Post("/api/shorten", handls.AddLongLinkJSON)

	Server := http.Server{}
	Server.Addr = conf.GetAddressServer()
	Server.Handler = router

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
