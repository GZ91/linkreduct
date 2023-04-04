package server

import (
	"errors"
	"github.com/GZ91/linkreduct/internal/config"
	"github.com/GZ91/linkreduct/internal/handlers"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strings"
	"sync"
)

func Start(conf *config.Config) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(r.(string))
		}
	}()
	handlers.InstallConfig(conf)
	WG := &sync.WaitGroup{}
	{
		router := chi.NewRouter()
		router.Post("/", handlers.MethodPost)

		WG.Add(1)
		go func(wg *sync.WaitGroup) {
			http.ListenAndServe(conf.GetAddressServer(), router)
			wg.Done()
		}(WG)
	}
	{
		router := chi.NewRouter()
		router.Get("/{id}", handlers.MethodGet)
		router.Post("/", handlers.MethodPost)
		StrNotPrefix := strings.TrimPrefix(conf.GetAddressServerURL(), "http://")
		adressServer := strings.Split(StrNotPrefix, "/")[0]
		WG.Add(1)
		go func(wg *sync.WaitGroup) {
			http.ListenAndServe(adressServer, router)
			wg.Done()
		}(WG)
	}
	WG.Wait()
	return nil
}
