package controllers

import (
	"BrawlPicks/webserver/controllers/requests"
	"BrawlPicks/webserver/inference"
	"BrawlPicks/webserver/internal/api"
	"BrawlPicks/webserver/services"

	"github.com/sirupsen/logrus"
)

type RecommendationController struct {
	svRecommendation services.RecommendationServiceInterface
}

func NewRecommendationController(svRecommendation services.RecommendationServiceInterface) *RecommendationController {
	return &RecommendationController{
		svRecommendation: svRecommendation,
	}
}

// RecommendBrawlers godoc
// @Summary Recommend brawlers for a partial draft
// @Description Proxies a recommendation request through the ML inference service
// @Tags recommendations
// @Accept json
// @Produce json
// @Param request body requests.RecommendBrawlers true "Recommendation request"
// @Success 200 {object} api.Response
// @Failure 400 {object} api.Response
// @Failure 500 {object} api.Response
// @Router /recommendations [post]
func (cl *RecommendationController) RecommendBrawlers(c *api.Context, r *requests.RecommendBrawlers) {
	logger := logrus.WithFields(logrus.Fields{
		"mapName": r.MapName,
		"mode":    r.Mode,
		"rank":    r.Rank,
	})

	res, err := cl.svRecommendation.Recommend(
		c.Request.Context(),
		r.MapName,
		r.Mode,
		r.Rank,
		r.AllyBrawlers,
		r.EnemyBrawlers,
		r.CandidateBrawlers,
		r.BannedBrawlers,
		r.TopK,
	)
	if err != nil {
		logger.WithError(err).Warn("failed-recommendation-api-call")
		if inference.IsBadRequest(err) {
			c.BadRequest(err)
			return
		}
		c.InternalServerError(err)
		return
	}

	c.OK(res)
}
