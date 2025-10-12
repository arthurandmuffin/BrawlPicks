package services

import (
	env "BrawlPicks/scraper/config"
	"BrawlPicks/scraper/repositories"
	"context"
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"
)

type MonitorServiceInterface interface {
	Start(ctx context.Context)
	AddGames(count int)
}

type MonitorService struct {
	heartbeat  time.Duration
	games      atomic.Int64
	rProcessed repositories.ProcessedRecordRepositoryInterface
}

func NewMonitorService(
	e *env.Env,
	rProcessed repositories.ProcessedRecordRepositoryInterface,
) *MonitorService {
	heartbeat := time.Duration(e.Monitor.HeartbeatInterval) * time.Second

	return &MonitorService{
		heartbeat:  heartbeat,
		rProcessed: rProcessed,
	}
}

func (s *MonitorService) Start(ctx context.Context) {
	if s == nil || s.rProcessed == nil || s.heartbeat <= 0 {
		return
	}

	queueLength, err := s.rProcessed.GetPlayerQueueLength(ctx)
	if err != nil {
		logrus.WithField("err", err).Warn("failed-to-read-startup-queue-length")
	} else {
		logrus.WithField("redis_queue_size", queueLength).Info("starting-scraper-monitor")
	}

	ticker := time.NewTicker(s.heartbeat)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logrus.Info("monitor-shutdown")
			return
		case <-ticker.C:
			queueLength, err := s.rProcessed.GetPlayerQueueLength(ctx)
			if err != nil {
				logrus.WithField("err", err).Warn("failed-to-read-monitor-queue-length")
				continue
			}

			logrus.WithFields(logrus.Fields{
				"redis_queue_size": queueLength,
				"games_added":      s.games.Swap(0),
			}).Info("scraper-heartbeat")
		}
	}
}

func (s *MonitorService) AddGames(count int) {
	if s == nil || count <= 0 {
		return
	}
	s.games.Add(int64(count))
}
