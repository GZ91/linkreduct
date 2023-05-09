package sizemiddleware

import (
	"net/http"
	"strconv"
)

type Size struct {
	http.ResponseWriter
}

func (s Size) Write(p []byte) (int, error) {
	len, err := s.ResponseWriter.Write(p)
	s.ResponseWriter.Header().Add("Content-Length", strconv.Itoa(len))
	return len, err
}

func CalculateSize(h http.Handler) http.Handler {
	CalSize := func(w http.ResponseWriter, r *http.Request) {
		s := Size{w}
		h.ServeHTTP(s, r)
	}
	return http.HandlerFunc(CalSize)
}
