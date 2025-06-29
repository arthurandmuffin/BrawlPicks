package services

import (
	"BrawlPicks/models"
	"BrawlPicks/repositories"
	"math"
	"sort"
)

type PLService struct {
	repo *repositories.PLRepository
}

func NewPLService(repo *repositories.PLRepository) *PLService {
	return &PLService{repo: repo}
}

// Adjusts WR upward based on use rate (UR), using log1p(3) base
func weight(ur float64) float64 {
	return 0.9 + 0.1*math.Log1p(ur)/math.Log1p(3)
}

// Returns winrate shrunk toward the prior based on use rate
func bayesianAdjustedWR(wr, ur, prior, m float64) float64 {
	return (wr*ur + prior*m) / (ur + m)
}

type weightedBrawler struct {
	models.BrawlerStats
	Score float64
}

func (s *PLService) GetTop5BrawlersByMap() (map[string][]models.BrawlerStats, error) {
	results, err := s.repo.FetchPLResults()
	if err != nil {
		return nil, err
	}

	topBrawlers := make(map[string][]models.BrawlerStats)

	for mapName, mapData := range results {
		var scored []weightedBrawler
		const priorWR = 50.0     // average winrate
		const priorWeight = 10.0 // adjust this to tune aggressiveness

		for _, b := range mapData.Individual {
			if b.UR < 2.0 {
				continue // Skip unreliable picks
			}

			// Apply Bayesian shrink
			adjustedWR := bayesianAdjustedWR(b.WR, b.UR, priorWR, priorWeight)

			// You can optionally apply an additional weight boost here based on UR
			// e.g., score := adjustedWR * weight(b.UR)
			score := adjustedWR

			scored = append(scored, weightedBrawler{
				BrawlerStats: b,
				Score:        score,
			})
		}

		// Sort by adjusted score
		sort.Slice(scored, func(i, j int) bool {
			return scored[i].Score > scored[j].Score
		})

		// Take top 5
		top := []models.BrawlerStats{}
		for i := 0; i < len(scored) && i < 5; i++ {
			top = append(top, scored[i].BrawlerStats)
		}

		topBrawlers[mapName] = top
	}

	return topBrawlers, nil
}
