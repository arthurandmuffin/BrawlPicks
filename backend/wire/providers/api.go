package providers

import (
	"BrawlPicks/internal/api"
	"BrawlPicks/internal/ctx"
	"BrawlPicks/internal/env"
	"BrawlPicks/repositories"
	"BrawlPicks/routes"
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
	wire.Bind(new(repositories.MapRankingRepositoryInterface), new(*repositories.MapRankingRepository)),
	MapRankingRepositorySet,
)

func NewApi(
	e *env.Env,
	ctx context.Context,
	swagger *routes.SwaggerRoute,
	mapRanking *routes.MapRankingRoute,
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
			mapRanking,
		},
	)
}
