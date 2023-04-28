package loggermiddleware

import (
	"github.com/GZ91/linkreduct/internal/app/logger"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type (
	responseData struct {
		status int
		size   int
	}

	loggingResoinseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

func (r *loggingResoinseWriter) WriteHeader(status int) {
	r.ResponseWriter.WriteHeader(status)
}

func (r *loggingResoinseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func WithLogging(h http.Handler) http.Handler {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		respData := &responseData{0, 0}
		lw := loggingResoinseWriter{responseData: respData, ResponseWriter: w}

		h.ServeHTTP(&lw, r)

		duration := time.Since(start)

		logger.Log.Info("logging middleware ", zap.Float32("duration", float32(duration)),
			zap.String("method", r.Method),
			zap.Int("status", lw.responseData.status),
			zap.Int("size", lw.responseData.size),
		)
	}
	return http.HandlerFunc(logFn)
}
