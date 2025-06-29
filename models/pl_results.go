package models

type PLResults map[string]MapData

type MapData struct {
	Individual []BrawlerStats `json:"individual"`
}

type BrawlerStats struct {
	Brawler string  `json:"brawler"` // Name or ID
	WR      float64 `json:"wr"`      // Win Rate
	UR      float64 `json:"ur"`      // Use Rate
	SR      float64 `json:"sr"`      // Star Rate
}
