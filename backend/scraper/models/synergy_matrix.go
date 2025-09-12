package models

type Count struct {
	Wins  float64 `json:"wins"`
	Total float64 `json:"total"`
}

type SynergyMatrix struct {
	Synergy map[int]map[int]Count `json:"synergy"`
	Counter map[int]map[int]Count `json:"counter"`
}

func NewSynergyMatrix() *SynergyMatrix {
	return &SynergyMatrix{
		Synergy: make(map[int]map[int]Count),
		Counter: make(map[int]map[int]Count),
	}
}

func (s *SynergyMatrix) IncrementSynergy(ids [2]int, winning bool, draw bool) {
	s.ensureCountMap(ids)
	cell1 := s.Synergy[ids[0]][ids[1]]
	cell2 := s.Synergy[ids[1]][ids[0]]

	cell1.Total++
	cell2.Total++
	if winning {
		cell1.Wins++
		cell2.Wins++
	} else if draw {
		cell1.Wins += 0.5
		cell2.Wins += 0.5
	}

	s.Synergy[ids[0]][ids[1]] = cell1
	s.Synergy[ids[1]][ids[0]] = cell2
}

func (s *SynergyMatrix) IncrementCounter(ids [2]int, draw bool) {
	s.ensureCountMap(ids)
	cellW := s.Counter[ids[0]][ids[1]]
	cellL := s.Counter[ids[1]][ids[0]]

	cellW.Total++
	cellL.Total++
	if draw {
		cellW.Wins += 0.5
		cellL.Wins += 0.5
	} else {
		cellW.Wins++
	}

	s.Counter[ids[0]][ids[1]] = cellW
	s.Counter[ids[1]][ids[0]] = cellL
}

func (s *SynergyMatrix) ensureCountMap(ids [2]int) {
	if s.Synergy[ids[0]] == nil {
		s.Synergy[ids[0]] = make(map[int]Count)
	}
	if s.Synergy[ids[1]] == nil {
		s.Synergy[ids[1]] = make(map[int]Count)
	}
	if s.Counter[ids[0]] == nil {
		s.Counter[ids[0]] = make(map[int]Count)
	}
	if s.Counter[ids[1]] == nil {
		s.Counter[ids[1]] = make(map[int]Count)
	}
}
