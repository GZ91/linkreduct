package main

import (
	"context"
	"net/http"
	_ "net/http/pprof"
	"strconv"
	"time"

	"github.com/GZ91/linkreduct/internal/app/config"
	"github.com/GZ91/linkreduct/internal/app/initializing"
	"github.com/GZ91/linkreduct/internal/models"
	"github.com/GZ91/linkreduct/internal/service"
	"github.com/GZ91/linkreduct/internal/service/genrunes"
	"github.com/GZ91/linkreduct/internal/storage/inmemory"
	"github.com/google/uuid"
)

func main() {
	conf := initializing.Configuration()
	go func() {
		for {
			testingFunction(conf)
			time.Sleep(time.Second * 1)
		}
	}()
	http.ListenAndServe(":8080", nil)
}

func testingFunction(conf *config.Config) {
	ctx := context.Background()

	genrun := genrunes.New()
	dbNode, err := inmemory.New(ctx, conf, genrun)
	if err != nil {
		panic(err)
	}
	serviceNode := service.New(ctx,
		service.AddDb(dbNode),
		service.AddChsURLForDel(ctx, make(chan []models.StructDelURLs)),
		service.AddConf(conf))

	var batchLink []models.IncomingBatchURL
	for i := 1; i <= 1000; i++ {
		link := "http://www." + uuid.New().String() + ".com"
		batchLink = append(batchLink, models.IncomingBatchURL{CorrelationID: strconv.Itoa(i), OriginalURL: link})
	}

	_, err = serviceNode.AddBatchLink(ctx, batchLink)
	if err != nil {
		panic(err)
	}

}
