package models

type Event struct {
	ID   int
	Name string
}

type BrawlerWinrate struct {
	Name    string
	Winrate float64
}

type EventWinrate struct {
	Event    Event
	Brawlers []BrawlerWinrate
}
