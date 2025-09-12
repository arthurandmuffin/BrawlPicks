package main

import (
	internal "BrawlPicks/internal/logging"
	"BrawlPicks/webserver/config"
	"BrawlPicks/webserver/wire"
	"flag"

	"github.com/sirupsen/logrus"
)

func main() {
	internal.SetupLogger()

	path := flag.String("d", "./webserver/config/default.yml", "")
	flag.Parse()

	env, err := env.Get(*path)
	if err != nil {
		logrus.Panic(err)
	}
	env.Info()

	app, err := wire.InitializeApp(env)
	if err != nil {
		logrus.Panic(err)
	}
	app.Run()
}
