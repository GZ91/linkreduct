package authenticationmiddleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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
	result := rec.Result()
	defer result.Body.Close()
	for _, val := range result.Cookies() {
		if val.Name == "Authorization" {
			userID = val.Value
			break
		}
	}

	assert.NotEqual(t, userID, "", "TEST middleware authentication")
}
