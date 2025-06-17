package app

import "BrawlPicks/internal/api"

type App struct {
	api *api.Api
}

func NewApp(api *api.Api) *App {
	return &App{
		api: api,
	}
}

func (a *App) Run() {
	a.api.Run()
}
