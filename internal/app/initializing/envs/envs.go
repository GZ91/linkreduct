package envs

import "github.com/caarlos0/env/v6"

type EnvVars struct {
	AddressServer       string `env:"SERVER_ADDRESS"`
	AddressServerForURL string `env:"BASE_URL"`
	LvlLogs             string `env:"LOG_LEVEL"`
	PathFileStorage     string `env:"FILE_STORAGE_PATH"`
}

func ReadEnv() (string, string, string, string) {
	envs := EnvVars{}
	if err := env.Parse(&envs); err != nil {
		panic(err)
	}

	return envs.AddressServer, envs.AddressServerForURL, envs.LvlLogs, envs.PathFileStorage
}
