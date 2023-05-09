package app

import (
	"fmt"
	"github.com/GZ91/linkreduct/internal/api/http/server"
	"github.com/GZ91/linkreduct/internal/app/config"
)

var appLink *App

type App struct {
	config *config.Config
}

func New(config *config.Config) *App {
	if appLink == nil {
		appLink = &App{
			config,
		}
		return appLink
	}
	return appLink
}

func (r App) Run() {
	if err := server.Start(r.config); err != nil {
		fmt.Printf("%v \n", err)
	}
}
