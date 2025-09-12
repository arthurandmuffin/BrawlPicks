package providers

import (
	"BrawlPicks/webserver/app"

	"github.com/google/wire"
)

var AppSet = wire.NewSet(
	app.NewApp,
	ApiSet,
)
