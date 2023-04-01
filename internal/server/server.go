package server

import (
	"errors"
	"github.com/GZ91/linkreduct/internal/config"
	"github.com/GZ91/linkreduct/internal/handlers"
	"net/http"
)

func Start(conf *config.Config) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(r.(string))
		}
	}()
	mux := &http.ServeMux{}
	handlers.InstallConfig(conf)
	mux.Handle("/", handlers.Conveyor(http.HandlerFunc(handlers.MethodGet), handlers.MethodPost))

	return http.ListenAndServe(conf.GetAddressServer(), mux)
}
