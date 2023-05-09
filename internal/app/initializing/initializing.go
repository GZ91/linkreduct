package initializing

import (
	"github.com/GZ91/linkreduct/internal/app/config"
	"github.com/GZ91/linkreduct/internal/app/initializing/envs"
	"github.com/GZ91/linkreduct/internal/app/initializing/flags"
	"github.com/GZ91/linkreduct/internal/app/logger"
	"go.uber.org/zap"
	"strings"
)

func Configuration() *config.Config {
	addressServer, addressServerForURL, logLvl, pathFileStorage := ReadParams()
	logger.Initializing(logLvl)
	return config.New(false, addressServer, addressServerForURL, 10, 5, pathFileStorage)
}

func ReadParams() (string, string, string, string) {

	envVars, err := envs.ReadEnv()
	if err != nil {
		logger.Log.Error("error when reading environment variables", zap.String("error", err.Error()))
	}

	var addressServer, addressServerForURL, logLvl, pathFileStorage string

	if envVars == nil {
		addressServer, addressServerForURL, logLvl, pathFileStorage = flags.ReadFlags()
	} else {
		addressServer, addressServerForURL, logLvl, pathFileStorage = envVars.AddressServer, envVars.AddressServerForURL, envVars.LvlLogs, envVars.PathFileStorage

		if addressServer == "" || addressServerForURL == "" || logLvl == "" || pathFileStorage == "" {
			addressServerFlag, addressServerForURLFlag, logLvlFlag, pathFileStorageFlag := flags.ReadFlags()
			if addressServer == "" {
				addressServer = addressServerFlag
			}
			if addressServerForURL == "" {
				addressServerForURL = addressServerForURLFlag
			}
			if logLvl == "" {
				logLvl = logLvlFlag
			}
			if pathFileStorage == "" {
				pathFileStorage = pathFileStorageFlag
			}
		}
	}

	addressServerForURL = CheckChangeBaseURL(addressServer, addressServerForURL)
	return addressServer, addressServerForURL, logLvl, pathFileStorage
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
