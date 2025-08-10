package upstream

type TopPlayersResponse struct {
	Players []TopPlayer `json:"items"`
}

type TopPlayer struct {
	Tag string `json:"tag"`
}
