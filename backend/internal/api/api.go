package api

import (
	"context"

	"github.com/gin-gonic/gin"
)

type Api struct {
	err            chan error
	ctx            context.Context
	port           string
	prefix         string
	cors           bool
	debug          bool
	baseRoutes     []Route
	prefixedRoutes []Route
}

func NewApi(
	ctx context.Context,
	port,
	prefix string,
	cors bool,
	debug bool,
	baseRoutes []Route,
	prefixedRoutes []Route,
) *Api {
	return &Api{
		make(chan error),
		ctx,
		port,
		prefix,
		cors,
		debug,
		baseRoutes,
		prefixedRoutes,
	}
}

func (a *Api) Run() {
	if !a.debug {
		gin.SetMode(gin.DebugMode)
	}

	engine := gin.Default()
	if a.cors {
		engine.Use(CorsMiddleware)
	}

	for _, route := range a.baseRoutes {
		route.Setup(&engine.RouterGroup)
	}

	group := engine.Group(a.prefix)
	for _, route := range a.prefixedRoutes {
		route.Setup(group)
	}

	go func() {
		a.err <- engine.Run(a.port)
	}()

	select {
	case <-a.err:
	case <-a.ctx.Done():
	}
}
