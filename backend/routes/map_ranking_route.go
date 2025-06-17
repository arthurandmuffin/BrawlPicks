package routes

import (
	"BrawlPicks/controllers"
	"BrawlPicks/internal/api"

	"github.com/gin-gonic/gin"
)

type MapRankingRoute struct {
	cl *controllers.MapRankingController
}

func NewMapRankingRoute(cl *controllers.MapRankingController) *MapRankingRoute {
	return &MapRankingRoute{
		cl: cl,
	}
}

func (r *MapRankingRoute) Setup(g *gin.RouterGroup) {
	mapRankings := g.Group("/map-rankings")
	{
		mapRankings.GET("/:rank/:mapName", api.UriHandler(r.cl.GetMapRanking))
	}
}
