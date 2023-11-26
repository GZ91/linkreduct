package config

import (
	"sync"

	"github.com/GZ91/linkreduct/internal/storage/postgresql/postgresqlconfig"
)

// Config представляет структуру для хранения конфигурационных параметров приложения.
type Config struct {
	debug             bool
	addressServer     string
	addressServerURL  string
	maxIterLen        int
	startLenShortLink int
	fileStorage       string
	configDB          *postgresqlconfig.ConfigDB
	mutex             sync.Mutex
}

// New создает и возвращает новый объект Config с заданными параметрами.
func New(debug bool, addressServer, addressServerURL string, maxIterRuneGen int, startLenShortLink int, fileStorage string) *Config {
	return &Config{
		debug:             debug,
		addressServer:     addressServer,
		maxIterLen:        maxIterRuneGen,
		addressServerURL:  addressServerURL,
		startLenShortLink: startLenShortLink,
		fileStorage:       fileStorage,
	}
}

// ConfigureDBPostgresql создает и настраивает конфигурацию для подключения к базе данных PostgreSQL.
func (r *Config) ConfigureDBPostgresql(StringServer string) {
	r.configDB = postgresqlconfig.New(StringServer)
}

// GetConfDB возвращает конфигурацию для подключения к базе данных PostgreSQL.
func (r *Config) GetConfDB() *postgresqlconfig.ConfigDB {
	return r.configDB
}

// GetAddressServerURL возвращает URL-адрес сервера.
func (r *Config) GetAddressServerURL() string {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	return r.addressServerURL
}

// GetDebug возвращает флаг отладки.
func (r *Config) GetDebug() bool {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	return r.debug
}

// GetAddressServer возвращает адрес сервера.
func (r *Config) GetAddressServer() string {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	return r.addressServer
}

// GetMaxIterLen возвращает максимальное количество итераций для генерации коротких ссылок.
func (r *Config) GetMaxIterLen() int {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	return r.maxIterLen
}

// GetStartLenShortLink возвращает начальную длину коротких ссылок.
func (r *Config) GetStartLenShortLink() int {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	return r.startLenShortLink
}

// GetNameFileStorage возвращает имя файла хранилища данных.
func (r *Config) GetNameFileStorage() string {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	return r.fileStorage
}
