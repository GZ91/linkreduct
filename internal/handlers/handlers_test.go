package handlers

import (
	"github.com/GZ91/linkreduct/internal/config"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type TestHandler struct{}

func (t TestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {}

func TestHandlers(t *testing.T) {
	InstallConfig(config.New(true, "localhost:8080", 5))

	targetLink := "google.com"
	var id string

	{ //POST 201

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(targetLink))

		methodpost := MethodPost(TestHandler{})
		methodpost.ServeHTTP(rec, req)

		res := rec.Result()

		assert.Equal(t, http.StatusCreated, res.StatusCode, "TEST POST 201")

		body, err := io.ReadAll(res.Body)
		res.Body.Close()
		if assert.NoErrorf(t, err, "error at the request body stage TEST POST 201") {
			assert.NotEqual(t, 0, len(string(body)), "TEST POST 201")
		}
		id = strings.TrimPrefix(string(body), "http://"+configHandler.GetAddressServer()+"/")
	}
	{ //POST 400

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(""))
		methodpost := MethodPost(TestHandler{})
		methodpost.ServeHTTP(rec, req)

		res := rec.Result()
		res.Body.Close()
		assert.Equal(t, http.StatusBadRequest, res.StatusCode, "TEST POST 400")
	}
	{ //GET 307

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/"+id, nil)

		MethodGet(rec, req)

		res := rec.Result()
		res.Body.Close()
		val := res.Header.Get("Location")

		assert.Equal(t, targetLink, val, "TEST GET 307")
		assert.Equal(t, http.StatusTemporaryRedirect, res.StatusCode, "TEST GET 307")
	}
	{ //GET 400 not found ID

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/"+"adsafwefgasgsgfasdfsdfasdsdafwvwe23dasdasd854@3e23K◘c☼", nil)

		MethodGet(rec, req)

		res := rec.Result()
		res.Body.Close()
		val := res.Header.Get("Location")

		assert.NotEqual(t, targetLink, val, "TEST GET 400 \"not found ID\" The ID exactly should not be found (Test entry of an unknown ID)")
		assert.Equal(t, http.StatusBadRequest, res.StatusCode, "TEST GET 400 \"not found ID\" The ID exactly should not be found (Test entry of an unknown ID)")
	}
	{ //GET 400 another method

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/"+id, nil)

		MethodGet(rec, req)

		res := rec.Result()
		res.Body.Close()
		val := res.Header.Get("Location")

		assert.NotEqual(t, targetLink, val, "TEST GET 400 another method")
		assert.Equal(t, http.StatusBadRequest, res.StatusCode, "TEST GET 400 another method")
	}
}
