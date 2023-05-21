package main

import (
	"context"
	"github.com/GZ91/linkreduct/internal/app"
	"github.com/GZ91/linkreduct/internal/app/initializing"
)

func main() {
	ctx := context.Background()
	appliction := app.New(initializing.Configuration())
	appliction.Run(ctx)

}
