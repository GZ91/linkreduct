package signalreception

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/GZ91/linkreduct/internal/app/logger"
	"go.uber.org/zap"
)

// CloserInterf представляет интерфейс для объектов с методом Close.
type CloserInterf interface {
	Close() error
}

// Stopper представляет структуру, реализующую интерфейс Closer и содержащую информацию о закрываемом объекте.
type Stopper struct {
	CloserInterf
	Name string
}

// GetName возвращает имя объекта Stopper.
func (s Stopper) GetName() string {
	return s.Name
}

// Closer представляет интерфейс для объектов с методами Close и GetName.
type Closer interface {
	Close() error
	GetName() string
}

// OnClose ожидает сигналы завершения работы и выполняет закрытие переданных объектов.
// Принимает функцию отмены контекста, список объектов, реализующих интерфейс Closer, и WaitGroup для ожидания завершения работы.
func OnClose(cancel context.CancelFunc, closers []Closer, wg *sync.WaitGroup) {
	// Добавление одной горутины к WaitGroup
	wg.Add(1)
	defer wg.Done()

	// Создание канала для сигналов завершения
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// Ожидание сигнала завершения
	for sig := range c {
		switch sig {
		case syscall.SIGINT:
			// Отмена контекста
			cancel()
			// Закрытие переданных объектов
			for _, closer := range closers {
				if err := closer.Close(); err != nil {
					logger.Log.Error(closer.GetName()+" close error", zap.String("error", err.Error()))
				} else {
					logger.Log.Info(closer.GetName()+" close", zap.String("status", sig.String()))
				}
			}
			return
		default:
			logger.Log.Info("there was a signal", zap.String("status", sig.String()))
		}
	}
}
