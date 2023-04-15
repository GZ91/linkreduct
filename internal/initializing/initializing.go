package initializing

import (
	"github.com/GZ91/linkreduct/internal/config"
	"github.com/GZ91/linkreduct/internal/initializing/envs"
	"github.com/GZ91/linkreduct/internal/initializing/flags"
	"strings"
)

func Configuration() *config.Config {
	addressServer, addressServerForURL := GetAddressesServer()
	return config.New(false, addressServer, addressServerForURL, 10)
}

func GetAddressesServer() (string, string) {

	addressServer, addressServerForURL := envs.ReadEnv()

	if addressServer == "" || addressServerForURL == "" {
		addressServerFlag, addressServerForURLFlag := flags.ReadFlags()
		if addressServer == "" {
			addressServer = addressServerFlag
		}
		if addressServerForURL == "" {
			addressServerForURL = addressServerForURLFlag
		}
	}

	addressServerForURL = CheckChangeBaseURL(addressServer, addressServerForURL)
	return addressServer, addressServerForURL
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
