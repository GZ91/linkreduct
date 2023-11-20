package loggermiddleware

import (
	"net/http"
	"time"

	"github.com/GZ91/linkreduct/internal/app/logger"
	"go.uber.org/zap"
)

// responseData представляет структуру для хранения данных о ответе HTTP.
type responseData struct {
	status     int
	returnData []byte
}

// loggingResoinseWriter представляет собой обертку для http.ResponseWriter с возможностью логирования.
// Реализует интерфейс io.Writer для записи данных ответа.
type loggingResoinseWriter struct {
	http.ResponseWriter
	responseData *responseData
}

// Write перенаправляет вызов метода Write на внутренний ResponseWriter.
// Записывает данные ответа в поле returnData структуры responseData.
// Возвращает количество записанных байт и ошибку (если есть).
func (r *loggingResoinseWriter) Write(b []byte) (int, error) {
	r.responseData.returnData = append(r.responseData.returnData, b...)
	size, err := r.ResponseWriter.Write(b)
	return size, err
}

// WriteHeader перенаправляет вызов метода WriteHeader на внутренний ResponseWriter.
// Записывает код состояния ответа в поле status структуры responseData.
func (r *loggingResoinseWriter) WriteHeader(statusCode int) {
	r.responseData.status = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

// WithLogging возвращает обработчик HTTP, добавляющий логирование длительности выполнения запроса,
// метода запроса, кода состояния ответа и тела ответа.
func WithLogging(h http.Handler) http.Handler {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		// Замеряем время начала обработки запроса
		start := time.Now()

		// Создаем объект responseData для хранения данных об ответе
		respData := &responseData{0, []byte{}}

		// Создаем loggingResoinseWriter для обертки ResponseWriter с возможностью логирования
		lw := loggingResoinseWriter{responseData: respData, ResponseWriter: w}

		// Передаем управление следующему обработчику в цепочке
		h.ServeHTTP(&lw, r)

		// Замеряем длительность выполнения запроса
		duration := time.Since(start)

		// Логируем информацию о запросе с использованием zap-логгера
		logger.Log.Info("logging middleware",
			zap.Float32("duration", float32(duration)),
			zap.String("method", r.Method),
			zap.Int("status", lw.responseData.status),
			zap.ByteString("body", lw.responseData.returnData),
		)
	}

	return http.HandlerFunc(logFn)
}
