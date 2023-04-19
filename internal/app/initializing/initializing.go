package initializing

import (
	"github.com/GZ91/linkreduct/internal/app/config"
	"github.com/GZ91/linkreduct/internal/app/initializing/envs"
	"github.com/GZ91/linkreduct/internal/app/initializing/flags"
	"github.com/GZ91/linkreduct/internal/app/logger"
	"strings"
)

func Configuration() *config.Config {
	addressServer, addressServerForURL, logLvl := ReadParams()
	logger.Initializing(logLvl)
	return config.New(false, addressServer, addressServerForURL, 10)
}

func ReadParams() (string, string, string) {

	addressServer, addressServerForURL, logLvl := envs.ReadEnv()

	if addressServer == "" || addressServerForURL == "" || logLvl == "" {
		addressServerFlag, addressServerForURLFlag, logLvlFlag := flags.ReadFlags()
		if addressServer == "" {
			addressServer = addressServerFlag
		}
		if addressServerForURL == "" {
			addressServerForURL = addressServerForURLFlag
		}
		if logLvl == "" {
			logLvl = logLvlFlag
		}
	}

	addressServerForURL = CheckChangeBaseURL(addressServer, addressServerForURL)
	return addressServer, addressServerForURL, logLvl
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
