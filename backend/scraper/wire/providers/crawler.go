package providers

import (
	"BrawlPicks/internal/ctx"
	ahttp "BrawlPicks/internal/http"
	"BrawlPicks/scraper/app"
	env "BrawlPicks/scraper/config"
	"BrawlPicks/scraper/repositories"
	"BrawlPicks/scraper/services"
	"context"
	"net/http"
	"time"

	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
)

var CrawlerSet = wire.NewSet(
	ctx.GetGracefulShutdownCtx,
	NewHttpClient,
	NewRedisClient,
	repositories.NewProcessedRecordRepository,
	NewBattleLogRepository,
	NewSynergyCounterRepository,
	NewCrawlerService,
	NewScraperApp,
)

func NewHttpClient() *ahttp.Client {
	return ahttp.NewClient(&http.Client{})
}

func NewRedisClient(e *env.Env) *redis.Client {
	if e.Redis != nil && e.Redis.Credentials != nil && e.Redis.Credentials.MasterName != "" {
		return redis.NewFailoverClient(&redis.FailoverOptions{
			MasterName:    e.Redis.Credentials.MasterName,
			SentinelAddrs: []string{e.Redis.Credentials.Address},
			Password:      e.Redis.Credentials.Password,
		})
	}

	address := ""
	password := ""
	if e.Redis != nil && e.Redis.Credentials != nil {
		address = e.Redis.Credentials.Address
		password = e.Redis.Credentials.Password
	}

	return redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
	})
}

func NewBattleLogRepository(e *env.Env) *repositories.BattleLogRepository {
	cfg := e.Storage.BattleLog
	return repositories.NewBattleLogRepository(
		cfg.MaxRows,
		time.Duration(cfg.FlushSeconds)*time.Second,
		cfg.Dir,
	)
}

func NewSynergyCounterRepository(e *env.Env) (*repositories.SynergyCounterRepository, error) {
	cfg := e.Storage.Synergy
	return repositories.NewSynergyCounterRepository(
		e,
		time.Duration(cfg.FlushSeconds)*time.Second,
		cfg.RetentionDays,
		cfg.Dir,
	)
}

func NewCrawlerService(
	e *env.Env,
	client *ahttp.Client,
	rProcessed *repositories.ProcessedRecordRepository,
	rBattleLog *repositories.BattleLogRepository,
	rSynergy *repositories.SynergyCounterRepository,
) *services.MatchDataCrawlerService {
	return services.NewMatchDataCrawlerService(e, client, rProcessed, rBattleLog, rSynergy)
}

func NewScraperApp(
	ctx context.Context,
	crawler *services.MatchDataCrawlerService,
	battleLog *repositories.BattleLogRepository,
	synergy *repositories.SynergyCounterRepository,
	redisClient *redis.Client,
) *app.App {
	return app.NewApp(ctx, crawler, battleLog, synergy, redisClient)
}
