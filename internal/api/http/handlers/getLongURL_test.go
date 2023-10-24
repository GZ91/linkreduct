package handlers

import (
	"context"
	"github.com/GZ91/linkreduct/internal/models"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGet400(t *testing.T) {
	SetupForTesting(t)
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
	SetupForTesting(t)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(""))
	handls.AddLongLink(rec, req)

	res := rec.Result()
	res.Body.Close()
	assert.Equal(t, http.StatusBadRequest, res.StatusCode, "TEST POST 400")
}
