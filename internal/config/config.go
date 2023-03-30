package config

import "sync"

type Config struct {
	debug         bool
	addressServer string
	maxIterLen    int
	mutex         sync.Mutex
}

func New(debug bool, addressServer string, maxIterRuneGen int) Config {
	return Config{
		debug:         debug,
		addressServer: addressServer,
		maxIterLen:    maxIterRuneGen,
	}
}

func (r Config) GetDebug() bool {
	return r.debug
}

func (r Config) GetAddressServer() string {
	return r.addressServer
}

func (r Config) GetMaxIterLen() int {
	r.mutex.Lock()
	maxLenIter := r.maxIterLen
	r.mutex.Unlock()
	return maxLenIter
}
