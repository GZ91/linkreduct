package main

import (
	"flag"
	"github.com/GZ91/linkreduct/internal/app"
	"github.com/GZ91/linkreduct/internal/config"
	"strings"
)

func main() {

	appliction := app.New(Configuration())
	appliction.Run()

}

func Configuration() *config.Config {
	addressServer := flag.String("a", "localhost:8080", "Run Address server")
	addressServerURL := flag.String("b", "http://localhost:8080/", "Address server for URL")

	strAddress := strings.Split(*addressServerURL, ":")
	port := strAddress[2]

	if len(port) == 0 {
		port = strings.Split(*addressServer, ":")[1]
	}
	if port[len(port)-1] != '/' {
		port = port + "/"
	}

	flag.Parse()
	return config.New(false, *addressServer, strAddress[0]+":"+strAddress[1]+":"+port, 10)
}
