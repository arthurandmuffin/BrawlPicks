package models

type Rank string

const (
	Diamond   Rank = "d1"
	Mythic1   Rank = "m1"
	Mythic3   Rank = "m3"
	Legendary Rank = "l1"
)

func Ranks() []Rank {
	return []Rank{Diamond, Mythic1, Mythic3, Legendary}
}

type MapData struct {
	MapName    string            `json:"map_name"`
	Mode       string            `json:"mode"`
	Brawlers   []*BrawlerMapData `json:"brawlers"`
	MatchCount int               `json:"match_count"`
}

type BrawlerMapData struct {
	Name     string  `json:"brawler"`
	WinRate  float64 `json:"wr"`
	UseRate  float64 `json:"ur"`
	StarRate float64 `json:"sr"`
}

type MapRanking struct {
	Name         string               `json:"map_name"`
	Mode         string               `json:"mode"`
	Brawlers     []*BrawlerMapRanking `json:"brawlers"`
	WinRateK     float64              `json:"winRateK"`
	WinRateMean  float64              `json:"winRateMean"`
	StarRateK    float64              `json:"starRateK"`
	StarRateMean float64              `json:"starRateMean"`
}

type BrawlerMapRanking struct {
	Name             string  `json:"brawler"`
	Score            float64 `json:"score"`
	AdjustedWinRate  float64 `json:"adjusted_wr"`
	AdjustedStarRate float64 `json:"adjusted_sr"`
}

type LastUpdated struct {
	Time int `json:"last_updated"`
}
