package handlers

import (
	"context"
	"github.com/GZ91/linkreduct/internal/app/config"
	"github.com/GZ91/linkreduct/internal/models"
	"github.com/GZ91/linkreduct/internal/service"
	"github.com/GZ91/linkreduct/internal/service/mocks"
	mock_test "github.com/stretchr/testify/mock"
	"testing"
)

var handls *handlers

func SetupForTesting(t *testing.T) *mocks.Storeger {
	conf := config.New(true, "localhost:8080", "http://localhost:8080/", 5, 5, "info.txt")

	NodeStorage := mocks.NewStoreger(t)
	NodeStorage.On("InitializingRemovalChannel", mock_test.Anything, mock_test.Anything).Return(nil).Maybe()
	NodeService := service.New(context.Background(), NodeStorage, conf, make(chan []models.StructDelURLs))
	handls = New(NodeService)
	return NodeStorage
}
