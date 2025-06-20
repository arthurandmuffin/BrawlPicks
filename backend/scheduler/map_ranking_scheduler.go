package scheduler

import (
	"BrawlPicks/services"
	"context"
	"time"

	"github.com/sirupsen/logrus"
)

type MapRankingScheduler struct {
	svMapRanking services.MapRankingDataServiceInterface
}

func NewMapRankingScheduler(svMapRanking services.MapRankingDataServiceInterface) *MapRankingScheduler {
	return &MapRankingScheduler{
		svMapRanking: svMapRanking,
	}
}

func (s *MapRankingScheduler) Start(c context.Context) {
	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()

	for {
		select {
		case time := <-ticker.C:
			err := s.CheckForNewData()
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"time": time,
					"err":  err,
				})
			}
		case <-c.Done():
			return
		}
	}
}

func (s *MapRankingScheduler) CheckForNewData() (err error) {
	_, err = s.svMapRanking.RefreshRankings(false)
	return
}
