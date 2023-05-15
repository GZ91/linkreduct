package config

import (
	"github.com/GZ91/linkreduct/internal/storage/postgresql/postgresqlconfig"
	"sync"
)

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

func (r *Config) ConfigureDBPostgresql(address, user, password, dbname string) {
	r.configDB = postgresqlconfig.New(address, user, password, dbname)
}

func (r *Config) GetConfDB() *postgresqlconfig.ConfigDB {
	return r.configDB
}

func (r *Config) GetAddressServerURL() string {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	return r.addressServerURL
}

func (r *Config) GetDebug() bool {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	return r.debug
}

func (r *Config) GetAddressServer() string {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	return r.addressServer
}

func (r *Config) GetMaxIterLen() int {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	return r.maxIterLen
}

func (r *Config) GetStartLenShortLink() int {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	return r.startLenShortLink
}

func (r *Config) GetNameFileStorage() string {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	return r.fileStorage
}
