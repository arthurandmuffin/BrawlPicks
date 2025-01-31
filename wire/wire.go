//go:build wireinject

package wire

import (
	"BrawlPicks/handlers"
	"BrawlPicks/repositories"
	"BrawlPicks/services"

	"github.com/google/wire"
)

func InitializeHandler() *handlers.Handler {
	panic(wire.Build(
		repositories.NewRepository,
		services.NewService,
		handlers.NewHandler,
	))
}
