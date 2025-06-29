package services

import (
	"BrawlPicks/models"
	//"BrawlPicks/repositories"
)

type WinrateService interface {
	GetTopWinrates() ([]models.EventWinrate, error)
}

/*
type winrateService struct {
	repo repositories.WinrateRepository
}

func NewWinrateService(repo repositories.WinrateRepository) WinrateService {
	return &winrateService{repo: repo}
}

func (s *winrateService) GetTopWinrates() ([]models.EventWinrate, error) {
	return s.repo.GetEventWinrates()
}
*/
