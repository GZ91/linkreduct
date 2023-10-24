package handlers

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_handlers_PingDataBase(t *testing.T) {

	SetupForTesting(t)

	{
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(""))
		handls.PingDataBase(rec, req)

		res := rec.Result()
		res.Body.Close()
		assert.Equal(t, http.StatusOK, res.StatusCode, "TEST GET ping DB")
	}

}
