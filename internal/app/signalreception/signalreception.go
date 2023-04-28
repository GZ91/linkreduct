package signalreception

import (
	"github.com/GZ91/linkreduct/internal/app/logger"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type closer interface {
	Close() error
}

func OnClose(server closer, wg *sync.WaitGroup, nameSystem string) {
	wg.Add(1)
	defer wg.Done()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	for sig := range c {
		switch sig {
		case syscall.SIGINT:
			if err := server.Close(); err != nil {
				logger.Log.Error(nameSystem+" stop error", zap.String("error", err.Error()))
				return
			}
			logger.Log.Info(nameSystem+" stoped", zap.String("status", sig.String()))
			return
		default:
			logger.Log.Info("there was a signal", zap.String("status", sig.String()))
		}

	}

}
