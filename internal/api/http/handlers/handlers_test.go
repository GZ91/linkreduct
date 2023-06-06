package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/GZ91/linkreduct/internal/app/config"
	"github.com/GZ91/linkreduct/internal/models"
	"github.com/GZ91/linkreduct/internal/service"
	"github.com/GZ91/linkreduct/internal/service/genrunes"
	"github.com/GZ91/linkreduct/internal/storage/inmemory"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var handls *handlers

func Test_handlers_AddLongLink(t *testing.T) {
	SetupForTesting()
	targetLink := "http://google.com"

	router := chi.NewRouter()
	router.Route("/", func(r chi.Router) {
		r.Get("/{id}", handls.GetLongURL)
		r.Post("/", handls.AddLongLink)
	})

	server := httptest.NewServer(router)
	defer server.Close()

	errRedirectBlocked := errors.New("HTTP redirect blocked")

	client := server.Client()

	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return errRedirectBlocked
	}

	result, err := client.Post(server.URL+"/", "text/html; charset=utf8", strings.NewReader(targetLink))
	if err != nil {
		return
	}
	body, _ := io.ReadAll(result.Body)
	result.Body.Close()
	strBody := string(body)
	id := strings.TrimPrefix(strBody, "http://localhost:8080/")

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

func TestGet400(t *testing.T) {
	SetupForTesting()
	targetLink := "http://google.com"

	{
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(""))
		handls.AddLongLink(rec, req)

		res := rec.Result()
		res.Body.Close()
		assert.Equal(t, http.StatusBadRequest, res.StatusCode, "TEST POST 400")
	}

	{

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/"+"adsafwefgasgsgfasdfsdfasdsdafwvwe23dasdasd854@3e23K◘c☼", nil)

		var userIDCTX models.CtxString = "userID"
		req = req.WithContext(context.WithValue(req.Context(), userIDCTX, "userID"))

		handls.GetLongURL(rec, req)

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
	handls.AddLongLink(rec, req)

	res := rec.Result()
	res.Body.Close()
	assert.Equal(t, http.StatusBadRequest, res.StatusCode, "TEST POST 400")
}

func Test_handlers_AddLongLinkJSON(t *testing.T) {
	SetupForTesting()
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

func SetupForTesting() {
	conf := config.New(true, "localhost:8080", "http://localhost:8080/", 5, 5, "C:\\Users\\Georgiy\\Desktop\\GO\\linkreduct\\info.txt")

	genrun := genrunes.New()
	NodeStorage := inmemory.New(context.Background(), conf, genrun)
	NodeService := service.New(NodeStorage, conf, make(chan []models.StructDelURLs))
	handls = New(NodeService)

}

func Test_handlers_PingDataBase(t *testing.T) {

	SetupForTesting()

	{
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(""))
		handls.PingDataBase(rec, req)

		res := rec.Result()
		res.Body.Close()
		assert.Equal(t, http.StatusOK, res.StatusCode, "TEST GET ping DB")
	}

}

func Test_handlers_AddBatchLinks(t *testing.T) {
	SetupForTesting()

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(`[
    {
        "correlation_id": "1",
        "original_url": "https://www.deepl.com"
    },
    {
        "correlation_id": "2",
        "original_url": "https://www.mail.ru"
    },
    {
        "correlation_id": "3",
        "original_url": "https://www.google.com"
    }
] `))
	var userIDCTX models.CtxString = "userID"
	req = req.WithContext(context.WithValue(req.Context(), userIDCTX, "userID"))

	handls.AddBatchLinks(rec, req)

	res := rec.Result()
	res.Body.Close()
	assert.Equal(t, http.StatusCreated, res.StatusCode, "TEST GET ping DB")

}

func Test_handlers_GetURLsUser(t *testing.T) {
	SetupForTesting()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(""))

	var userIDCTX models.CtxString = "userID"
	req = req.WithContext(context.WithValue(req.Context(), userIDCTX, "userID"))

	handls.GetURLsUser(rec, req)

	res := rec.Result()
	res.Body.Close()
	assert.Equal(t, http.StatusNoContent, res.StatusCode, "TEST GET ping DB")
}

func Test_handlers_GetURLsUser2(t *testing.T) {
	SetupForTesting()

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(`[
    {
        "correlation_id": "1",
        "original_url": "https://www.deepl.com"
    },
    {
        "correlation_id": "2",
        "original_url": "https://www.mail.ru"
    },
    {
        "correlation_id": "3",
        "original_url": "https://www.google.com"
    }
] `))
	var userIDCTX models.CtxString = "userID"
	req = req.WithContext(context.WithValue(req.Context(), userIDCTX, "userID"))

	handls.AddBatchLinks(rec, req)

	res := rec.Result()
	res.Body.Close()
	assert.Equal(t, http.StatusCreated, res.StatusCode, "TEST GET ping DB")

	rec = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/", strings.NewReader(""))

	req = req.WithContext(context.WithValue(req.Context(), userIDCTX, "userID"))

	handls.GetURLsUser(rec, req)

	res = rec.Result()
	res.Body.Close()
	assert.Equal(t, http.StatusOK, res.StatusCode, "TEST GET ping DB")
}

func Test_handlers_DeleteURLs(t *testing.T) {
	SetupForTesting()

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(`["VeJoV", "eCuqR", "oemJV"] `))

	var userIDCTX models.CtxString = "userID"
	req = req.WithContext(context.WithValue(req.Context(), userIDCTX, "userID"))

	handls.DeleteURLs(rec, req)

	res := rec.Result()
	res.Body.Close()
	assert.Equal(t, http.StatusAccepted, res.StatusCode, "TEST GET ping DB")

}
