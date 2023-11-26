package initializing

import (
	"strings"

	"github.com/GZ91/linkreduct/internal/app/config"
	"github.com/GZ91/linkreduct/internal/app/initializing/envs"
	"github.com/GZ91/linkreduct/internal/app/initializing/flags"
	"github.com/GZ91/linkreduct/internal/app/logger"
	"go.uber.org/zap"
)

// Configuration создает и возвращает объект конфигурации приложения.
// Считывает параметры из окружения или флагов командной строки.
func Configuration() *config.Config {
	// Чтение параметров из окружения или флагов
	addressServer, addressServerForURL, logLvl, pathFileStorage, connectionStringDB := ReadParams()
	// Инициализация логгера с уровнем логгирования
	logger.Initializing(logLvl)
	// Создание объекта конфигурации
	conf := config.New(false, addressServer, addressServerForURL, 10, 5, pathFileStorage)
	// Настройка подключения к базе данных PostgreSQL
	conf.ConfigureDBPostgresql(connectionStringDB)
	return conf
}

// ReadParams считывает параметры из окружения или флагов командной строки.
// Возвращает адрес сервера, адрес сервера для URL, уровень логгирования, путь к файлу хранилища и строку подключения к базе данных PostgreSQL.
func ReadParams() (string, string, string, string, string) {
	// Чтение параметров из окружения
	envVars, err := envs.ReadEnv()
	if err != nil {
		logger.Log.Error("error when reading environment variables", zap.String("error", err.Error()))
	}

	var addressServer, addressServerForURL, logLvl, pathFileStorage, connectionStringDB string

	// Если параметры не были считаны из окружения
	if envVars == nil {
		// Считывание параметров из флагов командной строки
		addressServer, addressServerForURL, logLvl, pathFileStorage, connectionStringDB = flags.ReadFlags()
	} else {
		// Использование параметров из окружения, если они были считаны
		addressServer, addressServerForURL, logLvl, pathFileStorage, connectionStringDB =
			envVars.AddressServer, envVars.AddressServerForURL, envVars.LvlLogs, envVars.PathFileStorage, envVars.ConnectionStringDB

		// Проверка и заполнение отсутствующих параметров, если они не были заданы
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

	// Проверка и изменение базового URL в соответствии с адресом сервера
	addressServerForURL = CheckChangeBaseURL(addressServer, addressServerForURL)
	return addressServer, addressServerForURL, logLvl, pathFileStorage, connectionStringDB
}

// CheckChangeBaseURL проверяет и изменяет базовый URL в соответствии с адресом сервера.
func CheckChangeBaseURL(addressServer, addressServerURL string) string {
	strAddress := strings.Split(addressServerURL, ":")
	var port string
	if len(strAddress) == 3 {
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
