package routes

import (
	"BrawlPicks/webserver/controllers"
	"BrawlPicks/webserver/internal/api"

	"github.com/gin-gonic/gin"
)

type RecommendationRoute struct {
	cl *controllers.RecommendationController
}

func NewRecommendationRoute(cl *controllers.RecommendationController) *RecommendationRoute {
	return &RecommendationRoute{
		cl: cl,
	}
}

func (r *RecommendationRoute) Setup(g *gin.RouterGroup) {
	recommendations := g.Group("/recommendations")
	{
		recommendations.POST("", api.JsonHandler(r.cl.RecommendBrawlers))
	}
}
