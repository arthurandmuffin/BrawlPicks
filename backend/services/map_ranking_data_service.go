package services

import (
	"BrawlPicks/internal/env"
	"BrawlPicks/internal/stats"
	"BrawlPicks/models"
	"BrawlPicks/repositories"
	"BrawlPicks/services/upstream"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"sync"

	"maps"

	"github.com/sirupsen/logrus"
)

type MapRankingDataServiceInterface interface {
	InitMapRankings() (err error)
	RefreshRankings(force bool) (updated bool, err error)
	GetMapRankings(rank models.Rank, mapNames []string) (rankings []*models.MapRanking, err error)
	NewDataAvailable() (newData bool, err error, newTime int)
	FetchRawMapData(rank models.Rank) (data []*models.MapData, err error)
}

type MapRankingDataService struct {
	e                      *env.Env
	client                 *http.Client
	rMap                   repositories.MapRankingRepositoryInterface
	mapRankings            map[models.Rank][]*models.MapRanking
	mapRankingsinitialized bool
	mu                     sync.RWMutex
}

func NewMapRankingDataService(e *env.Env, client *http.Client, rMap repositories.MapRankingRepositoryInterface) *MapRankingDataService {
	service := &MapRankingDataService{
		e:                      e,
		client:                 client,
		rMap:                   rMap,
		mapRankings:            make(map[models.Rank][]*models.MapRanking),
		mapRankingsinitialized: false,
	}
	return service
}

func (s *MapRankingDataService) InitMapRankings() (err error) {
	logger := logrus.WithField("op", "startup")

	_, err = s.RefreshRankings(true)
	if err != nil {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	for _, rank := range models.Ranks() {
		rankings, err := s.rMap.LoadMapRankings(rank)
		if err != nil {
			logger.WithError(err).Warn("failed-to-load-rankings-to-memory")
			return err
		}
		s.mapRankings[rank] = rankings
	}
	s.mapRankingsinitialized = true
	return nil
}

func (s *MapRankingDataService) RefreshRankings(force bool) (updated bool, err error) {
	newDataAvailable, err, newTimeStamp := s.NewDataAvailable()
	if err != nil {
		return
	}
	if !newDataAvailable && !force {
		return
	}
	// TODO: Archive present data when repo methods up

	newRawMapData, err := s.fetchNewRawMapData()
	if err != nil {
		return
	}

	for rank, rawMapData := range newRawMapData {
		if err = s.rMap.NewRawMapData(rank, rawMapData); err != nil {
			return
		}
	}

	newMapRankingData := s.computeMapRankings(newRawMapData)
	for rank, rankingData := range newMapRankingData {
		if err = s.rMap.NewMapRankings(rank, rankingData); err != nil {
			return
		}
	}

	s.mu.Lock()
	maps.Copy(s.mapRankings, newMapRankingData)
	s.mu.Unlock()

	err = s.rMap.UpdateLastUpdatedTime(newTimeStamp)
	if err != nil {
		return
	}

	return true, nil
}

func (s *MapRankingDataService) GetMapRankings(rank models.Rank, mapNames []string) (rankings []*models.MapRanking, err error) {
	if !s.mapRankingsinitialized {
		err = s.InitMapRankings()
		if err != nil {
			return
		}
	}
	for _, mapName := range mapNames {
		r, err := s.getMapRanking(rank, mapName)
		if err != nil {
			return nil, err
		}
		rankings = append(rankings, r)
	}
	return
}

func (s *MapRankingDataService) getMapRanking(rank models.Rank, mapName string) (ranking *models.MapRanking, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, mapRanking := range s.mapRankings[rank] {
		if mapRanking.Name == mapName {
			return mapRanking, nil
		}
	}
	return nil, fmt.Errorf("unknown map name `%s`", mapName)
}

func (s *MapRankingDataService) computeMapRankings(rawMapData map[models.Rank][]*models.MapData) (rankingMapData map[models.Rank][]*models.MapRanking) {
	rankingMapData = make(map[models.Rank][]*models.MapRanking)
	for rank, rawData := range rawMapData {
		for _, raw := range rawData {
			rankingMapData[rank] = append(rankingMapData[rank], s.computeMapRanking(raw))
		}
	}
	return
}

func (s *MapRankingDataService) computeMapRanking(rawMapData *models.MapData) (ranking *models.MapRanking) {
	var (
		winRateByBrawler    []float64
		starRateByBrawler   []float64
		matchCountByBrawler []float64
		brawlers            []*models.BrawlerMapRanking
	)
	winRateWeight := s.e.MapRanking.WinRateWeight
	starRateWeight := 1 - s.e.MapRanking.WinRateWeight

	for _, brawler := range rawMapData.Brawlers {
		winRateByBrawler = append(winRateByBrawler, brawler.WinRate)
		starRateByBrawler = append(starRateByBrawler, brawler.StarRate)
		matchCount := float64(rawMapData.MatchCount) * brawler.UseRate
		matchCountByBrawler = append(matchCountByBrawler, matchCount)
	}

	adjustedWinRateByBrawler, winRateK, winRateMean := stats.BayesianShrinkByVarianceMatching(winRateByBrawler, matchCountByBrawler)
	adjustedStarRateByBrawler, starRateK, starRateMean := stats.BayesianShrinkByVarianceMatching(starRateByBrawler, matchCountByBrawler)

	for i, brawler := range rawMapData.Brawlers {
		score := adjustedWinRateByBrawler[i]*winRateWeight + adjustedStarRateByBrawler[i]*starRateWeight
		brawlerRanking := &models.BrawlerMapRanking{
			Name:             brawler.Name,
			Score:            math.Round(score*100) / 100,
			AdjustedWinRate:  math.Round(adjustedWinRateByBrawler[i]*100) / 100,
			AdjustedStarRate: math.Round(adjustedStarRateByBrawler[i]*100) / 100,
		}
		brawlers = append(brawlers, brawlerRanking)
	}
	ranking = &models.MapRanking{
		Name:         rawMapData.MapName,
		Mode:         rawMapData.Mode,
		Brawlers:     brawlers,
		WinRateK:     winRateK,
		WinRateMean:  winRateMean,
		StarRateK:    starRateK,
		StarRateMean: starRateMean,
	}
	return
}

func (s *MapRankingDataService) fetchNewRawMapData() (mapData map[models.Rank][]*models.MapData, err error) {
	mapData = make(map[models.Rank][]*models.MapData)
	rawData := make(map[models.Rank]*upstream.RawMapDataResponse)

	for _, rank := range models.Ranks() {
		req, err := s.newUpstreamRawMapDataRequest(rank)
		if err != nil {
			return nil, err
		}
		resp, err := s.client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("unexpected status code %d for rank %s", resp.StatusCode, rank)
		}

		tmp := make(upstream.RawMapDataResponse)
		if err := json.NewDecoder(resp.Body).Decode(&tmp); err != nil {
			return nil, err
		}
		rawData[rank] = &tmp
	}

	for rank, data := range rawData {
		for _, rawMapData := range *data {

			tmp := &models.MapData{
				MapName:    rawMapData.MapName,
				Mode:       rawMapData.Mode,
				MatchCount: rawMapData.MatchCount,
			}

			for _, brawlerRawData := range rawMapData.Brawlers {
				if brawlerRawData.UseRate == 0 {
					continue
				}
				tmp.Brawlers = append(tmp.Brawlers,
					&models.BrawlerMapData{
						Name:     brawlerRawData.Name,
						WinRate:  brawlerRawData.WinRate,
						UseRate:  brawlerRawData.UseRate,
						StarRate: brawlerRawData.StarRate,
					},
				)
			}
			mapData[rank] = append(mapData[rank], tmp)
		}
	}
	return mapData, nil
}

func (s *MapRankingDataService) NewDataAvailable() (newData bool, err error, newTime int) {
	previousUpdated, err := s.fetchLastUpdatedTime()
	if err != nil {
		return
	}
	newDataTimestamp, err := s.fetchNewDataTimestamp()
	if err != nil {
		return
	}
	return newDataTimestamp > previousUpdated, nil, newDataTimestamp
}

func (s *MapRankingDataService) FetchRawMapData(rank models.Rank) (data []*models.MapData, err error) {
	return s.rMap.LoadRawMapData(rank)
}

func (s *MapRankingDataService) fetchNewDataTimestamp() (unixTime int, err error) {
	var (
		lastUpdated = &upstream.LastUpdated{}
	)
	req, err := s.newUpstreamDataTimestampRequest()
	if err != nil {
		return -1, err
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return
	}
	if err := json.NewDecoder(resp.Body).Decode(lastUpdated); err != nil {
		return -1, err
	}

	return lastUpdated.Time, nil
}

func (s *MapRankingDataService) fetchLastUpdatedTime() (unixTime int, err error) {
	return s.rMap.LoadLastUpdatedTime()
}

func (s *MapRankingDataService) newUpstreamRawMapDataRequest(rank models.Rank) (req *http.Request, err error) {
	endpoint := fmt.Sprintf(s.e.Upstream.MatchData.BasePath, s.e.Upstream.MatchData.Endpoints[string(rank)])
	return http.NewRequest(http.MethodGet, endpoint, nil)
}

func (s *MapRankingDataService) newUpstreamDataTimestampRequest() (req *http.Request, err error) {
	return http.NewRequest(http.MethodGet, s.e.Upstream.MatchData.LastUpdatedEndpoint, nil)
}
