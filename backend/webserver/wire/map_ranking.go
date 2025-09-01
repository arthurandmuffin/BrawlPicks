//go:build wireinject

package wire

import (
	env "BrawlPicks/webserver/config"
	"BrawlPicks/webserver/controllers"
	"BrawlPicks/webserver/repositories"
	"BrawlPicks/webserver/routes"
	"BrawlPicks/webserver/services"
	"BrawlPicks/webserver/wire/providers"

	"github.com/google/wire"
)

func InitializeMapRankingRoute(controller *controllers.MapRankingController) *routes.MapRankingRoute {
	panic(wire.Build(
		providers.MapRankingRouteSet,
	))
}

func InitializeMapRankingController(sv services.MapRankingDataServiceInterface) *controllers.MapRankingController {
	panic(wire.Build(
		providers.MapRankingControllerSet,
	))
}

func InitializeMapRankingService(
	e *env.Env,
	repo repositories.MapRankingRepositoryInterface,
) *services.MapRankingDataService {
	panic(wire.Build(
		providers.MapRankingServiceSet,
	))
}

func InitializeMapRankingRepository(e *env.Env) *repositories.MapRankingRepository {
	panic(wire.Build(
		providers.MapRankingRepositorySet,
	))
}
