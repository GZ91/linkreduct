package handlers

import (
	"github.com/GZ91/linkreduct/internal/config"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func SetupForTesting() {
	InstallConfig(config.New(true, "localhost:8080", "localhost:8080", 5))
}

func TestPostGet(t *testing.T) {
	SetupForTesting()
	targetLink := "google.com"

	router := chi.NewRouter()
	router.Route("/", func(r chi.Router) {
		r.Get("/{id}", MethodGet)
		r.Post("/", MethodPost)
	})

	server := httptest.NewServer(router)
	defer server.Close()
	client := server.Client()

	result, err := client.Post(server.URL+"/", "text/html; charset=utf8", strings.NewReader(targetLink))
	if err != nil {
		return
	}
	body, _ := io.ReadAll(result.Body)
	result.Body.Close()
	strBody := string(body)
	id := strings.TrimPrefix(strBody, configHandler.GetAddressServer())

	server.CloseClientConnections()
	resp, err := client.Get(server.URL + "/" + id)

	if err != nil {
		return
	}
	defer resp.Body.Close()

	val := resp.Header.Get("Location")
	bodyByte, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	b, _ := strconv.Atoi(string(bodyByte))

	assert.Equal(t, targetLink, val, "TEST GET 307")
	assert.Equal(t, http.StatusTemporaryRedirect, b, "TEST GET 307")
}

func TestGet400(t *testing.T) {
	SetupForTesting()
	targetLink := "google.com"

	{
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(""))
		MethodPost(rec, req)

		res := rec.Result()
		res.Body.Close()
		assert.Equal(t, http.StatusBadRequest, res.StatusCode, "TEST POST 400")
	}

	{

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/"+"adsafwefgasgsgfasdfsdfasdsdafwvwe23dasdasd854@3e23K◘c☼", nil)

		MethodGet(rec, req)

		res := rec.Result()
		res.Body.Close()
		val := res.Header.Get("Location")

		assert.NotEqual(t, targetLink, val, "TEST GET 400 \"not found ID\" The ID exactly should not be found (Test entry of an unknown ID)")
		assert.Equal(t, http.StatusBadRequest, res.StatusCode, "TEST GET 400 \"not found ID\" The ID exactly should not be found (Test entry of an unknown ID)")
	}
}

func TestPost400(t *testing.T) {
	SetupForTesting()

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(""))
	MethodPost(rec, req)

	res := rec.Result()
	res.Body.Close()
	assert.Equal(t, http.StatusBadRequest, res.StatusCode, "TEST POST 400")
}
