package main

import (
	"flag"
	"github.com/GZ91/linkreduct/internal/app"
	"github.com/GZ91/linkreduct/internal/config"
)

func main() {

	appliction := app.New(Configuration())
	appliction.Run()

}

func Configuration() *config.Config {
	addressServer := flag.String("a", "localhost:8080", "Run Address server")
	addressServerURL := flag.String("b", "localhost:8080", "Address server for URL")
	flag.Parse()
	return config.New(false, *addressServer, *addressServerURL, 10)
}
