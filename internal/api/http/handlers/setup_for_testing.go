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

var handls *Handlers

func SetupForTesting(t *testing.T) *mocks.Storeger {
	conf := config.New(true, "localhost:8080", "http://localhost:8080/", 5, 5, "info.txt")
	ctx := context.Background()
	NodeStorage := mocks.NewStoreger(t)
	NodeStorage.On("InitializingRemovalChannel", mock_test.Anything, mock_test.Anything).Return(nil).Maybe()
	NodeService := service.New(ctx,
		service.AddDB(NodeStorage),
		service.AddChsURLForDel(ctx, make(chan []models.StructDelURLs)),
		service.AddConf(conf))
	handls = New(NodeService)
	return NodeStorage
}
