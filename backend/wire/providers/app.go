package providers

import (
	"BrawlPicks/app"

	"github.com/google/wire"
)

var AppSet = wire.NewSet(
	app.NewApp,
	ApiSet,
)
