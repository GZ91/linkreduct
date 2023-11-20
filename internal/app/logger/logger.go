package logger

import "go.uber.org/zap"

// Log представляет глобальный логгер для приложения.
var Log *zap.Logger

// Initializing создает и настраивает логгер с заданным уровнем логгирования.
// Возвращает ошибку, если произошла ошибка при создании логгера.
func Initializing(level string) error {
	// Парсинг атомарного уровня логгирования
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err
	}

	// Создание конфигурации для логгера
	cfg := zap.NewProductionConfig()
	// Установка уровня логгирования
	cfg.Level = lvl

	// Построение логгера
	zl, err := cfg.Build()
	if err != nil {
		return err
	}

	// Установка глобального логгера
	Log = zl
	return nil
}
