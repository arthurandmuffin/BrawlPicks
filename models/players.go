package models

type PlayerResponse struct {
	Brawlers []Brawler `json:"brawlers"`
}

type Brawler struct {
	Name  string `json:"name"`
	Power int    `json:"power"`
}
