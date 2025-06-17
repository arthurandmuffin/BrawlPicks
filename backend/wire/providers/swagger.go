package providers

import (
	"BrawlPicks/internal/env"
	"BrawlPicks/routes"

	"github.com/google/wire"
)

var SwaggerRouteSet = wire.NewSet(
	NewSwaggerRoute,
)

func NewSwaggerRoute(e *env.Env) *routes.SwaggerRoute {
	return routes.NewSwaggerRoute(e.Api.Prefix)
}
