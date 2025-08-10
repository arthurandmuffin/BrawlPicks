package services

import (
	ahttp "BrawlPicks/internal/http"
	"BrawlPicks/internal/config"
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

const (
	defaultCrawlerQPS         = 5
	defaultCrawlerBurst       = 10
	defaultCrawlerIOWorkers   = 4
	defaultCrawlerCPUWorkers  = 2
	defaultCrawlerQueueBatch  = 25
	defaultCrawlerQueueLow    = 25
	defaultCrawlerQueueHigh   = 100
	defaultCrawlerChannelSize = 128
)

type MatchDataCrawlerServiceInterface interface {
	Crawl(ctx context.Context)
	SeedTopPlayers(ctx context.Context) error
	FetchTopPlayerTags(ctx context.Context) ([]string, error)
}

type MatchDataCrawlerService struct {
	e                             *env.Env
	client                        *ahttp.Client
	targetHost                    string
	token                         string
	qps, burst                    int
	ioWorkerCount, cpuWorkerCount int
	queueLimit                    int
	queueBatch                    int
	queueLow, queueHigh           int
	ioJobQueue                    chan string
	cpuJobQueue                   chan *models.Battle
	rProcessed                    repositories.ProcessedRecordRepositoryInterface
	rBattleLog                    repositories.BattleLogRepositoryInterface
	rSynergy                      repositories.SynergyCounterRepositoryInterface
}

func NewMatchDataCrawlerService(e *env.Env, client *ahttp.Client) *MatchDataCrawlerService {
	service := &MatchDataCrawlerService{
		e:          e,
		queueLimit: 10,
		client:     client,
		qps:        defaultCrawlerQPS,
		burst:      defaultCrawlerBurst,
		ioWorkerCount:  defaultCrawlerIOWorkers,
		cpuWorkerCount: defaultCrawlerCPUWorkers,
		queueBatch:     defaultCrawlerQueueBatch,
		queueLow:       defaultCrawlerQueueLow,
		queueHigh:      defaultCrawlerQueueHigh,
		ioJobQueue:     make(chan string, defaultCrawlerChannelSize),
		cpuJobQueue:    make(chan *models.Battle, defaultCrawlerChannelSize),
	}
	if e != nil && e.Brawl != nil {
		service.targetHost = e.Brawl.BattleLogEndpoint
		service.token = e.Brawl.Key
	}
	return service
}

func (s *MatchDataCrawlerService) Crawl(ctx context.Context) {
	lim := rate.NewLimiter(rate.Limit(s.qps), s.burst)

	for i := 0; i < s.ioWorkerCount; i++ {
		go s.ioWorker(ctx, lim, s.ioJobQueue, s.cpuJobQueue)
	}
	for i := 0; i < s.cpuWorkerCount; i++ {
		go s.cpuWorker(ctx, s.cpuJobQueue)
	}
	go s.jobFeeder(ctx, s.ioJobQueue)

	<-ctx.Done()
}

func (s *MatchDataCrawlerService) jobFeeder(ctx context.Context, ioJobQueue chan<- string) {
	for {
		if ctx.Err() != nil {
			return
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

func (s *MatchDataCrawlerService) ioWorker(ctx context.Context, lim *rate.Limiter, ioJobQueue <-chan string, cpuJobQueue chan<- *models.Battle) {
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

			for i, flag := range flags {
				if flag {
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
		}
	}
}

func (s *MatchDataCrawlerService) cpuWorker(ctx context.Context, cpuJobQueue <-chan *models.Battle) {
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

func (s *MatchDataCrawlerService) SeedTopPlayers(ctx context.Context) error {
	tags, err := s.FetchTopPlayerTags(ctx)
	if err != nil {
		return err
	}

	newTags, err := s.rProcessed.AddPlayersToBF(ctx, tags)
	if err != nil {
		return err
	}
	return s.rProcessed.AddPlayersToQueue(ctx, newTags)
}

func (s *MatchDataCrawlerService) FetchTopPlayerTags(ctx context.Context) ([]string, error) {
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

func (s *MatchDataCrawlerService) getBattleLog(ctx context.Context, tag string) (battles []*upstream.BattleLogBattle, err error) {
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
	for _, battle := range battleLog.Battles {
		if battle == nil || !battle.IsValid() {
			continue
		}
		battles = append(battles, battle)
	}
	return battles, nil
}

func (s *MatchDataCrawlerService) newBattleLogRequest(tag string) (req *http.Request, err error) {
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

func (s *MatchDataCrawlerService) newTopPlayersRequest() (req *http.Request, err error) {
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
