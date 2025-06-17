package api

import (
	"github.com/gin-contrib/cors"
)

var (
	CorsMiddleware = cors.New(cors.Config{
		AllowAllOrigins:        true,
		AllowCredentials:       true,
		AllowBrowserExtensions: true,
		AllowMethods:           []string{"GET", "POST", "DELETE", "OPTIONS", "PUT"},
		AllowHeaders: []string{
			"Authorization",
			"Content-Type",
			"Upgrade",
			"Origin",
			"Connection",
			"Accept-Encoding",
			"Accept-Language",
			"Host",
			"Access-Control-Request-Method",
			"Access-Control-Request-Headers",
		},
	})
)
