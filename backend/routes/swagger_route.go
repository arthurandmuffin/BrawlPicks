package routes

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"BrawlPicks/docs"
)

type SwaggerRoute struct {
	apiPrefix string
}

func NewSwaggerRoute(apiPrefix string) *SwaggerRoute {
	return &SwaggerRoute{
		apiPrefix: apiPrefix,
	}
}

func (r *SwaggerRoute) Setup(g *gin.RouterGroup) {
	docs.SwaggerInfo.BasePath = r.apiPrefix
	g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
