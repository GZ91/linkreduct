package main

import (
	"github.com/GZ91/linkreduct/internal/app"
	"github.com/GZ91/linkreduct/internal/initializing"
)

func main() {

	appliction := app.New(initializing.Configuration())
	appliction.Run()

}
