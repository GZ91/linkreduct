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

func Test_handlers_GetURLsUser(t *testing.T) {
	SetupForTesting(t)
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
	SetupForTesting(t)

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
