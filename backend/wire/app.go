//go:build wireinject

package wire

import (
	"BrawlPicks/app"
	"BrawlPicks/internal/env"
	"BrawlPicks/wire/providers"

	"github.com/google/wire"
)

func InitializeApp(*env.Env) (*app.App, error) {
	panic(wire.Build(
		providers.AppSet,
	))
}
