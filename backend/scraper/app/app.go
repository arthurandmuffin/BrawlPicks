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
	scraper services.ScraperServiceInterface
	closers []closer
}

func NewApp(ctx context.Context, scraper services.ScraperServiceInterface, closers ...closer) *App {
	return &App{
		ctx:     ctx,
		scraper: scraper,
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

	logrus.Info("starting-scraper")
	if err := a.scraper.SeedQueue(a.ctx); err != nil {
		return err
	}

	logrus.Info("starting-scrape-loop")
	a.scraper.Crawl(a.ctx)
	return nil
}
