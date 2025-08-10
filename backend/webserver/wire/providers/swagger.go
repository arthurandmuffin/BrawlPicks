package providers

import (
	env "BrawlPicks/internal/config"
	"BrawlPicks/webserver/docs"

	"github.com/google/wire"
)

var SwaggerRouteSet = wire.NewSet(
	NewSwaggerRoute,
)

func NewSwaggerRoute(e *env.Env) *docs.SwaggerRoute {
	return docs.NewSwaggerRoute(e.Api.Prefix)
}
