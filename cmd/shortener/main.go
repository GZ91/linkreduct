package main

import (
	"flag"
	"github.com/GZ91/linkreduct/internal/app"
	"github.com/GZ91/linkreduct/internal/config"
	"github.com/caarlos0/env/v6"
	"strings"
)

func main() {

	appliction := app.New(Configuration())
	appliction.Run()

}

func Configuration() *config.Config {
	var (
		addressServerFlag       string
		addressServerForURLFlag string

		addressServerEnv       string
		addressServerForURLEnv string

		addressServer       string
		addressServerForURL string
	)
	addressServerEnv, addressServerForURLEnv = ReadEnv()

	if addressServerEnv == "" || addressServerForURLEnv == "" {
		addressServerFlag, addressServerForURLFlag = ReadFlags()
		if addressServerEnv == "" {
			addressServer = addressServerFlag
		} else {
			addressServer = addressServerEnv
		}
		if addressServerForURLEnv == "" {
			addressServerForURL = addressServerForURLFlag
		} else {
			addressServerForURL = addressServerForURLEnv
		}
	} else {
		addressServer = addressServerEnv
		addressServerForURL = addressServerForURLEnv
	}
	addressServerForURL = CheckChangeBaseURL(addressServer, addressServerForURL)
	return config.New(false, addressServer, addressServerForURL, 10)

}

func ReadFlags() (string, string) {
	addressServer := flag.String("a", "localhost:8080", "Run Address server")
	addressServerURL := flag.String("b", "http://localhost:8080/", "Address server for URL")

	flag.Parse()
	return *addressServer, *addressServerURL
}

type EnvVars struct {
	AddressServer       string `env:"SERVER_ADDRESS"`
	AddressServerForURL string `env:"BASE_URL"`
}

func ReadEnv() (string, string) {
	envs := EnvVars{}
	if err := env.Parse(&envs); err != nil {
		panic(err)
	}

	return envs.AddressServer, envs.AddressServerForURL
}

func CheckChangeBaseURL(addressServer, addressServerURL string) string {
	strAddress := strings.Split(addressServerURL, ":")
	var port string
	if (len(strAddress)) == 3 {
		port = strAddress[2]
	} else {
		port = ""
	}

	if len(port) == 0 {
		port = strings.Split(addressServer, ":")[1]
	}
	if port[len(port)-1] != '/' {
		port = port + "/"
	}
	return strAddress[0] + ":" + strAddress[1] + ":" + port
}
