package main

import (
	"github.com/GZ91/linkreduct/internal/app"
	"github.com/GZ91/linkreduct/internal/config"
)

func main() {

	appliction := app.New(config.New(false, "localhost:8080", 10))
	appliction.Run()

}
