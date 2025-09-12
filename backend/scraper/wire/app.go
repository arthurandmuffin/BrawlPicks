//go:build wireinject

package wire

import (
	"BrawlPicks/scraper/app"
	env "BrawlPicks/scraper/config"
	"BrawlPicks/scraper/wire/providers"

	"github.com/google/wire"
)

func InitializeApp(*env.Env) (*app.App, error) {
	panic(wire.Build(
		providers.AppSet,
	))
}
