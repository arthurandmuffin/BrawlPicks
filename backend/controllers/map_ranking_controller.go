package controllers

import (
	"BrawlPicks/controllers/requests"
	"BrawlPicks/internal/api"
	"BrawlPicks/internal/e"
	"BrawlPicks/models"
	"BrawlPicks/services"
	"slices"

	"github.com/sirupsen/logrus"
)

type MapRankingController struct {
	svMapRanking services.MapRankingDataServiceInterface
}

func NewMapRankingController(svMapRanking services.MapRankingDataServiceInterface) *MapRankingController {
	return &MapRankingController{
		svMapRanking: svMapRanking,
	}
}

// GetMapRanking godoc
// @Summary Get map ranking for a specific map and rank
// @Description Returns the ranking of brawlers for a given map and rank
// @Tags map-rankings
// @Param rank path string true "Rank" Enums(d1, m1, m3, l1)
// @Param mapName path string true "Map Name"
// @Success 200 {object} models.MapRanking
// @Failure 400 {object} api.Response
// @Failure 500 {object} api.Response
// @Router /map-rankings/{rank}/{mapName} [get]
func (cl *MapRankingController) GetMapRanking(c *api.Context, r *requests.GetMapRanking) {
	logger := logrus.WithFields(logrus.Fields{
		"rank":    r.Rank,
		"mapName": r.MapName,
	})

	mapName := []string{r.MapName}
	mapRankings, err := cl.svMapRanking.GetMapRankings(models.Rank(r.Rank), mapName)
	if err != nil {
		logger.WithError(err).Warn("failed-ranking-api-call")
		c.InternalServerError(err)
		return
	} else if len(mapRankings) < 1 {
		logger.WithError(e.ErrEmptyResponse).Warn("failed-ranking-api-call")
		c.InternalServerError(e.ErrEmptyResponse)
		return
	}

	res := mapRankings[0]
	slices.SortStableFunc(res.Brawlers, func(a, b *models.BrawlerMapRanking) int {
		switch {
		case a.Score < b.Score:
			return 1
		case a.Score > b.Score:
			return -1
		default:
			return 0
		}
	})
	c.OK(res)
}
