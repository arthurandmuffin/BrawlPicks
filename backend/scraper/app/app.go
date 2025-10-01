package app

import (
	"BrawlPicks/scraper/services"
	"context"

	"github.com/sirupsen/logrus"
)

type closer interface {
	Close() error
}

type App struct {
	ctx     context.Context
	crawler services.MatchDataCrawlerServiceInterface
	closers []closer
}

func NewApp(ctx context.Context, crawler services.MatchDataCrawlerServiceInterface, closers ...closer) *App {
	return &App{
		ctx:     ctx,
		crawler: crawler,
		closers: closers,
	}
}

func (a *App) Run() error {
	defer func() {
		for _, c := range a.closers {
			if c == nil {
				continue
			}
			if err := c.Close(); err != nil {
				logrus.WithError(err).Warn("failed-to-close-resource")
			}
		}
	}()

	queueLength, err := a.crawler.GetQueueLength(a.ctx)
	if err != nil {
		return err
	}
	logrus.WithField("queue_length", queueLength).Info("starting-scraper")

	if err := a.crawler.SeedQueue(a.ctx); err != nil {
		return err
	}

	logrus.Info("starting-scrape-loop")
	a.crawler.Crawl(a.ctx)
	return nil
}
