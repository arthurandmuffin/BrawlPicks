package providers

import (
	"BrawlPicks/scheduler"

	"github.com/google/wire"
)

var MapRankingSchedulerSet = wire.NewSet(
	scheduler.NewMapRankingScheduler,
)
