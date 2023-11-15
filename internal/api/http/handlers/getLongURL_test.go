package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/GZ91/linkreduct/internal/models"
	"github.com/stretchr/testify/assert"
	mock_test "github.com/stretchr/testify/mock"
)

func TestGet400(t *testing.T) {
	mockStorager := SetupForTesting(t)
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
}

func TestPost400(t *testing.T) {
	SetupForTesting(t)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(""))
	handls.AddLongLink(rec, req)

	res := rec.Result()
	res.Body.Close()
	assert.Equal(t, http.StatusBadRequest, res.StatusCode, "TEST POST 400")
}

func BenchmarkGet400(b *testing.B) {
	for i := 0; i < b.N; i++ {
		t := &testing.T{}
		TestGet400(t)
	}
}

func BenchmarkPost400(b *testing.B) {
	for i := 0; i < b.N; i++ {
		t := &testing.T{}
		TestPost400(t)
	}
}
