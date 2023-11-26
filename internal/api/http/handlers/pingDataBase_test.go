package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	mock_test "github.com/stretchr/testify/mock"
)

func Test_handlers_PingDataBase(t *testing.T) {
	mockStorager := SetupForTesting(t)
	{
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(""))
		mockStorager.On("Ping", mock_test.Anything).Return(nil)
		handls.PingDataBase(rec, req)

		res := rec.Result()
		res.Body.Close()
		assert.Equal(t, http.StatusOK, res.StatusCode, "TEST GET ping DB")
	}
}
