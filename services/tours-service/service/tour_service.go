package service

import (
	"tours-service/domain"
	"tours-service/repository"
)

type TourService interface {
	Create(tour *domain.Tour) error
	GetByAuthorId(authorId string) ([]*domain.Tour, error) // <-- NOVA METODA
}

type tourService struct {
	repo repository.TourRepository
}

func NewTourService(repo repository.TourRepository) TourService {
	return &tourService{repo: repo}
}

func (s *tourService) Create(tour *domain.Tour) error {
	tour.Status = "draft"
	tour.Price = 0.0
	return s.repo.Create(tour)
}

// Implementacija nove metode
func (s *tourService) GetByAuthorId(authorId string) ([]*domain.Tour, error) {
	return s.repo.GetByAuthorId(authorId)
}
