package main

import (
	"fmt"
	"os"
	"path/filepath"

	"BrawlPicks/cli/app"
	"BrawlPicks/cli/brawlers"
	"BrawlPicks/cli/client"
	"BrawlPicks/cli/config"
)

func main() {
	configPath := filepath.Join("config", "default.yml")

	cfg, err := config.Load(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load cli config: %v\n", err)
		os.Exit(1)
	}

	brawlerCatalog, err := brawlers.Load(filepath.Join("..", "brawlers.json"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load brawler catalog: %v\n", err)
		os.Exit(1)
	}

	httpClient := client.New(cfg.Server.BaseURL, cfg.Server.RecommendPath)
	terminalApp := app.New(cfg, httpClient, brawlerCatalog)

	if err := terminalApp.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "cli exited with error: %v\n", err)
		os.Exit(1)
	}
}
