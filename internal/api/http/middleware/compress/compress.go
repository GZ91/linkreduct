package compressmiddleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/GZ91/linkreduct/internal/app/logger"
	"go.uber.org/zap"
)

// Compression представляет собой обертку для http.ResponseWriter с поддержкой сжатия.
// Реализует интерфейс io.Writer для записи данных сжатия.
type Compression struct {
	http.ResponseWriter
	Writer io.Writer
}

// Write перенаправляет вызов метода Write на внутренний Writer Compression.
// Возвращает количество записанных байт и ошибку (если есть).
func (c Compression) Write(p []byte) (int, error) {
	return c.Writer.Write(p)
}

// Compress возвращает обработчик HTTP, добавляющий поддержку сжатия для передачи данных.
// Извлекает информацию о поддерживаемых методах сжатия из заголовков запроса.
// При наличии метода сжатия gzip в Accept-Encoding и подходящем Content-Type,
// добавляет сжатие gzip к ResponseWriter. Если запрос содержит Content-Encoding: gzip,
// добавляет декомпрессию gzip к Body запроса.
// Логирует соответствующие действия с использованием zap-логгера.
func Compress(h http.Handler) http.Handler {
	compMid := func(w http.ResponseWriter, r *http.Request) {
		// Создаем объект Compression для обертки ResponseWriter
		comp := Compression{w, w}

		// Проверяем наличие Content-Encoding: gzip в заголовке запроса
		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			logger.Log.Info("middleware_Compress", zap.String("status", "decompression added - gzip"))

			// Если присутствует, добавляем декомпрессию gzip к Body запроса
			reader, err := gzip.NewReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
				logger.Log.Error("middleware_Compress", zap.String("error", err.Error()))
				return
			}
			defer reader.Close()
			r.Body = reader
		}

		// Проверяем наличие gzip в Accept-Encoding заголовка и соответствующем Content-Type
		AcceptEncoding := false
		TypesEncoding := r.Header.Values("Accept-Encoding")

		for index := range TypesEncoding {
			if strings.Contains(TypesEncoding[index], "gzip") {
				AcceptEncoding = true
				break
			}
		}

		AcceptEncodingType := strings.Contains(r.Header.Get("Content-Type"), "application/json") ||
			strings.Contains(r.Header.Get("Content-Type"), "text/html")

		// Если поддерживается gzip и Content-Type соответствует, добавляем сжатие gzip к ResponseWriter
		if AcceptEncoding && AcceptEncodingType {
			logger.Log.Info("middleware_Compress", zap.String("status", "compression added - gzip"))

			// Создаем объект gzip.Writer для записи сжатых данных в ResponseWriter
			gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
				logger.Log.Error("middleware_Compress", zap.String("error", err.Error()))
				return
			}
			defer gz.Close()

			// Устанавливаем заголовок Content-Encoding: gzip
			w.Header().Add("Content-Encoding", "gzip")
			comp.Writer = gz
		}

		// Передаем управление следующему обработчику в цепочке
		h.ServeHTTP(comp, r)
	}

	return http.HandlerFunc(compMid)
}
