package initializing

import (
	"strings"

	"github.com/GZ91/linkreduct/internal/app/config"
	"github.com/GZ91/linkreduct/internal/app/initializing/envs"
	"github.com/GZ91/linkreduct/internal/app/initializing/flags"
	"github.com/GZ91/linkreduct/internal/app/logger"
	"go.uber.org/zap"
)

func Configuration() *config.Config {
	addressServer, addressServerForURL, logLvl, pathFileStorage, connectionStringDB := ReadParams()
	logger.Initializing(logLvl)
	conf := config.New(false, addressServer, addressServerForURL, 10, 5, pathFileStorage)
	conf.ConfigureDBPostgresql(connectionStringDB)
	return conf
}

func ReadParams() (string, string, string, string, string) {

	envVars, err := envs.ReadEnv()
	if err != nil {
		logger.Log.Error("error when reading environment variables", zap.String("error", err.Error()))
	}

	var addressServer, addressServerForURL, logLvl, pathFileStorage, connectionStringDB string

	if envVars == nil {
		addressServer, addressServerForURL, logLvl, pathFileStorage, connectionStringDB = flags.ReadFlags()
	} else {
		addressServer, addressServerForURL, logLvl, pathFileStorage, connectionStringDB =
			envVars.AddressServer, envVars.AddressServerForURL, envVars.LvlLogs, envVars.PathFileStorage, envVars.ConnectionStringDB

		if addressServer == "" || addressServerForURL == "" || logLvl == "" || pathFileStorage == "" || connectionStringDB == "" {
			addressServerFlag, addressServerForURLFlag, logLvlFlag, pathFileStorageFlag, connectionStringDBFlag := flags.ReadFlags()
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
			if connectionStringDB == "" {
				connectionStringDB = connectionStringDBFlag
			}
		}
	}

	addressServerForURL = CheckChangeBaseURL(addressServer, addressServerForURL)
	return addressServer, addressServerForURL, logLvl, pathFileStorage, connectionStringDB
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
