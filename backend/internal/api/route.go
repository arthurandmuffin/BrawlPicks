package api

import (
	"github.com/gin-gonic/gin"
)

type Route interface {
	Setup(g *gin.RouterGroup)
}
