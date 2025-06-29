package main

import (
	"BrawlPicks/wire"
	"log"

	"github.com/gin-gonic/gin"
)

type BrawlerStat struct {
	Brawler string  `json:"brawler"`
	WR      float64 `json:"wr"`
	UR      float64 `json:"ur"`
	SR      float64 `json:"sr"`
}

type MapStat struct {
	Individual []BrawlerStat `json:"individual"`
}

func main() {

	apiKey := "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzUxMiIsImtpZCI6IjI4YTMxOGY3LTAwMDAtYTFlYi03ZmExLTJjNzQzM2M2Y2NhNSJ9.eyJpc3MiOiJzdXBlcmNlbGwiLCJhdWQiOiJzdXBlcmNlbGw6Z2FtZWFwaSIsImp0aSI6ImYzOTA3ZjEzLWM1NzAtNDViMS1iYjYwLWRlMzZmNjA4ZWM5MiIsImlhdCI6MTczODIwMjIxNiwic3ViIjoiZGV2ZWxvcGVyLzRhNmMzMjcyLTIzYWMtZDIxYi0zY2NlLTUzYzkxNDNkYjAxNCIsInNjb3BlcyI6WyJicmF3bHN0YXJzIl0sImxpbWl0cyI6W3sidGllciI6ImRldmVsb3Blci9zaWx2ZXIiLCJ0eXBlIjoidGhyb3R0bGluZyJ9LHsiY2lkcnMiOlsiMTczLjE3Ni4xMzcuMTA1IiwiNTQuODYuNTAuMTM5Il0sInR5cGUiOiJjbGllbnQifV19.ptYJLGyJKL6TMsCK_I-fAT6shsPI-Uf6QUuqXwP5P_dz4pKQOry-khaKxEYUOrKkdhMdjk5lKRee96RtneVtzA"
	if apiKey == "" {
		log.Fatal("Missing BRAWLSTARS_API_KEY environment variable")
	}

	handler := wire.InitializeBrawlStarsHandler(apiKey)
	plHandler := wire.InitializePLHandler()

	r := gin.Default()
	r.GET("/player/:tag", handler.GetPlayer)
	r.GET("/power11/:tag", handler.GetPower11Brawlers)
	r.GET("/winrates", plHandler.GetTopWinrates)

	log.Println("Server running on port 8080")

	r.Run(":8080")

}
