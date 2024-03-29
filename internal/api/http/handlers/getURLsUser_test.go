package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/GZ91/linkreduct/internal/models"
	"github.com/stretchr/testify/assert"
	mock_test "github.com/stretchr/testify/mock"
)

func Test_handlers_GetURLsUser(t *testing.T) {
	mockStorager := SetupForTesting(t)
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

func Test_handlers_GetURLsUser2(t *testing.T) {
	mockStorager := SetupForTesting(t)

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
	returnData = append(returnData, models.ReleasedBatchURL{CorrelationID: "1", ShortURL: "sdgsg"})
	returnData = append(returnData, models.ReleasedBatchURL{CorrelationID: "2", ShortURL: "sdfg"})
	returnData = append(returnData, models.ReleasedBatchURL{CorrelationID: "3", ShortURL: "sgrgrw"})

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
	returnedData = append(returnedData, models.ReturnedStructURL{ShortURL: "sdgsg", OriginalURL: "https://www.deepl.com"})
	returnedData = append(returnedData, models.ReturnedStructURL{ShortURL: "sdfg", OriginalURL: "https://www.mail.ru"})
	returnedData = append(returnedData, models.ReturnedStructURL{ShortURL: "sgrgrw", OriginalURL: "https://www.google.com"})

	mockStorager.On("GetLinksUser", mock_test.Anything, "userID").Return(returnedData, nil)

	handls.GetURLsUser(rec, req)

	res = rec.Result()
	res.Body.Close()
	assert.Equal(t, http.StatusOK, res.StatusCode, "TEST GET ping DB")
}
