package main

import (
	internal "BrawlPicks/internal/logging"
	env "BrawlPicks/scraper/config"
	"BrawlPicks/scraper/wire"
	"flag"

	"github.com/sirupsen/logrus"
)

func main() {
	internal.SetupLogger()

	path := flag.String("d", "./scraper/config/default.yml", "")
	flag.Parse()

	cfg, err := env.Get(*path)
	if err != nil {
		logrus.Panic(err)
	}
	cfg.Info()

	app, err := wire.InitializeApp(cfg)
	if err != nil {
		logrus.Panic(err)
	}
	if err := app.Run(); err != nil {
		logrus.Panic(err)
	}
}
