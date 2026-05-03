package models

import (
	"time"
)

type Battle struct {
	Timestamp time.Time
	MapName   string
	Mode      string
	Rank      int
	Draw      bool
	TeamW     []int
	TeamL     []int
}

func (b *Battle) SynergyCombinations(winning bool) [][2]int {
	var team []int
	switch winning {
	case true:
		team = b.TeamW
	default:
		team = b.TeamL
	}
	out := make([][2]int, 0, 3)
	for i := 0; i < len(team); i++ {
		for j := i + 1; j < len(team); j++ {
			out = append(out, [2]int{team[i], team[j]})
		}
	}
	return out
}

func (b *Battle) CounterCombinations() [][2]int {
	out := make([][2]int, 0, 9)
	for _, wid := range b.TeamW {
		for _, lid := range b.TeamL {
			out = append(out, [2]int{wid, lid})
		}
	}
	return out
}
