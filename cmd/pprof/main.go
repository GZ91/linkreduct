package main

import (
	"context"
	"net/http"
	_ "net/http/pprof"

	"github.com/GZ91/linkreduct/internal/app/initializing"
	"github.com/GZ91/linkreduct/internal/models"
	"github.com/GZ91/linkreduct/internal/service"
	"github.com/GZ91/linkreduct/internal/service/genrunes"
	"github.com/GZ91/linkreduct/internal/storage/postgresql"
)

func main() {
	conf := initializing.Configuration()
	go func() {
		ctx := context.Background()

		genrun := genrunes.New()
		dbNode, err := postgresql.New(ctx, conf, genrun)
		if err != nil {
			panic(err)
		}
		serviceNode := service.New(ctx,
			service.AddDB(dbNode),
			service.AddChsURLForDel(ctx, make(chan []models.StructDelURLs)),
			service.AddConf(conf))
		for {
			testingFunction(ctx, serviceNode)
		}
	}()
	http.ListenAndServe(":8080", nil)
}

func testingFunction(ctx context.Context, serviceNode *service.NodeService) {
	_, err := serviceNode.GetURLsUser(ctx, "")
	if err != nil {
		panic(err)
	}

}
