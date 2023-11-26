package envs

import (
	"github.com/caarlos0/env/v6"
)

// EnvVars представляет структуру для хранения параметров окружения.
type EnvVars struct {
	AddressServer       string `env:"SERVER_ADDRESS"`
	AddressServerForURL string `env:"BASE_URL"`
	LvlLogs             string `env:"LOG_LEVEL"`
	PathFileStorage     string `env:"FILE_STORAGE_PATH"`
	ConnectionStringDB  string `env:"DATABASE_DSN"`
}

// ReadEnv считывает параметры из переменных окружения и возвращает объект EnvVars.
func ReadEnv() (*EnvVars, error) {
	envs := EnvVars{}

	// Считывание значений переменных окружения в объект EnvVars
	if err := env.Parse(&envs); err != nil {
		return nil, err
	}

	return &envs, nil
}
