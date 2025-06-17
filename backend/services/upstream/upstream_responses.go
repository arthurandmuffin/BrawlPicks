package upstream

type RawMapDataResponse map[string]*RawMapData

type RawMapData struct {
	Brawlers   []*BrawlerRawMapData `json:"individual"`
	MapName    string               `json:"brawlify_map_id"`
	Mode       string               `json:"mode"`
	MatchCount int                  `json:"match_count"`
}

type BrawlerRawMapData struct {
	Name     string  `json:"brawler"`
	WinRate  float64 `json:"wr"`
	UseRate  float64 `json:"ur"`
	StarRate float64 `json:"sr"`
}

type LastUpdated struct {
	Time int `json:"last_updated"`
}
