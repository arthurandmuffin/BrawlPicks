package main

import (
	"BrawlPicks/wire"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	handler := wire.InitializeHandler()

	r.GET("/test", handler.HandlerTest)
	log.Println("Server running on port 8080")
	r.Run(":8080")
}
