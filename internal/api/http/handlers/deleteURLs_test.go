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

func Test_handlers_DeleteURLs(t *testing.T) {
	SetupForTesting(t)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(`["VeJoV", "eCuqR", "oemJV"] `))

	var userIDCTX models.CtxString = "userID"
	req = req.WithContext(context.WithValue(req.Context(), userIDCTX, "userID"))

	handls.DeleteURLs(rec, req)

	res := rec.Result()
	res.Body.Close()
	assert.Equal(t, http.StatusAccepted, res.StatusCode, "TEST GET ping DB")

}
