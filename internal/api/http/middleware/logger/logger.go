package loggermiddleware

import (
	"github.com/GZ91/linkreduct/internal/app/logger"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type (
	responseData struct {
		status     int
		returnData []byte
	}

	loggingResoinseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

func (r *loggingResoinseWriter) Write(b []byte) (int, error) {
	r.responseData.returnData = append(r.responseData.returnData, b...)
	size, err := r.ResponseWriter.Write(b)
	return size, err
}

func (r *loggingResoinseWriter) WriteHeader(statusCode int) {
	r.responseData.status = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func WithLogging(h http.Handler) http.Handler {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		respData := &responseData{0, []byte{}}
		lw := loggingResoinseWriter{responseData: respData, ResponseWriter: w}

		h.ServeHTTP(&lw, r)

		duration := time.Since(start)

		logger.Log.Info("logging middleware ", zap.Float32("duration", float32(duration)),
			zap.String("method", r.Method),
			zap.Int("status", lw.responseData.status),
			zap.ByteString("body", lw.responseData.returnData),
		)
	}
	return http.HandlerFunc(logFn)
}
