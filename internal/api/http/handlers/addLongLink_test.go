package handlers

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	mock_test "github.com/stretchr/testify/mock"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_handlers_AddLongLink(t *testing.T) {
	mockStorager := SetupForTesting(t)
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

func BenchmarkAddLongLink(b *testing.B) {
	for i := 0; i < b.N; i++ {
		t := &testing.T{}
		Test_handlers_AddLongLink(t)
	}
}
