package services

import repository "BrawlPicks/repositories"

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
