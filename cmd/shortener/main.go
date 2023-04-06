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

	addressServer, addressServerForURL := ReadFlags()
	return config.New(false, addressServer, addressServerForURL, 10)

}

func ReadFlags() (string, string) {
	addressServer := flag.String("a", "localhost:8080", "Run Address server")
	addressServerURL := flag.String("b", "http://localhost:8080/", "Address server for URL")

	flag.Parse()

	strAddress := strings.Split(*addressServerURL, ":")
	var port string
	if (len(strAddress)) == 3 {
		port = strAddress[2]
	} else {
		port = ""
	}

	if len(port) == 0 {
		port = strings.Split(*addressServer, ":")[1]
	}
	if port[len(port)-1] != '/' {
		port = port + "/"
	}

	addressServerForURL := strAddress[0] + ":" + strAddress[1] + ":" + port
	return *addressServer, addressServerForURL
}