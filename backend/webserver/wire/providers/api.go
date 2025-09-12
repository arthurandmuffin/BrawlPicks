package providers

import (
	"BrawlPicks/webserver/internal/api"
	"BrawlPicks/internal/ctx"
	env "BrawlPicks/webserver/config"
	"BrawlPicks/webserver/docs"
	"BrawlPicks/webserver/repositories"
	"BrawlPicks/webserver/routes"
	"BrawlPicks/webserver/scheduler"
	"BrawlPicks/webserver/services"
	"context"

	"github.com/google/wire"
)

var ApiSet = wire.NewSet(
	NewApi,
	ctx.GetGracefulShutdownCtx,
	SwaggerRouteSet,
	MapRankingRouteSet,
	MapRankingControllerSet,
	MapRankingServiceSet,
	wire.Bind(new(services.MapRankingDataServiceInterface), new(*services.MapRankingDataService)),
	MapRankingRepositorySet,
	wire.Bind(new(repositories.MapRankingRepositoryInterface), new(*repositories.MapRankingRepository)),
	MapRankingSchedulerSet,
)

func NewApi(
	e *env.Env,
	ctx context.Context,
	swagger *docs.SwaggerRoute,
	mapRankingRoute *routes.MapRankingRoute,
	mapRankingScheduler *scheduler.MapRankingScheduler,
) *api.Api {
	return api.NewApi(
		ctx,
		e.Api.Port,
		e.Api.Prefix,
		e.Api.Cors,
		e.Api.Debug,
		[]api.Route{
			swagger,
		},
		[]api.Route{
			mapRankingRoute,
		},
		[]scheduler.Scheduler{
			mapRankingScheduler,
		},
	)
}
