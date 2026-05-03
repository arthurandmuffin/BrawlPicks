package responses

type RecommendBrawlers struct {
	ModelID         string                   `json:"modelId"`
	Recommendations []*BrawlerRecommendation `json:"recommendations"`
}

type BrawlerRecommendation struct {
	BrawlerID int     `json:"brawlerId"`
	Score     float64 `json:"score"`
}
