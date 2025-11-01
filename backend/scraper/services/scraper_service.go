package services

import (
	ahttp "BrawlPicks/internal/http"
	env "BrawlPicks/scraper/config"
	"BrawlPicks/scraper/models"
	"BrawlPicks/scraper/repositories"
	"BrawlPicks/scraper/services/upstream"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
)

type ScraperServiceInterface interface {
	Crawl(ctx context.Context)
	SeedQueue(ctx context.Context) error
	FetchTopPlayerTags(ctx context.Context) ([]string, error)
}

type ScraperService struct {
	e                             *env.Env
	client                        *ahttp.Client
	monitor                       MonitorServiceInterface
	targetHost                    string
	token                         string
	qps, burst                    int
	ioWorkerCount, cpuWorkerCount int
	queueBatch                    int
	queueLow, queueHigh           int
	seedThreshold                 int64
	seedCooldown                  time.Duration
	maxBattleAge                  time.Duration
	lastSeededAt                  time.Time
	ioJobQueue                    chan string
	cpuJobQueue                   chan *models.Battle
	rProcessed                    repositories.ProcessedRecordRepositoryInterface
	rBattleLog                    repositories.BattleLogRepositoryInterface
	rSynergy                      repositories.SynergyCounterRepositoryInterface
}

func NewScraperService(
	e *env.Env,
	client *ahttp.Client,
	monitor MonitorServiceInterface,
	rProcessed repositories.ProcessedRecordRepositoryInterface,
	rBattleLog repositories.BattleLogRepositoryInterface,
	rSynergy repositories.SynergyCounterRepositoryInterface,
) *ScraperService {
	service := &ScraperService{
		e:          e,
		client:     client,
		monitor:    monitor,
		rProcessed: rProcessed,
		rBattleLog: rBattleLog,
		rSynergy:   rSynergy,
	}
	if e != nil && e.Brawl != nil {
		service.targetHost = e.Brawl.BattleLogEndpoint
		service.token = e.Brawl.Key
	}
	if e != nil && e.Scraper != nil {
		if e.Scraper.RateLimit != nil {
			service.qps = e.Scraper.RateLimit.QPS
			service.burst = e.Scraper.RateLimit.Burst
		}
		if e.Scraper.Workers != nil {
			service.ioWorkerCount = e.Scraper.Workers.IO
			service.cpuWorkerCount = e.Scraper.Workers.CPU
		}
		if e.Scraper.Queue != nil {
			service.queueBatch = e.Scraper.Queue.Batch
			service.queueLow = e.Scraper.Queue.Low
			service.queueHigh = e.Scraper.Queue.High
			service.ioJobQueue = make(chan string, e.Scraper.Queue.ChannelSize)
			service.cpuJobQueue = make(chan *models.Battle, e.Scraper.Queue.ChannelSize)
		}
		if e.Scraper.Processing != nil {
			service.maxBattleAge = time.Duration(e.Scraper.Processing.MaxBattleAgeDays) * 24 * time.Hour
		}
		if e.Scraper.Seeding != nil {
			service.seedThreshold = e.Scraper.Seeding.Threshold
			service.seedCooldown = time.Duration(e.Scraper.Seeding.CooldownSeconds) * time.Second
		}
	}
	return service
}

func (s *ScraperService) Crawl(ctx context.Context) {
	lim := rate.NewLimiter(rate.Limit(s.qps), s.burst)

	if s.monitor != nil {
		go s.monitor.Start(ctx)
	}
	for i := 0; i < s.ioWorkerCount; i++ {
		go s.ioWorker(ctx, lim, s.ioJobQueue, s.cpuJobQueue)
	}
	for i := 0; i < s.cpuWorkerCount; i++ {
		go s.cpuWorker(ctx, s.cpuJobQueue)
	}
	go s.jobFeeder(ctx, s.ioJobQueue)

	<-ctx.Done()
}

func (s *ScraperService) jobFeeder(ctx context.Context, ioJobQueue chan<- string) {
	for {
		if ctx.Err() != nil {
			return
		}
		if err := s.SeedQueue(ctx); err != nil {
			logrus.WithField("err", err).Warn("failed-to-seed-top-players")
			if !sleepOrDone(ctx, 5*time.Second) {
				return
			}
			continue
		}

		if len(ioJobQueue) > s.queueLow {
			select {
			case <-time.After(time.Second):
				continue
			case <-ctx.Done():
				return
			}
		}
		for len(ioJobQueue) < s.queueHigh {
			res, err := s.rProcessed.PopPlayersFromQueue(ctx, s.queueBatch)
			if err != nil {
				logrus.WithField("err", err).Warn("failed-to-popPlayersFromQueue")
				if !sleepOrDone(ctx, 5*time.Second) {
					return
				}
				continue
			}
			if len(res) == 0 {
				if err := s.SeedQueue(ctx); err != nil {
					logrus.WithField("err", err).Warn("failed-to-reseed-top-players")
				}
				if !sleepOrDone(ctx, time.Second) {
					return
				}
				break
			}

			tags, err := s.rProcessed.AddPlayersToBF(ctx, res)
			if err != nil {
				logrus.WithField("err", err).Warn("failed-to-AddPlayersToBF")
				continue
			}
			for _, tag := range tags {
				select {
				case ioJobQueue <- tag:
				case <-ctx.Done():
					return
				}
			}
		}
	}
}

func (s *ScraperService) SeedQueue(ctx context.Context) error {
	if s.rProcessed == nil {
		return nil
	}

	queueLength, err := s.rProcessed.GetPlayerQueueLength(ctx)
	if err != nil {
		return err
	}
	if queueLength >= s.seedThreshold {
		return nil
	}
	if !s.lastSeededAt.IsZero() && time.Since(s.lastSeededAt) < s.seedCooldown {
		logrus.WithFields(logrus.Fields{
			"queue_length": queueLength,
			"threshold":    s.seedThreshold,
			"cooldown":     s.seedCooldown.String(),
		}).Info("skip-seeding-cooldown-active")
		return nil
	}

	logrus.WithFields(logrus.Fields{
		"queue_length": queueLength,
		"threshold":    s.seedThreshold,
	}).Info("seeding-triggered")
	if err := s.seedTopPlayers(ctx); err != nil {
		return err
	}
	s.lastSeededAt = time.Now()
	return nil
}

func (s *ScraperService) ioWorker(ctx context.Context, lim *rate.Limiter, ioJobQueue <-chan string, cpuJobQueue chan<- *models.Battle) {
	for {
		select {
		case <-ctx.Done():
			return
		case tag := <-ioJobQueue:
			if err := lim.Wait(ctx); err != nil {
				return
			}

			battles, err := s.getBattleLog(ctx, tag)
			if err != nil {
				logrus.WithField("err", err).Warn("failed-to-getBattleLog")
				continue
			}

			flags, err := s.rProcessed.AddGamesToBF(ctx, battles)
			if err != nil {
				logrus.WithField("err", err).Warn("failed-to-AddGamesToBF")
				continue
			}
			newGames := 0

			for i, flag := range flags {
				if flag {
					newGames++
					battle := battles[i]
					newTags := battle.GetTags(tag)
					if err := s.rProcessed.AddPlayersToQueue(ctx, newTags); err != nil {
						logrus.WithField("err", err).Warn("failed-to-AddPlayersToQueue")
					}
					transformed := battle.Transform(tag)
					if transformed == nil {
						continue
					}
					select {
					case cpuJobQueue <- transformed:
					case <-ctx.Done():
						return
					}
				}
			}
			if s.monitor != nil {
				s.monitor.AddGames(newGames)
			}
		}
	}
}

func (s *ScraperService) cpuWorker(ctx context.Context, cpuJobQueue <-chan *models.Battle) {
	for {
		select {
		case <-ctx.Done():
			return
		case battle := <-cpuJobQueue:
			if battle == nil {
				continue
			}
			s.rSynergy.RecordBattle(battle)
			if err := s.rBattleLog.WriteBattleLog(*battle); err != nil {
				logrus.WithField("err", err).Warn("failed-to-WriteBattleLog")
			}
		}
	}
}

func (s *ScraperService) seedTopPlayers(ctx context.Context) error {
	tags, err := s.FetchTopPlayerTags(ctx)
	if err != nil {
		return err
	}
	logrus.Info(fmt.Sprintf("Seeding top players: fetched %d", len(tags)))
	return s.rProcessed.AddPlayersToQueue(ctx, tags)
}

func (s *ScraperService) FetchTopPlayerTags(ctx context.Context) ([]string, error) {
	logrus.Info("fetching-top-players")
	topPlayers := new(upstream.TopPlayersResponse)
	req, err := s.newTopPlayersRequest()
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := s.client.Do(req, topPlayers); err != nil {
		return nil, err
	}

	tags := make([]string, 0, len(topPlayers.Players))
	for _, player := range topPlayers.Players {
		if player.Tag == "" {
			continue
		}
		tags = append(tags, player.Tag)
	}
	return tags, nil
}

func (s *ScraperService) getBattleLog(ctx context.Context, tag string) (battles []*upstream.BattleLogBattle, err error) {
	battleLog := new(upstream.BattleLogResponse)
	req, err := s.newBattleLogRequest(tag)
	if err != nil {
		return
	}
	req = req.WithContext(ctx)
	if err = s.client.Do(req, battleLog); err != nil {
		return
	}

	battles = make([]*upstream.BattleLogBattle, 0, len(battleLog.Battles))
	cutoffTime := time.Time{}
	if s.maxBattleAge > 0 {
		cutoffTime = time.Now().UTC().Add(-s.maxBattleAge)
	}
	for _, battle := range battleLog.Battles {
		if battle == nil || !battle.IsValid() {
			continue
		}
		if !cutoffTime.IsZero() && battle.Time.Before(cutoffTime) {
			continue
		}
		battles = append(battles, battle)
	}
	return battles, nil
}

func (s *ScraperService) newBattleLogRequest(tag string) (req *http.Request, err error) {
	endpoint := s.targetHost
	if strings.Contains(endpoint, "%s") {
		endpoint = fmt.Sprintf(endpoint, url.QueryEscape(tag))
	} else {
		endpoint = fmt.Sprintf("%s/players/%s/battlelog", strings.TrimRight(endpoint, "/"), url.QueryEscape(tag))
	}

	return ahttp.NewBearerAuthRequest(
		http.MethodGet,
		endpoint,
		s.token, nil,
	)
}

func (s *ScraperService) newTopPlayersRequest() (req *http.Request, err error) {
	return ahttp.NewBearerAuthRequest(
		http.MethodGet,
		s.e.Brawl.TopPlayersEndpoint,
		s.token, nil,
	)
}

func sleepOrDone(ctx context.Context, wait time.Duration) bool {
	select {
	case <-time.After(wait):
		return true
	case <-ctx.Done():
		return false
	}
}
