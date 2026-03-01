package providers

import (
	"net/http"
	"time"

	env "BrawlPicks/webserver/config"
	"BrawlPicks/webserver/controllers"
	"BrawlPicks/webserver/inference"
	"BrawlPicks/webserver/routes"
	"BrawlPicks/webserver/services"

	"github.com/google/wire"
)

var RecommendationRouteSet = wire.NewSet(
	routes.NewRecommendationRoute,
)

var RecommendationControllerSet = wire.NewSet(
	controllers.NewRecommendationController,
)

var RecommendationServiceSet = wire.NewSet(
	services.NewRecommendationService,
)

var InferenceClientSet = wire.NewSet(
	NewInferenceHTTPClient,
	inference.NewClient,
)

func NewInferenceHTTPClient(e *env.Env) *http.Client {
	timeout := 5 * time.Second
	if e.Inference != nil && e.Inference.TimeoutSeconds > 0 {
		timeout = time.Duration(e.Inference.TimeoutSeconds) * time.Second
	}
	return &http.Client{
		Timeout: timeout,
	}
}
