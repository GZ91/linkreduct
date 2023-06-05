package authenticationmiddleware

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAuthentication(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(``))

	handle := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello world"))
	})
	resHandle := Authentication(handle)

	resHandle.ServeHTTP(rec, req)

	var userID string
	for _, val := range rec.Result().Cookies() {
		if val.Name == "Authorization" {
			userID = val.Value
			break
		}
	}
	assert.NotEqual(t, userID, "", "TEST middleware authentication")
}
