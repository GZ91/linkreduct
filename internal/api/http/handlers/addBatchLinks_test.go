package handlers

import (
	"context"
	"github.com/GZ91/linkreduct/internal/models"
	"github.com/stretchr/testify/assert"
	mock_test "github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_handlers_AddBatchLinks(t *testing.T) {
	mockStorager := SetupForTesting(t)

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
	batch = append(batch, models.IncomingBatchURL{CorrelationID: "1", OriginalURL: "https://www.deepl.com"})
	batch = append(batch, models.IncomingBatchURL{CorrelationID: "2", OriginalURL: "https://www.mail.ru"})
	batch = append(batch, models.IncomingBatchURL{CorrelationID: "3", OriginalURL: "https://www.google.com"})

	var retBatch []models.ReleasedBatchURL
	retBatch = append(retBatch, models.ReleasedBatchURL{CorrelationID: "1", ShortURL: "uygh"})
	retBatch = append(retBatch, models.ReleasedBatchURL{CorrelationID: "2", ShortURL: "uasdygh"})
	retBatch = append(retBatch, models.ReleasedBatchURL{CorrelationID: "3", ShortURL: "usdfger4h"})

	mockStorager.On("AddBatchLink", mock_test.Anything, batch).Return(retBatch, nil)
	handls.AddBatchLinks(rec, req)

	res := rec.Result()
	res.Body.Close()
	assert.Equal(t, http.StatusCreated, res.StatusCode, "TEST GET ping DB")

}

func BenchmarkAddBatchLinks(b *testing.B) {
	for i := 0; i < b.N; i++ {
		t := &testing.T{}
		Test_handlers_AddBatchLinks(t)
	}
}
