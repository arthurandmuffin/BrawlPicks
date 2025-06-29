//go:build wireinject

package wire

import (
	"BrawlPicks/handlers"
	"BrawlPicks/repositories"
	"BrawlPicks/services"

	"github.com/google/wire"
)

func InitializeBrawlStarsHandler(apiKey string) *handlers.BrawlStarsHandler {
	wire.Build(
		repositories.NewBrawlStarsRepository,
		services.NewBrawlStarsService,
		handlers.NewBrawlStarsHandler,
	)
	return nil
}

func InitializePLHandler() *handlers.PLHandler {
	wire.Build(
		repositories.NewPLRepository,
		services.NewPLService,
		handlers.NewPLHandler,
	)
	return nil
}
