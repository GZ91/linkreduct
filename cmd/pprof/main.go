package main

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/GZ91/linkreduct/internal/api/http/handlers"
	"github.com/GZ91/linkreduct/internal/app/config"
	"github.com/GZ91/linkreduct/internal/models"
	"github.com/GZ91/linkreduct/internal/service"
	"github.com/GZ91/linkreduct/internal/service/mocks"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	mock_test "github.com/stretchr/testify/mock"
	"io"
	"net/http"
	"net/http/httptest"
	_ "net/http/pprof" // подключаем пакет pprof
	"strings"
	"testing"
)

func main() {
	go func() {
		for {
			testingFunction()
		}
	}()
	http.ListenAndServe(":8080", nil)
}

func testingFunction() {
	t := &testing.T{}

	conf := config.New(true, "localhost:8080", "http://localhost:8080/", 5, 5, "info.txt")
	NodeStorage := mocks.NewStoreger(t)
	NodeStorage.On("InitializingRemovalChannel", mock_test.Anything, mock_test.Anything).Return(nil).Maybe()
	NodeService := service.New(context.Background(), NodeStorage, conf, make(chan []models.StructDelURLs))
	handls := handlers.New(NodeService)

	mockStorager := NodeStorage

	{

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

		var batch []models.IncomingBatchURL
		batch = append(batch, models.IncomingBatchURL{"1", "https://www.deepl.com"})
		batch = append(batch, models.IncomingBatchURL{"2", "https://www.mail.ru"})
		batch = append(batch, models.IncomingBatchURL{"3", "https://www.google.com"})

		var retBatch []models.ReleasedBatchURL
		retBatch = append(retBatch, models.ReleasedBatchURL{"1", "uygh"})
		retBatch = append(retBatch, models.ReleasedBatchURL{"2", "uasdygh"})
		retBatch = append(retBatch, models.ReleasedBatchURL{"3", "usdfger4h"})

		mockStorager.On("AddBatchLink", mock_test.Anything, batch).Return(retBatch, nil)
		handls.AddBatchLinks(rec, req)

		res := rec.Result()
		res.Body.Close()
		assert.Equal(t, http.StatusCreated, res.StatusCode, "TEST GET ping DB")
	}

	{
		mockStorager := handlers.SetupForTesting(t)
		targetLink := "http://google.com"
		mockStorager.On("FindLongURL", mock_test.Anything, targetLink).Return("nhjsdf", true, nil)

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
		mockStorager.On("GetURL", mock_test.Anything, id).Return(targetLink, true, nil)
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

	{
		targetLink := "https://practicum.yandex.ru"
		mockStorager.On("FindLongURL", mock_test.Anything, targetLink).Return("nhjsdf", true, nil)
		mockStorager.On("GetURL", mock_test.Anything, "nhjsdf").Return(targetLink, true, nil)

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

	{
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(`["VeJoV", "eCuqR", "oemJV"] `))

		var userIDCTX models.CtxString = "userID"
		req = req.WithContext(context.WithValue(req.Context(), userIDCTX, "userID"))

		handls.DeleteURLs(rec, req)

		res := rec.Result()
		res.Body.Close()
		assert.Equal(t, http.StatusAccepted, res.StatusCode, "TEST GET ping DB")
	}
	{
		targetLink := "http://google.com"
		mockStorager.On("GetURL", mock_test.Anything, "").Return(targetLink, true, nil)
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
			req := httptest.NewRequest(http.MethodGet, "/", nil)

			var userIDCTX models.CtxString = "userID"
			req = req.WithContext(context.WithValue(req.Context(), userIDCTX, "userID"))

			handls.GetLongURL(rec, req)

			res := rec.Result()
			res.Body.Close()
			val := res.Header.Get("Location")

			assert.Equal(t, targetLink, val, "TEST GET 307")
			assert.Equal(t, http.StatusTemporaryRedirect, res.StatusCode, "TEST GET 307")
		}

		{
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(""))
			handls.AddLongLink(rec, req)

			res := rec.Result()
			res.Body.Close()
			assert.Equal(t, http.StatusBadRequest, res.StatusCode, "TEST POST 400")
		}
		{
			mockStorager.On("GetLinksUser", mock_test.Anything, "userID").Return(nil, nil)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(""))

			var userIDCTX models.CtxString = "userID"
			req = req.WithContext(context.WithValue(req.Context(), userIDCTX, "userID"))

			handls.GetURLsUser(rec, req)

			res := rec.Result()
			res.Body.Close()
			assert.Equal(t, http.StatusNoContent, res.StatusCode, "TEST GET ping DB")
		}

		{
			textBody := `[
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
] `
			var data []models.IncomingBatchURL
			json.Unmarshal([]byte(textBody), &data)
			var returnData []models.ReleasedBatchURL
			returnData = append(returnData, models.ReleasedBatchURL{"1", "sdgsg"})
			returnData = append(returnData, models.ReleasedBatchURL{"2", "sdfg"})
			returnData = append(returnData, models.ReleasedBatchURL{"3", "sgrgrw"})

			mockStorager.On("AddBatchLink", mock_test.Anything, data).Return(returnData, nil)
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(textBody))
			var userIDCTX models.CtxString = "userID"
			req = req.WithContext(context.WithValue(req.Context(), userIDCTX, "userID"))

			handls.AddBatchLinks(rec, req)

			res := rec.Result()
			res.Body.Close()
			assert.Equal(t, http.StatusCreated, res.StatusCode, "TEST GET ping DB")

			rec = httptest.NewRecorder()
			req = httptest.NewRequest(http.MethodGet, "/", strings.NewReader(""))

			req = req.WithContext(context.WithValue(req.Context(), userIDCTX, "userID"))

			var returnedData []models.ReturnedStructURL
			returnedData = append(returnedData, models.ReturnedStructURL{"sdgsg", "https://www.deepl.com"})
			returnedData = append(returnedData, models.ReturnedStructURL{"sdfg", "https://www.mail.ru"})
			returnedData = append(returnedData, models.ReturnedStructURL{"sgrgrw", "https://www.google.com"})

			mockStorager.On("GetLinksUser", mock_test.Anything, "userID").Return(returnedData, nil)

			handls.GetURLsUser(rec, req)

			res = rec.Result()
			res.Body.Close()
			assert.Equal(t, http.StatusOK, res.StatusCode, "TEST GET ping DB")
		}
	}
}
