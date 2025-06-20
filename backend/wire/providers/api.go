package providers

import (
	"BrawlPicks/internal/api"
	"BrawlPicks/internal/ctx"
	"BrawlPicks/internal/env"
	"BrawlPicks/repositories"
	"BrawlPicks/routes"
	"BrawlPicks/scheduler"
	"BrawlPicks/services"
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
	swagger *routes.SwaggerRoute,
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
