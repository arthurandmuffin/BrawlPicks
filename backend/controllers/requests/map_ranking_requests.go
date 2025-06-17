package requests

type GetMapRanking struct {
	Rank    string `uri:"rank" binding:"required,oneof=d1 m1 m3 l1"`
	MapName string `uri:"mapName" binding:"required"`
}
