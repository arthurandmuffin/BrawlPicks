package providers

import (
	"BrawlPicks/webserver/scheduler"

	"github.com/google/wire"
)

var MapRankingSchedulerSet = wire.NewSet(
	scheduler.NewMapRankingScheduler,
)
