package handlers

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_handlers_AddLongLinkJSON(t *testing.T) {
	SetupForTesting(t)
	targetLink := "https://practicum.yandex.ru"

	router := chi.NewRouter()
	router.Route("/", func(r chi.Router) {
		r.Get("/{id}", handls.GetLongURL)
		r.Post("/api/shorten", handls.AddLongLinkJSON)
	})

	server := httptest.NewServer(router)
	defer server.Close()

	errRedirectBlocked := errors.New("HTTP redirect blocked")

	client := server.Client()

	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return errRedirectBlocked
	}

	result, err := client.Post(server.URL+"/api/shorten", "text/html; charset=utf8", strings.NewReader(`{"url": "https://practicum.yandex.ru"} `))
	if err != nil {
		return
	}
	body, _ := io.ReadAll(result.Body)
	type resType struct {
		Result string `json:"result"`
	}
	var res resType
	json.Unmarshal(body, &res)

	result.Body.Close()

	id := strings.TrimPrefix(res.Result, "http://localhost:8080/")

	server.CloseClientConnections()
	resp, err := client.Get(server.URL + "/" + id)

	if err != nil {
		assert.Equal(t, true, errors.Is(err, errRedirectBlocked))
	}

	defer resp.Body.Close()

	val := resp.Header.Get("Location")
	io.Copy(io.Discard, resp.Body)

	assert.Equal(t, targetLink, val, "TEST GET 307")
	assert.Equal(t, http.StatusTemporaryRedirect, resp.StatusCode, "TEST GET 307")
}
