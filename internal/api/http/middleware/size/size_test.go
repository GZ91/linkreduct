package sizemiddleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateSize(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(``))

	handle := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello world"))
	})
	resHandle := CalculateSize(handle)

	resHandle.ServeHTTP(rec, req)

	len := rec.Header().Get("Content-Length")
	assert.Equal(t, len, "11", "TEST middleware size")
}
