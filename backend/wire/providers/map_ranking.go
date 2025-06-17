package providers

import (
	"BrawlPicks/controllers"
	"BrawlPicks/repositories"
	"BrawlPicks/routes"
	"BrawlPicks/services"
	"net/http"

	"github.com/google/wire"
)

var MapRankingRouteSet = wire.NewSet(
	routes.NewMapRankingRoute,
)

var MapRankingControllerSet = wire.NewSet(
	controllers.NewMapRankingController,
)

var MapRankingServiceSet = wire.NewSet(
	NewHttpClient,
	services.NewMapRankingDataService,
)

var MapRankingRepositorySet = wire.NewSet(
	repositories.NewMapRankingRepository,
)

func NewHttpClient() *http.Client {
	return &http.Client{}
}
