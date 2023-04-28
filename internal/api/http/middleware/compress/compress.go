package compress

import (
	"compress/gzip"
	"github.com/GZ91/linkreduct/internal/app/logger"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strings"
)

type Compression struct {
	http.ResponseWriter
	Writer io.Writer
}

func (c Compression) Write(p []byte) (int, error) {
	return c.Writer.Write(p)

}

func Compress(h http.Handler) http.Handler {
	compMid := func(w http.ResponseWriter, r *http.Request) {

		comp := Compression{w, w}

		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {

			logger.Log.Info("middleware_Compress", zap.String("status", "decompression added - gzip"))

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

		if AcceptEncoding && AcceptEncodingType {

			logger.Log.Info("middleware_Compress", zap.String("status", "compression added - gzip"))
			gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
				logger.Log.Error("middleware_Compress", zap.String("error", err.Error()))
				return
			}
			defer gz.Close()

			w.Header().Add("Content-Encoding", "gzip")
			comp.Writer = gz
		}
		h.ServeHTTP(comp, r)

	}

	return http.HandlerFunc(compMid)
}
