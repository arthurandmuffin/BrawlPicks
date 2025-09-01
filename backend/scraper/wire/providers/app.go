package providers

import "github.com/google/wire"

var AppSet = wire.NewSet(
	CrawlerSet,
)
