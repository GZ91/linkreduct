package config

import "sync"

type Config struct {
	debug            bool
	addressServer    string
	addressServerURL string
	maxIterLen       int
	mutex            sync.Mutex
}

func New(debug bool, addressServer, addressServerURL string, maxIterRuneGen int) *Config {
	return &Config{
		debug:            debug,
		addressServer:    addressServer,
		maxIterLen:       maxIterRuneGen,
		addressServerURL: addressServerURL,
	}
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
