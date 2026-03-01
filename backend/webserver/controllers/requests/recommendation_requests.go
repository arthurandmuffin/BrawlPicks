package requests

type RecommendBrawlers struct {
	MapName           string `json:"mapName" binding:"required"`
	Mode              string `json:"mode" binding:"required"`
	Rank              int    `json:"rank" binding:"required"`
	AllyBrawlers      []int  `json:"allyBrawlers"`
	EnemyBrawlers     []int  `json:"enemyBrawlers"`
	CandidateBrawlers []int  `json:"candidateBrawlers"`
	BannedBrawlers    []int  `json:"bannedBrawlers"`
	TopK              *int   `json:"topK"`
}
