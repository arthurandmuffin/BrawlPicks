//go:build wireinject

package wire

import (
	"BrawlPicks/webserver/app"
	env "BrawlPicks/webserver/config"
	"BrawlPicks/webserver/wire/providers"

	"github.com/google/wire"
)

func InitializeApp(*env.Env) (*app.App, error) {
	panic(wire.Build(
		providers.AppSet,
	))
}
