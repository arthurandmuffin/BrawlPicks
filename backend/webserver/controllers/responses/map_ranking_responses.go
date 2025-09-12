package responses

type GetMapRanking struct {
	Rank     string               `json:"rank"`
	MapName  string               `json:"mapName"`
	Mode     string               `json:"mode"`
	Brawlers []*BrawlerMapRanking `json:"brawlers"`
}

type BrawlerMapRanking struct {
	Name     string  `json:"name"`
	Score    float64 `json:"score"`
	WinRate  float64 `json:"winRate"`
	StarRate float64 `json:"starRate"`
}
