package services

import (
	"BrawlPicks/models"
	repository "BrawlPicks/repositories"
)

type BrawlStarsService struct {
	repo *repository.BrawlStarsRepository
}

func NewBrawlStarsService(repo *repository.BrawlStarsRepository) *BrawlStarsService {
	return &BrawlStarsService{repo: repo}
}

func (s *BrawlStarsService) GetPlayer(tag string) (map[string]interface{}, error) {
	return s.repo.GetPlayer(tag)
}

func (s *BrawlStarsService) GetPower11Brawlers(tag string) ([]models.Brawler, error) {
	player, err := s.repo.GetPlayerData(tag)
	if err != nil {
		return nil, err
	}

	var power11 []models.Brawler
	for _, b := range player.Brawlers {
		if b.Power == 11 {
			power11 = append(power11, b)
		}
	}
	return power11, nil
}

/*
type Service struct {
	repo *repository.Repository
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		repo: repo,
	}
}


func (s *Service) ServiceTest() string {
	return s.repo.RepositoryTest()
}
*/
