package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	mock_test "github.com/stretchr/testify/mock"
)

func Exapmle() {

	req := httptest.NewRequest(http.MethodGet, "/ShortURL", nil)
	w := httptest.NewRecorder()

	t := &testing.T{}
	Storager := SetupForTesting(t)
	Storager.On("GetURL", mock_test.Anything, "ShortURL").Return("LongURL", true, nil)

	// handls инициализируется и объявлена в функции SetupForTesting
	handls.GetLongURL(w, req)

}
