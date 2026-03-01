package providers

import (
	"BrawlPicks/internal/ctx"
	env "BrawlPicks/webserver/config"
	"BrawlPicks/webserver/docs"
	"BrawlPicks/webserver/internal/api"
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
	RecommendationRouteSet,
	RecommendationControllerSet,
	RecommendationServiceSet,
	wire.Bind(new(services.RecommendationServiceInterface), new(*services.RecommendationService)),
	InferenceClientSet,
	MapRankingRepositorySet,
	wire.Bind(new(repositories.MapRankingRepositoryInterface), new(*repositories.MapRankingRepository)),
	MapRankingSchedulerSet,
)

func NewApi(
	e *env.Env,
	ctx context.Context,
	swagger *docs.SwaggerRoute,
	mapRankingRoute *routes.MapRankingRoute,
	recommendationRoute *routes.RecommendationRoute,
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
			recommendationRoute,
		},
		[]scheduler.Scheduler{
			mapRankingScheduler,
		},
	)
}
