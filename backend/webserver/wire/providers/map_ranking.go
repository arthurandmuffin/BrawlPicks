package providers

import (
	"net/http"

	"github.com/google/wire"

	ahttp "BrawlPicks/internal/http"
	"BrawlPicks/webserver/controllers"
	"BrawlPicks/webserver/repositories"
	"BrawlPicks/webserver/routes"
	"BrawlPicks/webserver/services"
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

func NewHttpClient() *ahttp.Client {
	return ahttp.NewClient(&http.Client{})
}
