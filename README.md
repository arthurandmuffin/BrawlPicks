# BrawlPicks

Simple draft-assistant passion project for Brawl Stars.

![System Diagram](/Users/arthurhuang/BrawlPicks/BrawlPicks/assets/system_diagram.png)

## What It Is

- scraper collects battle logs and derives stats
- ML pipeline transforms data, trains models, and serves inference
- webserver exposes the product-facing API
- CLI is the first thin user interface over the backend
- Full frontend in progress, to come in the future!

## Repo Roadmap

- `backend/scraper`: data collection and persistence
- `backend/webserver`: API and product orchestration
- `ml/`: transformer, trainer, and inference service
- `cli/`: terminal interface
