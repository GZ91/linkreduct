package server

import (
	"context"
	"net/http"
	"sync"

	compressmiddleware "github.com/GZ91/linkreduct/internal/api/http/middleware/compress"
	sizemiddleware "github.com/GZ91/linkreduct/internal/api/http/middleware/size"
	"github.com/GZ91/linkreduct/internal/models"

	"github.com/GZ91/linkreduct/internal/api/http/handlers"
	authenticationmiddleware "github.com/GZ91/linkreduct/internal/api/http/middleware/authentication"
	"github.com/GZ91/linkreduct/internal/api/http/middleware/logger"
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
)

// NodeStorager представляет интерфейс, расширяющий интерфейс service.Storeger и добавляющий метод Close().
type NodeStorager interface {
	service.Storeger
	Close() error
}

// Утверждения типов, проверяющие реализацию интерфейсов NodeStorager различными типами.

var (
	_ NodeStorager               = (*postgresql.DB)(nil)
	_ NodeStorager               = (*infile.DB)(nil)
	_ NodeStorager               = (*inmemory.DB)(nil)
	_ handlers.HandlerserService = (*service.NodeService)(nil)
)

// Start запускает сервер с использованием заданного контекста и конфигурации.
// Создает хранилище и сервис, инициализирует HTTP-сервер и слушает сигналы для корректного завершения работы.
func Start(ctx context.Context, conf *config.Config) (er error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Получаем хранилище данных в зависимости от конфигурации
	NodeStorage, err := getNodeStorage(ctx, conf)
	if err != nil {
		return err
	}

	// Создаем сервис с использованием хранилища данных
	NodeService := service.New(ctx,
		service.AddDB(NodeStorage),
		service.AddChsURLForDel(ctx, make(chan []models.StructDelURLs)),
		service.AddConf(conf))

	// Инициализируем HTTP-сервер
	Server := http.Server{}
	Server.Addr = conf.GetAddressServer()
	Server.Handler = routing(handlers.New(NodeService))

	// Создаем WaitGroup для синхронизации горутин
	wg := sync.WaitGroup{}

	// Запускаем слушатель сигналов для корректного завершения работы
	go signalreception.OnClose(cancel, []signalreception.Closer{
		&signalreception.Stopper{CloserInterf: &Server, Name: "server"},
		&signalreception.Stopper{CloserInterf: NodeStorage, Name: "node storage"},
		&signalreception.Stopper{CloserInterf: NodeService, Name: "node service"}},
		&wg)

	// Запускаем HTTP-сервер
	if err := Server.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			logger.Log.Error("server startup error", zap.String("error", err.Error()))
		}
	}

	// Ожидаем завершения работы всех горутин
	wg.Wait()
	return
}

// getNodeStorage создает и возвращает хранилище данных в зависимости от конфигурации.
func getNodeStorage(ctx context.Context, conf *config.Config) (NodeStorager, error) {
	GeneratorRunes := genrunes.New()
	if !conf.GetConfDB().Empty() {
		return postgresql.New(ctx, conf, GeneratorRunes)
	} else if conf.GetNameFileStorage() != "" {
		return infile.New(ctx, conf, GeneratorRunes)
	} else {
		return inmemory.New(ctx, conf, GeneratorRunes)
	}
}

// routing создает и настраивает маршрутизатор с использованием заданных обработчиков.
func routing(handls *handlers.Handlers) *chi.Mux {
	router := chi.NewRouter()
	router.Use(authenticationmiddleware.Authentication)
	router.Use(sizemiddleware.CalculateSize)
	router.Use(loggermiddleware.WithLogging)
	router.Use(compressmiddleware.Compress)

	// Задаем обработчики для различных маршрутов
	router.Get("/ping", handls.PingDataBase)
	router.Get("/{id}", handls.GetLongURL)
	router.Get("/api/user/urls", handls.GetURLsUser)
	router.Post("/", handls.AddLongLink)
	router.Post("/api/shorten/batch", handls.AddBatchLinks)
	router.Post("/api/shorten", handls.AddLongLinkJSON)
	router.Delete("/api/user/urls", handls.DeleteURLs)
	return router
}
